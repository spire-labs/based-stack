// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "src/libraries/BatchInboxErrors.sol";

enum DataType {
    CALLDATA,
    BLOB
}

/// @title BatchInbox
/// @notice The BatchInbox is a contract responsible for authentication of the current
///         Sequencer Election winner.
contract BatchInbox {
    /// @notice Event emitted when a batch was successfully submitted.
    ///
    /// @param _submitter The address of the submitter.
    event BatchSubmitted(address indexed _submitter, DataType _type);

    /// @notice Submits a new batch.
    function submitBlob(uint256 _targetTimestamp) public {
        if (_targetTimestamp != block.timestamp) {
            revert InvalidTargetTimestamp(_targetTimestamp, block.timestamp);
        }
        emit BatchSubmitted(msg.sender, DataType.BLOB);
    }

    function submitCalldata(uint256 _targetTimestamp, bytes calldata) public {
        if (_targetTimestamp != block.timestamp) {
            revert InvalidTargetTimestamp(_targetTimestamp, block.timestamp);
        }
        emit BatchSubmitted(msg.sender, DataType.CALLDATA);
    }
}
