// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { Script } from "forge-std/Script.sol";
import { console2 as console } from "forge-std/console2.sol";
import { stdJson } from "forge-std/StdJson.sol";
import { Executables } from "scripts/libraries/Executables.sol";
import { Process } from "scripts/libraries/Process.sol";
import { Chains } from "scripts/libraries/Chains.sol";
import { Config, Fork, ForkUtils } from "scripts/libraries/Config.sol";
import { ElectionTickets } from "src/L2/ElectionTickets.sol";
import { ElectionSystemConfig } from "src/L1/ElectionSystemConfig.sol";

/// @title DeployConfig
/// @notice Represents the configuration required to deploy the system. It is expected
///         to read the file from JSON. A future improvement would be to have fallback
///         values if they are not defined in the JSON themselves.
contract DeployConfig is Script {
    using stdJson for string;
    using ForkUtils for Fork;

    /// @notice Represents an unset offset value, as opposed to 0, which denotes no-offset.
    uint256 constant NULL_OFFSET = type(uint256).max;

    string internal _json;

    address public finalSystemOwner;
    address public superchainConfigGuardian;
    uint256 public l1ChainID;
    uint256 public l2ChainID;
    uint256 public l2BlockTime;
    uint256 public l2GenesisDeltaTimeOffset;
    uint256 public l2GenesisEcotoneTimeOffset;
    uint256 public l2GenesisFjordTimeOffset;
    uint256 public l2GenesisGraniteTimeOffset;
    uint256 public maxSequencerDrift;
    uint256 public sequencerWindowSize;
    uint256 public channelTimeout;
    address public p2pSequencerAddress;
    address public batchSenderAddress;
    uint256 public l2OutputOracleSubmissionInterval;
    int256 internal _l2OutputOracleStartingTimestamp;
    uint256 public l2OutputOracleStartingBlockNumber;
    address public l2OutputOracleProposer;
    address public l2OutputOracleChallenger;
    uint256 public finalizationPeriodSeconds;
    bool public fundDevAccounts;
    address public proxyAdminOwner;
    address public baseFeeVaultRecipient;
    uint256 public baseFeeVaultMinimumWithdrawalAmount;
    uint256 public baseFeeVaultWithdrawalNetwork;
    address public l1FeeVaultRecipient;
    uint256 public l1FeeVaultMinimumWithdrawalAmount;
    uint256 public l1FeeVaultWithdrawalNetwork;
    address public sequencerFeeVaultRecipient;
    uint256 public sequencerFeeVaultMinimumWithdrawalAmount;
    uint256 public sequencerFeeVaultWithdrawalNetwork;
    string public governanceTokenName;
    string public governanceTokenSymbol;
    address public governanceTokenOwner;
    uint256 public l2GenesisBlockGasLimit;
    uint32 public basefeeScalar;
    uint32 public blobbasefeeScalar;
    bool public enableGovernance;
    uint256 public eip1559Denominator;
    uint256 public eip1559Elasticity;
    uint256 public faultGameAbsolutePrestate;
    uint256 public faultGameGenesisBlock;
    bytes32 public faultGameGenesisOutputRoot;
    uint256 public faultGameMaxDepth;
    uint256 public faultGameSplitDepth;
    uint256 public faultGameClockExtension;
    uint256 public faultGameMaxClockDuration;
    uint256 public faultGameWithdrawalDelay;
    uint256 public preimageOracleMinProposalSize;
    uint256 public preimageOracleChallengePeriod;
    uint256 public systemConfigStartBlock;
    uint256 public requiredProtocolVersion;
    uint256 public recommendedProtocolVersion;
    uint256 public proofMaturityDelaySeconds;
    uint256 public disputeGameFinalityDelaySeconds;
    uint256 public respectedGameType;
    bool public useFaultProofs;
    bool public useAltDA;
    string public daCommitmentType;
    uint256 public daChallengeWindow;
    uint256 public daResolveWindow;
    uint256 public daBondSize;
    uint256 public daResolverRefundPercentage;
    bytes32 public electionFallbackList;

    bool public useCustomGasToken;
    address public customGasTokenAddress;

    bool public useInterop;

    // Made internal to create a view function that returns the proper types
    ElectionSystemConfig.SequencerRule[] internal _sequencerRules;
    ElectionTickets.GenesisAllocation[] internal _genesisAllocation;

    function read(string memory _path) public {
        console.log("DeployConfig: reading file %s", _path);
        try vm.readFile(_path) returns (string memory data) {
            _json = data;
        } catch {
            require(false, string.concat("Cannot find deploy config file at ", _path));
        }

        finalSystemOwner = stdJson.readAddress(_json, "$.finalSystemOwner");
        superchainConfigGuardian = stdJson.readAddress(_json, "$.superchainConfigGuardian");
        l1ChainID = stdJson.readUint(_json, "$.l1ChainID");
        l2ChainID = stdJson.readUint(_json, "$.l2ChainID");
        l2BlockTime = stdJson.readUint(_json, "$.l2BlockTime");

        l2GenesisDeltaTimeOffset = _readOr(_json, "$.l2GenesisDeltaTimeOffset", NULL_OFFSET);
        l2GenesisEcotoneTimeOffset = _readOr(_json, "$.l2GenesisEcotoneTimeOffset", NULL_OFFSET);
        l2GenesisFjordTimeOffset = _readOr(_json, "$.l2GenesisFjordTimeOffset", NULL_OFFSET);
        l2GenesisGraniteTimeOffset = _readOr(_json, "$.l2GenesisGraniteTimeOffset", NULL_OFFSET);

        maxSequencerDrift = stdJson.readUint(_json, "$.maxSequencerDrift");
        sequencerWindowSize = stdJson.readUint(_json, "$.sequencerWindowSize");
        channelTimeout = stdJson.readUint(_json, "$.channelTimeout");
        p2pSequencerAddress = stdJson.readAddress(_json, "$.p2pSequencerAddress");
        batchSenderAddress = stdJson.readAddress(_json, "$.batchSenderAddress");
        l2OutputOracleSubmissionInterval = stdJson.readUint(_json, "$.l2OutputOracleSubmissionInterval");
        _l2OutputOracleStartingTimestamp = stdJson.readInt(_json, "$.l2OutputOracleStartingTimestamp");
        l2OutputOracleStartingBlockNumber = stdJson.readUint(_json, "$.l2OutputOracleStartingBlockNumber");
        l2OutputOracleProposer = stdJson.readAddress(_json, "$.l2OutputOracleProposer");
        l2OutputOracleChallenger = stdJson.readAddress(_json, "$.l2OutputOracleChallenger");
        finalizationPeriodSeconds = stdJson.readUint(_json, "$.finalizationPeriodSeconds");
        fundDevAccounts = _readOr(_json, "$.fundDevAccounts", false);
        proxyAdminOwner = stdJson.readAddress(_json, "$.proxyAdminOwner");
        baseFeeVaultRecipient = stdJson.readAddress(_json, "$.baseFeeVaultRecipient");
        baseFeeVaultMinimumWithdrawalAmount = stdJson.readUint(_json, "$.baseFeeVaultMinimumWithdrawalAmount");
        baseFeeVaultWithdrawalNetwork = stdJson.readUint(_json, "$.baseFeeVaultWithdrawalNetwork");
        l1FeeVaultRecipient = stdJson.readAddress(_json, "$.l1FeeVaultRecipient");
        l1FeeVaultMinimumWithdrawalAmount = stdJson.readUint(_json, "$.l1FeeVaultMinimumWithdrawalAmount");
        l1FeeVaultWithdrawalNetwork = stdJson.readUint(_json, "$.l1FeeVaultWithdrawalNetwork");
        sequencerFeeVaultRecipient = stdJson.readAddress(_json, "$.sequencerFeeVaultRecipient");
        sequencerFeeVaultMinimumWithdrawalAmount = stdJson.readUint(_json, "$.sequencerFeeVaultMinimumWithdrawalAmount");
        sequencerFeeVaultWithdrawalNetwork = stdJson.readUint(_json, "$.sequencerFeeVaultWithdrawalNetwork");
        governanceTokenName = stdJson.readString(_json, "$.governanceTokenName");
        governanceTokenSymbol = stdJson.readString(_json, "$.governanceTokenSymbol");
        governanceTokenOwner = stdJson.readAddress(_json, "$.governanceTokenOwner");
        l2GenesisBlockGasLimit = stdJson.readUint(_json, "$.l2GenesisBlockGasLimit");
        basefeeScalar = uint32(_readOr(_json, "$.gasPriceOracleBaseFeeScalar", 1368));
        blobbasefeeScalar = uint32(_readOr(_json, "$.gasPriceOracleBlobBaseFeeScalar", 810949));

        enableGovernance = stdJson.readBool(_json, "$.enableGovernance");
        eip1559Denominator = stdJson.readUint(_json, "$.eip1559Denominator");
        eip1559Elasticity = stdJson.readUint(_json, "$.eip1559Elasticity");
        systemConfigStartBlock = stdJson.readUint(_json, "$.systemConfigStartBlock");
        requiredProtocolVersion = stdJson.readUint(_json, "$.requiredProtocolVersion");
        recommendedProtocolVersion = stdJson.readUint(_json, "$.recommendedProtocolVersion");

        useFaultProofs = _readOr(_json, "$.useFaultProofs", false);
        proofMaturityDelaySeconds = _readOr(_json, "$.proofMaturityDelaySeconds", 0);
        disputeGameFinalityDelaySeconds = _readOr(_json, "$.disputeGameFinalityDelaySeconds", 0);
        respectedGameType = _readOr(_json, "$.respectedGameType", 0);

        faultGameAbsolutePrestate = stdJson.readUint(_json, "$.faultGameAbsolutePrestate");
        faultGameMaxDepth = stdJson.readUint(_json, "$.faultGameMaxDepth");
        faultGameSplitDepth = stdJson.readUint(_json, "$.faultGameSplitDepth");
        faultGameClockExtension = stdJson.readUint(_json, "$.faultGameClockExtension");
        faultGameMaxClockDuration = stdJson.readUint(_json, "$.faultGameMaxClockDuration");
        faultGameGenesisBlock = stdJson.readUint(_json, "$.faultGameGenesisBlock");
        faultGameGenesisOutputRoot = stdJson.readBytes32(_json, "$.faultGameGenesisOutputRoot");
        faultGameWithdrawalDelay = stdJson.readUint(_json, "$.faultGameWithdrawalDelay");

        preimageOracleMinProposalSize = stdJson.readUint(_json, "$.preimageOracleMinProposalSize");
        preimageOracleChallengePeriod = stdJson.readUint(_json, "$.preimageOracleChallengePeriod");

        useAltDA = _readOr(_json, "$.useAltDA", false);
        daCommitmentType = _readOr(_json, "$.daCommitmentType", "KeccakCommitment");
        daChallengeWindow = _readOr(_json, "$.daChallengeWindow", 1000);
        daResolveWindow = _readOr(_json, "$.daResolveWindow", 1000);
        daBondSize = _readOr(_json, "$.daBondSize", 1000000000);
        daResolverRefundPercentage = _readOr(_json, "$.daResolverRefundPercentage", 0);

        useCustomGasToken = _readOr(_json, "$.useCustomGasToken", false);
        customGasTokenAddress = _readOr(_json, "$.customGasTokenAddress", address(0));

        useInterop = _readOr(_json, "$.useInterop", false);

        // Election system config
        electionFallbackList = stdJson.readBytes32(_json, "$.electionFallbackList");

        address[] memory targets = _readOr(_json, "$.genesisAllocation.targets", new address[](0));
        uint256[] memory amounts = _readOr(_json, "$.genesisAllocation.amounts", new uint256[](0));

        // Sanity check that the arrays are parallel
        if (targets.length != amounts.length) {
            revert("GenesisAllocation: targets and amounts must be the same length");
        }

        for (uint256 i; i < targets.length; i++) {
            _genesisAllocation.push(ElectionTickets.GenesisAllocation({ target: targets[i], amount: amounts[i] }));
        }

        // Read sequencer rules
        uint256[] memory assertionType = _readOr(_json, "$.sequencerRules.assertionType", new uint256[](0));
        bytes32[] memory desiredRetdata = _readOr(_json, "$.sequencerRules.desiredRetdata", new bytes32[](0));
        bytes[] memory configCalldata = _readOr(_json, "$.sequencerRules.configCalldata", new bytes[](0));
        address[] memory target = _readOr(_json, "$.sequencerRules.targets", new address[](0));

        if (
            assertionType.length != desiredRetdata.length || desiredRetdata.length != configCalldata.length
                || configCalldata.length != target.length
        ) {
            revert("SequencerRules: assertionType, desiredRetdata, configCalldata, and target must be the same length");
        }

        // TODO(spire): Support addressOffsets, should be added in the enhance json config pr
        for (uint256 i; i < assertionType.length; i++) {
            _sequencerRules.push(
                ElectionSystemConfig.SequencerRule({
                    assertionType: ElectionSystemConfig.SequencerAssertion(assertionType[i]),
                    desiredRetdata: desiredRetdata[i],
                    configCalldata: configCalldata[i],
                    target: target[i],
                    addressOffsets: new uint256[](0)
                })
            );
        }
    }

    function fork() public view returns (Fork fork_) {
        // let env var take precedence
        fork_ = Config.fork();
        if (fork_ == Fork.NONE) {
            // Will revert if no deploy config can be found either.
            fork_ = latestGenesisFork();
            console.log("DeployConfig: using deploy config fork: %s", fork_.toString());
        } else {
            console.log("DeployConfig: using env var fork: %s", fork_.toString());
        }
    }

    function l1StartingBlockTag() public returns (bytes32) {
        try vm.parseJsonBytes32(_json, "$.l1StartingBlockTag") returns (bytes32 tag) {
            return tag;
        } catch {
            try vm.parseJsonString(_json, "$.l1StartingBlockTag") returns (string memory tag) {
                return _getBlockByTag(tag);
            } catch {
                try vm.parseJsonUint(_json, "$.l1StartingBlockTag") returns (uint256 tag) {
                    return _getBlockByTag(vm.toString(tag));
                } catch { }
            }
        }
        revert("l1StartingBlockTag must be a bytes32, string or uint256 or cannot fetch l1StartingBlockTag");
    }

    function l2OutputOracleStartingTimestamp() public returns (uint256) {
        if (_l2OutputOracleStartingTimestamp < 0) {
            bytes32 tag = l1StartingBlockTag();
            string[] memory cmd = new string[](3);
            cmd[0] = Executables.bash;
            cmd[1] = "-c";
            cmd[2] = string.concat("cast block ", vm.toString(tag), " --json | ", Executables.jq, " .timestamp");
            bytes memory res = Process.run(cmd);
            return stdJson.readUint(string(res), "");
        }
        return uint256(_l2OutputOracleStartingTimestamp);
    }

    /// @notice Allow the `useAltDA` config to be overridden in testing environments
    function setUseAltDA(bool _useAltDA) public {
        useAltDA = _useAltDA;
    }

    /// @notice Allow the `useFaultProofs` config to be overridden in testing environments
    function setUseFaultProofs(bool _useFaultProofs) public {
        useFaultProofs = _useFaultProofs;
    }

    /// @notice Allow the `useInterop` config to be overridden in testing environments
    function setUseInterop(bool _useInterop) public {
        useInterop = _useInterop;
    }

    /// @notice Allow the `fundDevAccounts` config to be overridden.
    function setFundDevAccounts(bool _fundDevAccounts) public {
        fundDevAccounts = _fundDevAccounts;
    }

    /// @notice Allow the `useCustomGasToken` config to be overridden in testing environments
    function setUseCustomGasToken(address _token) public {
        useCustomGasToken = true;
        customGasTokenAddress = _token;
    }

    function genesisAllocation() public view returns (ElectionTickets.GenesisAllocation[] memory) {
        return _genesisAllocation;
    }

    function sequencerRules() public view returns (ElectionSystemConfig.SequencerRule[] memory) {
        return _sequencerRules;
    }

    function latestGenesisFork() internal view returns (Fork) {
        if (l2GenesisGraniteTimeOffset == 0) {
            return Fork.GRANITE;
        } else if (l2GenesisFjordTimeOffset == 0) {
            return Fork.FJORD;
        } else if (l2GenesisEcotoneTimeOffset == 0) {
            return Fork.ECOTONE;
        } else if (l2GenesisDeltaTimeOffset == 0) {
            return Fork.DELTA;
        }
        revert("DeployConfig: no supported fork active at genesis");
    }

    function _getBlockByTag(string memory _tag) internal returns (bytes32) {
        string[] memory cmd = new string[](3);
        cmd[0] = Executables.bash;
        cmd[1] = "-c";
        cmd[2] = string.concat("cast block ", _tag, " --json | ", Executables.jq, " -r .hash");
        bytes memory res = Process.run(cmd);
        return abi.decode(res, (bytes32));
    }

    function _readOr(string memory json, string memory key, bool defaultValue) internal view returns (bool) {
        return vm.keyExistsJson(json, key) ? json.readBool(key) : defaultValue;
    }

    function _readOr(string memory json, string memory key, uint256 defaultValue) internal view returns (uint256) {
        return (vm.keyExistsJson(json, key) && !_isNull(json, key)) ? json.readUint(key) : defaultValue;
    }

    function _readOr(string memory json, string memory key, address defaultValue) internal view returns (address) {
        return vm.keyExistsJson(json, key) ? json.readAddress(key) : defaultValue;
    }

    function _readOr(
        string memory json,
        string memory key,
        address[] memory defaultValue
    )
        internal
        view
        returns (address[] memory)
    {
        return vm.keyExistsJson(json, key) ? json.readAddressArray(key) : defaultValue;
    }

    function _readOr(
        string memory json,
        string memory key,
        uint256[] memory defaultValue
    )
        internal
        view
        returns (uint256[] memory)
    {
        return vm.keyExistsJson(json, key) ? json.readUintArray(key) : defaultValue;
    }

    function _isNull(string memory json, string memory key) internal pure returns (bool) {
        string memory value = json.readString(key);
        return (keccak256(bytes(value)) == keccak256(bytes("null")));
    }

    function _readOr(
        string memory json,
        string memory key,
        string memory defaultValue
    )
        internal
        view
        returns (string memory)
    {
        return vm.keyExists(json, key) ? json.readString(key) : defaultValue;
    }

    function _readOr(
        string memory json,
        string memory key,
        bytes32[] memory defaultValue
    )
        internal
        view
        returns (bytes32[] memory)
    {
        return vm.keyExists(json, key) ? json.readBytes32Array(key) : defaultValue;
    }

    function _readOr(
        string memory json,
        string memory key,
        bytes[] memory defaultValue
    )
        internal
        view
        returns (bytes[] memory)
    {
        return vm.keyExists(json, key) ? json.readBytesArray(key) : defaultValue;
    }
}
