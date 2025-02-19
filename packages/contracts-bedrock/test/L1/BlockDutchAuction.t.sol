// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { ICrossDomainMessenger } from "src/universal/interfaces/ICrossDomainMessenger.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { BlockDutchAuction } from "src/L1/BlockDutchAuction.sol";
import { ElectionTickets } from "src/L2/ElectionTickets.sol";

import "src/libraries/BlockAuctionErrors.sol";

contract TestBlockDutchAuction is BlockDutchAuction {
    constructor(
        uint216 _startBlock,
        uint8 _durationBlocks,
        uint256 _startPrice,
        uint8 _discountRate,
        SystemConfig _systemConfig
    )
        BlockDutchAuction(_startBlock, _durationBlocks, _startPrice, _discountRate, _systemConfig)
    { }

    function findStartBlock(
        uint256 _predictedEndBlock,
        uint256 _durationBlocks
    )
        external
        view
        returns (uint256 _newStartBlock)
    {
        _newStartBlock = _findStartBlock(_predictedEndBlock, _durationBlocks);
    }

    function forTest_setTicketsLeft(uint8 _amount) external {
        _ticketsLeft = _amount;
    }

    function forTest_setPendingDiscountRate(uint8 _discountRate) external {
        pendingDiscountRate = _discountRate;
    }

    function forTest_setPendingStartPrice(uint256 _startPrice) external {
        pendingStartPrice = _startPrice;
    }

    function forTest_setPendingDurationBlocks(uint8 _durationBlocks) external {
        pendingDurationBlocks = _durationBlocks;
    }

    function forTest_setDurationBlocks(uint8 _durationBlocks) external {
        durationBlocks = _durationBlocks;
    }
}

contract BlockDutchAuction_Test is Test {
    TestBlockDutchAuction public auction;
    SystemConfig public systemConfig = SystemConfig(makeAddr("systemConfig"));
    address public messenger = makeAddr("messenger");
    address public owner = makeAddr("owner");
    address public electionTicket = 0x4200000000000000000000000000000000000028;

    function setUp() public virtual {
        vm.prank(owner);
        auction = new TestBlockDutchAuction(1, 32, 1e18, 10, systemConfig);
    }
}

