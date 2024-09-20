// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "src/libraries/ElectionTicketErrors.sol";

contract ElectionTickets is ERC721 {
    /// @notice The address of the Election contract
    address internal immutable ELECTION;

    /// @notice The address of the BatchInbox contract
    address internal immutable BATCH_INBOX;

    /// @notice The token id of the most recently minted ticket
    uint256 public tokenId;

    /// @notice Constructs the ElectionTickets contract
    ///
    /// @param _election The address of the Election contract
    /// @param _batchInbox The address of the BatchInbox contract
    constructor(address _election, address _batchInbox) ERC721("ElectionTickets", "ET") {
        ELECTION = _election;
        BATCH_INBOX = _batchInbox;
    }

    /// @notice Mints a new ticket
    ///
    /// @param _to The address to mint the ticket to
    function mint(address _to) external {
        if (msg.sender != ELECTION) revert NotElection();

        // Not feasible for this to ever overflow
        unchecked {
            ++tokenId;
        }

        _mint(_to, tokenId);
    }

    /// @notice Burns a ticket
    ///
    /// @param _tokenId The token id of the ticket to burn
    function burn(uint256 _tokenId) external {
        // This check might change in the future
        if (msg.sender != BATCH_INBOX) revert NotBatchInbox();

        _burn(_tokenId);
    }

    /// @notice Returns the address of the Election contract
    ///
    /// @return election_ The address of the Election contract
    function election() external view returns (address election_) {
        election_ = ELECTION;
    }

    /// @notice Returns the address of the BatchInbox contract
    ///
    /// @return batchInbox_ The address of the BatchInbox contract
    function batchInbox() external view returns (address batchInbox_) {
        batchInbox_ = BATCH_INBOX;
    }
}
