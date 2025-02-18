package derivation

import (
	"math/big"
	"math/rand"
	"testing"

	actionsHelpers "github.com/ethereum-optimism/optimism/op-e2e/actions/helpers"
	upgradesHelpers "github.com/ethereum-optimism/optimism/op-e2e/actions/upgrades/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-node/rollup/sync"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
)

// TestSystemConfigBatchType run each system config-related test case in singular batch mode and span batch mode.
func TestSystemConfigBatchType(t *testing.T) {
	t.Skip("TODO(spire): Reenable these tests")
	tests := []struct {
		name string
		f    func(gt *testing.T, deltaTimeOffset *hexutil.Uint64)
	}{
		{"GPOParamsChange", GPOParamsChange},
		{"GasLimitChange", GasLimitChange},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name+"_SingularBatch", func(t *testing.T) {
			test.f(t, nil)
		})
	}

	deltaTimeOffset := hexutil.Uint64(0)
	for _, test := range tests {
		test := test
		t.Run(test.name+"_SpanBatch", func(t *testing.T) {
			test.f(t, &deltaTimeOffset)
		})
	}
}

// GPOParamsChange tests that the GPO params can be updated to adjust fees of L2 transactions,
// and that the L1 data fees to the L2 transaction are applied correctly before, during and after the GPO update in L2.
func GPOParamsChange(gt *testing.T, deltaTimeOffset *hexutil.Uint64) {
	t := actionsHelpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, actionsHelpers.DefaultRollupTestParams)
	upgradesHelpers.ApplyDeltaTimeOffset(dp, deltaTimeOffset)

	// activating Delta only, not Ecotone and further:
	// the GPO change assertions here all apply only for the Delta transition.
	// Separate tests cover Ecotone GPO changes.
	dp.DeployConfig.L2GenesisEcotoneTimeOffset = nil
	dp.DeployConfig.L2GenesisFjordTimeOffset = nil
	dp.DeployConfig.L2GenesisGraniteTimeOffset = nil

	sd := e2eutils.Setup(t, dp, actionsHelpers.DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)
	miner, seqEngine, sequencer := actionsHelpers.SetupSequencerTest(t, sd, dp, log)
	batcher := actionsHelpers.NewL2Batcher(log, sd.RollupCfg, actionsHelpers.DefaultBatcherCfg(dp),
		sequencer.RollupClient(), miner.EthClient(), seqEngine.EthClient(), seqEngine.EngineClient(t, sd.RollupCfg))

	alice := actionsHelpers.NewBasicUser[any](log, dp.Secrets.Alice, rand.New(rand.NewSource(1234)))
	alice.SetUserEnv(&actionsHelpers.BasicUserEnv[any]{
		EthCl:  seqEngine.EthClient(),
		Signer: types.LatestSigner(sd.L2Cfg.Config),
	})

	sequencer.ActL2PipelineFull(t)

	// new L1 block, with new L2 chain
	miner.ActEmptyBlock(t)
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1Head(t)
	basefee := miner.L1Chain().CurrentBlock().BaseFee

	// alice makes a L2 tx, sequencer includes it
	alice.ActResetTxOpts(t)
	alice.ActMakeTx(t)
	sequencer.ActL2StartBlock(t)
	seqEngine.ActL2IncludeTx(dp.Addresses.Alice)(t)
	sequencer.ActL2EndBlock(t)

	receipt := alice.LastTxReceipt(t)
	require.Equal(t, basefee, receipt.L1GasPrice, "L1 gas price matches basefee of L1 origin")
	require.NotZero(t, receipt.L1GasUsed, "L2 tx uses L1 data")
	require.Equal(t,
		new(big.Float).Mul(
			new(big.Float).SetInt(basefee),
			new(big.Float).Mul(new(big.Float).SetInt(receipt.L1GasUsed), receipt.FeeScalar),
		),
		new(big.Float).SetInt(receipt.L1Fee), "fee field in receipt matches gas used times scalar times basefee")
	// receipt.L1GasUsed includes the overhead already, so subtract that before passing it into the L1 cost func
	l1Cost := types.L1Cost(receipt.L1GasUsed.Uint64()-2100, basefee, big.NewInt(2100), big.NewInt(1000_000))
	require.Equal(t, l1Cost, receipt.L1Fee, "L1 fee is computed with standard GPO params")
	require.Equal(t, "1", receipt.FeeScalar.String(), "1000_000 divided by 6 decimals = float(1)")

	// confirm L2 chain on L1
	batcher.ActSubmitAll(t)
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTx(dp.Addresses.Batcher)(t)
	miner.ActL1EndBlock(t)

	sysCfgContract, err := bindings.NewSystemConfig(sd.RollupCfg.L1SystemConfigAddress, miner.EthClient())
	require.NoError(t, err)

	sysCfgOwner, err := bind.NewKeyedTransactorWithChainID(dp.Secrets.SysCfgOwner, sd.RollupCfg.L1ChainID)
	require.NoError(t, err)

	// overhead changes from 2100 (default) to 1000
	// scalar changes from 1_000_000 (default) to 2_300_000
	// e.g. if system operator determines that l2 txs need to be more expensive, but small ones less
	_, err = sysCfgContract.SetGasConfig(sysCfgOwner, big.NewInt(1000), big.NewInt(2_300_000))
	require.NoError(t, err)

	// include the GPO change tx in L1
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTx(dp.Addresses.SysCfgOwner)(t)
	miner.ActL1EndBlock(t)
	basefeeGPOUpdate := miner.L1Chain().CurrentBlock().BaseFee

	// build empty L2 chain, up to but excluding the L2 block with the L1 origin that processes the GPO change
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1HeadExcl(t)

	engCl := seqEngine.EngineClient(t, sd.RollupCfg)
	envelope, err := engCl.PayloadByLabel(t.Ctx(), eth.Unsafe)
	require.NoError(t, err)
	sysCfg, err := derive.PayloadToSystemConfig(sd.RollupCfg, envelope.ExecutionPayload)
	require.NoError(t, err)
	require.Equal(t, sd.RollupCfg.Genesis.SystemConfig, sysCfg, "still have genesis system config before we adopt the L1 block with GPO change")

	// Now alice makes another transaction, which gets included in the same block that adopts the L1 origin with GPO change
	alice.ActResetTxOpts(t)
	alice.ActMakeTx(t)
	sequencer.ActL2StartBlock(t)
	seqEngine.ActL2IncludeTx(dp.Addresses.Alice)(t)
	sequencer.ActL2EndBlock(t)

	envelope, err = engCl.PayloadByLabel(t.Ctx(), eth.Unsafe)
	require.NoError(t, err)
	sysCfg, err = derive.PayloadToSystemConfig(sd.RollupCfg, envelope.ExecutionPayload)
	require.NoError(t, err)
	require.Equal(t, eth.Bytes32(common.BigToHash(big.NewInt(1000))), sysCfg.Overhead, "overhead changed")
	require.Equal(t, eth.Bytes32(common.BigToHash(big.NewInt(2_300_000))), sysCfg.Scalar, "scalar changed")

	receipt = alice.LastTxReceipt(t)
	require.Equal(t, basefeeGPOUpdate, receipt.L1GasPrice, "L1 gas price matches basefee of L1 origin")
	require.NotZero(t, receipt.L1GasUsed, "L2 tx uses L1 data")
	// subtract overhead from L1GasUsed receipt field, types.L1Cost applies it again
	l1Cost = types.L1Cost(receipt.L1GasUsed.Uint64()-1000, basefeeGPOUpdate, big.NewInt(1000), big.NewInt(2_300_000))
	require.Equal(t, l1Cost, receipt.L1Fee, "L1 fee is computed with updated GPO params")
	require.Equal(t, "2.3", receipt.FeeScalar.String(), "2_300_000 divided by 6 decimals = float(2.3)")

	// build more L2 blocks, with new L1 origin
	miner.ActEmptyBlock(t)
	basefee = miner.L1Chain().CurrentBlock().BaseFee
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1Head(t)
	// and Alice makes a tx again
	alice.ActResetTxOpts(t)
	alice.ActMakeTx(t)
	sequencer.ActL2StartBlock(t)
	seqEngine.ActL2IncludeTx(dp.Addresses.Alice)(t)
	sequencer.ActL2EndBlock(t)

	// and verify the new GPO params are persistent, even though the L1 origin and L2 chain have progressed
	receipt = alice.LastTxReceipt(t)
	require.Equal(t, basefee, receipt.L1GasPrice, "L1 gas price matches basefee of L1 origin")
	require.NotZero(t, receipt.L1GasUsed, "L2 tx uses L1 data")
	// subtract overhead from L1GasUsed receipt field, types.L1Cost applies it again
	l1Cost = types.L1Cost(receipt.L1GasUsed.Uint64()-1000, basefee, big.NewInt(1000), big.NewInt(2_300_000))
	require.Equal(t, l1Cost, receipt.L1Fee, "L1 fee is computed with updated GPO params")
	require.Equal(t, "2.3", receipt.FeeScalar.String(), "2_300_000 divided by 6 decimals = float(2.3)")
}

