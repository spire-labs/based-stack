// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title ElectionSystemConfig
/// @notice This config file is abstract and inherited by the SystemConfig contracts
///         It contains all the logic to manage election related system configuration
abstract contract ElectionSystemConfig {
    /// @notice Throws when a fallback list fails the sanity check
    error InvalidFallbackList();

    /// @notice Enum representing different fallback rules
    ///
    /// @custom:value NO_FALLBACK                       Indicates there is no fallback left
    /// @custom:value CURRENT_PROPOSER                  Indicates that the current fallback is the current L1 proposer
    /// @custom:value CURRENT_PROPOSER_WITH_CONFIG      Indicates that the current fallback is the current L1 proposer,
    /// with enforced config
    /// @custom:value NEXT_PROPOSER                     Indicates that the current fallback is the next l1 proposer in
    /// the lookahead
    /// @custom:value NEXT_PROPOSER_WITH_CONFIG         Indicates that the current fallback is the next l1 proposer in
    /// the lookahead, with enforced config
    /// @custom:value RANDOM_TICKET_HOLDER              Indicates that the fallback will be a pseudo random ticket
    /// holder
    /// @custom:value PERMISSIONLESS                    Indicates tat the fallback is completely permissionless
    enum ElectionFallback {
        NO_FALLBACK, // 0x00
        CURRENT_PROPOSER, // 0x01
        CURRENT_PROPOSER_WITH_CONFIG, // 0x02
        NEXT_PROPOSER, // 0x03
        NEXT_PROPOSER_WITH_CONFIG, // 0x04
        RANDOM_TICKET_HOLDER, // 0x05
        PERMISSIONLESS // 0x06

    }

    // NOTE: The following structs are single values
    // but are defined as structs because we will likely add more to them in future iterations

    /// @notice Struct for storing the rules of election precedence
    /// @dev This is a lower level struct meant to be used in sync with ElectionConfigRules
    ///
    /// @param electionFallbackList A list going through the configured order of precedence, each byte represents a
    /// fallback rule
    struct ElectionPrecedence {
        bytes32 electionFallbackList;
    }

    /// @notice The available rules to enforce on the elected sequencers
    /// @dev This is a lower level struct meant to be used in sync with ElectionPrecedence
    ///
    /// @param minimumPreconfirmationCollateral The minimum amount of collateral required
    struct ElectionConfigRules {
        uint256 minimumPreconfirmationCollateral;
    }

    /// @notice The configuration for the election
    ///
    /// @param rules The defined election rules
    /// @param precedence The defined order of precedence
    struct ElectionConfig {
        ElectionConfigRules rules;
        ElectionPrecedence precedence;
    }

    /// @notice The storage slot that the election config is stored at
    ///
    /// @dev Set as internal and exposes a getter function to make it return a struct instead of tuple
    ElectionConfig internal _electionConfig;

    /// @notice Storage gap for future upgrades
    uint256[50] private __gap;

    /// @notice Fetches the minimum preconfirmation collateral that is set
    ///
    /// @return minimumPreconfirmationCollateral_ The minimum preconfirmation collateral
    function minimumPreconfirmationCollateral() external view returns (uint256 minimumPreconfirmationCollateral_) {
        minimumPreconfirmationCollateral_ = _electionConfig.rules.minimumPreconfirmationCollateral;
    }

    event log2(uint256 _index, bytes1 _byte);
    /// @notice Fetches the election fallback list that is set
    ///
    /// @return electionFallbackList_ The election fallback list

    function electionFallbackList() external view returns (ElectionFallback[] memory electionFallbackList_) {
        // The list is intended to be a right padded hexadecimal string
        // Each byte represents an ElectionFallback enum value
        bytes32 _fallbackList = _electionConfig.precedence.electionFallbackList;

        if (_fallbackList == bytes32(0)) return electionFallbackList_;

        bytes memory _listAsBytes = abi.encode(_fallbackList);

        uint256 _byte;
        uint256 _val;

        assembly {
            // Allocate memory for the array
            electionFallbackList_ := mload(0x40)
        }

        // If we encounter byte 00 (NO_FALLBACK) we know we've reached the end of the list
        while (uint256(uint8(_listAsBytes[_byte])) != uint256(ElectionFallback.NO_FALLBACK)) {
            _val = uint256(uint8(_listAsBytes[_byte]));

            unchecked {
                ++_byte;
            }

            // Dynamically resize the array
            assembly {
                // Store the length of the array
                mstore(electionFallbackList_, _byte)

                // Get the location to store the value in and store it
                let _dataLocation := add(electionFallbackList_, mul(_byte, 0x20))
                mstore(_dataLocation, _val)

                // Update the free memory pointer
                mstore(0x40, add(_dataLocation, 0x20))
            }
        }
    }

    /// @notice Fetches the election config that is set
    ///
    /// @return electionConfig_ The election config
    function electionConfig() external view returns (ElectionConfig memory electionConfig_) {
        electionConfig_ = _electionConfig;
    }

    /// @notice Updates the election queried by the offchain node for computing the election
    ///
    /// @param _config The config to update to
    function _setElectionConfig(ElectionConfig memory _config) internal {
        _sanitzeFallbackList(_config.precedence.electionFallbackList);
        _electionConfig = _config;
    }

    /// @notice Sanitzes a fallback list before it is set
    ///
    /// @param _fallbackListAsBytes The fallback list to sanitze
    function _sanitzeFallbackList(bytes32 _fallbackListAsBytes) internal pure {
        // The list is intended to be a right padded hexadecimal string
        // Each byte represents an ElectionFallback enum value
        bytes memory _listAsBytes = abi.encode(_fallbackListAsBytes);

        uint256 _val;
        uint256 _byte;
        bool _didLoop;

        // If we encounter byte 00 (NO_FALLBACK) we know we've reached the end of the list
        while (uint256(uint8(_listAsBytes[_byte])) != uint256(ElectionFallback.NO_FALLBACK)) {
            if (!_didLoop) _didLoop = true;

            _val = uint256(uint8(_listAsBytes[_byte]));

            // The list contains an invalid enum
            if (_val > uint256(ElectionFallback.PERMISSIONLESS)) revert InvalidFallbackList();

            unchecked {
                ++_byte;
            }
        }

        // If we did not loop and the list is not empty, this means the list is not right padded
        // Meaning its an invalid format
        if (!_didLoop && uint256(_fallbackListAsBytes) != 0) revert InvalidFallbackList();
    }
}
