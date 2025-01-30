// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

import {BatchCheckSeqConfig, ISystemConfig} from 'src/BatchCheckSeqConfig.sol';
import {Utils} from 'test/Utils.sol';

contract Unit_BatchCheckSeqConfig is Utils {
  address public sysConfig = makeAddr('SystemConfig');

  /// @dev Test that the BatchCheckSeqConfig contract returns true if the sequencer supports the config
  function test_supportsConfig_succeeds() public {
    address[] memory addrs = new address[](1);
    addrs[0] = makeAddr('addr');

    _mockAndExpect(sysConfig, abi.encodeWithSignature('checkSequencerRules(address)', addrs[0]), abi.encode(true));
    BatchCheckSeqConfig batchCheckSeqConfig = new BatchCheckSeqConfig(sysConfig, addrs);

    bytes memory retdata = address(batchCheckSeqConfig).code;
    bool[] memory results = abi.decode(retdata, (bool[]));

    assertTrue(results[0]);
    assertEq(results.length, 1);
  }

  /// @dev Test that the BatchCheckSeqConfig contract returns false if the sequencer does not support the config
  function test_doesNotSupportConfig_succeeds() public {
    address[] memory addrs = new address[](1);
    addrs[0] = makeAddr('addr');

    _mockAndExpect(sysConfig, abi.encodeWithSignature('checkSequencerRules(address)', addrs[0]), abi.encode(false));
    BatchCheckSeqConfig batchCheckSeqConfig = new BatchCheckSeqConfig(sysConfig, addrs);

    bytes memory retdata = address(batchCheckSeqConfig).code;
    bool[] memory results = abi.decode(retdata, (bool[]));

    assertFalse(results[0]);
    assertEq(results.length, 1);
  }

  /// @dev Test that the BatchCheckSeqConfig contract returns varying results based on the config
  function test_multipleAddresses_succeeds() public {
    address[] memory addrs = new address[](2);
    addrs[0] = makeAddr('addr1');
    addrs[1] = makeAddr('addr2');

    _mockAndExpect(sysConfig, abi.encodeWithSignature('checkSequencerRules(address)', addrs[0]), abi.encode(true));
    _mockAndExpect(sysConfig, abi.encodeWithSignature('checkSequencerRules(address)', addrs[1]), abi.encode(false));
    BatchCheckSeqConfig batchCheckSeqConfig = new BatchCheckSeqConfig(sysConfig, addrs);

    bytes memory retdata = address(batchCheckSeqConfig).code;
    bool[] memory results = abi.decode(retdata, (bool[]));

    assertTrue(results[0]);
    assertFalse(results[1]);
    assertEq(results.length, 2);
  }

  /// @dev Tests that the BatchCheckSeqConfig contract returns false if no address is passed in
  function test_zeroAddressIsSkipped_succeeds() public {
    address[] memory addrs = new address[](1);
    addrs[0] = address(0);

    BatchCheckSeqConfig batchCheckSeqConfig = new BatchCheckSeqConfig(sysConfig, addrs);

    bytes memory retdata = address(batchCheckSeqConfig).code;
    bool[] memory results = abi.decode(retdata, (bool[]));

    assertFalse(results[0]);
    assertEq(results.length, 1);
  }
}
