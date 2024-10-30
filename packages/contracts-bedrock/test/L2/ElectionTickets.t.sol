// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Test } from "forge-std/Test.sol";
import { ElectionTickets } from "src/L2/ElectionTickets.sol";
import { Constants } from "src/libraries/Constants.sol";
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
    /// @dev Tests that the `mint` function reverts when called by a non-Election address.
    function test_mint_onlyAuction_reverts() public {
        vm.expectRevert(NotAuction.selector);
        electionTicket.mint(to);
    }

    /// @dev Tests that the `mint` function succeeds when called by the Election contract.
    function test_mint_succeeds() public {
        vm.prank(election);
        electionTicket.mint(to);
    }
}

contract ElectionTickets_burn_Test is ElectionTickets_Test {
    function setUp() public override {
        super.setUp();
        vm.prank(election);
        electionTicket.mint(to);
    }

    /// @dev Tests that the `burn` function reverts when called by a non-BatchInbox address.
    function test_burn_onlyBatchInbox_reverts() public {
        vm.expectRevert(NotSystemBurn.selector);
        electionTicket.burn(1);
    }

    /// @dev Tests that the `burn` function reverts when the token id is invalid.
    function test_burn_succeeds() public {
        vm.prank(Constants.DEPOSITOR_ACCOUNT);
        electionTicket.burn(1);
    }
}
