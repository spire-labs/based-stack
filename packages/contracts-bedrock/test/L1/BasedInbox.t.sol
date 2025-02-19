// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { IBasedInbox } from "src/L1/interfaces/IBasedInbox.sol";
import { BasedInbox } from "src/L1/BasedInbox.sol";

contract BasedInbox_BroadcastTx_Test is Test {
    event TransactionBroadcasted(uint256 indexed _chainId, uint256 _count);

    IBasedInbox inbox;

    function setUp() public {
        inbox = IBasedInbox(address(new BasedInbox()));
    }

    /// @dev Test that broadcasting a transaction increases count
    function testFuzz_increasesCount_succeeds(uint256 _chainId, bytes calldata _tx) public {
        assertEq(inbox.count(), 0);

        inbox.broadcastTx(_tx, _chainId);

        assertEq(inbox.count(), 1);
    }

    /// @dev Test that broadcasting a transaction emits an event
    function testFuzz_emitsEvent_succeeds(uint256 _chainId, bytes calldata _tx) public {
        vm.expectEmit(true, true, true, true);
        emit TransactionBroadcasted(_chainId, 1);

        inbox.broadcastTx(_tx, _chainId);
    }
}
