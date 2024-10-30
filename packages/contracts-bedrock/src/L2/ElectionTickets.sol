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

    /// @notice The sentinel ticket id is used to traverse the ticket stack, so that
    ///         ticketQueue[msg.sender][SENTINEL_TICKET_ID] is the top of the stack
    ///         and ticketQueue[msg.sender][bottom_ticket_id] points to SENTINEL_TICKET_ID, meaning you are at the bottom of the stack
    uint256 internal constant SENTINEL_TICKET_ID = 0;

    /// @notice The token id of the most recently minted ticket
    uint256 public tokenId;

    /// @notice The ticket stack for each given address
    /// @dev Implemented as a linked list
    mapping(address => mapping(uint256 => uint256)) public ticketStack;

    /// @notice The number of tickets for each address
    mapping(address => uint256) public ticketCount;


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

        uint256 _tokenId;
        // Not feasible for this to ever overflow
        unchecked {
            _tokenId =  ++tokenId;
        }

        uint256 _top = _top(_to);

        if (_top == 0) {
            // This is the first ticket for this address
            // Set the top to the token id
            ticketStack[_to][SENTINEL_TICKET_ID] = _tokenId;
        } else {
            // This is not the first ticket for this address
            // Move the previous top down the linked list
            ticketStack[_to][_tokenId] = _top;
            // Set the top to the token id that just got minted
            ticketStack[_to][SENTINEL_TICKET_ID] = _tokenId;
        }

        ticketCount[_to]++;

        _mint(_to, tokenId);
    }

    /// @notice Burns a ticket
    ///
    /// @param _target The address to burn the ticket for
    function burn(address _target) external {
        // This check might change in the future
        if (msg.sender != Constants.DEPOSITOR_ACCOUNT) revert NotSystemBurn();

        uint256 _top = _top(_target);

        // Sanity check that there are tickets left
        if (_top == SENTINEL_TICKET_ID) revert NoTicketsLeft();

        // Update the top pointer to the next ticket in the stack
        uint256 _nextTicket = ticketStack[_target][_top];
        ticketStack[_target][SENTINEL_TICKET_ID] = _nextTicket;

        // Clear the reference for the burned ticket
        ticketStack[_target][_top] = 0;

        ticketCount[_target]--;

        // Remove the ticket from the stack
        _burn(_top);
    }

    /// @notice Returns the top of the ticket stack for a given address
    ///
    /// @param _to The address to get the top of the stack for
    /// @return top_ The top of the stack
    function top(address _to) external view returns (uint256 top_) {
        top_ = _top(_to);
    }

    function traverseTicketStack(address _to) external view returns (uint256[] memory stack_) {
        uint256 _count = ticketCount[_to];
        stack_ = new uint256[](_count);

        uint256 _lastTicketId = SENTINEL_TICKET_ID;
        for (uint256 i = 0; i < _count; i++) {
            uint256 _ticketId = ticketStack[_to][_lastTicketId];
            stack_[i] = _ticketId;
            _lastTicketId = _ticketId;
        }
    }

    /// @notice Returns the address of the BlockDutchAuction contract
    ///
    /// @return auction_ The address of the BlockDutchAuction contract
    function auction() external view returns (address auction_) {
        auction_ = AUCTION;
    }

    function _top(address _to) internal view returns (uint256 top_) {
        top_ = ticketStack[_to][SENTINEL_TICKET_ID];
    }
}
