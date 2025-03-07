// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { OPStackManager } from "src/L1/OPStackManager.sol";
import { SuperchainConfig } from "src/L1/SuperchainConfig.sol";
import { ProtocolVersions } from "src/L1/ProtocolVersions.sol";
import { ResourceMetering } from "src/L1/ResourceMetering.sol";
import { SystemConfig } from "src/L1/SystemConfig.sol";
import { SystemConfigInterop } from "src/L1/SystemConfigInterop.sol";

/// @custom:proxied true
contract OPStackManagerInterop is OPStackManager {
    constructor(
        SuperchainConfig _superchainConfig,
        ProtocolVersions _protocolVersions
    )
        OPStackManager(_superchainConfig, _protocolVersions)
    { }

    // The `SystemConfigInterop` contract has an extra `address _dependencyManager` argument
    // that we must account for.
    function encodeSystemConfigInitializer(
        bytes4 selector,
        DeployInput memory _input,
        DeployOutput memory _output
    )
        internal
        view
        virtual
        override
        returns (bytes memory)
    {
        (ResourceMetering.ResourceConfig memory referenceResourceConfig, SystemConfig.Addresses memory opChainAddrs) =
            defaultSystemConfigParams(selector, _input, _output);

        // TODO For now we assume that the dependency manager is the same as the proxy admin owner.
        // This is currently undefined since it's not part of the standard config, so we may need
        // to update where this value is pulled from in the future. To support a different dependency
        // manager in this contract without an invasive change of redefining the `Roles` struct,
        // we will make the change described in https://github.com/ethereum-optimism/optimism/issues/11783.
        return _getCalldata(selector, _input, referenceResourceConfig, opChainAddrs);
    }

    /// Using a new internal function to avoid stack too deep errors.
    function _getCalldata(
        bytes4 selector,
        DeployInput memory _input,
        ResourceMetering.ResourceConfig memory referenceResourceConfig,
        SystemConfig.Addresses memory opChainAddrs
    )
        internal
        pure
        returns (bytes memory)
    {
        return abi.encodeWithSelector(
            selector,
            _input.roles.systemConfigOwner,
            _input.basefeeScalar,
            _input.blobBasefeeScalar,
            30_000_000, // gasLimit TODO make this configurable?
            _input.roles.unsafeBlockSigner,
            referenceResourceConfig,
            chainIdToBatchInboxAddress(_input.l2ChainId),
            opChainAddrs,
            _input.electionFallbackList,
            _input.roles.opChainProxyAdminOwner,
            _input.sequencerRules
        );
    }
}
