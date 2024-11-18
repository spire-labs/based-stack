// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @notice Thrown when the caller is not the election contract
error NotAuction();

/// @notice Thrown when the burn is not coming from the system level
error NotSystemBurn();

/// @notice Thrown when there are no tickets left
error NoTicketsLeft();

/// @notice Thrown when a ticket is attempted to be transferred
error Untransferable();

/// @notice Thrown when an invalid genesis allocation is provided
error InvalidGenesisAllocation();
