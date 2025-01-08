// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import { IResourceMetering } from "src/L1/interfaces/IResourceMetering.sol";
import { ElectionSystemConfig } from "src/L1/ElectionSystemConfig.sol";

interface ISystemConfig {
    enum UpdateType {
        BATCHER,
        GAS_CONFIG,
        GAS_LIMIT,
        UNSAFE_BLOCK_SIGNER,
        ELECTION_CONFIG
    }

    enum ElectionFallback {
        NO_FALLBACK,
        CURRENT_PROPOSER,
        CURRENT_PROPOSER_WITH_CONFIG,
        NEXT_PROPOSER,
        NEXT_PROPOSER_WITH_CONFIG,
        RANDOM_TICKET_HOLDER,
        PERMISSIONLESS
    }

    struct Addresses {
        address l1CrossDomainMessenger;
        address l1ERC721Bridge;
        address l1StandardBridge;
        address disputeGameFactory;
        address optimismPortal;
        address optimismMintableERC20Factory;
        address gasPayingToken;
    }

    error InvalidFallbackList();

    event ConfigUpdate(uint256 indexed version, UpdateType indexed updateType, bytes data);
    event Initialized(uint8 version);
    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    function BATCH_INBOX_SLOT() external view returns (bytes32);
    function DISPUTE_GAME_FACTORY_SLOT() external view returns (bytes32);
    function L1_CROSS_DOMAIN_MESSENGER_SLOT() external view returns (bytes32);
    function L1_ERC_721_BRIDGE_SLOT() external view returns (bytes32);
    function L1_STANDARD_BRIDGE_SLOT() external view returns (bytes32);
    function OPTIMISM_MINTABLE_ERC20_FACTORY_SLOT() external view returns (bytes32);
    function OPTIMISM_PORTAL_SLOT() external view returns (bytes32);
    function START_BLOCK_SLOT() external view returns (bytes32);
    function UNSAFE_BLOCK_SIGNER_SLOT() external view returns (bytes32);
    function VERSION() external view returns (uint256);
    function basefeeScalar() external view returns (uint32);
    function batchInbox() external view returns (address addr_);
    function batcherHash() external view returns (bytes32);
    function blobbasefeeScalar() external view returns (uint32);
    function disputeGameFactory() external view returns (address addr_);
    function gasLimit() external view returns (uint64);
    function gasPayingToken() external view returns (address addr_, uint8 decimals_);
    function gasPayingTokenName() external view returns (string memory name_);
    function gasPayingTokenSymbol() external view returns (string memory symbol_);
    function initialize(
        address _owner,
        uint32 _basefeeScalar,
        uint32 _blobbasefeeScalar,
        bytes32 _batcherHash,
        uint64 _gasLimit,
        address _unsafeBlockSigner,
        IResourceMetering.ResourceConfig memory _config,
        address _batchInbox,
        Addresses memory _addresses,
        ElectionSystemConfig.ElectionConfig memory _eConfig
    )
        external;
    function isCustomGasToken() external view returns (bool);
    function l1CrossDomainMessenger() external view returns (address addr_);
    function l1ERC721Bridge() external view returns (address addr_);
    function l1StandardBridge() external view returns (address addr_);
    function maximumGasLimit() external pure returns (uint64);
    function minimumGasLimit() external view returns (uint64);
    function optimismMintableERC20Factory() external view returns (address addr_);
    function optimismPortal() external view returns (address addr_);
    function overhead() external view returns (uint256);
    function owner() external view returns (address);
    function renounceOwnership() external;
    function resourceConfig() external view returns (IResourceMetering.ResourceConfig memory);
    function scalar() external view returns (uint256);
    function setBatcherHash(bytes32 _batcherHash) external;
    function setGasConfig(uint256 _overhead, uint256 _scalar) external;
    function setGasConfigEcotone(uint32 _basefeeScalar, uint32 _blobbasefeeScalar) external;
    function setGasLimit(uint64 _gasLimit) external;
    function setUnsafeBlockSigner(address _unsafeBlockSigner) external;
    function startBlock() external view returns (uint256 startBlock_);
    function transferOwnership(address newOwner) external;
    function unsafeBlockSigner() external view returns (address addr_);
    function version() external pure returns (string memory);

    // ElectionSystemConfig
    function setElectionConfig(ElectionSystemConfig.ElectionConfig memory _config) external;
    function minimumPreconfirmationCollateral() external view returns (uint256 minimumPreconfirmationCollateral_);
    function electionFallbackList()
        external
        view
        returns (ElectionSystemConfig.ElectionFallback[] memory electionFallbackList_);
    function electionConfig() external view returns (ElectionSystemConfig.ElectionConfig memory electionConfig_);
}
