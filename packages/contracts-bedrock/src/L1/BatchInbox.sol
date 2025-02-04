// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "src/libraries/BatchInboxErrors.sol";

/// @title BatchInbox
/// @notice The BatchInbox is a contract responsible for authentication of the current
///         Sequencer Election winner.
contract BatchInbox {
    /// @notice Event emitted when a batch was successfully submitted.
    ///
    /// @param _submitter The address of the submitter.
    event BatchSubmitted(address indexed _submitter);

    /// @notice Submits a new batch.
    function submit(uint256 _targetTimestamp) public {
        if (_targetTimestamp != block.timestamp) {
            revert InvalidTargetTimestamp(_targetTimestamp, block.timestamp);
        }
        emit BatchSubmitted(msg.sender);
    }
}
