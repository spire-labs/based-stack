// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

import {BatchRandomTicketInstruction, IElectionTickets, Retdata} from 'src/BatchRandomTicketInstruction.sol';
import {Utils} from 'test/Utils.sol';

/// @dev No `setUp` function needed because all tests deploy the batch contract
contract Unit_BatchRandomTicketInstruction is Utils {
  address public constant ELECTION_TICKETS = 0x4200000000000000000000000000000000000028;

  /// @dev Tests that when called with empty data no data is returned
  function test_empty_returns_empty_succeeds(
    uint256 _totalTickets
  ) public {
    uint256[] memory timestamps = new uint256[](0);

    _mockAndExpect(
      ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.tokenId.selector), abi.encode(_totalTickets)
    );
    BatchRandomTicketInstruction batchRandomTicketInstruction = new BatchRandomTicketInstruction(timestamps);

    // Test the return type and the data matches whats expected of the batch
    bytes memory retdata = address(batchRandomTicketInstruction).code;
    Retdata[] memory expected = abi.decode(retdata, (Retdata[]));
    assertEq(expected.length, 0);
  }

  /// @dev Tests that when called with one total token, the winner is always the same
  function test_oneTotalTokenAlwaysWins_succeeds(uint256 _timestamp, address _randomAddr) public {
    vm.assume(_randomAddr != address(0));

    uint256 _totalTickets = 1;

    // tokenId 1 should always win if total tickets is 1
    uint256 _tokenId = 1;

    uint256[] memory timestamps = new uint256[](1);
    timestamps[0] = _timestamp;

    _mockAndExpect(
      ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.tokenId.selector), abi.encode(_totalTickets)
    );
    _mockAndExpect(
      ELECTION_TICKETS,
      abi.encodeWithSelector(IElectionTickets.balanceOf.selector, _randomAddr),
      abi.encode(_totalTickets)
    );
    _mockAndExpect(
      ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.ownerOf.selector, _tokenId), abi.encode(_randomAddr)
    );
    BatchRandomTicketInstruction batchRandomTicketInstruction = new BatchRandomTicketInstruction(timestamps);

    // Test the return type and the data matches whats expected of the batch
    bytes memory retdata = address(batchRandomTicketInstruction).code;
    Retdata[] memory expected = abi.decode(retdata, (Retdata[]));
    assertEq(expected.length, 1);
    assertEq(expected[0].timestamp, _timestamp);
    assertEq(expected[0].winner, _randomAddr);
  }

  /// @dev Tests that if the winner is out of tickets, the slot is missed
  function test_missedSlotsIfWinnerOutOfTickets_succeeds(
    uint256 _timestamp1,
    uint256 _timestamp2,
    uint256 _timestamp3,
    address _randomAddr
  ) public {
    vm.assume(_randomAddr != address(0));
    vm.assume(_timestamp1 < _timestamp2 && _timestamp2 < _timestamp3);

    // Make it so that there are two tickets and the winner has all of them
    uint256 _totalTickets = 2;

    // Three timestamps to guarantee the third one should be missed
    uint256[] memory timestamps = new uint256[](3);
    timestamps[0] = _timestamp1;
    timestamps[1] = _timestamp2;
    timestamps[2] = _timestamp3;

    _mockAndExpect(
      ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.tokenId.selector), abi.encode(_totalTickets)
    );
    _mockAndExpect(
      ELECTION_TICKETS,
      abi.encodeWithSelector(IElectionTickets.balanceOf.selector, _randomAddr),
      abi.encode(_totalTickets)
    );

    // We only mock these because the randomness could choose either, and mocking this in a test environment is a potential improvement for later
    // But this is fine for now as which tickets get chosen are covered in other scenarios
    vm.mockCall(ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.ownerOf.selector, 1), abi.encode(_randomAddr));
    vm.mockCall(ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.ownerOf.selector, 2), abi.encode(_randomAddr));

    BatchRandomTicketInstruction batchRandomTicketInstruction = new BatchRandomTicketInstruction(timestamps);

    // Test the return type and the data matches whats expected of the batch
    bytes memory retdata = address(batchRandomTicketInstruction).code;
    Retdata[] memory expected = abi.decode(retdata, (Retdata[]));
    assertEq(expected.length, timestamps.length);
    assertEq(expected[0].timestamp, _timestamp1);
    assertEq(expected[0].winner, _randomAddr);
    assertEq(expected[1].timestamp, _timestamp2);
    assertEq(expected[1].winner, _randomAddr);
    assertEq(expected[2].winner, address(0));
    assertEq(expected[2].timestamp, _timestamp3);
  }

  /// @dev Tests that if the winner is out of tickets, the slot is missed
  function test_ifTicketBurnedIsMissedSlot_succeeds(address _randomAddr, uint256 _timestamp) public {
    vm.assume(_randomAddr != address(0));

    uint256 _totalTickets = 1;

    // tokenId 1 should always win if total tickets is 1
    uint256 _tokenId = 1;

    uint256[] memory timestamps = new uint256[](1);
    timestamps[0] = _timestamp;

    _mockAndExpect(
      ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.tokenId.selector), abi.encode(_totalTickets)
    );
    _mockAndExpect(
      ELECTION_TICKETS, abi.encodeWithSelector(IElectionTickets.ownerOf.selector, _tokenId), abi.encode(address(0))
    );
    BatchRandomTicketInstruction batchRandomTicketInstruction = new BatchRandomTicketInstruction(timestamps);

    // Test the return type and the data matches whats
    bytes memory retdata = address(batchRandomTicketInstruction).code;
    Retdata[] memory expected = abi.decode(retdata, (Retdata[]));
    assertEq(expected.length, 1);
    assertEq(expected[0].timestamp, _timestamp);
    assertEq(expected[0].winner, address(0));
  }
}
