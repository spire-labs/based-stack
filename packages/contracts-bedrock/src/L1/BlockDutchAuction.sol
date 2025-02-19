// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";
import { ICrossDomainMessenger } from "src/universal/interfaces/ICrossDomainMessenger.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { ElectionTickets } from "src/L2/ElectionTickets.sol";
import "src/libraries/BlockAuctionErrors.sol";

contract BlockDutchAuction is Ownable {
    /// @notice The minimum number of validators in the look ahead
    uint256 public constant VALIDATORS_IN_LOOKAHEAD = 32;

    /// @notice The address of the ElectionTicket contract on L2
    address public constant ELECTION_TICKET = 0x4200000000000000000000000000000000000028;

    /// @notice The address of the SystemConfig contract on L1
    SystemConfig public immutable SYSTEM_CONFIG;

    /// @notice The start block of the current running auction
    uint216 public startBlock;

    /// @notice The duration the auction will run for in blocks
    uint8 public durationBlocks;

    /// @notice The pending new block duration of the auction
    uint8 public pendingDurationBlocks;

    /// @notice The discount rate of the auction
    uint8 public discountRate;

    /// @notice The pending new discount rate of the auction
    uint8 public pendingDiscountRate;

    /// @notice The number of tickets left in the current auction
    uint8 internal _ticketsLeft;

    /// @notice The starting price of the auction
    uint256 public startPrice;

    /// @notice The pending new start price of the auction
    uint256 public pendingStartPrice;

    /// @notice Emitted when a new pending start price is set
    ///
    /// @param newStartPrice The new start price
    event PendingStartPriceSet(uint256 newStartPrice);

    /// @notice Emitted when a new pending discount rate is
    ///
    /// @param newDiscountRate The new discount rate
    event PendingDiscountRateSet(uint8 newDiscountRate);

    /// @notice Emitted when a new pending duration blocks is set
    /// @param newDurationBlocks The new duration blocks
    event PendingDurationBlocksSet(uint8 newDurationBlocks);

    /// @notice Emitted when a ticket is bought
    ///
    /// @param buyer The address of the buyer
    /// @param startBlock The start block of the current auction the ticket was bought in
    /// @param price The price of the ticket
    event TicketBought(address indexed buyer, uint256 indexed startBlock, uint256 price, uint8 ticketsLeft);

    /// @notice Constructs and initilizes the BlockDutchAuction
    ///
    /// @param _startBlock The start block of the auction
    /// @param _durationBlocks The duration of the auction in blocks
    /// @param _startPrice The starting price of the auction
    /// @param _discountRate The discount rate of the auction
    /// @param _systemConfig The address of the SystemConfig contract on L1
    constructor(
        uint216 _startBlock,
        uint8 _durationBlocks,
        uint256 _startPrice,
        uint8 _discountRate,
        SystemConfig _systemConfig
    ) {
        if (_discountRate >= 100 || _discountRate == 0) revert InvalidDiscountRate();
        if (_durationBlocks > VALIDATORS_IN_LOOKAHEAD) revert InvalidBlockDuration();
        if (_startPrice < 1e3 || _startPrice > type(uint256).max / 100) revert InvalidStartPrice();

        startBlock = _startBlock;
        durationBlocks = _durationBlocks;
        startPrice = _startPrice;
        discountRate = _discountRate;
        SYSTEM_CONFIG = _systemConfig;

        // Set the starting amount of tickets
        _ticketsLeft = _durationBlocks;
    }

    /*///////////////////////////////////////////////////////////////
                              ADMIN FUNCTIONS
    ///////////////////////////////////////////////////////////////*/

    /// @notice Sets the starting price of the auction
    ///
    /// @param _newStartPrice The new starting price
    function setStartPrice(uint256 _newStartPrice) external onlyOwner {
        if (_newStartPrice < 1e3 || _newStartPrice > type(uint256).max / 100) revert InvalidStartPrice();

        pendingStartPrice = _newStartPrice;

        emit PendingStartPriceSet(_newStartPrice);
    }

    /// @notice Sets the discount rate of the auction
    ///
    /// @param _newDiscountRate The new discount rate
    function setDiscountRate(uint8 _newDiscountRate) external onlyOwner {
        if (_newDiscountRate >= 100 || _newDiscountRate == 0) revert InvalidDiscountRate();

        pendingDiscountRate = _newDiscountRate;

        emit PendingDiscountRateSet(_newDiscountRate);
    }

    /// @notice Sets the block duration of the auction
    ///
    /// @param _newDurationBlocks The new block duration
    function setDurationBlocks(uint8 _newDurationBlocks) external onlyOwner {
        if (_newDurationBlocks > VALIDATORS_IN_LOOKAHEAD) revert InvalidBlockDuration();

        pendingDurationBlocks = _newDurationBlocks;

        emit PendingDurationBlocksSet(_newDurationBlocks);
    }

    /*///////////////////////////////////////////////////////////////
                              CORE LOGIC
    ///////////////////////////////////////////////////////////////*/

    /// @notice Buys a ticket in the auction
    ///
    /// @param _amount The amount of tickets to buy
    function buy(uint8 _amount) external payable {
        // TODO: Batch minting/buying
        uint256 _startBlock = startBlock;
        uint256 _durationBlocks = durationBlocks;

        // If the current block number is past the end of the auction
        // This branch will trigger the start of a new auction
        // View functions have this logic baked in to calculate based on what the auction would be at
        if (block.number > _startBlock + _durationBlocks) {
            // Need to define this up here incase duration blocks is changed
            uint256 _predictedEndBlock = _startBlock + _durationBlocks;

            // If there is a change in duration, we need to update this value before finding the start block
            if (pendingDurationBlocks != 0) {
                _durationBlocks = pendingDurationBlocks;
                durationBlocks = uint8(_durationBlocks);
                pendingDurationBlocks = 0;
            }

            // Find the start block of the next auction
            _startBlock = _findStartBlock(_predictedEndBlock, _durationBlocks);
            startBlock = uint216(_startBlock);

            if (pendingStartPrice != 0) {
                startPrice = pendingStartPrice;
                pendingStartPrice = 0;
            }

            if (pendingDiscountRate != 0) {
                discountRate = pendingDiscountRate;
                pendingDiscountRate = 0;
            }

            // Amount of tickets sold in an auction should always equal the block duration
            _ticketsLeft = uint8(_durationBlocks);
        }

        // This check needs to come second incase its a new auction
        uint8 __ticketsLeft = _ticketsLeft;
        if (__ticketsLeft < _amount) revert NoTicketsLeft();

        uint256 _price = _getPrice(uint256(discountRate), startPrice, _startBlock);

        if (_price > msg.value) {
            revert InsufficientFunds();
        } else if (msg.value > _price) {
            unchecked {
                (bool _success,) = payable(msg.sender).call{ value: msg.value - _price }("");
                if (!_success) revert FailedLowLevelCall();
            }
        }

        // Send the message to L2 to mint the ticket
        ICrossDomainMessenger _messenger = ICrossDomainMessenger(SYSTEM_CONFIG.l1CrossDomainMessenger());

        // Mint theoretical max gas limit is ~130_000
        // We send as 150_000 to allow for some buffer
        // TODO(spire): A potential improvement here is to allow for mass minting in one call
        //              This would require a much larger refactor, can be handled as an optimization in a future PR
        for (uint256 i; i < _amount; i++) {
            _messenger.sendMessage(ELECTION_TICKET, abi.encodeCall(ElectionTickets.mint, (msg.sender)), 150_000);
        }

        unchecked {
            __ticketsLeft -= _amount;
            _ticketsLeft = __ticketsLeft;
        }

        emit TicketBought(msg.sender, _startBlock, _price, __ticketsLeft);
    }

    /*///////////////////////////////////////////////////////////////
                              VIEW FUNCTIONS
    ///////////////////////////////////////////////////////////////*/

    /// @notice Returns the current price of a ticket in the auction
    ///
    /// @return price_ The current price of a ticket
    function getPrice() external view returns (uint256 price_) {
        uint256 _startBlock = startBlock;
        uint256 _durationBlocks = durationBlocks;
        uint256 _startPrice = startPrice;
        uint256 _discountRate = discountRate;

        if (block.number > _startBlock + _durationBlocks) {
            // Need to define this up here incase duration blocks is changed
            uint256 _predictedEndBlock = _startBlock + _durationBlocks;

            // If there is a change in duration, we need to update this value before finding the start block
            if (pendingDurationBlocks != 0) {
                _durationBlocks = pendingDurationBlocks;
            }

            // Find the expected start block
            _startBlock = _findStartBlock(_predictedEndBlock, _durationBlocks);

            // Use the pending values to calculate the price
            if (pendingDiscountRate != 0) {
                _discountRate = pendingDiscountRate;
            }

            if (pendingStartPrice != 0) {
                _startPrice = pendingStartPrice;
            }
        }

        price_ = _getPrice(_discountRate, _startPrice, _startBlock);
    }

    /// @notice Returns the amount of tickets left in the current auction
    ///
    /// @return amount_ The amount of tickets left in the current auction
    function ticketsLeft() external view returns (uint256 amount_) {
        uint256 _durationBlocks = durationBlocks;
        if (block.number > startBlock + _durationBlocks) {
            amount_ = _durationBlocks;
        } else {
            amount_ = _ticketsLeft;
        }
    }

    /// @notice Returns the current price of a ticket in the auction
    ///
    /// @param _discountRate The discount rate of the auction
    /// @param _startPrice The starting price of the auction
    /// @param _startBlock The start block of the auction
    ///
    /// @return price_ The current price of a ticket
    function _getPrice(
        uint256 _discountRate,
        uint256 _startPrice,
        uint256 _startBlock
    )
        internal
        view
        returns (uint256 price_)
    {
        uint256 blocksSinceStart = block.number - _startBlock;
        price_ = _startPrice;

        // NOTE: Is this the math we want to go with?
        for (uint256 i; i < blocksSinceStart; i++) {
            price_ = price_ - ((price_ * _discountRate) / 100);
        }
    }

    /// @notice Recursively finds the start block of the next auction
    ///
    /// @param _predictedEndBlock The predicted end block of the auction
    /// @param _durationBlocks The duration of the auction in blocks
    ///
    /// @return newStartBlock_ The start block of the next auction
    function _findStartBlock(
        uint256 _predictedEndBlock,
        uint256 _durationBlocks
    )
        internal
        view
        returns (uint256 newStartBlock_)
    {
        uint256 _difference = block.number - _predictedEndBlock;
        uint256 _result = _difference % (_durationBlocks);

        if (_result == 0) {
            newStartBlock_ = block.number;
        } else {
            newStartBlock_ = block.number - _result + 1;
        }
    }
}
