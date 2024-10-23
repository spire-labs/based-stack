// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @notice Error when the batch's target block is not the current block.
error InvalidTargetBlock(uint256 targetBlock, uint256 currentBlock);
