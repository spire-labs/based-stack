// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { ElectionTickets } from "src/L2/ElectionTickets.sol";
import { Predeploys } from "src/libraries/Predeploys.sol";
import { Constants } from "src/libraries/Constants.sol";
import { ICrossDomainMessenger } from "src/universal/interfaces/ICrossDomainMessenger.sol";
import "src/libraries/ElectionTicketErrors.sol";

contract ElectionTickets_Test is Test {
    ElectionTickets public electionTicket;
    address public election = makeAddr("election");
    address public to = makeAddr("to");

    function setUp() public virtual {
        electionTicket = new ElectionTickets(election);
    }
}

contract ElectionTickets_mint_Test is ElectionTickets_Test {
    /// @dev Tests that the `mint` function reverts when called by a non-messenger address.
    function test_mint_onlyAuction_wrongCaller_reverts() public {
        vm.expectRevert(NotAuction.selector);
        electionTicket.mint(to);
    }

    /// @dev Tests that the `mint` function reverts when l1 sender is not the auction contract.
    function test_mint_onlyAuction_wrongL1Sender_reverts() public {
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(address(0))
        );

        vm.expectRevert(NotAuction.selector);
        vm.prank(Predeploys.L1_BLOCK_ATTRIBUTES);
        electionTicket.mint(to);
    }

    /// @dev Tests that the `mint` function succeeds when called by the Election contract.
    function test_mint_succeeds() public {
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );
        vm.expectCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector)
        );
        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        electionTicket.mint(to);
    }

    /// @dev Tests that the mint function correctly updates the stack
    function test_mint_updatesTicketStack() public {
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );

        for (uint256 i = 0; i < 10; i++) {
            vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
            electionTicket.mint(to);
        }

        uint256[] memory stack = electionTicket.traverseTicketStack(to);

        // Top should be the last minted ticket
        assertEq(electionTicket.top(to), electionTicket.tokenId());

        for (uint256 i; i < 10; i++) {
            // Top of stack should be the last minted ticket so id 10
            assertEq(stack[i], electionTicket.tokenId() - i);
        }

        assertEq(electionTicket.ticketCount(to), 10);
    }
}

contract ElectionTickets_burn_Test is ElectionTickets_Test {
    function setUp() public override {
        super.setUp();
        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );
        electionTicket.mint(to);
    }

    /// @dev Tests that the `burn` function reverts when called by a non-BatchInbox address.
    function test_burn_onlyBatchInbox_reverts() public {
        vm.expectRevert(NotSystemBurn.selector);
        electionTicket.burn(to);
    }

    /// @dev Tests that the `burn` function reverts when there are no tickets left.
    function test_burn_noTicketsLeft_reverts() public {
        // Burn the ticket minted in the setup
        vm.prank(Constants.DEPOSITOR_ACCOUNT);
        electionTicket.burn(to);

        // Attempt to burn the ticket again
        vm.expectRevert(NoTicketsLeft.selector);
        vm.prank(Constants.DEPOSITOR_ACCOUNT);
        electionTicket.burn(to);
    }

    /// @dev Tests that the `burn` function reverts when the token id is invalid.
    function test_burn_succeeds() public {
        vm.prank(Constants.DEPOSITOR_ACCOUNT);
        electionTicket.burn(to);
    }

    /// @dev Test that the burn function correctly updates the stack
    function test_burn_traversesStack() public {
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );

        // Including the mint in setup this mints 10 tickets to "to"
        for (uint256 i = 0; i < 9; i++) {
            vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
            electionTicket.mint(to);
        }

        for (uint256 i = 0; i < 10; i++) {
            uint256 top = electionTicket.top(to);

            // Assert that the top of the stack is the last minted ticket
            assertEq(top, electionTicket.tokenId() - i);

            vm.prank(Constants.DEPOSITOR_ACCOUNT);
            electionTicket.burn(to);
        }

        // Stack should be empty
        assertEq(electionTicket.top(to), 0);
        assertEq(electionTicket.ticketCount(to), 0);
    }
}
