package derivation

import (
	"fmt"
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

func TestElectionType(gt *testing.T) {
	t := actionsHelpers.NewDefaultTesting(gt)
	p := &e2eutils.TestParams{
		MaxSequencerDrift:   20, // larger than L1 block time we simulate in this test (12)
		SequencerWindowSize: 24,
		ChannelTimeout:      20,
		L1BlockTime:         12,
	}
	dp := e2eutils.MakeDeployParams(t, p) //actionsHelpers.DefaultRollupTestParams)
	sd := e2eutils.Setup(t, dp, actionsHelpers.DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)
	miner, l2Engine, sequencer, verifier, verifierEng, batcher := actionsHelpers.SetupElectionTest(t, dp, sd, log)

	fmt.Println(verifierEng, l2Engine)

	// TODO (wintercode) create a batch tx from the election winner. In testing this is hardcoded to be
	// 0x7c60541eB6f54f0F3c8B34D0a00De9045d2f5534
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
	// TODO: hardcode alice's address as election winner for testing
	require.NoError(t, cl.SendTransaction(t.Ctx(), tx))

	// Sequence the tx
	sequencer.ActL2PipelineFull(t)
	verifier.ActL2PipelineFull(t)

	// Make L2 block
	sequencer.ActL2StartBlock(t)
	l2Engine.ActL2IncludeTx(dp.Addresses.Alice)(t)
	sequencer.ActL2EndBlock(t)

	// batch submit to L1
	batcher.ActL2BatchBuffer(t)
	batcher.ActL2ChannelClose(t)
	batcher.ActL2BatchSubmit(t)

	// confirm batch on L1
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTx(dp.Addresses.Batcher)(t)
	miner.ActL1EndBlock(t)
	bl := miner.L1Chain().CurrentBlock()
	log.Info("bl", "txs", len(miner.L1Chain().GetBlockByHash(bl.Hash()).Transactions()))

	// Now make enough L1 blocks that the verifier will have to derive a L2 block
	// It will also eagerly derive the block from the batcher
	// for i := uint64(0); i < sd.RollupCfg.SeqWindowSize; i++ {
	// 	miner.ActL1StartBlock(12)(t)
	// 	miner.ActL1EndBlock(t)
	// }

	// sync verifier from L1 batch in otherwise empty sequence window
	verifier.ActL1HeadSignal(t)
	verifier.ActL2PipelineFull(t)

	// check that the tx from alice made it into the L2 chain, meaning the batch included by the election winner was processed
	verifCl := verifierEng.EthClient()
	vTx, isPending, err := verifCl.TransactionByHash(t.Ctx(), tx.Hash())
	require.NoError(t, err)
	require.False(t, isPending)
	require.NotNil(t, vTx)
}
