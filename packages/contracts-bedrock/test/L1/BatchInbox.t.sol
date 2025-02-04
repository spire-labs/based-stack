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
        uint256 targetTimestamp = block.timestamp;

        vm.expectEmit(true, true, true, true);
        emit BatchSubmitted(mockWinner);

        vm.prank(mockWinner);
        batchInbox.submit(targetTimestamp);
    }

    function test_submit_invalidTargetTimestamp_reverts() public {
        uint256 currentTimestamp = block.timestamp;
        uint256 invalidTargetTimestamp = currentTimestamp + 1;

        vm.expectRevert(
            abi.encodeWithSelector(InvalidTargetTimestamp.selector, invalidTargetTimestamp, currentTimestamp)
        );
        vm.prank(mockWinner);
        batchInbox.submit(invalidTargetTimestamp);
    }
}