contract BlockDutchAuction_buy_Test is BlockDutchAuction_Test {
    function setUp() public override {
        super.setUp();
        vm.mockCall(
            address(systemConfig), abi.encodeCall(SystemConfig.l1CrossDomainMessenger, ()), abi.encode(messenger)
        );
        vm.mockCall(
            messenger,
            abi.encodeWithSelector(
                ICrossDomainMessenger.sendMessage.selector,
                address(electionTicket),
                abi.encodeCall(ElectionTickets.mint, (owner)),
                150_000
            ),
            abi.encode()
        );
    }

    event TicketBought(address indexed _buyer, uint256 indexed _startBlock, uint256 _price, uint8 _ticketsLeft);

    /// @dev Tests that the `buy` function reverts when there are no tickets left.
    function test_buy_noTicketsLeft_reverts() public {
        auction.forTest_setTicketsLeft(0);

        vm.expectRevert(NoTicketsLeft.selector);
        auction.buy(1);
    }

    /// @dev Tests that the `buy` function reverts when attempting to buy more tickets then available
    function testFuzz_buy_overflow_reverts(uint8 _amount) public {
        // Validators in lookahead is the absolute maximum amount of tickets to potentially be avilable
        // It is possible for certain cases and auction setups for the max to be < VALIDATORS_IN_LOOKAHEAD based on
        // config
        vm.assume(_amount > auction.VALIDATORS_IN_LOOKAHEAD());

        vm.expectRevert(NoTicketsLeft.selector);
        auction.buy(_amount);
    }

    /// @dev Tests that the `buy` function reverts when the caller does not send enough funds.
    function test_buy_insufficientFunds_reverts() public {
        vm.expectRevert(InsufficientFunds.selector);
        auction.buy(1);
    }

    /// @dev Tests that the `buy` function correctly buys a ticket if the auction round auto increments
    function test_buy_differentAuctionRound_succeeds(uint256 _n, uint256 _randomIncrement) public {
        uint256 _firstStartBlock = auction.startBlock();
        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < auction.durationBlocks());

        vm.roll(auction.durationBlocks() * _n + _randomIncrement);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{ value: _price }(1);

        // Tickets left should be durationBlocks - 1 becuase `durationBlocks` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.durationBlocks() - 1);
        assertTrue(auction.startBlock() != _firstStartBlock);
    }

    /// @dev Tests that the `buy` function correctly buys a ticket if the auction round auto increments
    ///      and the pending start price is set.
    function test_buy_differentAuctionRoundWithPendingStartPrice_succeeds(
        uint256 _n,
        uint256 _randomIncrement,
        uint256 _pendingStartPrice
    )
        public
    {
        uint256 _firstStartBlock = auction.startBlock();
        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < auction.durationBlocks());
        _pendingStartPrice = bound(_pendingStartPrice, 1e3, type(uint256).max / 100);

        vm.roll(auction.durationBlocks() * _n + _randomIncrement);

        auction.forTest_setPendingStartPrice(_pendingStartPrice);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{ value: _price }(1);

        // Tickets left should be durationBlocks - 1 becuase `durationBlocks` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.durationBlocks() - 1);
        assertEq(auction.startPrice(), _pendingStartPrice);
        assertEq(auction.pendingStartPrice(), 0);
        assertTrue(auction.startBlock() != _firstStartBlock);
    }

    /// @dev Test that the `buy` function correctly buys a ticket if the auction round auto increments
    ///      and the pending discount rate is set.
    function test_buy_differentAuctionRoundWithPendingDiscountRate_succeeds(
        uint256 _n,
        uint256 _randomIncrement,
        uint8 _pendingDiscountRate
    )
        public
    {
        uint256 _firstStartBlock = auction.startBlock();

        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < auction.durationBlocks());
        _pendingDiscountRate = uint8(bound(_pendingDiscountRate, 1, 100));

        vm.roll(auction.durationBlocks() * _n + _randomIncrement);

        auction.forTest_setPendingDiscountRate(_pendingDiscountRate);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{ value: _price }(1);

        // Tickets left should be durationBlocks - 1 becuase `durationBlocks` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.durationBlocks() - 1);
        assertEq(auction.discountRate(), _pendingDiscountRate);
        assertEq(auction.pendingDiscountRate(), 0);
        assertTrue(auction.startBlock() != _firstStartBlock);
    }

    /// @dev Tests that the `buy` function correctly buys a ticket when the block duration is changed.
    function test_buy_differentAuctionRoundWithPendingDurationBlocks_succeeds(
        uint8 _durationBlocks,
        uint256 _n,
        uint256 _randomIncrement
    )
        public
    {
        uint256 _firstStartBlock = auction.startBlock();

        _durationBlocks = uint8(bound(_durationBlocks, 1, auction.VALIDATORS_IN_LOOKAHEAD()));
        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < _durationBlocks);

        auction.forTest_setPendingDurationBlocks(_durationBlocks);

        vm.roll(auction.durationBlocks() * _n + _randomIncrement);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{ value: _price }(1);

        // Tickets left should be durationBlocks - 1 becuase `durationBlocks` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.durationBlocks() - 1);
        assertEq(auction.durationBlocks(), _durationBlocks);
        assertEq(auction.pendingDurationBlocks(), 0);
        assertTrue(auction.startBlock() != _firstStartBlock);
    }

    /// @dev Tests that the `buy` function correctly buys a ticket.
    function test_buy_succeeds() public {
        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);

        vm.expectCall(
            address(messenger),
            abi.encodeWithSelector(
                ICrossDomainMessenger.sendMessage.selector,
                electionTicket,
                abi.encodeCall(ElectionTickets.mint, (owner)),
                150_000
            )
        );
        vm.prank(owner);
        auction.buy{ value: _price }(1);

        // Tickets left should be durationBlocks - 1 becuase `durationBlocks` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.durationBlocks() - 1);
    }

    /// @dev Tests that the `buy` function correctly buys multiple tickets
    function testFuzz_buy_multipleTickets_succeeds(uint8 _amount) public {
        vm.assume(_amount > 0);
        vm.assume(_amount <= auction.VALIDATORS_IN_LOOKAHEAD());

        uint256 _price = auction.getPrice();
        uint256 _totalPrice = _price * _amount;
        vm.deal(owner, _totalPrice);

        vm.expectCall(
            address(messenger),
            abi.encodeWithSelector(
                ICrossDomainMessenger.sendMessage.selector,
                electionTicket,
                abi.encodeCall(ElectionTickets.mint, (owner)),
                150_000
            )
        );

        vm.prank(owner);
        auction.buy{ value: _totalPrice }(_amount);

        assertEq(auction.ticketsLeft(), auction.durationBlocks() - _amount);
    }

    /// @dev Tests that the `buy` function correctly buys a ticket when the caller overpays
    ///      and the overpayment is refunded.
    function test_buy_overpayment_succeeds(uint256 _overpayment) public {
        uint256 _price = auction.getPrice();

        _overpayment = bound(_overpayment, 1, type(uint128).max);

        vm.deal(owner, _price + _overpayment);
        vm.prank(owner);
        auction.buy{ value: _price + _overpayment }(1);

        assertEq(owner.balance, _overpayment);
    }

    /// @dev Tests that the `buy` function emits the correct event.
    function test_buy_emitsEvent_succeeds() public {
        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);

        vm.expectEmit(true, true, true, true);
        emit TicketBought(owner, auction.startBlock(), _price, uint8(auction.ticketsLeft() - 1));
        vm.prank(owner);
        auction.buy{ value: _price }(1);
    }
}

