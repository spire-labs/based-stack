// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { BatchInbox } from "src/L1/BatchInbox.sol";
import { Election } from "src/L1/Election.sol";

import "src/libraries/BatchInboxErrors.sol";

contract MockElection is Election {
    address public winner;

    function setWinner(address _winner) external {
        winner = _winner;
    }

    function electionWinner() external view override returns (address) {
        return winner;
    }
}

contract BatchInboxTest is Test {
    BatchInbox public batchInbox;
    MockElection public mockElection;
    address public winner = address(0x123);
    address public nonWinner = address(0x456);

    event BatchSubmitted();

    function setUp() public {
        mockElection = new MockElection();
        batchInbox = new BatchInbox(Election(address(mockElection)));
        mockElection.setWinner(winner);
    }

    function testSubmitOnTargetBlockSuccess() public {
        vm.prank(winner);
        uint256 targetBlock = 1337;
        vm.roll(targetBlock);
        batchInbox.submit(targetBlock);
        vm.expectEmit(true, true, true, true);
        emit BatchSubmitted();
    }

    function testSubmitOnTargetBlockRevert_InvalidBlock() public {
        vm.prank(winner);
        uint256 targetBlock = 1337;
        vm.roll(targetBlock + 1);
        batchInbox.submit(targetBlock);
        vm.expectRevert("InvalidTargetBlock");
        revert InvalidTargetBlock(targetBlock, targetBlock + 1);
    }
}
