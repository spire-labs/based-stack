// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import { Constants } from "src/libraries/Constants.sol";
import "src/libraries/ElectionTicketErrors.sol";

/// @custom:proxied false
/// @custom:predeploy 0x4200000000000000000000000000000000000028
/// @title ElectionTickets
/// @notice The ERC721 token representing a ticket for sequencing rights of the L2
contract ElectionTickets is ERC721 {
    /// @notice The address of the auction contract on L1
    /// @dev This is used to check that a message being received is sent from the correct contract
    address internal immutable AUCTION;

    /// @notice The token id of the most recently minted ticket
    uint256 public tokenId;

    /// @notice Constructs the ElectionTickets contract
    ///
    /// @param _auction The address of the Election contract
    constructor(address _auction) ERC721("ElectionTickets", "ET") {
        AUCTION = _auction;
    }

    /// @notice Mints a new ticket
    ///
    /// @param _to The address to mint the ticket to
    function mint(address _to) external {
        if (msg.sender != AUCTION) revert NotAuction();

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
        if (msg.sender != Constants.DEPOSITOR_ACCOUNT) revert NotSystemBurn();

        _burn(_tokenId);
    }

    /// @notice Returns the address of the BlockDutchAuction contract
    ///
    /// @return auction_ The address of the BlockDutchAuction contract
    function auction() external view returns (address auction_) {
        auction_ = AUCTION;
    }
}
