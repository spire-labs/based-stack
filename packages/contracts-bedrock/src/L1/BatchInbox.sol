// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "./Election.sol";

contract BatchInbox {
    Election public electionAddress;

    constructor(Election _election) {
        electionAddress = _election;
    }

    // There's potential problem with transaction inclusion
    // in the specified block IF the sequencer is not the L1 proposer.
    modifier onlyElectionWinner() {
        address electionWinner = electionAddress.electionWinner();
        require(msg.sender == electionWinner, "Only election winner");
        _;
    }

    function submit() public onlyElectionWinner { }
}
