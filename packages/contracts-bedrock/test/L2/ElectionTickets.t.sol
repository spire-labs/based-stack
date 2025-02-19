// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Testing
import { Test } from "forge-std/Test.sol";
import { Proxy } from "src/universal/Proxy.sol";

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

contract ElectionTickets_constructor_Test is ElectionTickets_Test {
    /// @dev Tests that the constructor sets the auction address.
    function test_constructor_auction_succeeds() public {
        ElectionTickets electionTicket = new ElectionTickets(election);
        assertEq(address(electionTicket.auction()), election);
    }
}

contract ElectionTickets_initialize_Test is ElectionTickets_Test {
    /// @dev Tests that the constructor mints the genesis tickets amount.
    function test_initialize_mintGenesisTickets_succeeds() public {
        ElectionTickets.GenesisAllocation[] memory _genAlloc = new ElectionTickets.GenesisAllocation[](3);

        _genAlloc[0] = ElectionTickets.GenesisAllocation(1, address(1));
        _genAlloc[1] = ElectionTickets.GenesisAllocation(2, address(2));
        _genAlloc[2] = ElectionTickets.GenesisAllocation(3, address(3));

        for (uint256 i; i < _genAlloc.length; i++) {
            vm.assume(_genAlloc[i].target != address(0));
            // hardcoded to make test runtime fast
            _genAlloc[i].amount = 3;
        }

        ElectionTickets electionTicket = new ElectionTickets(election);
        Proxy proxy = new Proxy(address(to));

        vm.prank(to);
        proxy.upgradeTo(address(electionTicket));

        electionTicket = ElectionTickets(address(proxy));

        electionTicket.initialize(_genAlloc);

        uint256 _amountMinted;

        for (uint256 i; i < _genAlloc.length; i++) {
            for (uint256 j; j < _genAlloc[i].amount; j++) {
                assertEq(electionTicket.ownerOf(_amountMinted + j + 1), _genAlloc[i].target);
            }

            _amountMinted += _genAlloc[i].amount;
        }

        assertEq(electionTicket.tokenId(), _amountMinted);
    }
}

contract ElectionTickets_mint_Test is ElectionTickets_Test {
    /// @dev Tests that the `mint` function reverts when called by a non-messenger address.
    function test_mint_onlyAuction_wrongCaller_reverts() public {
        vm.expectRevert(NotAuction.selector);
        electionTicket.mint(to, 1);
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
        electionTicket.mint(to, 1);
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
        electionTicket.mint(to, 1);
    }

    /// @dev Tests that `mint` correctly mints multiple tickets
    function test_mint_multipleTickets_succeeds(uint256 _amount) public {
        // Theoretical max based on the current validator in lookahead is 32
        vm.assume(_amount <= 32);

        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );
        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        electionTicket.mint(to, _amount);

        assertEq(electionTicket.tokenId(), _amount);

        for (uint256 i; i < _amount; i++) {
            assertEq(electionTicket.ownerOf(i + 1), to);
        }
    }

    /// @dev Tests that the mint function correctly updates the stack
    function test_mint_updatesTicketStack() public {
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );

        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        electionTicket.mint(to, 10);

        uint256[] memory stack = electionTicket.traverseTicketStack(to);

        // Top should be the last minted ticket
        assertEq(electionTicket.top(to), electionTicket.tokenId());

        for (uint256 i; i < 10; i++) {
            // Top of stack should be the last minted ticket so id 10
            assertEq(stack[i], electionTicket.tokenId() - i);
        }

        assertEq(electionTicket.balanceOf(to), 10);
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
        electionTicket.mint(to, 1);
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
        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        electionTicket.mint(to, 9);

        for (uint256 i = 0; i < 10; i++) {
            uint256 top = electionTicket.top(to);

            // Assert that the top of the stack is the last minted ticket
            assertEq(top, electionTicket.tokenId() - i);

            vm.prank(Constants.DEPOSITOR_ACCOUNT);
            electionTicket.burn(to);
        }

        // Stack should be empty
        assertEq(electionTicket.top(to), 0);
        assertEq(electionTicket.balanceOf(to), 0);
    }
}

contract ElectionTickets_Untransferable_Test is ElectionTickets_Test {
    address random = makeAddr("random");

    function setUp() public override {
        super.setUp();
        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );
        electionTicket.mint(to, 1);
    }

    /// @notice Test that the transferFrom function reverts
    function test_transferFrom_reverts() public {
        vm.expectRevert(abi.encodeWithSelector(Untransferable.selector));
        vm.prank(to);
        electionTicket.transferFrom(to, random, 1);
    }

    /// @notice Test that the safeTransferFrom function reverts
    function test_safeTransferFrom_reverts() public {
        vm.expectRevert(abi.encodeWithSelector(Untransferable.selector));
        vm.prank(to);
        electionTicket.safeTransferFrom(to, random, 1);
    }

    /// @notice Test that the safeTransferFrom function with data reverts
    function test_safeTransferFromWithData_reverts() public {
        vm.expectRevert(abi.encodeWithSelector(Untransferable.selector));
        vm.prank(to);
        electionTicket.safeTransferFrom(to, random, 1, "");
    }
}

contract ElectionTickets_OwnerOf_Test is ElectionTickets_Test {
    /// @dev Tests that the `ownerOf` function returns the correct owner.
    function test_ownerOf_succeeds() public {
        vm.mockCall(
            Predeploys.L2_CROSS_DOMAIN_MESSENGER,
            abi.encodeWithSelector(ICrossDomainMessenger.xDomainMessageSender.selector),
            abi.encode(election)
        );
        vm.prank(Predeploys.L2_CROSS_DOMAIN_MESSENGER);
        electionTicket.mint(to, 1);

        assertEq(electionTicket.ownerOf(1), to);
    }

    /// @dev Tests that the `ownerOf` function returns address(0) when the token does not exist.
    function test_ownerOf_doesNotExist_returnsZeroAddress() public view {
        assertEq(electionTicket.ownerOf(69), address(0));
    }
}
