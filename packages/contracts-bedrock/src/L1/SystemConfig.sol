// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Config inheritance
import { ElectionSystemConfig } from "src/L1/ElectionSystemConfig.sol";

// Contracts
import { OwnableUpgradeable } from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { ERC20 } from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

// Libraries
import { Storage } from "src/libraries/Storage.sol";
import { Constants } from "src/libraries/Constants.sol";
import { GasPayingToken, IGasToken } from "src/libraries/GasPayingToken.sol";

// Interfaces
import { ISemver } from "src/universal/interfaces/ISemver.sol";
import { IOptimismPortal } from "src/L1/interfaces/IOptimismPortal.sol";
import { IResourceMetering } from "src/L1/interfaces/IResourceMetering.sol";

/// @custom:proxied true
/// @title SystemConfig
/// @notice The SystemConfig contract is used to manage configuration of an Optimism network.
///         All configuration is stored on L1 and picked up by L2 as part of the derviation of
///         the L2 chain.
contract SystemConfig is OwnableUpgradeable, ElectionSystemConfig, ISemver, IGasToken {
    /// @notice Throws when a fallback list fails the sanity check
    error InvalidFallbackList();

    /// @notice Throws when a call is not an eth call
    error NotEthCall();

    /// @notice Throws when a given offset is out of bounds
    error OffsetOOB();

    /// @notice Throws when a given rule has exceeded the maximum number of rules
    error RuleOOB();

    /// @notice Throws when a given rule is not found
    error ConfigRuleNotFound();

    /// @notice Enum representing different types of updates.
    /// @custom:value BATCHER               Represents an update to the batcher hash.
    /// @custom:value GAS_CONFIG            Represents an update to txn fee config on L2.
    /// @custom:value GAS_LIMIT             Represents an update to gas limit on L2.
    /// @custom:value UNSAFE_BLOCK_SIGNER   Represents an update to the signer key for unsafe
    ///                                     block distrubution.
    /// @custom:value ELECTION_CONFIG       Represents an update to the election system config
    /// @custom:value DELETE_SEQUENCER_RULE Represents an update to delete a sequencer rule
    /// @custom:value INSERT_SEQUENCER_RULE Represents an update to insert a sequencer rule
    enum UpdateType {
        BATCHER,
        GAS_CONFIG,
        GAS_LIMIT,
        UNSAFE_BLOCK_SIGNER,
        ELECTION_CONFIG,
        DELETE_SEQUENCER_RULE,
        INSERT_SEQUENCER_RULE
    }

    /// @notice Struct representing the addresses of L1 system contracts. These should be the
    ///         contracts that users interact with (not implementations for proxied contracts)
    ///         and are network specific.
    struct Addresses {
        address l1CrossDomainMessenger;
        address l1ERC721Bridge;
        address l1StandardBridge;
        address disputeGameFactory;
        address optimismPortal;
        address optimismMintableERC20Factory;
        address gasPayingToken;
    }

    /// @notice Version identifier, used for upgrades.
    uint256 public constant VERSION = 0;

    /// @notice Storage slot that the unsafe block signer is stored at.
    ///         Storing it at this deterministic storage slot allows for decoupling the storage
    ///         layout from the way that `solc` lays out storage. The `op-node` uses a storage
    ///         proof to fetch this value.
    /// @dev    NOTE: this value will be migrated to another storage slot in a future version.
    ///         User input should not be placed in storage in this contract until this migration
    ///         happens. It is unlikely that keccak second preimage resistance will be broken,
    ///         but it is better to be safe than sorry.
    bytes32 public constant UNSAFE_BLOCK_SIGNER_SLOT = keccak256("systemconfig.unsafeblocksigner");

    /// @notice Storage slot that the L1CrossDomainMessenger address is stored at.
    bytes32 public constant L1_CROSS_DOMAIN_MESSENGER_SLOT =
        bytes32(uint256(keccak256("systemconfig.l1crossdomainmessenger")) - 1);

    /// @notice Storage slot that the L1ERC721Bridge address is stored at.
    bytes32 public constant L1_ERC_721_BRIDGE_SLOT = bytes32(uint256(keccak256("systemconfig.l1erc721bridge")) - 1);

    /// @notice Storage slot that the L1StandardBridge address is stored at.
    bytes32 public constant L1_STANDARD_BRIDGE_SLOT = bytes32(uint256(keccak256("systemconfig.l1standardbridge")) - 1);

    /// @notice Storage slot that the OptimismPortal address is stored at.
    bytes32 public constant OPTIMISM_PORTAL_SLOT = bytes32(uint256(keccak256("systemconfig.optimismportal")) - 1);

    /// @notice Storage slot that the OptimismMintableERC20Factory address is stored at.
    bytes32 public constant OPTIMISM_MINTABLE_ERC20_FACTORY_SLOT =
        bytes32(uint256(keccak256("systemconfig.optimismmintableerc20factory")) - 1);

    /// @notice Storage slot that the batch inbox address is stored at.
    bytes32 public constant BATCH_INBOX_SLOT = bytes32(uint256(keccak256("systemconfig.batchinbox")) - 1);

    /// @notice Storage slot for block at which the op-node can start searching for logs from.
    bytes32 public constant START_BLOCK_SLOT = bytes32(uint256(keccak256("systemconfig.startBlock")) - 1);

    /// @notice Storage slot for the DisputeGameFactory address.
    bytes32 public constant DISPUTE_GAME_FACTORY_SLOT =
        bytes32(uint256(keccak256("systemconfig.disputegamefactory")) - 1);

    /// @notice The maximum amount of allowed sequencer config rules
    uint256 public constant MAX_SEQUENCER_RULES = 32;

    /// @notice The number of decimals that the gas paying token has.
    uint8 internal constant GAS_PAYING_TOKEN_DECIMALS = 18;

    /// @notice The maximum gas limit that can be set for L2 blocks. This limit is used to enforce that the blocks
    ///         on L2 are not too large to process and prove. Over time, this value can be increased as various
    ///         optimizations and improvements are made to the system at large.
    uint64 internal constant MAX_GAS_LIMIT = 200_000_000;

    /// @notice Fixed L2 gas overhead. Used as part of the L2 fee calculation.
    ///         Deprecated since the Ecotone network upgrade
    uint256 public overhead;

    /// @notice Dynamic L2 gas overhead. Used as part of the L2 fee calculation.
    ///         The most significant byte is used to determine the version since the
    ///         Ecotone network upgrade.
    uint256 public scalar;

    /// @notice L2 block gas limit.
    uint64 public gasLimit;

    /// @notice Basefee scalar value. Part of the L2 fee calculation since the Ecotone network upgrade.
    uint32 public basefeeScalar;

    /// @notice Blobbasefee scalar value. Part of the L2 fee calculation since the Ecotone network upgrade.
    uint32 public blobbasefeeScalar;

    /// @notice The configuration for the deposit fee market.
    ///         Used by the OptimismPortal to meter the cost of buying L2 gas on L1.
    ///         Set as internal with a getter so that the struct is returned instead of a tuple.
    IResourceMetering.ResourceConfig internal _resourceConfig;

    /// @notice The storage slot that the election config is stored at
    ///
    /// @dev Set as internal and exposes a getter function to make it return a struct instead of tuple
    ElectionConfig internal _electionConfig;

    /// @notice Emitted when configuration is updated.
    /// @param version    SystemConfig version.
    /// @param updateType Type of update.
    /// @param data       Encoded update data.
    event ConfigUpdate(uint256 indexed version, UpdateType indexed updateType, bytes data);

    /// @notice Semantic version.
    /// @custom:semver 2.3.1-beta.12
    function version() public pure virtual returns (string memory) {
        return "2.3.1-beta.12";
    }

    /// @notice Constructs the SystemConfig contract. Cannot set
    ///         the owner to `address(0)` due to the Ownable contract's
    ///         implementation, so set it to `address(0xdEaD)`
    /// @dev    START_BLOCK_SLOT is set to type(uint256).max here so that it will be a dead value
    ///         in the singleton and is skipped by initialize when setting the start block.
    constructor() {
        Storage.setUint(START_BLOCK_SLOT, type(uint256).max);
        initialize({
            _owner: address(0xdEaD),
            _basefeeScalar: 0,
            _blobbasefeeScalar: 0,
            _gasLimit: 1,
            _unsafeBlockSigner: address(0),
            _config: IResourceMetering.ResourceConfig({
                maxResourceLimit: 1,
                elasticityMultiplier: 1,
                baseFeeMaxChangeDenominator: 2,
                minimumBaseFee: 0,
                systemTxMaxGas: 0,
                maximumBaseFee: 0
            }),
            _batchInbox: address(0),
            _addresses: SystemConfig.Addresses({
                l1CrossDomainMessenger: address(0),
                l1ERC721Bridge: address(0),
                l1StandardBridge: address(0),
                disputeGameFactory: address(0),
                optimismPortal: address(0),
                optimismMintableERC20Factory: address(0),
                gasPayingToken: address(0)
            }),
            _fallbackList: bytes32(0),
            _sequencerRules: new ElectionSystemConfig.SequencerRule[](0)
        });
    }

    /// @notice Initializer.
    ///         The resource config must be set before the require check.
    /// @param _owner             Initial owner of the contract.
    /// @param _basefeeScalar     Initial basefee scalar value.
    /// @param _blobbasefeeScalar Initial blobbasefee scalar value.
    /// @param _gasLimit          Initial gas limit.
    /// @param _unsafeBlockSigner Initial unsafe block signer address.
    /// @param _config            Initial ResourceConfig.
    /// @param _batchInbox        Batch inbox address. An identifier for the op-node to find
    ///                           canonical data.
    /// @param _addresses         Set of L1 contract addresses. These should be the proxies.
    /// @param _fallbackList      The defined election fallbacklist
    /// @param _sequencerRules    The sequencer rules to be initialized with
    function initialize(
        address _owner,
        uint32 _basefeeScalar,
        uint32 _blobbasefeeScalar,
        uint64 _gasLimit,
        address _unsafeBlockSigner,
        IResourceMetering.ResourceConfig memory _config,
        address _batchInbox,
        SystemConfig.Addresses memory _addresses,
        bytes32 _fallbackList,
        SequencerRule[] memory _sequencerRules
    )
        public
        initializer
    {
        __Ownable_init();
        transferOwnership(_owner);

        // These are set in ascending order of their UpdateTypes.
        _setGasConfigEcotone({ _basefeeScalar: _basefeeScalar, _blobbasefeeScalar: _blobbasefeeScalar });
        _setGasLimit(_gasLimit);

        Storage.setAddress(UNSAFE_BLOCK_SIGNER_SLOT, _unsafeBlockSigner);
        Storage.setAddress(BATCH_INBOX_SLOT, _batchInbox);
        Storage.setAddress(L1_CROSS_DOMAIN_MESSENGER_SLOT, _addresses.l1CrossDomainMessenger);
        Storage.setAddress(L1_ERC_721_BRIDGE_SLOT, _addresses.l1ERC721Bridge);
        Storage.setAddress(L1_STANDARD_BRIDGE_SLOT, _addresses.l1StandardBridge);
        Storage.setAddress(DISPUTE_GAME_FACTORY_SLOT, _addresses.disputeGameFactory);
        Storage.setAddress(OPTIMISM_PORTAL_SLOT, _addresses.optimismPortal);
        Storage.setAddress(OPTIMISM_MINTABLE_ERC20_FACTORY_SLOT, _addresses.optimismMintableERC20Factory);

        _setStartBlock();
        _setGasPayingToken(_addresses.gasPayingToken);
        if (!_sanitizeFallbackList(_fallbackList)) revert InvalidFallbackList();
        _setElectionFallbackList(_fallbackList);

        for (uint256 i; i < _sequencerRules.length; i++) {
            _setSequencerConfigRule(_sequencerRules[i]);
        }

        _setResourceConfig(_config);
        require(_gasLimit >= minimumGasLimit(), "SystemConfig: gas limit too low");
    }

    /// @notice Updates the election queried by the offchain node for computing the election
    ///
    /// @param _fallbackList The config to update to
    function setElectionFallbackList(bytes32 _fallbackList) external onlyOwner {
        bool _success = _sanitizeFallbackList(_fallbackList);

        if (!_success) revert InvalidFallbackList();

        _setElectionFallbackList(_fallbackList);

        bytes memory data = abi.encode(_fallbackList);
        emit ConfigUpdate(VERSION, UpdateType.ELECTION_CONFIG, data);
    }

    /// @notice Sets a sequencer rule in the config
    ///
    /// @param _rule The rule to be set
    function setSequencerConfigRule(SequencerRule memory _rule) external onlyOwner {
        _setSequencerConfigRule(_rule);

        emit ConfigUpdate(VERSION, UpdateType.INSERT_SEQUENCER_RULE, abi.encode(_rule));
    }

    /// @notice Removes a sequencer rule from the config
    ///
    /// @param _index The index of the rule to remove
    function removeSequencerConfigRule(uint256 _index) external onlyOwner {
        bytes32 _layout = _electionConfig.config.sequencerRulesLayout;

        // Extract the target byte
        bytes32 _shifted = _layout >> ((31 - _index) * 8);
        uint8 _targetByte = uint8(uint256(_shifted) & 0xFF);
        if (_targetByte == 0) revert ConfigRuleNotFound();

        // Clear the index
        bytes32 _mask = ~(bytes32(uint256(0xFF)) << ((31 - _index) * 8));
        _layout = _layout & _mask;

        // Update the state
        _electionConfig.config.sequencerRulesLayout = _layout;
        delete _electionConfig.config.rules[_index];

        emit ConfigUpdate(VERSION, UpdateType.DELETE_SEQUENCER_RULE, abi.encode(_index));
    }

    /// @notice Returns the sequencer rule at the given index
    ///
    /// @param _index The index of the rule to return
    ///
    /// @return SequencerRule The sequencer rule at the given index
    function getSequencerRuleAtIndex(uint256 _index) external view returns (SequencerRule memory) {
        bytes32 _layout = _electionConfig.config.sequencerRulesLayout;
        bytes32 _shifted = _layout >> ((31 - _index) * 8);
        uint8 _targetByte = uint8(uint256(_shifted) & 0xFF);
        if (_targetByte == 0) revert ConfigRuleNotFound();

        return _electionConfig.config.rules[_index];
    }

    /// @notice Returns the sequencer rules layout
    ///
    /// @return bytes32 The sequencer rules layout
    function sequencerRulesLayout() external view returns (bytes32) {
        return _electionConfig.config.sequencerRulesLayout;
    }

    /// @notice Returns the minimum L2 gas limit that can be safely set for the system to
    ///         operate. The L2 gas limit must be larger than or equal to the amount of
    ///         gas that is allocated for deposits per block plus the amount of gas that
    ///         is allocated for the system transaction.
    ///         This function is used to determine if changes to parameters are safe.
    /// @return uint64 Minimum gas limit.
    function minimumGasLimit() public view returns (uint64) {
        return uint64(_resourceConfig.maxResourceLimit) + uint64(_resourceConfig.systemTxMaxGas);
    }

    /// @notice Returns the maximum L2 gas limit that can be safely set for the system to
    ///         operate. This bound is used to prevent the gas limit from being set too high
    ///         and causing the system to be unable to process and/or prove L2 blocks.
    /// @return uint64 Maximum gas limit.
    function maximumGasLimit() public pure returns (uint64) {
        return MAX_GAS_LIMIT;
    }

    /// @notice High level getter for the unsafe block signer address.
    ///         Unsafe blocks can be propagated across the p2p network if they are signed by the
    ///         key corresponding to this address.
    /// @return addr_ Address of the unsafe block signer.
    function unsafeBlockSigner() public view returns (address addr_) {
        addr_ = Storage.getAddress(UNSAFE_BLOCK_SIGNER_SLOT);
    }

    /// @notice Getter for the L1CrossDomainMessenger address.
    function l1CrossDomainMessenger() external view returns (address addr_) {
        addr_ = Storage.getAddress(L1_CROSS_DOMAIN_MESSENGER_SLOT);
    }

    /// @notice Getter for the L1ERC721Bridge address.
    function l1ERC721Bridge() external view returns (address addr_) {
        addr_ = Storage.getAddress(L1_ERC_721_BRIDGE_SLOT);
    }

    /// @notice Getter for the L1StandardBridge address.
    function l1StandardBridge() external view returns (address addr_) {
        addr_ = Storage.getAddress(L1_STANDARD_BRIDGE_SLOT);
    }

    /// @notice Getter for the DisputeGameFactory address.
    function disputeGameFactory() external view returns (address addr_) {
        addr_ = Storage.getAddress(DISPUTE_GAME_FACTORY_SLOT);
    }

    /// @notice Getter for the OptimismPortal address.
    function optimismPortal() public view returns (address addr_) {
        addr_ = Storage.getAddress(OPTIMISM_PORTAL_SLOT);
    }

    /// @notice Checks the sequencer rules
    ///
    /// @dev This function is not marked as view because it can used as simulations
    ///      Due to this we restrict to only be callable through the context of an eth_call
    function checkSequencerRules() external returns (bool) {
        return _checkSequencerRules(address(0));
    }

    /// @notice Checks the sequencer
    ///
    /// @param _optionalInjectee The address to inject into the calldata
    ///
    /// @dev This function is not marked as view because it can used as simulations
    ///      Due to this we restrict to only be callable through the context of an eth_call
    function checkSequencerRules(address _optionalInjectee) external returns (bool) {
        return _checkSequencerRules(_optionalInjectee);
    }

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

    /// @notice Getter for the OptimismMintableERC20Factory address.
    function optimismMintableERC20Factory() external view returns (address addr_) {
        addr_ = Storage.getAddress(OPTIMISM_MINTABLE_ERC20_FACTORY_SLOT);
    }

    /// @notice Getter for the BatchInbox address.
    function batchInbox() external view returns (address addr_) {
        addr_ = Storage.getAddress(BATCH_INBOX_SLOT);
    }

    /// @notice Getter for the StartBlock number.
    function startBlock() external view returns (uint256 startBlock_) {
        startBlock_ = Storage.getUint(START_BLOCK_SLOT);
    }

    /// @notice Getter for the gas paying asset address.
    function gasPayingToken() public view returns (address addr_, uint8 decimals_) {
        (addr_, decimals_) = GasPayingToken.getToken();
    }

    /// @notice Getter for custom gas token paying networks. Returns true if the
    ///         network uses a custom gas token.
    function isCustomGasToken() public view returns (bool) {
        (address token,) = gasPayingToken();
        return token != Constants.ETHER;
    }

    /// @notice Getter for the gas paying token name.
    function gasPayingTokenName() external view returns (string memory name_) {
        name_ = GasPayingToken.getName();
    }

    /// @notice Getter for the gas paying token symbol.
    function gasPayingTokenSymbol() external view returns (string memory symbol_) {
        symbol_ = GasPayingToken.getSymbol();
    }

    /// @notice Returns the modified calldata with the injected address
    /// TODO(spire): Made public for now to make testing easier and might be used offchain, if not needed refactor to
    /// internal
    ///
    /// @param _calldata The calldata to inject the address into
    /// @param _offsets The offsets of the addresses to be dynamically injected inside the calldata
    /// @param _injectee The address to inject into the calldata
    ///
    /// @dev The offsets should be inputed to insert the padded address so for example the balanceOf selector
    ///      Would have an offset of 4
    function injectAddressIntoCalldata(
        bytes memory _calldata,
        uint256[] memory _offsets,
        address _injectee
    )
        public
        pure
        returns (bytes memory _newCalldata)
    {
        uint256 len = _offsets.length;
        _newCalldata = _calldata; // Initialize with original calldata

        // Convert the address to a 32-byte word (right-padded as Solidity handles addresses)
        uint256 injecteeWord = uint256(uint160(_injectee));

        for (uint256 i; i < len; i++) {
            uint256 offset = _offsets[i];

            if (offset + 20 > _calldata.length) revert OffsetOOB();

            // Inject the 20-byte address into calldata at the specified offset
            assembly {
                // Calculate the memory position
                let dataPtr := add(add(_newCalldata, 0x20), offset)

                // Load the existing word at the offset
                let existingWord := mload(dataPtr)

                // Mask the word to clear the bytes where the address will be injected
                let clearedWord := and(existingWord, not(0xffffffffffffffffffffffffffffffffffffffff))

                // Combine the cleared word with the address (injecteeWord) and store it back
                mstore(dataPtr, or(clearedWord, injecteeWord))
            }
        }
    }

    /// @notice Checks the sequencer rules
    ///
    /// @dev This function is not marked as view because it can used as simulations
    ///      Due to this we restrict to only be callable through the context of an eth_call
    function _checkSequencerRules(address _optionalInjectee) internal returns (bool) {
        // eth_call from field needs to be address(0)
        // We use tx.origin to allow for contracts to call this function
        if (tx.origin != address(0)) revert NotEthCall();

        SequencerRule memory _rule;
        bool _success;
        bytes memory _returnData;
        bytes memory _calldata;
        address _target;
        bytes32 _layout = _electionConfig.config.sequencerRulesLayout;

        for (uint256 i; i < MAX_SEQUENCER_RULES; i++) {
            if (_layout[i] == 0) continue;
            _rule = _electionConfig.config.rules[i];
            _target = _rule.target;

            // To avoid an extra storage read we check the injectee address instead
            // If the injectee is address zero we can assume there is no injectee
            if (_optionalInjectee != address(0)) {
                _calldata = injectAddressIntoCalldata(_rule.configCalldata, _rule.addressOffsets, _optionalInjectee);
            } else {
                _calldata = _rule.configCalldata;
            }

            (_success, _returnData) = _target.call(_calldata);

            // Check the desired assertion type
            if (_rule.assertionType == SequencerAssertion.SUCCESS) {
                if (!_success) return false;
            } else if (_rule.assertionType == SequencerAssertion.REVERT) {
                if (_success) return false;
            } else if (_rule.assertionType == SequencerAssertion.LT) {
                if (bytes32(_returnData) >= _rule.desiredRetdata) return false;
            } else if (_rule.assertionType == SequencerAssertion.GT) {
                if (bytes32(_returnData) <= _rule.desiredRetdata) return false;
            } else if (_rule.assertionType == SequencerAssertion.LTE) {
                if (bytes32(_returnData) > _rule.desiredRetdata) return false;
            } else if (_rule.assertionType == SequencerAssertion.GTE) {
                if (bytes32(_returnData) < _rule.desiredRetdata) return false;
            } else if (_rule.assertionType == SequencerAssertion.EQ) {
                if (bytes32(_returnData) != _rule.desiredRetdata) return false;
            } else if (_rule.assertionType == SequencerAssertion.NEQ) {
                if (bytes32(_returnData) == _rule.desiredRetdata) return false;
            }
        }

        return true;
    }

    /// @notice Sets a sequencer rule in the config
    ///
    /// @param _rule The rule to be set
    function _setSequencerConfigRule(SequencerRule memory _rule) internal {
        uint256 _i;
        bytes32 _layout = _electionConfig.config.sequencerRulesLayout;

        // If the layout is empty we dont need to check for anything
        if (_layout != bytes32(0)) {
            // We still need to start at zero incase of a deletion at zero
            for (_i; _i < MAX_SEQUENCER_RULES; _i++) {
                // Unoccupied slot was found
                if (_layout[_i] == bytes1(0)) {
                    break;
                }
            }

            // If we reach the max amount, not slots are available to be set into
            if (_i == MAX_SEQUENCER_RULES) revert RuleOOB();
        }

        // Mark the byte as occupied
        bytes32 _value = bytes32(1 << ((31 - _i) * 8));
        _layout = _layout | _value;

        _electionConfig.config.sequencerRulesLayout = _layout;

        _electionConfig.config.rules[_i] = _rule;
    }

    /// @notice Internal setter for the gas paying token address, includes validation.
    ///         The token must not already be set and must be non zero and not the ether address
    ///         to set the token address. This prevents the token address from being changed
    ///         and makes it explicitly opt-in to use custom gas token.
    /// @param _token Address of the gas paying token.
    function _setGasPayingToken(address _token) internal virtual {
        if (_token != address(0) && _token != Constants.ETHER && !isCustomGasToken()) {
            require(
                ERC20(_token).decimals() == GAS_PAYING_TOKEN_DECIMALS, "SystemConfig: bad decimals of gas paying token"
            );
            bytes32 name = GasPayingToken.sanitize(ERC20(_token).name());
            bytes32 symbol = GasPayingToken.sanitize(ERC20(_token).symbol());

            // Set the gas paying token in storage and in the OptimismPortal.
            GasPayingToken.set({ _token: _token, _decimals: GAS_PAYING_TOKEN_DECIMALS, _name: name, _symbol: symbol });
            IOptimismPortal(payable(optimismPortal())).setGasPayingToken({
                _token: _token,
                _decimals: GAS_PAYING_TOKEN_DECIMALS,
                _name: name,
                _symbol: symbol
            });
        }
    }

    /// @notice Updates the unsafe block signer address. Can only be called by the owner.
    /// @param _unsafeBlockSigner New unsafe block signer address.
    function setUnsafeBlockSigner(address _unsafeBlockSigner) external onlyOwner {
        _setUnsafeBlockSigner(_unsafeBlockSigner);
    }

    /// @notice Updates the unsafe block signer address.
    /// @param _unsafeBlockSigner New unsafe block signer address.
    function _setUnsafeBlockSigner(address _unsafeBlockSigner) internal {
        Storage.setAddress(UNSAFE_BLOCK_SIGNER_SLOT, _unsafeBlockSigner);

        bytes memory data = abi.encode(_unsafeBlockSigner);
        emit ConfigUpdate(VERSION, UpdateType.UNSAFE_BLOCK_SIGNER, data);
    }

    /// @notice Updates gas config. Can only be called by the owner.
    ///         Deprecated in favor of setGasConfigEcotone since the Ecotone upgrade.
    /// @param _overhead New overhead value.
    /// @param _scalar   New scalar value.
    function setGasConfig(uint256 _overhead, uint256 _scalar) external onlyOwner {
        _setGasConfig(_overhead, _scalar);
    }

    /// @notice Internal function for updating the gas config.
    /// @param _overhead New overhead value.
    /// @param _scalar   New scalar value.
    function _setGasConfig(uint256 _overhead, uint256 _scalar) internal {
        require((uint256(0xff) << 248) & _scalar == 0, "SystemConfig: scalar exceeds max.");

        overhead = _overhead;
        scalar = _scalar;

        bytes memory data = abi.encode(_overhead, _scalar);
        emit ConfigUpdate(VERSION, UpdateType.GAS_CONFIG, data);
    }

    /// @notice Updates gas config as of the Ecotone upgrade. Can only be called by the owner.
    /// @param _basefeeScalar     New basefeeScalar value.
    /// @param _blobbasefeeScalar New blobbasefeeScalar value.
    function setGasConfigEcotone(uint32 _basefeeScalar, uint32 _blobbasefeeScalar) external onlyOwner {
        _setGasConfigEcotone(_basefeeScalar, _blobbasefeeScalar);
    }

    /// @notice Internal function for updating the fee scalars as of the Ecotone upgrade.
    /// @param _basefeeScalar     New basefeeScalar value.
    /// @param _blobbasefeeScalar New blobbasefeeScalar value.
    function _setGasConfigEcotone(uint32 _basefeeScalar, uint32 _blobbasefeeScalar) internal {
        basefeeScalar = _basefeeScalar;
        blobbasefeeScalar = _blobbasefeeScalar;

        scalar = (uint256(0x01) << 248) | (uint256(_blobbasefeeScalar) << 32) | _basefeeScalar;

        bytes memory data = abi.encode(overhead, scalar);
        emit ConfigUpdate(VERSION, UpdateType.GAS_CONFIG, data);
    }

    /// @notice Updates the L2 gas limit. Can only be called by the owner.
    /// @param _gasLimit New gas limit.
    function setGasLimit(uint64 _gasLimit) external onlyOwner {
        _setGasLimit(_gasLimit);
    }

    /// @notice Internal function for updating the L2 gas limit.
    /// @param _gasLimit New gas limit.
    function _setGasLimit(uint64 _gasLimit) internal {
        require(_gasLimit >= minimumGasLimit(), "SystemConfig: gas limit too low");
        require(_gasLimit <= maximumGasLimit(), "SystemConfig: gas limit too high");
        gasLimit = _gasLimit;

        bytes memory data = abi.encode(_gasLimit);
        emit ConfigUpdate(VERSION, UpdateType.GAS_LIMIT, data);
    }

    /// @notice Sets the start block in a backwards compatible way. Proxies
    ///         that were initialized before the startBlock existed in storage
    ///         can have their start block set by a user provided override.
    ///         A start block of 0 indicates that there is no override and the
    ///         start block will be set by `block.number`.
    /// @dev    This logic is used to patch legacy deployments with new storage values.
    ///         Use the override if it is provided as a non zero value and the value
    ///         has not already been set in storage. Use `block.number` if the value
    ///         has already been set in storage
    function _setStartBlock() internal {
        if (Storage.getUint(START_BLOCK_SLOT) == 0) {
            Storage.setUint(START_BLOCK_SLOT, block.number);
        }
    }

    /// @notice Updates the election fallback list queried by the offchain node for computing the election
    ///
    /// @param _fallbackList The config to update to
    function _setElectionFallbackList(bytes32 _fallbackList) internal {
        _electionConfig.precedence.electionFallbackList = _fallbackList;
    }

    /// @notice A getter for the resource config.
    ///         Ensures that the struct is returned instead of a tuple.
    /// @return ResourceConfig
    function resourceConfig() external view returns (IResourceMetering.ResourceConfig memory) {
        return _resourceConfig;
    }

    /// @notice An internal setter for the resource config.
    ///         Ensures that the config is sane before storing it by checking for invariants.
    ///         In the future, this method may emit an event that the `op-node` picks up
    ///         for when the resource config is changed.
    /// @param _config The new resource config.
    function _setResourceConfig(IResourceMetering.ResourceConfig memory _config) internal {
        // Min base fee must be less than or equal to max base fee.
        require(
            _config.minimumBaseFee <= _config.maximumBaseFee, "SystemConfig: min base fee must be less than max base"
        );
        // Base fee change denominator must be greater than 1.
        require(_config.baseFeeMaxChangeDenominator > 1, "SystemConfig: denominator must be larger than 1");
        // Max resource limit plus system tx gas must be less than or equal to the L2 gas limit.
        // The gas limit must be increased before these values can be increased.
        require(_config.maxResourceLimit + _config.systemTxMaxGas <= gasLimit, "SystemConfig: gas limit too low");
        // Elasticity multiplier must be greater than 0.
        require(_config.elasticityMultiplier > 0, "SystemConfig: elasticity multiplier cannot be 0");
        // No precision loss when computing target resource limit.
        require(
            ((_config.maxResourceLimit / _config.elasticityMultiplier) * _config.elasticityMultiplier)
                == _config.maxResourceLimit,
            "SystemConfig: precision loss with target resource limit"
        );

        _resourceConfig = _config;
    }
}
