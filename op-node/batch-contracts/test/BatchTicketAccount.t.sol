// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

// External dependencies
import {Test} from "forge-std/Test.sol";

// Internal dependencies
import {BatchTicketAccounting, IElectionTickets} from "src/BatchTicketAccounting.sol";

/// @dev No `setUp` function needed because all tests deploy the batch contract
contract Unit_BatchTicketAccounting is Test {
  address public constant ELECTION_TICKETS = 0x4200000000000000000000000000000000000028;

  function testBatchTicketAccounting() public {
    address[] memory addresses = new address[](1);
    addresses[0] = address(0x42);

    BatchTicketAccounting accounting = new BatchTicketAccounting(addresses);

    uint256[] memory ticketCountPerValidator = abi.decode(accounting.getAccounting(), (uint256[]));

    assertEq(ticketCountPerValidator[0], 0);
  }
}