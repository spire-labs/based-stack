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
        uint256 _blockDuration,
        uint256 n
    )
        external
        view
        returns (uint256 _newStartBlock)
    {
        _newStartBlock = _findStartBlock(_currentStartBlock, _blockDuration, n);
    }

    function forTest_setStartBlock(uint256 _startBlock) external {
        startBlock = _startBlock;
    }

    function forTest_setTicketsLeft(uint256 _amount) external {
        _ticketsLeft = _amount;
    }
}

contract BlockDutchAuction_Test is Test {
    ForTestBlockDutchAuction public auction;
    ElectionTickets public electionTicket = ElectionTickets(address(0x1));

    function setUp() public {
        auction = new ForTestBlockDutchAuction(1, 32, 1e18, 10, electionTicket);
    }
}

contract BlockDutchAuction_setStartPrice_Test is BlockDutchAuction_Test {
    event PendingStartPriceSet(uint256 _newStartPrice);

    function testSetStartPrice_state_success() public {
        uint256 _newStartPrice = auction.startingPrice() + 1;
        auction.setStartPrice(_newStartPrice);

        assertEq(auction.pendingStartPrice(), _newStartPrice);
    }

    function testSetStartPrice_emitsEvent_success() public {
        uint256 _newStartPrice = auction.startingPrice() + 1;

        vm.expectEmit(true, true, true, true);
        emit PendingStartPriceSet(_newStartPrice);
        auction.setStartPrice(_newStartPrice);
    }
}

contract BlockDutchAuction_setDiscountRate_Test is BlockDutchAuction_Test {
    event PendingDiscountRateSet(uint256 _newDiscountRate);

    function testSetDiscountRate_InvalidDiscountRate_reverts() public {
        uint256 _newDiscountRate = 101;

        vm.expectRevert(InvalidDiscountRate.selector);
        auction.setDiscountRate(_newDiscountRate);
    }

    function testSetDiscountRate_state_success() public {
        uint256 _newDiscountRate = auction.discountRate() + 1;
        auction.setDiscountRate(_newDiscountRate);

        assertEq(auction.pendingDiscountRate(), _newDiscountRate);
    }

    function testSetDiscountRate_emitsEvent_success() public {
        uint256 _newDiscountRate = auction.discountRate() + 1;

        vm.expectEmit(true, true, true, true);
        emit PendingDiscountRateSet(_newDiscountRate);
        auction.setDiscountRate(_newDiscountRate);
    }
}

contract BlockDutchAuction_findStartBlock_Test is BlockDutchAuction_Test {
    function testFindStartBlock_success(uint256 _randomBlock) public {
        vm.assume(_randomBlock > auction.startBlock() + (auction.blockDuration() * 2) && _randomBlock < type(uint8).max);

        vm.roll(_randomBlock);
        uint256 _newStartBlock = auction.findStartBlock(auction.startBlock(), auction.blockDuration(), 3);

        assertGe(_newStartBlock, _randomBlock - auction.blockDuration());
        assertLe(_newStartBlock, _randomBlock);
    }

    function testFindStartBlock_AtExactMultipleOfN_success(uint256 _n) public {
        vm.assume(_n > 0 && _n < 10);

        vm.roll(auction.startBlock() + (auction.blockDuration() * _n + 1));
        uint256 _newStartBlock = auction.findStartBlock(auction.startBlock(), auction.blockDuration(), _n);

        assertEq(_newStartBlock, block.number);
    }
}

contract BlockDutchAuction_getPrice_Test is BlockDutchAuction_Test {
    function testGetPrice_AtStartingPrice_success() public {
        vm.roll(auction.startBlock());
        uint256 _price = auction.getPrice();

        assertEq(_price, auction.startingPrice());
    }

    function testGetPrice_AtDiscontedPrice_success() public {
        vm.roll(auction.startBlock() + 3);
        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    function testGetPrice_AtDiscountedPriceInFutureAuction_success() public {
        vm.roll(auction.startBlock() + (auction.blockDuration() + 1) + 3);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
        assertEq(_price, _expectedPrice);
    }

    function testGetPrice_AtDiscountedPriceInMoreThenOneFutureAuction_success(uint256 _n) public {
        vm.assume(_n > 2 && _n < 10);

        vm.roll(auction.startBlock() + (auction.blockDuration() * _n + 1) + 3);

        uint256 _price = auction.getPrice();

        // Current discount rate is 10% so it takes 10% of the new price each block
        // start block + 1 = 1000000000000000000 * 0.1 = 900000000000000000
        // start block + 2 = 900000000000000000 * 0.1 = 810000000000000000
        // start block + 3 = 810000000000000000 * 0.1 = 729000000000000000
        uint256 _expectedPrice = 729000000000000000;
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
