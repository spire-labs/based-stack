// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @notice Error for when the discount rate is greater than 100%
error InvalidDiscountRate();

/// @notice Error for when the buyer does not have enough funds to buy the ticket
error InsufficientFunds();

/// @notice Error for when the block duration is greater than 32
error InvalidBlockDuration();

/// @notice Error for when the starting price is invalid
error InvalidStartPrice();

/// @notice Error when there are no tickets left to buy
error NoTicketsLeft();

/// @notice Error when a low level call fails
error FailedLowLevelCall();
