// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { ERC721 } from "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract ElectionTickets is ERC721 {
    constructor() ERC721("ElectionTickets", "ET") { }
}
