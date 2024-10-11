// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { BatchInbox } from "src/L1/BatchInbox.sol";
import { Election } from "src/L1/Election.sol";
import { Test } from "forge-std/Test.sol";
import "src/libraries/BatchInboxErrors.sol";

contract BatchInboxTest is Test {
    BatchInbox public batchInbox;
    Election public election = Election(makeAddr("election"));

    address public mockWinner = address(0x123);

    event BatchSubmitted();

    function setUp() public {
        batchInbox = new BatchInbox(election);
    }

    function test_submit_emitsEvent_succeeds() public {
        bytes memory data = abi.encodeWithSelector(Election.electionWinner.selector);
        vm.mockCall(address(election), data, abi.encode(mockWinner));

        vm.expectCall(address(election), data);

        vm.prank(mockWinner);
        uint256 targetBlock = block.number;

        vm.expectEmit(true, true, true, true);
        emit BatchSubmitted();

        batchInbox.submit(targetBlock);
    }

    function test_submit_notElectionWinner_reverts() public {
        bytes memory data = abi.encodeWithSelector(Election.electionWinner.selector);
        vm.mockCall(address(election), data, abi.encode(mockWinner));

        address nonWinner = address(0x456);
        vm.prank(nonWinner);

        uint256 targetBlock = block.number;

        vm.expectRevert(abi.encodeWithSelector(OnlyElectionWinner.selector, nonWinner, mockWinner));

        batchInbox.submit(targetBlock);
    }

    function test_submit_invalidTargetBlock_reverts() public {
        bytes memory data = abi.encodeWithSelector(Election.electionWinner.selector);
        vm.mockCall(address(election), data, abi.encode(mockWinner));

        vm.prank(mockWinner);

        uint256 currentBlock = block.number;
        uint256 invalidTargetBlock = currentBlock + 1;

        vm.expectRevert(abi.encodeWithSelector(InvalidTargetBlock.selector, invalidTargetBlock, currentBlock));

        batchInbox.submit(invalidTargetBlock);
    }
}
