// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Election } from "src/L1/Election.sol";
import { Test } from "forge-std/Test.sol";

// TODO: Add this more seamlessly into the op test suite
contract ElectionTest is Test {
  Election public election;

  function setUp() public {
    election = new Election();
  }

  function testReturnsWinner() public view {
    address _expectedWinner = election.sequencers(0);
    address _winner = election.electionWinner();

    vm.assertTrue(_winner == _expectedWinner);
  }
}
