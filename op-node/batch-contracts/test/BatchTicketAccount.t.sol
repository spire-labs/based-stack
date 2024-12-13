// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

import {BatchTicketAccounting, IElectionTickets} from 'src/BatchTicketAccounting.sol';
import {Utils} from 'test/Utils.sol';

/// @dev No `setUp` function needed because all tests deploy the batch contract
contract Unit_BatchTicketAccounting is Utils {
  address public constant ELECTION_TICKETS = 0x4200000000000000000000000000000000000028;

  /// @dev Tests that when called with empty data no data is returned
  function test_returnsListOfAddresses_succeeds(uint256 _amountOfTickets, address _randomAddr) public {
    vm.assume(_randomAddr != address(0));

    address[] memory addrs = new address[](1);
    addrs[0] = _randomAddr;

    _mockAndExpect(
      ELECTION_TICKETS,
      abi.encodeWithSelector(IElectionTickets.balanceOf.selector, _randomAddr),
      abi.encode(_amountOfTickets)
    );
    BatchTicketAccounting batchTicketAccounting = new BatchTicketAccounting(addrs);

    // Test the return type and the data matches whats expected of the batch
    bytes memory retdata = address(batchTicketAccounting).code;
    uint256[] memory expected = abi.decode(retdata, (uint256[]));
    assertEq(expected.length, 1);
    assertEq(expected[0], _amountOfTickets);
  }

  /// @dev Tests that when called with multiple addresses, the list is returned
  function test_returnsListWithMultipleValues_succeeds(
    uint256 _amountOfTickets1,
    uint256 _amountOfTickets2,
    uint256 _amountOfTickets3,
    address _randomAddr1,
    address _randomAddr2,
    address _randomAddr3
  ) public {
    vm.assume(_randomAddr1 != address(0));
    vm.assume(_randomAddr2 != address(0));
    vm.assume(_randomAddr3 != address(0));

    vm.assume(_randomAddr1 != _randomAddr2 && _randomAddr1 != _randomAddr3 && _randomAddr2 != _randomAddr3);

    address[] memory addrs = new address[](3);
    addrs[0] = _randomAddr1;
    addrs[1] = _randomAddr2;
    addrs[2] = _randomAddr3;

    _mockAndExpect(
      ELECTION_TICKETS,
      abi.encodeWithSelector(IElectionTickets.balanceOf.selector, _randomAddr1),
      abi.encode(_amountOfTickets1)
    );
    _mockAndExpect(
      ELECTION_TICKETS,
      abi.encodeWithSelector(IElectionTickets.balanceOf.selector, _randomAddr2),
      abi.encode(_amountOfTickets2)
    );
    _mockAndExpect(
      ELECTION_TICKETS,
      abi.encodeWithSelector(IElectionTickets.balanceOf.selector, _randomAddr3),
      abi.encode(_amountOfTickets3)
    );

    BatchTicketAccounting batchTicketAccounting = new BatchTicketAccounting(addrs);

    // Test the return type and the data matches whats expected of the batch
    bytes memory retdata = address(batchTicketAccounting).code;
    uint256[] memory expected = abi.decode(retdata, (uint256[]));
    assertEq(expected.length, 3);
    assertEq(expected[0], _amountOfTickets1);
    assertEq(expected[1], _amountOfTickets2);
    assertEq(expected[2], _amountOfTickets3);
  }
}
