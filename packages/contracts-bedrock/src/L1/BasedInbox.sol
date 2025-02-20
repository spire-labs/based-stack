// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title BasedInbox
/// @notice The BasedInbox is a contract responsible for routing messages to corresponding L2s through a chain id
///         the contract does minimal work on-chain and relies on off-chain infrastructure.
///
/// @dev This contract is still experimental and not currently hooked up to the rest of the stack
contract BasedInbox {
    /// @notice Event emitted when a transaction is broadcasted.
    ///
    /// @param _chainId The chain id of the L2 the transaction to broadcast to.
    /// @param _count The internal count of broadcast transactions
    event TransactionBroadcasted(uint256 indexed _chainId, uint256 _count);

    uint256 public count;

    /// @notice Broadcasts a transaction to the chain with the specified chain id.
    ///
    /// @param _tx The signed transaction to broadcast.
    /// @param _chainId The chain id of the L2 to broadcast to.
    ///
    /// @dev The intended purpose of this is for the off-chain infrastructure to extract the tx from calldata
    ///      The _tx parameter must represent a signed transaction this is because offchain will call
    /// eth_sendRawTransaction
    ///      on the L2 chain specified by chainId.
    function broadcastTx(bytes calldata _tx, uint256 _chainId) external {
        uint256 _count;

        unchecked {
            _count = ++count;
        }

        emit TransactionBroadcasted(_chainId, _count);
    }
}