contract BlockDutchAuction_setStartPrice_Test is BlockDutchAuction_Test {
    event PendingStartPriceSet(uint256 _newStartPrice);

    /// @dev Tests that the `setStartPrice` function reverts when called by a non-owner.
    function test_setStartPrice_onlyOwner_reverts(uint256 _newStartPrice) public {
        vm.expectRevert("Ownable: caller is not the owner");
        auction.setStartPrice(_newStartPrice);
    }

    /// @dev Tests that the `setStartPrice` function reverts when the starting price is too low.
    function test_setStartPrice_invalidStartPrice_TooLow_reverts(uint256 _newStartPrice) public {
        vm.assume(_newStartPrice < 1e3);

        vm.expectRevert(InvalidStartPrice.selector);
        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);
    }

    /// @dev Tests that the `setStartPrice` function reverts when the starting price is too high.
    function test_setStartPrice_InvalidStartPrice_TooHigh_reverts(uint256 _newStartPrice) public {
        vm.assume(_newStartPrice > type(uint256).max / 100);

        vm.expectRevert(InvalidStartPrice.selector);
        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);
    }

    /// @dev Tests that the `setStartPrice` function correctly sets the pending starting price.
    function test_setStartPrice_state_succeeds(uint256 _newStartPrice) public {
        _newStartPrice = bound(_newStartPrice, 1e3, type(uint256).max / 100);

        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);

        assertEq(auction.pendingStartPrice(), _newStartPrice);
    }

    /// @dev Tests that the `setStartPrice` function emits the correct event.
    function test_setStartPrice_emitsEvent_succeeds(uint256 _newStartPrice) public {
        _newStartPrice = bound(_newStartPrice, 1e3, type(uint256).max / 100);

        vm.expectEmit(true, true, true, true);
        emit PendingStartPriceSet(_newStartPrice);
        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);
    }
}

