package derivation

import (
	"math/big"
	"testing"

	actionsHelpers "github.com/ethereum-optimism/optimism/op-e2e/actions/helpers"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

// Tests the simplest election winner scenario.
// Batcher is the election winner. It receives an L2 block from Alice, and includes it in the L1 block.
// That block is then picked up by the verifier, and successfully processed in the derivation pipeline, with the tx from Alice included.
// The second batcher is not an election winner, so the batch it creates and sends
// is not included in the L2 chain since it does not pass through the derivation pipeline.
func TestElectionType(gt *testing.T) {
	t := actionsHelpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, actionsHelpers.DefaultRollupTestParams)
	sd := e2eutils.Setup(t, dp, actionsHelpers.DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)
	miner, l2Engine, sequencer, verifier, _, batcher := actionsHelpers.SetupElectionTest(t, dp, sd, log)

	cl := l2Engine.EthClient()
	n, err := cl.PendingNonceAt(t.Ctx(), dp.Addresses.Alice)
	require.NoError(t, err)
	signer := types.LatestSigner(sd.L2Cfg.Config)
	tx := types.MustSignNewTx(dp.Secrets.Alice, signer, &types.DynamicFeeTx{
		ChainID:   sd.L2Cfg.Config.ChainID,
		Nonce:     n,
		GasTipCap: big.NewInt(2 * params.GWei),
		GasFeeCap: new(big.Int).Add(miner.L1Chain().CurrentBlock().BaseFee, big.NewInt(2*params.GWei)),
		Gas:       params.TxGas,
		To:        &dp.Addresses.Bob,
		Value:     e2eutils.Ether(2),
	})
	require.NoError(t, cl.SendTransaction(t.Ctx(), tx))

	// Sequence the tx
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	// build empty L1 block
	miner.ActEmptyBlock(t)
	// finalize it, so the L1 geth blob pool doesn't log errors about missing finality
	miner.ActL1SafeNext(t)
	miner.ActL1FinalizeNext(t)
	// Create L2 blocks, and reference the L1 head as origin
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1Head(t)
	// submit all new L2 blocks
	batcher.ActSubmitAll(t)
	batchTx := batcher.LastSubmitted
	require.Equal(t, uint8(types.BlobTxType), batchTx.Type(), "batch tx must be blob-tx")
	// new L1 block with L2 batch
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTxByHash(batchTx.Hash())(t)
	miner.ActL1EndBlock(t)

	// verifier picks up the L2 chain that was submitted
	verifier.ActL1HeadSignal(t)
	verifier.ActL2PipelineFull(t)

	require.Equal(t, verifier.L2Safe(), sequencer.L2Unsafe(), "verifier syncs from sequencer via L1")
	require.NotEqual(t, sequencer.L2Safe(), sequencer.L2Unsafe(), "sequencer has not processed L1 yet")

	// check that the tx from alice made it into the L2 chain
	vTx, _, err := cl.TransactionByHash(t.Ctx(), tx.Hash())
	require.NoError(t, err)
	require.NotNil(t, vTx)

	// at this point, the batch from the election winner was successfully processed.
	// This next batch is sent from a batcher who is not the election winner, so their batch is not included in the L2 chain.
	// next step, send a batch from someone other than the election winner, and confirm that
	// its l2 transactions are not included in the L2 chain
	miner.ActL1SafeNext(t)
	miner.ActL1FinalizeNext(t)
	newBatcher := setupBatcher(t, log, sd, dp, miner, sequencer, l2Engine)

	n2, err := cl.PendingNonceAt(t.Ctx(), dp.Addresses.Mallory)
	require.NoError(t, err)
	signer2 := types.LatestSigner(sd.L2Cfg.Config)

	tx2 := types.MustSignNewTx(dp.Secrets.Mallory, signer2, &types.DynamicFeeTx{
		ChainID:   sd.L2Cfg.Config.ChainID,
		Nonce:     n2,
		GasTipCap: big.NewInt(2 * params.GWei),
		GasFeeCap: new(big.Int).Add(miner.L1Chain().CurrentBlock().BaseFee, big.NewInt(2*params.GWei)),
		Gas:       params.TxGas,
		To:        &dp.Addresses.Bob,
		Value:     e2eutils.Ether(2),
	})
	require.NoError(t, cl.SendTransaction(t.Ctx(), tx2))

	// build empty L1 block
	miner.ActEmptyBlock(t)

	// Create L2 blocks, and reference the L1 head as origin
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1Head(t)

	// submit all new L2 blocks
	newBatcher.ActSubmitAll(t)
	batchTx = newBatcher.LastSubmitted

	// new L1 block with L2 batch
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTxByHash(batchTx.Hash())(t)
	miner.ActL1EndBlock(t)

	require.Equal(t, uint8(types.BlobTxType), batchTx.Type(), "batch tx must be blob-tx")

	verifier.ActL1HeadSignal(t)
	verifier.ActL2PipelineFull(t)

	// verifier does not pick up the L2 chain that was submitted
	require.NotEqual(t, verifier.L2Safe(), sequencer.L2Unsafe(), "verifier does not sync from sequencer via L1")
}

func setupBatcher(t actionsHelpers.Testing, log log.Logger, sd *e2eutils.SetupData, dp *e2eutils.DeployParams, miner *actionsHelpers.L1Miner,
	sequencer *actionsHelpers.L2Sequencer, engine *actionsHelpers.L2Engine) *actionsHelpers.L2Batcher {
	return actionsHelpers.NewL2Batcher(log, sd.RollupCfg, &actionsHelpers.BatcherCfg{
		MinL1TxSize:          0,
		MaxL1TxSize:          128_000,
		BatcherKey:           dp.Secrets.Bob,
		DataAvailabilityType: "blobs",
		L1BlockTime:          dp.DeployConfig.L1BlockTime,
	}, sequencer.RollupClient(), miner.EthClient(), engine.EthClient(), engine.EngineClient(t, sd.RollupCfg))
}
