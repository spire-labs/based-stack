// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { ElectionTickets } from "src/L1/ElectionTickets.sol";
import { BlockDutchAuction } from "src/L1/BlockDutchAuction.sol";

import "src/libraries/BlockAuctionErrors.sol";

contract ForTestBlockDutchAuction is BlockDutchAuction {
    constructor(
        uint256 _startBlock,
        uint256 _blockDuration,
        uint256 _startingPrice,
        uint256 _discountRate,
        ElectionTickets _electionTicket
    )
        BlockDutchAuction(_startBlock, _blockDuration, _startingPrice, _discountRate, _electionTicket)
    { }

    function findStartBlock(
        uint256 _currentStartBlock,
        uint256 _blockDuration
    )
        external
        view
        returns (uint256 _newStartBlock)
    {
        _newStartBlock = _findStartBlock(_currentStartBlock, _blockDuration);
    }

    function forTest_setTicketsLeft(uint256 _amount) external {
        _ticketsLeft = _amount;
    }

    function forTest_setPendingDiscountRate(uint256 _discountRate) external {
        pendingDiscountRate = _discountRate;
    }

    function forTest_setPendingStartPrice(uint256 _startingPrice) external {
        pendingStartPrice = _startingPrice;
    }
}

contract BlockDutchAuction_Test is Test {
    ForTestBlockDutchAuction public auction;
    ElectionTickets public electionTicket = ElectionTickets(makeAddr("electionTicket"));
    address public owner = makeAddr("owner");

    function setUp() public {
        vm.prank(owner);
        auction = new ForTestBlockDutchAuction(1, 32, 1e18, 10, electionTicket);
    }
}

contract BlockDutchAuction_buy_Test is BlockDutchAuction_Test {
    event TicketBought(address indexed _buyer, uint256 indexed _startBlock, uint256 _price, uint256 _ticketsLeft);

    function testBuy_NoTicketsLeft_reverts() public {
        auction.forTest_setTicketsLeft(0);

        vm.expectRevert(NoTicketsLeft.selector);
        auction.buy();
    }

    function testBuy_InsufficientFunds_reverts() public {
        vm.expectRevert(InsufficientFunds.selector);
        auction.buy();
    }

    function testBuy_DifferentAuctionRound_success(uint256 _n, uint256 _randomIncrement) public {
        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < auction.blockDuration());

        vm.roll(auction.blockDuration() * _n + _randomIncrement);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{value: _price}();

        // Tickets left should be blockDuration - 1 becuase `blockDuration` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.blockDuration() - 1);
    }

    function testBuy_DifferentAuctionRoundWithPendingStartPrice_success(uint256 _n, uint256 _randomIncrement, uint256 _pendingStartPrice) public {
        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < auction.blockDuration());
        _pendingStartPrice = bound(_pendingStartPrice, 1e3, type(uint256).max / 100);

        vm.roll(auction.blockDuration() * _n + _randomIncrement);

        auction.forTest_setPendingStartPrice(_pendingStartPrice);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{value: _price}();

        // Tickets left should be blockDuration - 1 becuase `blockDuration` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.blockDuration() - 1);
        assertEq(auction.startingPrice(), _pendingStartPrice);
        assertEq(auction.pendingStartPrice(), 0);
    }

    function testBuy_DifferentAuctionRoundWithPendingDiscountRate_success(uint256 _n, uint256 _randomIncrement, uint256 _pendingDiscountRate) public {
        _n = bound(_n, 2, type(uint64).max);
        vm.assume(_randomIncrement < auction.blockDuration());
        _pendingDiscountRate = bound(_pendingDiscountRate, 1, 100);

        vm.roll(auction.blockDuration() * _n + _randomIncrement);

        auction.forTest_setPendingDiscountRate(_pendingDiscountRate);

        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{value: _price}();

        // Tickets left should be blockDuration - 1 becuase `blockDuration` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.blockDuration() - 1);
        assertEq(auction.discountRate(), _pendingDiscountRate);
        assertEq(auction.pendingDiscountRate(), 0);
    }

    function testBuy_success() public {
        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);
        auction.buy{value: _price}();

        // Tickets left should be blockDuration - 1 becuase `blockDuration` tickets are in an auction round
        assertEq(auction.ticketsLeft(), auction.blockDuration() - 1);
    }

    function testBuy_overpayment_success(uint256 _overpayment) public {

        uint256 _price = auction.getPrice();

        _overpayment = bound(_overpayment, 1, type(uint128).max);

        vm.deal(owner, _price + _overpayment);
        vm.prank(owner);
        auction.buy{value: _price + _overpayment}();

        assertEq(owner.balance, _overpayment);
    }

    function testBuy_emitsEvent_success() public {
        uint256 _price = auction.getPrice();
        vm.deal(owner, _price);
        vm.prank(owner);

        vm.expectEmit(true, true, true, true);
        emit TicketBought(owner, auction.startBlock(), _price, auction.ticketsLeft() - 1);
        vm.prank(owner);
        auction.buy{value: _price}();
    }
}

