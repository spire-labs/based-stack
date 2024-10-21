// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { BatchInbox } from "src/L1/BatchInbox.sol";
import { Test } from "forge-std/Test.sol";
import "src/libraries/BatchInboxErrors.sol";

contract BatchInboxTest is Test {
    BatchInbox public batchInbox;

    address public mockWinner = makeAddr("winner");

    event BatchSubmitted(address indexed _submitter);

    function setUp() public {
        batchInbox = new BatchInbox();
    }

    function test_submit_emitsEvent_succeeds() public {
        uint256 targetBlock = block.number;

        vm.expectEmit(true, true, true, true);
        emit BatchSubmitted(mockWinner);

        vm.prank(mockWinner);
        batchInbox.submit(targetBlock);
    }

    function test_submit_invalidTargetBlock_reverts() public {
        uint256 currentBlock = block.number;
        uint256 invalidTargetBlock = currentBlock + 1;

        vm.expectRevert(abi.encodeWithSelector(InvalidTargetBlock.selector, invalidTargetBlock, currentBlock));
        vm.prank(mockWinner);
        batchInbox.submit(invalidTargetBlock);
    }
}
