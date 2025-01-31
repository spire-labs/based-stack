// SPDX-License-Identifier: MIT
pragma solidity 0.8.25;

interface ISystemConfig {
  function checkSequencerRules() external returns (bool);
  function checkSequencerRules(
    address _optionalInjectee
  ) external returns (bool);
}

/// @title BatchCheckSeqConfig
///
/// @dev This contract is used to batch check the sequencer rules for a list of potential sequencers.
///      The ticket accounting is expected to be handled externally, this batch just returns if the sequencer
///      supports the required rules.
contract BatchCheckSeqConfig {
  constructor(address _sysConfig, address[] memory _addrs) {
    bool[] memory results = new bool[](_addrs.length);

    for (uint256 i = 0; i < _addrs.length; i++) {
      // If we passed in zero address we mark it as false and skip
      if (_addrs[i] == address(0)) {
        results[i] = false;
        continue;
      }

      results[i] = ISystemConfig(_sysConfig).checkSequencerRules(_addrs[i]);
    }

    bytes memory data = abi.encode(results);

    assembly {
      let dataStart := add(data, 32)
      let dataEnd := sub(msize(), dataStart)
      return(dataStart, dataEnd)
    }
  }
}