contract BlockDutchAuction_setStartPrice_Test is BlockDutchAuction_Test {
    event PendingStartPriceSet(uint256 _newStartPrice);

    function testSetStartPrice_OwnerOnly_reverts(uint256 _newStartPrice) public {
        vm.expectRevert("Ownable: caller is not the owner");
        auction.setStartPrice(_newStartPrice);
    }

    function testSetStartPrice_InvalidStartingPrice_TooLow_reverts(uint256 _newStartPrice) public {
        vm.assume(_newStartPrice < 1e3);

        vm.expectRevert(InvalidStartingPrice.selector);
        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);
    }

    function testSetStartPrice_InvalidStartingPrice_TooHigh_reverts(uint256 _newStartPrice) public {
        vm.assume(_newStartPrice > type(uint256).max / 100);

        vm.expectRevert(InvalidStartingPrice.selector);
        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);
    }

    function testSetStartPrice_state_success(uint256 _newStartPrice) public {
        _newStartPrice = bound(_newStartPrice, 1e3, type(uint256).max / 100);

        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);

        assertEq(auction.pendingStartPrice(), _newStartPrice);
    }

    function testSetStartPrice_emitsEvent_success(uint256 _newStartPrice) public {
        _newStartPrice = bound(_newStartPrice, 1e3, type(uint256).max / 100);

        vm.expectEmit(true, true, true, true);
        emit PendingStartPriceSet(_newStartPrice);
        vm.prank(owner);
        auction.setStartPrice(_newStartPrice);
    }
}

contract BlockDutchAuction_setDiscountRate_Test is BlockDutchAuction_Test {
    event PendingDiscountRateSet(uint256 _newDiscountRate);

    function testSetDiscountRate_OnlyOwner_reverts(uint256 _newDiscountRate) public {
        vm.expectRevert("Ownable: caller is not the owner");
        auction.setDiscountRate(_newDiscountRate);
    }

    function testSetDiscountRate_InvalidDiscountRate_reverts(uint256 _newDiscountRate) public {
        vm.assume(_newDiscountRate >= 100 || _newDiscountRate == 0);

        vm.expectRevert(InvalidDiscountRate.selector);
        vm.prank(owner);
        auction.setDiscountRate(_newDiscountRate);
    }

    function testSetDiscountRate_state_success(uint256 _newDiscountRate) public {
        _newDiscountRate = bound(_newDiscountRate, 1, 99);

        vm.prank(owner);
        auction.setDiscountRate(_newDiscountRate);

        assertEq(auction.pendingDiscountRate(), _newDiscountRate);
    }

    function testSetDiscountRate_emitsEvent_success(uint256 _newDiscountRate) public {
        _newDiscountRate = bound(_newDiscountRate, 1, 99);

        vm.expectEmit(true, true, true, true);
        emit PendingDiscountRateSet(_newDiscountRate);
        vm.prank(owner);
        auction.setDiscountRate(_newDiscountRate);
    }
}