contract BlockDutchAuction_setDiscountRate_Test is BlockDutchAuction_Test {
    event PendingDiscountRateSet(uint8 _newDiscountRate);

    /// @dev Tests that the `setDiscountRate` function reverts when called by a non-owner.
    function test_setDiscountRate_onlyOwner_reverts(uint8 _newDiscountRate) public {
        vm.expectRevert("Ownable: caller is not the owner");
        auction.setDiscountRate(_newDiscountRate);
    }

    /// @dev Tests that the `setDiscountRate` function reverts when the discount rate is invalid.
    function test_setDiscountRate_invalidDiscountRate_reverts(uint8 _newDiscountRate) public {
        vm.assume(_newDiscountRate >= 100 || _newDiscountRate == 0);

        vm.expectRevert(InvalidDiscountRate.selector);
        vm.prank(owner);
        auction.setDiscountRate(_newDiscountRate);
    }

    /// @dev Tests that the `setDiscountRate` function correctly sets the pending discount rate.
    function test_setDiscountRate_state_succeeds(uint8 _newDiscountRate) public {
        _newDiscountRate = uint8(bound(_newDiscountRate, 1, 99));

        vm.prank(owner);
        auction.setDiscountRate(_newDiscountRate);

        assertEq(auction.pendingDiscountRate(), _newDiscountRate);
    }

    /// @dev Tests that the `setDiscountRate` function emits the correct event.
    function test_setDiscountRate_emitsEvent_succeeds(uint8 _newDiscountRate) public {
        _newDiscountRate = uint8(bound(_newDiscountRate, 1, 99));

        vm.expectEmit(true, true, true, true);
        emit PendingDiscountRateSet(_newDiscountRate);
        vm.prank(owner);
        auction.setDiscountRate(_newDiscountRate);
    }
}

contract BlockDutchAuction_setDurationBlocks_Test is BlockDutchAuction_Test {
    event PendingDurationBlocksSet(uint8 _newDurationBlocks);

    /// @dev Tests that the `setDurationBlocks` function reverts when called by a non-owner.
    function test_setDurationBlocks_onlyOwner_reverts(uint8 _newDurationBlocks) public {
        vm.expectRevert("Ownable: caller is not the owner");
        auction.setDurationBlocks(_newDurationBlocks);
    }

    /// @dev Tests that the `setDurationBlocks` function reverts when the duration blocks is invalid.
    function test_setDurationBlocks_invalidDurationBlocks_reverts(uint8 _newDurationBlocks) public {
        vm.assume(_newDurationBlocks > auction.VALIDATORS_IN_LOOKAHEAD());

        vm.expectRevert(InvalidBlockDuration.selector);
        vm.prank(owner);
        auction.setDurationBlocks(_newDurationBlocks);
    }

    /// @dev Tests that the `setDurationBlocks` function correctly sets the pending duration blocks.
    function test_setDurationBlocks_state_succeeds(uint8 _newDurationBlocks) public {
        _newDurationBlocks = uint8(bound(_newDurationBlocks, 1, auction.VALIDATORS_IN_LOOKAHEAD()));

        vm.prank(owner);
        auction.setDurationBlocks(_newDurationBlocks);

        assertEq(auction.pendingDurationBlocks(), _newDurationBlocks);
    }

    /// @dev Tests that the `setDurationBlocks` function emits the correct event.
    function test_setDurationBlocks_emitsEvent_succeeds(uint8 _newDurationBlocks) public {
        _newDurationBlocks = uint8(bound(_newDurationBlocks, 1, auction.VALIDATORS_IN_LOOKAHEAD()));

        vm.expectEmit(true, true, true, true);
        emit PendingDurationBlocksSet(_newDurationBlocks);
        vm.prank(owner);
        auction.setDurationBlocks(_newDurationBlocks);
    }
}

