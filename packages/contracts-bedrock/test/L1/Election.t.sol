// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Election } from "src/L1/Election.sol";
import { ElectionTickets } from "src/L1/ElectionTickets.sol";
import { Test } from "forge-std/Test.sol";

// TODO: Add this more seamlessly into the op test suite
contract ElectionTest is Test {
    Election public election;
    uint216 public startBlock = 1;
    uint8 public durationBlocks = 32;
    uint256 public startPrice = 1e18;
    uint8 public discountRate = 10;
    ElectionTickets public tickets = ElectionTickets(makeAddr("electionTicket"));

    function setUp() public {
        election = new Election(startBlock, durationBlocks, startPrice, discountRate, tickets);
    }

    function testReturnsWinner(uint256 _blockNumber) public {
        vm.assume(_blockNumber > 0);

        vm.roll(_blockNumber);

        address _expectedWinner = election.sequencers(0);
        address _winner = election.electionWinner();

        vm.assertTrue(_winner == _expectedWinner);
    }
}
