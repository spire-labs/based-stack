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
    uint256 public startBlock;

    /// @notice The duration the auction will run for in blocks
    uint256 public blockDuration;

    /// @notice The starting price of the auction
    uint256 public startingPrice;

    /// @notice The discount rate of the auction
    uint256 public discountRate;

    /// @notice The pending new start price of the auction
    uint256 public pendingStartPrice;

    /// @notice The pending new discount rate of the auction
    uint256 public pendingDiscountRate;

    /// @notice The number of tickets left in the current auction
    uint256 internal _ticketsLeft;

    /// @notice Emitted when a new pending start price is set
    /// @param _newStartPrice The new start price
    event PendingStartPriceSet(uint256 _newStartPrice);

    /// @notice Emitted when a new pending discount rate is set
    /// @param _newDiscountRate The new discount rate
    event PendingDiscountRateSet(uint256 _newDiscountRate);

    /// @notice Emitted when a ticket is bought
    /// @param _buyer The address of the buyer
    /// @param _startBlock The start block of the current auction the ticket was bought in
    /// @param _price The price of the ticket
    event TicketBought(address indexed _buyer, uint256 indexed _startBlock, uint256 _price, uint256 _ticketsLeft);

    /// @notice Constructs and initilizes the BlockDutchAuction
    /// @param _startBlock The start block of the auction
    /// @param _blockDuration The duration of the auction in blocks
    /// @param _startingPrice The starting price of the auction
    /// @param _discountRate The discount rate of the auction
    /// @param _electionTicket The address of the ElectionTicket contract
    constructor(
        uint256 _startBlock,
        uint256 _blockDuration,
        uint256 _startingPrice,
        uint256 _discountRate,
        ElectionTickets _electionTicket
    ) {
        if (_discountRate >= 100 || _discountRate == 0) revert InvalidDiscountRate();
        if (_blockDuration > VALIDATORS_IN_LOOKAHEAD) revert InvalidBlockDuration();
        if (_startingPrice < 1e3 || _startingPrice > type(uint256).max / 100) revert InvalidStartingPrice();

        startBlock = _startBlock;
        blockDuration = _blockDuration;
        startingPrice = _startingPrice;
        discountRate = _discountRate;
        ELECTION_TICKET = _electionTicket;

        // Set the starting amount of tickets
        _ticketsLeft = _blockDuration;
    }

    /*///////////////////////////////////////////////////////////////
                              ADMIN FUNCTIONS
    ///////////////////////////////////////////////////////////////*/

    function setStartPrice(uint256 _newStartPrice) external onlyOwner {
        if (_newStartPrice < 1e3 || _newStartPrice > type(uint256).max / 100) revert InvalidStartingPrice();

        pendingStartPrice = _newStartPrice;

        emit PendingStartPriceSet(_newStartPrice);
    }

    function setDiscountRate(uint256 _newDiscountRate) external onlyOwner {
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
        uint256 _blockDuration = blockDuration;

        // If the current block number is past the end of the auction
        // This branch will trigger the start of a new auction
        // View functions have this logic baked in to calculate based on what the auction would be at
        if (block.number > _startBlock + _blockDuration) {
            // Find the start block of the next auction
            startBlock = _findStartBlock(_startBlock, _blockDuration);

            // is warm storage cheaper then calling findStartBlock again?
            _startBlock = startBlock;

            if (pendingStartPrice != 0) {
                startingPrice = pendingStartPrice;
                pendingStartPrice = 0;
            }

            if (pendingDiscountRate != 0) {
                discountRate = pendingDiscountRate;
                pendingDiscountRate = 0;
            }

            // Amount of tickets sold in an auction should always equal the block duration
            _ticketsLeft = _blockDuration;
        }

        // This check needs to come second incase its a new auction
        uint256 __ticketsLeft = _ticketsLeft;
        if (__ticketsLeft == 0) revert NoTicketsLeft();

        uint256 _price = _getPrice(discountRate, startingPrice, _startBlock);

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
        uint256 _blockDuration = blockDuration;
        uint256 _startingPrice = startingPrice;
        uint256 _discountRate = discountRate;

        if (block.number > _startBlock + _blockDuration) {
            // Find the expected start block
            _startBlock = _findStartBlock(_startBlock, _blockDuration);

            // Use the pending values to calculate the price
            if (pendingDiscountRate != 0) {
                _discountRate = pendingDiscountRate;
            }

            if (pendingStartPrice != 0) {
                _startingPrice = pendingStartPrice;
            }
        }

        _price = _getPrice(_discountRate, _startingPrice, _startBlock);
    }

    /// @notice Returns the amount of tickets left in the current auction
    /// @return _amount The amount of tickets left in the current auction
    function ticketsLeft() external view returns (uint256 _amount) {
        uint256 _blockDuration = blockDuration;
        if (block.number > startBlock + _blockDuration) {
            _amount = _blockDuration;
        } else {
            _amount = _ticketsLeft;
        }
    }

    /// @notice Returns the current price of a ticket in the auction
    /// @param _discountRate The discount rate of the auction
    /// @param _startingPrice The starting price of the auction
    /// @param _startBlock The start block of the auction
    /// @return _price The current price of a ticket
    function _getPrice(
        uint256 _discountRate,
        uint256 _startingPrice,
        uint256 _startBlock
    )
        internal
        view
        returns (uint256 _price)
    {
        uint256 blocksSinceStart = block.number - _startBlock;
        _price = _startingPrice;

        // NOTE: Is this the math we want to go with?
        for (uint256 i; i < blocksSinceStart; i++) {
            _price = _price - ((_price * _discountRate) / 100);
        }
    }

    /// @notice Recursively finds the start block of the next auction
    /// @param _currentStartBlock The current start block of the auction
    /// @param _blockDuration The duration of the auction in blocks
    /// @return _newStartBlock The start block of the next auction
    function _findStartBlock(
        uint256 _currentStartBlock,
        uint256 _blockDuration
    )
        internal
        view
        returns (uint256 _newStartBlock)
    {
        uint256 _predictedEndBlock = _currentStartBlock + _blockDuration;
        uint256 _difference = block.number - _predictedEndBlock;
        uint256 _result = _difference % (_blockDuration);

        if (_result == 0) {
            _newStartBlock = block.number;
        } else {
            _newStartBlock = block.number - _result + 1;
        }
    }
}