contract BlockDutchAuction_getPrice_Test is BlockDutchAuction_Test {
    /// @dev Tests that the `getPrice` function returns the correct price at the starting price.
    function test_getPrice_atstartPrice_succeeds() public {
        vm.roll(auction.startBlock());
        uint256 _price = auction.getPrice();

        assertEq(_price, auction.startPrice());
    }

    /// @dev Tests that the `getPrice` function returns the correct price after some decay has occurred.
    function test_getPrice_atDiscountedPrice_succeeds() public {
        vm.roll(auction.startBlock() + 3);
        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    /// @dev Tests that the `getPrice` function returns the correct price after some decay has occurred
    ///      in a future auction round.
    function test_getPrice_atDiscountedPriceInFutureAuction_succeeds() public {
        vm.roll(auction.startBlock() + (auction.durationBlocks() + 1) + 3);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    /// @dev Tests that the `getPrice` function returns the correct price after some decay has occurred
    ///      in a future auction round far into the future.
    function test_getPrice_atDiscountedPriceInMoreThenOneFutureAuction_succeeds(uint256 _n) public {
        _n = bound(_n, 2, type(uint64).max);

        vm.roll(auction.startBlock() + (auction.durationBlocks() * _n + 1) + 3);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    /// @dev Tests that the `getPrice` function returns the correct price after some decay has occurred
    ///      in a future auction round with a pending start price.
    function test_getPrice_atDiscountedPriceInFutureAuctionWithPendingStartPrice_succeeds() public {
        vm.roll(auction.startBlock() + (auction.durationBlocks() + 1) + 3);

        auction.forTest_setPendingStartPrice(10e18);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 10000000000000000000 - 10000000000000000000 * 0.1 = 9000000000000000000
        // start block + 2 = 9000000000000000000 -  9000000000000000000 * 0.1 = 8100000000000000000
        // start block + 3 = 8100000000000000000 -  8100000000000000000 * 0.1 = 7290000000000000000
        uint256 _expectedPrice = 7290000000000000000;
        assertEq(_price, _expectedPrice);
    }

    /// @dev Tests that the `getPrice` function returns the correct price after some decay has occurred
    ///      in a future auction round with a pending discount rate.
    function test_getPrice_atDiscountedPriceInFutureAuctionWithPendingDiscountRate_succeeds() public {
        vm.roll(auction.startBlock() + (auction.durationBlocks() + 1) + 3);

        auction.forTest_setPendingDiscountRate(5);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.05 = 9500000000000000000
        // start block + 2 = 9500000000000000000 -  950000000000000000 * 0.05 = 902500000000000000
        // start block + 3 = 902500000000000000 - 902500000000000000 = 857375000000000000
        uint256 _expectedPrice = 857375000000000000;
        assertEq(_price, _expectedPrice);
    }

    /// @dev Tests that the `getPrice` function returns the correct price after some decay has occurred
    ///      in a future auction round with a pending duration blocks.
    function test_getPrice_atDiscountedPriceInFutureAuctionWithPendingDurationBlocks_succeeds() public {
        vm.roll(auction.startBlock() + (auction.durationBlocks() + 1) + 3);

        auction.forTest_setPendingDurationBlocks(10);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }
}

contract BlockDutchAuction_ticketsLeft_Test is BlockDutchAuction_Test {
    /// @dev Tests that the `ticketsLeft` function returns the correct number of tickets at the start of an auction.
    function test_ticketsLeft_atAuctionStart_succeeds() public {
        vm.roll(auction.startBlock());
        uint256 _ticketsLeft = auction.ticketsLeft();

        assertEq(_ticketsLeft, auction.durationBlocks());
    }

    /// @dev Tests that the `ticketsLeft` function returns the correct number of tickets during an auction.
    function test_ticketsLeft_duringAuction_succeeds(uint8 _fakeTickets) public {
        vm.roll(auction.startBlock() + 1);
        vm.assume(_fakeTickets < auction.durationBlocks());

        auction.forTest_setTicketsLeft(_fakeTickets);
        uint256 _ticketsLeft = auction.ticketsLeft();

        assertEq(_ticketsLeft, _fakeTickets);
    }

    /// @dev Tests that the `ticketsLeft` function returns the correct number of tickets after an auction.
    function test_ticketsLeft_afterAuction_succeeds(uint8 _fakeTickets) public {
        vm.roll(auction.startBlock() + (auction.durationBlocks() + 1));
        vm.assume(_fakeTickets < auction.durationBlocks());

        auction.forTest_setTicketsLeft(_fakeTickets);
        uint256 _ticketsLeft = auction.ticketsLeft();

        // View function should ignore the state as it assumes
        // if block.number > startBlock + durationBlocks then it is in a new auction
        assertEq(_ticketsLeft, auction.durationBlocks());
    }
}

contract BlockDutchAuction_findStartBlock_Test is BlockDutchAuction_Test {
    /// @dev Tests that the `findStartBlock` function returns the correct start block.
    function test_findStartBlock_succeeds(uint256 _randomBlock) public {
        // Find start block is only intended to be used if the current auction is over
        vm.assume(_randomBlock > auction.startBlock() + auction.durationBlocks() && _randomBlock < type(uint64).max);

        vm.roll(_randomBlock);
        uint256 _predictedEndBlock = auction.startBlock() + auction.durationBlocks();
        uint256 _newStartBlock = auction.findStartBlock(_predictedEndBlock, auction.durationBlocks());

        assertGe(_newStartBlock, _randomBlock - auction.durationBlocks());
        assertLe(_newStartBlock, _randomBlock);
    }

    /// @dev Tests that the `findStartBlock` function returns the correct start block
    ///      when the current block is an exact multiple of the block duration.
    function test_findStartBlock_atExactMultipleOfN_succeeds(uint256 _n) public {
        vm.assume(_n > 0 && _n < type(uint64).max);

        vm.roll(auction.startBlock() + (auction.durationBlocks() * _n + 1));
        uint256 _predictedEndBlock = auction.startBlock() + auction.durationBlocks();
        uint256 _newStartBlock = auction.findStartBlock(_predictedEndBlock, auction.durationBlocks());

        // Should recognize that the current block is the start block
        assertEq(_newStartBlock, block.number);
    }

    /// @dev Tests that the `findStartBlock` function returns the correct start block when the block duration is
    /// changed.
    function test_findStartBlock_changedDurationBlocks_succeeds(uint256 _randomBlock, uint8 _durationBlocks) public {
        _durationBlocks = uint8(bound(_durationBlocks, 1, auction.VALIDATORS_IN_LOOKAHEAD()));

        auction.forTest_setDurationBlocks(_durationBlocks);

        // Find start block is only intended to be used if the current auction is over
        vm.assume(_randomBlock > auction.startBlock() + auction.durationBlocks() && _randomBlock < type(uint64).max);

        vm.roll(_randomBlock);
        uint256 _predictedEndBlock = auction.startBlock() + auction.durationBlocks();
        uint256 _newStartBlock = auction.findStartBlock(_predictedEndBlock, auction.durationBlocks());

        assertGe(_newStartBlock, _randomBlock - auction.durationBlocks());
        assertLe(_newStartBlock, _randomBlock);
    }

    /// @dev Tests that the `findStartBlock` function returns the correct start block when the block duration is changed
    ///      at an exact multiple of the block duration.
    function test_findStartBlock_changedDurationBlocks_atExactMultipleOfN_succeeds(
        uint256 _n,
        uint8 _durationBlocks
    )
        public
    {
        _durationBlocks = uint8(bound(_durationBlocks, 1, auction.VALIDATORS_IN_LOOKAHEAD()));

        auction.forTest_setDurationBlocks(_durationBlocks);

        vm.assume(_n > 0 && _n < type(uint64).max);

        vm.roll(auction.startBlock() + (auction.durationBlocks() * _n + 1));
        uint256 _predictedEndBlock = auction.startBlock() + auction.durationBlocks();
        uint256 _newStartBlock = auction.findStartBlock(_predictedEndBlock, auction.durationBlocks());

        // Should recognize that the current block is the start block
        assertEq(_newStartBlock, block.number);
    }
}
