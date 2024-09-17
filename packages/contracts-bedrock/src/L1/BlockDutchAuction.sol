// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "src/libraries/BlockAuctionErrors.sol";

abstract contract BlockDutchAuction {
    /// @notice The minimum number of validators in the look ahead
    uint256 public constant VALIDATORS_IN_LOOKAHEAD = 32;

    /// @notice The address of the ElectionTicket contract
    address public immutable ELECTION_TICKET;

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
        address _electionTicket
    ) {
        if (_discountRate > 100 || _discountRate == 0) revert InvalidDiscountRate();
        if (_blockDuration > VALIDATORS_IN_LOOKAHEAD) revert InvalidBlockDuration();

        startBlock = _startBlock;
        blockDuration = _blockDuration;
        startingPrice = _startingPrice;
        discountRate = _discountRate;
        ELECTION_TICKET = _electionTicket;
    }

    /*///////////////////////////////////////////////////////////////
                              ADMIN FUNCTIONS
    ///////////////////////////////////////////////////////////////*/

    // TODO: Write admin functions that connect to the larger stack through permissioned channels
    function setStartPrice(uint256 _newStartPrice) external {
        pendingStartPrice = _newStartPrice;

        emit PendingStartPriceSet(_newStartPrice);
    }

    function setDiscountRate(uint256 _newDiscountRate) external {
        if (_newDiscountRate > 100) revert InvalidDiscountRate();

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

        // If the current block number is past the end of the auction, set the start block to the end of the auction
        if (block.number > _startBlock + _blockDuration) {
            // In the rare case where an auction went by without any bidders we need to recursively find a start block
            if (block.number > _startBlock + (_blockDuration * 2)) {
                // NOTE: This is not a full safe solution as it can fail
                // with stack overflow if auctions are ignored for too long
                // TODO: Improve this
                // perhaps if such a state is reached we can force manual intervention from the appchain
                startBlock = _findStartBlock(_startBlock, _blockDuration, 3);
                _startBlock = startBlock;
            } else {
                startBlock = _startBlock + _blockDuration + 1;
                _startBlock = _startBlock + _blockDuration + 1;
            }

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

        uint256 _price = _getPrice(discountRate, startingPrice, _startBlock);

        if (_price > msg.value) revert InsufficientFunds();
        else if (msg.value > _price) payable(msg.sender).transfer(_price - msg.value);

        // TODO: Mint ticket

        uint256 __ticketsLeft = _ticketsLeft;
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

        // In most cases it will be hitting the else block
        // In the rare case where an auction went by without any bidders we need to recursively find a start block
        if (block.number > _startBlock + (_blockDuration * 2)) {
            // NOTE: This is not a full safe solution as it can fail with stack overflow if auctions are ignored for too
            // long
            // TODO: Improve this, perhaps if such a state is reached we can force manual intervention from the appchain
            _startBlock = _findStartBlock(_startBlock, _blockDuration, 3);
        } else if (block.number > _startBlock + _blockDuration) {
            _startBlock = _startBlock + _blockDuration + 1;
        }

        _price = _getPrice(discountRate, startingPrice, _startBlock);
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
    /// @param _n The number of `_blockDuration`s to check for the next start block
    /// @return _newStartBlock The start block of the next auction
    function _findStartBlock(
        uint256 _currentStartBlock,
        uint256 _blockDuration,
        uint256 _n
    )
        internal
        view
        returns (uint256 _newStartBlock)
    {
        if (block.number > _currentStartBlock + (_blockDuration * _n)) {
            _newStartBlock = _findStartBlock(_currentStartBlock, _blockDuration, _n + 1);
        } else if (block.number == _currentStartBlock + (_blockDuration * _n) + 1) {
            _newStartBlock = _currentStartBlock + (_blockDuration * (_n)) + 1;
        } else {
            _newStartBlock = _currentStartBlock + (_blockDuration * (_n - 1)) + 1;
        }
    }
}
