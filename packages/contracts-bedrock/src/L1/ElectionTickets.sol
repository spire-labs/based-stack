// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "src/libraries/ElectionTicketErrors.sol";

contract ElectionTickets is ERC721 {
    /// @notice The address of the Election contract
    address public immutable ELECTION;

    /// @notice The address of the BatchInbox contract
    address public immutable BATCH_INBOX;

    /// @notice The token id of the most recently minted ticket
    uint256 public tokenId;

    /// @notice The block number a specific ticket was minted at
    /// @dev `tokenId => block.number`
    mapping(uint256 => uint256) public blockOfMint;

    constructor(address _election, address _batchInbox) ERC721("ElectionTickets", "ET") {
        ELECTION = _election;
        BATCH_INBOX = _batchInbox;
     }

     /// @notice Mints a new ticket
     /// @param _to The address to mint the ticket to
     function mint(address _to) external {
        if (msg.sender != ELECTION) revert NotElection();

        // Not feasible for this to ever overflow
        unchecked {
            ++tokenId;
        }

        _mint(_to, tokenId);

        blockOfMint[tokenId] = block.number;
     }

    /// @notice Burns a ticket
    /// @param _tokenId The token id of the ticket to burn
     function burn(uint256 _tokenId) external {
        if (msg.sender != BATCH_INBOX) revert NotBatchInbox();

        _burn(_tokenId);

        blockOfMint[_tokenId] = 0;
     }




}
