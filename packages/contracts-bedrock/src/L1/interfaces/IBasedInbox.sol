// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

interface IBasedInbox {
    event TransactionBroadcasted(uint256 indexed _chainId, uint256 _count);

    function broadcastTx(bytes calldata _tx, uint256 _chainId) external;
    function count() external view returns (uint256);
}
