// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";
import { ElectionTickets } from "src/L1/ElectionTickets.sol";
import "src/libraries/BlockAuctionErrors.sol";

abstract contract BlockDutchAuction is Ownable {
    /// @notice The minimum number of validators in the look ahead
    uint256 public constant VALIDATORS_IN_LOOKAHEAD = 32;

    /// @notice The address of the ElectionTicket contract
    ElectionTickets public immutable ELECTION_TICKET;

    /// @notice The start block of the current running auction
    uint224 public startBlock;

    /// @notice The duration the auction will run for in blocks
    uint8 public durationBlocks;

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
    /// @param _newStartPrice The new start price
    event PendingStartPriceSet(uint256 _newStartPrice);

    /// @notice Emitted when a new pending discount rate is set
    /// @param _newDiscountRate The new discount rate
    event PendingDiscountRateSet(uint8 _newDiscountRate);

    /// @notice Emitted when a ticket is bought
    /// @param _buyer The address of the buyer
    /// @param _startBlock The start block of the current auction the ticket was bought in
    /// @param _price The price of the ticket
    event TicketBought(address indexed _buyer, uint256 indexed _startBlock, uint256 _price, uint8 _ticketsLeft);

    /// @notice Constructs and initilizes the BlockDutchAuction
    /// @param _startBlock The start block of the auction
    /// @param _durationBlocks The duration of the auction in blocks
    /// @param _startPrice The starting price of the auction
    /// @param _discountRate The discount rate of the auction
    /// @param _electionTicket The address of the ElectionTicket contract
    constructor(
        uint224 _startBlock,
        uint8 _durationBlocks,
        uint256 _startPrice,
        uint8 _discountRate,
        ElectionTickets _electionTicket
    ) {
        if (_discountRate >= 100 || _discountRate == 0) revert InvalidDiscountRate();
        if (_durationBlocks > VALIDATORS_IN_LOOKAHEAD) revert InvalidBlockDuration();
        if (_startPrice < 1e3 || _startPrice > type(uint256).max / 100) revert InvalidStartPrice();

        startBlock = _startBlock;
        durationBlocks = _durationBlocks;
        startPrice = _startPrice;
        discountRate = _discountRate;
        ELECTION_TICKET = _electionTicket;

        // Set the starting amount of tickets
        _ticketsLeft = _durationBlocks;
    }

    /*///////////////////////////////////////////////////////////////
                              ADMIN FUNCTIONS
    ///////////////////////////////////////////////////////////////*/

    function setStartPrice(uint256 _newStartPrice) external onlyOwner {
        if (_newStartPrice < 1e3 || _newStartPrice > type(uint256).max / 100) revert InvalidStartPrice();

        pendingStartPrice = _newStartPrice;

        emit PendingStartPriceSet(_newStartPrice);
    }

    function setDiscountRate(uint8 _newDiscountRate) external onlyOwner {
        if (_newDiscountRate >= 100 || _newDiscountRate == 0) revert InvalidDiscountRate();

        pendingDiscountRate = _newDiscountRate;

        emit PendingDiscountRateSet(_newDiscountRate);
    }

    /*///////////////////////////////////////////////////////////////
                              CORE LOGIC
    ///////////////////////////////////////////////////////////////*/

    /// @notice Buys a ticket in the auction
    function buy() external payable {
        uint256 _startBlock = startBlock;
        uint256 _durationBlocks = durationBlocks;

        // If the current block number is past the end of the auction
        // This branch will trigger the start of a new auction
        // View functions have this logic baked in to calculate based on what the auction would be at
        if (block.number > _startBlock + _durationBlocks) {
            // Find the start block of the next auction
            _startBlock = _findStartBlock(_startBlock, _durationBlocks);
            startBlock = uint224(_startBlock);

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
        if (__ticketsLeft == 0) revert NoTicketsLeft();

        uint256 _price = _getPrice(uint256(discountRate), startPrice, _startBlock);

        if (_price > msg.value) {
            revert InsufficientFunds();
        } else if (msg.value > _price) {
            (bool _success,) = payable(msg.sender).call{ value: msg.value - _price }("");
            if (!_success) revert FailedLowLevelCall();
        }

        // TODO: Mint ticket

        _ticketsLeft = __ticketsLeft - 1;
        emit TicketBought(msg.sender, _startBlock, _price, __ticketsLeft - 1);
    }

    /*///////////////////////////////////////////////////////////////
                              VIEW FUNCTIONS
    ///////////////////////////////////////////////////////////////*/

    /// @notice Returns the current price of a ticket in the auction
    /// @return _price The current price of a ticket
    function getPrice() external view returns (uint256 _price) {
        uint256 _startBlock = startBlock;
        uint256 _durationBlocks = durationBlocks;
        uint256 _startPrice = startPrice;
        uint256 _discountRate = discountRate;

        if (block.number > _startBlock + _durationBlocks) {
            // Find the expected start block
            _startBlock = _findStartBlock(_startBlock, _durationBlocks);

            // Use the pending values to calculate the price
            if (pendingDiscountRate != 0) {
                _discountRate = pendingDiscountRate;
            }

            if (pendingStartPrice != 0) {
                _startPrice = pendingStartPrice;
            }
        }

        _price = _getPrice(_discountRate, _startPrice, _startBlock);
    }

    /// @notice Returns the amount of tickets left in the current auction
    /// @return _amount The amount of tickets left in the current auction
    function ticketsLeft() external view returns (uint256 _amount) {
        uint256 _durationBlocks = durationBlocks;
        if (block.number > startBlock + _durationBlocks) {
            _amount = _durationBlocks;
        } else {
            _amount = _ticketsLeft;
        }
    }

    /// @notice Returns the current price of a ticket in the auction
    /// @param _discountRate The discount rate of the auction
    /// @param _startPrice The starting price of the auction
    /// @param _startBlock The start block of the auction
    /// @return _price The current price of a ticket
    function _getPrice(
        uint256 _discountRate,
        uint256 _startPrice,
        uint256 _startBlock
    )
        internal
        view
        returns (uint256 _price)
    {
        uint256 blocksSinceStart = block.number - _startBlock;
        _price = _startPrice;

        // NOTE: Is this the math we want to go with?
        for (uint256 i; i < blocksSinceStart; i++) {
            _price = _price - ((_price * _discountRate) / 100);
        }
    }

    /// @notice Recursively finds the start block of the next auction
    /// @param _currentStartBlock The current start block of the auction
    /// @param _durationBlocks The duration of the auction in blocks
    /// @return _newStartBlock The start block of the next auction
    function _findStartBlock(
        uint256 _currentStartBlock,
        uint256 _durationBlocks
    )
        internal
        view
        returns (uint256 _newStartBlock)
    {
        uint256 _predictedEndBlock = _currentStartBlock + _durationBlocks;
        uint256 _difference = block.number - _predictedEndBlock;
        uint256 _result = _difference % (_durationBlocks);

        if (_result == 0) {
            _newStartBlock = block.number;
        } else {
            _newStartBlock = block.number - _result + 1;
        }
    }
}