contract BlockDutchAuction_getPrice_Test is BlockDutchAuction_Test {
    function testGetPrice_AtStartingPrice_success() public {
        vm.roll(auction.startBlock());
        uint256 _price = auction.getPrice();

        assertEq(_price, auction.startingPrice());
    }

    function testGetPrice_AtDiscountedPrice_success() public {
        vm.roll(auction.startBlock() + 3);
        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    function testGetPrice_AtDiscountedPriceInFutureAuction_success() public {
        vm.roll(auction.startBlock() + (auction.blockDuration() + 1) + 3);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    function testGetPrice_AtDiscountedPriceInMoreThenOneFutureAuction_success(uint256 _n) public {
       _n = bound(_n, 2, type(uint64).max);

        vm.roll(auction.startBlock() + (auction.blockDuration() * _n + 1) + 3);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 -  900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 -  810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    function testGetPrice_AtDiscountedPriceInFutureAuctionWithPendingStartPrice_success() public {
        vm.roll(auction.startBlock() + (auction.blockDuration() + 1) + 3);

        auction.forTest_setPendingStartPrice(10e18);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 10000000000000000000 - 10000000000000000000 * 0.1 = 9000000000000000000
        // start block + 2 = 9000000000000000000 -  9000000000000000000 * 0.1 = 8100000000000000000
        // start block + 3 = 8100000000000000000 -  8100000000000000000 * 0.1 = 7290000000000000000
        uint256 _expectedPrice = 7290000000000000000;
        assertEq(_price, _expectedPrice);
    }

    function testGetPrice_AtDiscountedPriceInFutureAuctionWithPendingDiscountRate_success() public {
        vm.roll(auction.startBlock() + (auction.blockDuration() + 1) + 3);

        auction.forTest_setPendingDiscountRate(5);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 - 1000000000000000000 * 0.05 = 9500000000000000000
        // start block + 2 = 9500000000000000000 -  950000000000000000 * 0.05 = 902500000000000000
        // start block + 3 = 902500000000000000 - 902500000000000000 = 857375000000000000
        uint256 _expectedPrice = 857375000000000000;
        assertEq(_price, _expectedPrice);
    }
}

contract BlockDutchAuction_ticketsLeft_Test is BlockDutchAuction_Test {
    function testTicketsLeft_AtStartingPrice_success() public {
        vm.roll(auction.startBlock());
        uint256 _ticketsLeft = auction.ticketsLeft();

        assertEq(_ticketsLeft, auction.blockDuration());
    }

    function testTicketsLeft_DuringAuction_success(uint256 _fakeTickets) public {
        vm.roll(auction.startBlock() + 1);
        vm.assume(_fakeTickets < auction.blockDuration());

        auction.forTest_setTicketsLeft(_fakeTickets);
        uint256 _ticketsLeft = auction.ticketsLeft();

        assertEq(_ticketsLeft, _fakeTickets);
    }

    function testTicketsLeft_AfterAuction_success(uint256 _fakeTickets) public {
        vm.roll(auction.startBlock() + (auction.blockDuration() + 1));
        vm.assume(_fakeTickets < auction.blockDuration());

        auction.forTest_setTicketsLeft(_fakeTickets);
        uint256 _ticketsLeft = auction.ticketsLeft();

        // View function should ignore the state as it assumes
        // if block.number > startBlock + blockDuration then it is in a new auction
        assertEq(_ticketsLeft, auction.blockDuration());
    }
}

contract BlockDutchAuction_findStartBlock_Test is BlockDutchAuction_Test {
    function testFindStartBlock_success(uint256 _randomBlock) public {
        // Find start block is only intended to be used if the current auction is over
        vm.assume(_randomBlock > auction.startBlock() + auction.blockDuration() && _randomBlock < type(uint64).max);

        vm.roll(_randomBlock);
        uint256 _newStartBlock = auction.findStartBlock(auction.startBlock(), auction.blockDuration());

        assertGe(_newStartBlock, _randomBlock - auction.blockDuration());
        assertLe(_newStartBlock, _randomBlock);
    }

    function testFindStartBlock_AtExactMultipleOfN_success(uint256 _n) public {
        vm.assume(_n > 0 && _n < type(uint64).max);

        vm.roll(auction.startBlock() + (auction.blockDuration() * _n + 1));
        uint256 _newStartBlock = auction.findStartBlock(auction.startBlock(), auction.blockDuration());

        // Should recognize that the current block is the start block
        assertEq(_newStartBlock, block.number);
    }
}
