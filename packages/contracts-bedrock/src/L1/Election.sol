// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { BlockDutchAuction } from "src/L1/BlockDutchAuction.sol";
import { ElectionTickets } from "src/L1/ElectionTickets.sol";

// Ultra minimalistic and dirt simple temproary election contract to allow us to build
// The multiple batcher and sequencer ontop of
contract Election is BlockDutchAuction {
    // TODO: Remove this from being hardcoded
    // Currently set to the batcher address that runs in devnet
    // Set as an array for now as we are going to add more
    /// @notice The list of sequencers that are allowed to submit batches
    address[1] public sequencers = [0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC];

    /// @notice Constructs the Election contract
    ///
    /// @param _startBlock The block number to start the auction at
    /// @param _durationBlocks The number of blocks to run the auction for
    /// @param _startPrice The starting price of the auction
    /// @param _discountRate The discount rate of the auction
    /// @param _electionTicket The address of the ElectionTickets contract
    constructor(
        uint216 _startBlock,
        uint8 _durationBlocks,
        uint256 _startPrice,
        uint8 _discountRate,
        ElectionTickets _electionTicket
    )
        BlockDutchAuction(_startBlock, _durationBlocks, _startPrice, _discountRate, _electionTicket)
    { }

    /// @notice The winner of the election
    ///
    /// @return winner_ The address of the winner
    function electionWinner() external view returns (address winner_) {
        winner_ = sequencers[block.number % sequencers.length];
    }
}
