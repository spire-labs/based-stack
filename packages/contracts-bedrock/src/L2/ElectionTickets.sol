// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import { Initializable } from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import { Constants } from "src/libraries/Constants.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { ICrossDomainMessenger } from "src/universal/interfaces/ICrossDomainMessenger.sol";
import "src/libraries/ElectionTicketErrors.sol";

/// @custom:proxied true
/// @custom:predeploy 0x4200000000000000000000000000000000000028
/// @title ElectionTickets
/// @notice The ERC721 token representing a ticket for sequencing rights of the L2
contract ElectionTickets is ERC721, Initializable {
    /// @notice The struct for the genesis allocation
    ///
    /// @dev This structs parameters need to be kept in alphabetical order due to us using it for JSON deserialization
    ///
    /// @param amount The amount of tickets to mint
    /// @param target The address to mint the ticket to
    struct GenesisAllocation {
        uint256 amount;
        address target;
    }

    /// @notice The address of the auction contract on L1
    /// @dev This is used to check that a message being received is sent from the correct contract
    address internal immutable AUCTION;

    /// @notice The sentinel ticket id is used to traverse the ticket stack, so that
    ///         ticketStack[msg.sender][SENTINEL_TICKET_ID] is the top of the stack
    ///         and ticketStack[msg.sender][bottom_ticket_id] points to SENTINEL_TICKET_ID, meaning you are at the
    ///         bottom of the stack
    uint256 internal constant SENTINEL_TICKET_ID = 0;

    /// @notice The token id of the most recently minted ticket
    uint256 public tokenId;

    /// @notice The ticket stack for each given address
    /// @dev Implemented as a linked list
    mapping(address => mapping(uint256 => uint256)) public ticketStack;

    /// @notice Constructs the ElectionTickets contract
    ///
    /// @param _auction The address of the Election contract
    constructor(address _auction) ERC721("ElectionTickets", "ET") {
        AUCTION = _auction;

        initialize(new GenesisAllocation[](0));
    }

    /// @notice Initializes the ElectionTickets contract
    ///
    /// @param _genesisAllocation The array of allocation details to mint the genesis tickets to
    function initialize(GenesisAllocation[] memory _genesisAllocation) public initializer {
        uint256 _amountMinted;
        uint256 _genesisTicketsAmount = _genesisAllocation.length;
        for (uint256 i; i < _genesisTicketsAmount; i++) {
            for (uint256 j; j < _genesisAllocation[i].amount; j++) {
                _mintTo(_genesisAllocation[i].target, _amountMinted + j + 1);
            }

            _amountMinted += _genesisAllocation[i].amount;
        }

        unchecked {
            tokenId += _amountMinted;
        }
    }

    /// @notice Mints a new ticket
    ///
    /// @param _to The address to mint the ticket to
    function mint(address _to) external {
        // Mint is supposed to be called through a cross chain message
        // We need to also check that the l1 sender is the auction contract
        if (
            msg.sender != Predeploys.L2_CROSS_DOMAIN_MESSENGER
                || ICrossDomainMessenger(msg.sender).xDomainMessageSender() != AUCTION
        ) revert NotAuction();

        uint256 _tokenId;

        // Not feasible for this to ever overflow
        unchecked {
            _tokenId = ++tokenId;
        }

        _mintTo(_to, _tokenId);
    }

    /// @notice Burns a ticket
    ///
    /// @param _target The address to burn the ticket for
    function burn(address _target) external {
        // This check might change in the future
        if (msg.sender != Constants.DEPOSITOR_ACCOUNT) revert NotSystemBurn();

        uint256 _topTicket = _top(_target);

        // Sanity check that there are tickets left
        if (_topTicket == SENTINEL_TICKET_ID) revert NoTicketsLeft();

        // Update the top pointer to the next ticket in the stack
        uint256 _nextTicket = ticketStack[_target][_topTicket];
        ticketStack[_target][SENTINEL_TICKET_ID] = _nextTicket;

        // Clear the reference for the burned ticket
        ticketStack[_target][_topTicket] = 0;

        // Remove the ticket from the stack
        _burn(_topTicket);
    }

    /// @notice Returns the top of the ticket stack for a given address
    ///
    /// @param _to The address to get the top of the stack for
    /// @return top_ The top of the stack
    function top(address _to) external view returns (uint256 top_) {
        top_ = _top(_to);
    }

    /// @notice Returns the ticket stack for a given address
    ///
    /// @param _to The address to get the stack for
    /// @return stack_ The ticket stack
    function traverseTicketStack(address _to) external view returns (uint256[] memory stack_) {
        uint256 _count = balanceOf(_to);
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

    /// @notice Overrides the ownerOf function to make it not revert on address(0)
    ///
    /// @param _tokenId The token id to get the owner of
    /// @return owner_ The owner of the token
    function ownerOf(uint256 _tokenId) public view override returns (address owner_) {
        // This can be optimized in future versions of openzeppelin, but that is a large dependency across the repo
        if (_exists(_tokenId)) {
            owner_ = super.ownerOf(_tokenId);
        }
    }

    /// @notice Returns the top of the ticket stack for a given address
    ///
    /// @param _to The address to get the top of the stack for
    /// @return top_ The top of the stack
    function _top(address _to) internal view returns (uint256 top_) {
        top_ = ticketStack[_to][SENTINEL_TICKET_ID];
    }

    /// @notice Mints a ticket to a given address
    /// @dev This function does not update the tokenId of this contract
    ///
    /// @param _to The address to mint the ticket to
    /// @param _tokenId The token id to mint
    function _mintTo(address _to, uint256 _tokenId) internal {
        uint256 _topTicket = _top(_to);

        if (_topTicket == 0) {
            // This is the first ticket for this address
            // Set the top to the token id
            ticketStack[_to][SENTINEL_TICKET_ID] = _tokenId;
        } else {
            // This is not the first ticket for this address
            // Move the previous top down the linked list
            ticketStack[_to][_tokenId] = _topTicket;
            // Set the top to the token id that just got minted
            ticketStack[_to][SENTINEL_TICKET_ID] = _tokenId;
        }

        _mint(_to, _tokenId);
    }

    /// @notice Overrides the transfer function to prevent tickets from being transferred
    function _transfer(address, address, uint256) internal pure override {
        revert Untransferable();
    }
}