// GasLimitChange tests that the gas limit can be configured to L1,
// and that the L2 changes the gas limit instantly at the exact block that adopts the L1 origin with
// the gas limit change event. And checks if a verifier node can reproduce the same gas limit change.
func GasLimitChange(gt *testing.T, deltaTimeOffset *hexutil.Uint64) {
	t := actionsHelpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, actionsHelpers.DefaultRollupTestParams)
	upgradesHelpers.ApplyDeltaTimeOffset(dp, deltaTimeOffset)
	sd := e2eutils.Setup(t, dp, actionsHelpers.DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)
	miner, seqEngine, sequencer := actionsHelpers.SetupSequencerTest(t, sd, dp, log)
	batcher := actionsHelpers.NewL2Batcher(log, sd.RollupCfg, actionsHelpers.DefaultBatcherCfg(dp),
		sequencer.RollupClient(), miner.EthClient(), seqEngine.EthClient(), seqEngine.EngineClient(t, sd.RollupCfg))

	sequencer.ActL2PipelineFull(t)
	miner.ActEmptyBlock(t)
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1Head(t)

	oldGasLimit := seqEngine.L2Chain().CurrentBlock().GasLimit
	require.Equal(t, oldGasLimit, uint64(dp.DeployConfig.L2GenesisBlockGasLimit))

	// change gas limit on L1 to triple what it was
	sysCfgContract, err := bindings.NewSystemConfig(sd.RollupCfg.L1SystemConfigAddress, miner.EthClient())
	require.NoError(t, err)

	sysCfgOwner, err := bind.NewKeyedTransactorWithChainID(dp.Secrets.SysCfgOwner, sd.RollupCfg.L1ChainID)
	require.NoError(t, err)

	_, err = sysCfgContract.SetGasLimit(sysCfgOwner, oldGasLimit*3)
	require.NoError(t, err)

	// include the gaslimit update on L1
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTx(dp.Addresses.SysCfgOwner)(t)
	miner.ActL1EndBlock(t)

	// build to latest L1, excluding the block that adopts the L1 block with the gaslimit change
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1HeadExcl(t)

	require.Equal(t, oldGasLimit, seqEngine.L2Chain().CurrentBlock().GasLimit)
	require.Equal(t, uint64(1), sequencer.SyncStatus().UnsafeL2.L1Origin.Number)

	// now include the L1 block with the gaslimit change, and see if it changes as expected
	sequencer.ActBuildToL1Head(t)
	require.Equal(t, oldGasLimit*3, seqEngine.L2Chain().CurrentBlock().GasLimit)
	require.Equal(t, uint64(2), sequencer.SyncStatus().UnsafeL2.L1Origin.Number)

	// now submit all this to L1, and see if a verifier can sync and reproduce it
	batcher.ActSubmitAll(t)
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTx(dp.Addresses.Batcher)(t)
	miner.ActL1EndBlock(t)
	l1Cl := miner.L1Client(t, sd.RollupCfg)

	_, verifier := actionsHelpers.SetupVerifier(t, sd, log, l1Cl, l1Cl.EthClient, miner.BlobStore(), miner.BeaconClient(), &sync.Config{})
	verifier.ActL2PipelineFull(t)

	require.Equal(t, sequencer.L2Unsafe(), verifier.L2Safe(), "verifier stays in sync, even with gaslimit changes")
}
