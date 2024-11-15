// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

interface IElectionTickets {
  function ticketCount(address _target) external view returns (uint256);
}

/// @title BatchTicketAccounting
/// @dev This contract is used to fetch the accounting of the tickets
///      and is only supposed to be used within the context of an eth_call
contract BatchTicketAccounting {
  /// @notice Fetches the ticket accounting for the given addresses
  ///
  /// @param addresses The addresses to fetch the ticket accounting for
  constructor(address[] memory addresses) {
    IElectionTickets tickets = IElectionTickets(0x4200000000000000000000000000000000000028);
    uint256[] memory ticketCountPerValidator = new uint256[](addresses.length);

    for (uint256 i = 0; i < addresses.length; i++) {
      ticketCountPerValidator[i] = tickets.ticketCount(addresses[i]);
    }

    bytes memory _data = abi.encode(ticketCountPerValidator);

    assembly {
      let _dataStart := add(_data, 32)
      let _dataEnd := sub(msize(), _dataStart)
      return(_dataStart, _dataEnd)
    }
  }
}