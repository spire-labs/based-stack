// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Election } from "src/L1/Election.sol";
import "src/libraries/BatchInboxErrors.sol";

/// @title BatchInbox
/// @notice The BatchInbox is a contract responsible for authentication of the current
///         Sequencer Election winner.
contract BatchInbox {
    /// @notice Address of Election contract.
    Election public election;

    /// @notice Event emitted when a batch was successfully submitted.
    event BatchSubmitted();

    /// @notice Constructs the BatchInboc contract.
    ///
    /// @param _election Election contract address
    constructor(Election _election) {
        election = _election;
    }

    /// @notice Authenticates the batch sender. There's potential problem with transaction inclusion
    //          in the specified block IF the sequencer is not the L1 proposer.
    modifier onlyElectionWinner() {
        address electionWinner = election.electionWinner();
        if (msg.sender != electionWinner) {
            revert OnlyElectionWinner(msg.sender, electionWinner);
        }
        _;
    }

    /// @notice Submits a new batch.
    function submit(uint256 _targetBlock) public onlyElectionWinner {
        if (_targetBlock != block.number) {
            revert InvalidTargetBlock(_targetBlock, block.number);
        }
        emit BatchSubmitted();
    }
}
