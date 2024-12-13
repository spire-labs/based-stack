// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

interface IElectionTickets {
  function balanceOf(address _target) external view returns (uint256);
  function tokenId() external view returns (uint256);
  function ownerOf(uint256 _tokenId) external view returns (address);
}

struct Retdata {
  uint256 timestamp;
  address winner;
}

/// @title BatchRandomTicketInstruction
/// @dev This contract is used to determine the winner of a certain slot in the election
///      if the random ticket instruction is used.
///
/// @dev This contract is only supposed to be used within the context of an eth_call
contract BatchRandomTicketInstruction {
  // TODO(spire): Use in memory mappings or transient storage to avoid using SSTORE
  mapping(address => uint256) public balances;
  mapping(address => bool) public seenWinners;
  mapping(uint256 => bool) public usedTickets;

  /// @notice Determines certain timestamps slot winners based on the random ticket instruction
  ///
  /// @param timestamps The timestamps to determine the winner for
  constructor(uint256[] memory timestamps) {
    IElectionTickets tickets = IElectionTickets(0x4200000000000000000000000000000000000028);
    uint256 totalTickets = tickets.tokenId();
    Retdata[] memory retdatas = new Retdata[](timestamps.length);


    for (uint256 i; i < timestamps.length; i++) {
      // With a low "totalTickets" it is feasible to have a collision where the same tokenId wins multiple times
      // For this scenario we will choose the same winner, assuming they have enough total tickets to win
      uint256 pseudoRandomTokenId = uint256(keccak256(abi.encodePacked(timestamps[i], block.prevrandao))) % (totalTickets + 1);

      address winner = tickets.ownerOf(pseudoRandomTokenId);
      uint256 balanceOf = tickets.balanceOf(winner);

      bool hasSeenWinner = seenWinners[winner];

      if (hasSeenWinner && balances[winner] == 0) {
        // TODO(spire): Maybe we can do some form of search to avoid a missed slot
        // For now if the winner is out of tickets from previous iterations, we will assume its a missed slot
        winner = address(0);
      }

      // If this is the first time we see this winner, set the balance
      if (!hasSeenWinner) {
        // They won one ticket and it will be burn so we do - 1 for future potential iterations
        balances[winner] = balanceOf - 1;
        seenWinners[winner] = true;
      }

      retdatas[i] = Retdata(timestamps[i], tickets.ownerOf(pseudoRandomTokenId));
    }

    bytes memory data = abi.encode(retdatas);

    assembly {
      let dataStart := add(data, 32)
      let dataEnd := sub(msize(), dataStart)
      return(dataStart, dataEnd)
    }
  }
}