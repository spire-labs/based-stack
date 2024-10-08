// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

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
        NO_FALLBACK,
        CURRENT_PROPOSER,
        CURRENT_PROPOSER_WITH_CONFIG,
        NEXT_PROPOSER,
        NEXT_PROPOSER_WITH_CONFIG,
        RANDOM_TICKET_HOLDER,
        PERMISSIONLESS
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

    /// @notice Fetches the minimum preconfirmation collateral that is set
    ///
    /// @return minimumPreconfirmationCollateral_ The minimum preconfirmation collateral
    function minimumPreconfirmationCollateral() external view returns (uint256 minimumPreconfirmationCollateral_) {
        minimumPreconfirmationCollateral_ = _electionConfig.rules.minimumPreconfirmationCollateral;
    }

    /// @notice Fetches the election fallback list that is set
    ///
    /// @return electionFallbackList_ The election fallback list
    function electionFallbackList() external view returns (bytes32 electionFallbackList_) {
        electionFallbackList_ = _electionConfig.precedence.electionFallbackList;
    }

    /// @notice Fetches the election config that is set
    ///
    /// @return electionConfig_ The election config
    function electionConfig() external view returns (ElectionConfig memory electionConfig_) {
        electionConfig_ = _electionConfig;
    }

    /// @notice Updates the election queried by the offchain node for computing the election
    /// @param _config The config to update to
    function _setElectionConfig(ElectionConfig memory _config) internal {
        _electionConfig = _config;
    }
}
