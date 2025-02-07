// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { BatchInbox, DataType } from "src/L1/BatchInbox.sol";
import { Test } from "forge-std/Test.sol";
import "src/libraries/BatchInboxErrors.sol";

contract BatchInboxTest is Test {
    BatchInbox public batchInbox;

    address public mockWinner = makeAddr("winner");

    event BatchSubmitted(address indexed _submitter, DataType _type);

    function setUp() public {
        batchInbox = new BatchInbox();
    }

    function test_submitCalldata_validTimestamp_succeds() public {
        bytes memory data = abi.encode("some data");

        uint256 targetTimestamp = 1234567;
        vm.warp(targetTimestamp);

        vm.expectEmit(true, false, false, true);
        emit BatchSubmitted(address(this), DataType.CALLDATA);

        batchInbox.submitCalldata(targetTimestamp, data);
    }

    function test_submitCalldata_invalidTimestamp_reverts() public {
        bytes memory data = abi.encode("some data");

        uint256 targetTimestamp = 1234567;
        vm.warp(targetTimestamp + 1);

        vm.expectRevert(abi.encodeWithSelector(InvalidTargetTimestamp.selector, targetTimestamp, targetTimestamp + 1));

        batchInbox.submitCalldata(targetTimestamp, data);
    }

    function test_submitBlob_emitsEvent_succeeds() public {
        uint256 targetTimestamp = block.timestamp;

        vm.expectEmit(true, false, false, true);
        emit BatchSubmitted(mockWinner, DataType.BLOB);

        vm.prank(mockWinner);
        batchInbox.submitBlob(targetTimestamp);
    }

    function test_submitBlob_invalidTargetTimestamp_reverts() public {
        uint256 currentTimestamp = block.timestamp;
        uint256 invalidTargetTimestamp = currentTimestamp + 1;

        vm.expectRevert(
            abi.encodeWithSelector(InvalidTargetTimestamp.selector, invalidTargetTimestamp, currentTimestamp)
        );
        vm.prank(mockWinner);
        batchInbox.submitBlob(invalidTargetTimestamp);
    }
}
