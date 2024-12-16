// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

// External dependencies
import {Test} from 'forge-std/Test.sol';

contract Utils is Test {
  function _mockAndExpect(address target, bytes memory input, bytes memory retdata) internal {
    vm.mockCall(target, input, retdata);
    vm.expectCall(target, input);
  }
}
