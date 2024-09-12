// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Ultra minimalistic and dirt simple temproary election contract to allow us to build
// The multiple batcher and sequencer ontop of
contract Election {
  // TODO: Remove this from being hardcoded
  // Currently set to the batcher address that runs in devnet
  // Set as an array for now as we are going to add more addresses
  address[1] public sequencers = [0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC];

  function electionWinner() external view returns(address _winner) {
    _winner = sequencers[sequencers.length % block.number];
  }
}