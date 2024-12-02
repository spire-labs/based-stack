package derivation

import (
	"fmt"
	"testing"

	actionsHelpers "github.com/ethereum-optimism/optimism/op-e2e/actions/helpers"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func TestElectionType(gt *testing.T) {
	t := actionsHelpers.NewDefaultTesting(gt)
	dp := e2eutils.MakeDeployParams(t, actionsHelpers.DefaultRollupTestParams)
	sd := e2eutils.Setup(t, dp, actionsHelpers.DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)
	miner, l2Engine, sequencer, verifier, verifierEng, batcher := actionsHelpers.SetupElectionTest(t, dp, sd, log)

	fmt.Println(verifierEng, l2Engine)

	sequencer.ActL2PipelineFull(t)

	// build empty L1 block
	miner.ActEmptyBlock(t)

	// Create L2 blocks, and reference the L1 head as origin
	sequencer.ActL1HeadSignal(t)
	sequencer.ActBuildToL1Head(t)

	// submit all new L2 blocks
	batcher.ActSubmitAll(t)

	// new L1 block with L2 batch
	miner.ActL1StartBlock(12)(t)
	miner.ActL1IncludeTx(sd.RollupCfg.Genesis.SystemConfig.BatcherAddr)(t)
	/*batchTx*/ _ = miner.L1Transactions[0]
	miner.ActL1EndBlock(t)

	// verifier picks up the L2 chain that was submitted
	verifier.ActL1HeadSignal(t)
	verifier.ActL2PipelineFull(t)
	require.NotEqual(t, sequencer.L2Safe(), sequencer.L2Unsafe(), "sequencer has not processed L1 yet")

	// l1unsafe block event -> time to process new L1 block

	// call out to beacon client to get epoch number
	// todo: mock beacon to return lookahead addresses

	// calls to get winners at epoch
	// calls to get batch ticket accounting from contract

	// test that ElectionWinnerEvent is emitted
	// test that NextElectionWinnerEvent is emitted

}
