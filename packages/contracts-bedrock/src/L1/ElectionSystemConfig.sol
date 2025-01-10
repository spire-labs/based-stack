// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

/// @title ElectionSystemConfig
/// @notice This config file is abstract and inherited by the SystemConfig contracts
///         It contains all the datatypes and documentation for election related system configuration
abstract contract ElectionSystemConfig {
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

    /// @notice Sanitzes a fallback list before it is set
    ///
    /// @param _fallbackListAsBytes The fallback list to sanitze
    /// @return bool Whether the list is valid or not
    function _sanitzeFallbackList(bytes32 _fallbackListAsBytes) internal pure returns (bool) {
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
            if (_val > uint256(ElectionFallback.PERMISSIONLESS)) return false;

            unchecked {
                ++_byte;
            }
        }

        // If we did not loop and the list is not empty, this means the list is not right padded
        // Meaning its an invalid format
        if (!_didLoop && uint256(_fallbackListAsBytes) != 0) return false;

        return true;
    }
}
