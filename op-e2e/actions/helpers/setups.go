package helpers

import (
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-e2e/actions/upgrades/helpers"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-node/rollup/election"
	"github.com/ethereum-optimism/optimism/op-node/rollup/election_store"
	"github.com/ethereum-optimism/optimism/op-node/rollup/sync"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"
)

func SetupElectionTest(t Testing, dp *e2eutils.DeployParams, sd *e2eutils.SetupData, log log.Logger, opts ...SequencerOpt) (*L1Miner, *L2Engine, *L2Sequencer, *L2Verifier, *L2Engine, *L2Batcher) {
	miner, seqEngine, sequencer := SetupSequencerTest(t, sd, dp, log)
	l1Cl := miner.L1Client(t, sd.RollupCfg)
	verifEngine, verifier := SetupVerifier(t, sd, log, l1Cl, l1Cl.EthClient, miner.BlobStore(), miner.BeaconClient(), &sync.Config{})
	rollupSeqCl := sequencer.RollupClient()
	batcher := NewL2Batcher(log, sd.RollupCfg, BlobBatcherCfg(dp),
		rollupSeqCl, miner.EthClient(), seqEngine.EthClient(), seqEngine.EngineClient(t, sd.RollupCfg))
	return miner, seqEngine, sequencer, verifier, verifEngine, batcher
}

func SetupSequencerTest(t Testing, sd *e2eutils.SetupData, dp *e2eutils.DeployParams, log log.Logger, opts ...SequencerOpt) (*L1Miner, *L2Engine, *L2Sequencer) {
	jwtPath := e2eutils.WriteDefaultJWT(t)
	cfg := DefaultSequencerConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	miner := NewL1Miner(t, log.New("role", "l1-miner"), sd.L1Cfg, dp.DeployConfig)

	l1F, err := sources.NewL1Client(miner.RPCClient(), log, nil, sources.L1ClientDefaultConfig(sd.RollupCfg, false, sources.RPCKindStandard))
	require.NoError(t, err)
	engine := NewL2Engine(t, log.New("role", "sequencer-engine"), sd.L2Cfg, sd.RollupCfg.Genesis.L1, jwtPath, EngineWithP2P())
	l2Cl, err := sources.NewEngineClient(engine.RPCClient(), log, nil, sources.EngineClientDefaultConfig(sd.RollupCfg))
	require.NoError(t, err)

	sequencer := NewL2Sequencer(t, log.New("role", "sequencer"), l1F, miner.BlobStore(), miner.BeaconClient(), altda.Disabled, l2Cl, l1F.EthClient, sd.RollupCfg, 0, cfg.InteropBackend)
	return miner, engine, sequencer
}

func SetupVerifier(t Testing, sd *e2eutils.SetupData, log log.Logger,
	l1F derive.L1Fetcher, l1Client L1API, blobSrc derive.L1BlobsFetcher, beaconClient election.BeaconClient, syncCfg *sync.Config, opts ...VerifierOpt) (*L2Engine, *L2Verifier) {
	cfg := DefaultVerifierCfg()
	for _, opt := range opts {
		opt(cfg)
	}
	jwtPath := e2eutils.WriteDefaultJWT(t)
	engine := NewL2Engine(t, log.New("role", "verifier-engine"), sd.L2Cfg, sd.RollupCfg.Genesis.L1, jwtPath, EngineWithP2P())
	engCl := engine.EngineClient(t, sd.RollupCfg)
	electionStore := election_store.NewElectionStore(log)
	verifier := NewL2Verifier(t, log.New("role", "verifier"), l1F, blobSrc, beaconClient, altda.Disabled, engCl, l1Client, sd.RollupCfg, syncCfg, cfg.SafeHeadListener, cfg.InteropBackend, electionStore)
	return engine, verifier
}

func SetupVerifierOnlyTest(t Testing, sd *e2eutils.SetupData, dp *e2eutils.DeployParams, log log.Logger) (*L1Miner, *L2Engine, *L2Verifier) {
	miner := NewL1Miner(t, log, sd.L1Cfg, dp.DeployConfig)
	l1Cl := miner.L1Client(t, sd.RollupCfg)
	engine, verifier := SetupVerifier(t, sd, log, l1Cl, l1Cl.EthClient, miner.BlobStore(), miner.BeaconClient(), &sync.Config{})
	return miner, engine, verifier
}

func SetupReorgTest(t Testing, config *e2eutils.TestParams, deltaTimeOffset *hexutil.Uint64) (*e2eutils.SetupData, *e2eutils.DeployParams, *L1Miner, *L2Sequencer, *L2Engine, *L2Verifier, *L2Engine, *L2Batcher) {
	dp := e2eutils.MakeDeployParams(t, config)
	helpers.ApplyDeltaTimeOffset(dp, deltaTimeOffset)

	sd := e2eutils.Setup(t, dp, DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)

	return SetupReorgTestActors(t, dp, sd, log)
}

func SetupReorgTestActors(t Testing, dp *e2eutils.DeployParams, sd *e2eutils.SetupData, log log.Logger) (*e2eutils.SetupData, *e2eutils.DeployParams, *L1Miner, *L2Sequencer, *L2Engine, *L2Verifier, *L2Engine, *L2Batcher) {
	miner, seqEngine, sequencer := SetupSequencerTest(t, sd, dp, log)
	miner.ActL1SetFeeRecipient(common.Address{'A'})
	sequencer.ActL2PipelineFull(t)
	client := miner.L1Client(t, sd.RollupCfg)
	verifEngine, verifier := SetupVerifier(t, sd, log, client, client.EthClient, miner.BlobStore(), miner.BeaconClient(), &sync.Config{})
	rollupSeqCl := sequencer.RollupClient()
	batcher := NewL2Batcher(log, sd.RollupCfg, DefaultBatcherCfg(dp),
		rollupSeqCl, miner.EthClient(), seqEngine.EthClient(), seqEngine.EngineClient(t, sd.RollupCfg))
	return sd, dp, miner, sequencer, seqEngine, verifier, verifEngine, batcher
}

func SetupReorgTestBlob(t Testing, config *e2eutils.TestParams, deltaTimeOffset *hexutil.Uint64) (*e2eutils.SetupData, *e2eutils.DeployParams, *L1Miner, *L2Sequencer, *L2Engine, *L2Verifier, *L2Engine, *L2Batcher) {
	dp := e2eutils.MakeDeployParams(t, config)
	helpers.ApplyDeltaTimeOffset(dp, deltaTimeOffset)

	sd := e2eutils.Setup(t, dp, DefaultAlloc)
	log := testlog.Logger(t, log.LevelDebug)

	return SetupReorgTestActorsBlob(t, dp, sd, log)
}

func SetupReorgTestActorsBlob(t Testing, dp *e2eutils.DeployParams, sd *e2eutils.SetupData, log log.Logger) (*e2eutils.SetupData, *e2eutils.DeployParams, *L1Miner, *L2Sequencer, *L2Engine, *L2Verifier, *L2Engine, *L2Batcher) {
	miner, seqEngine, sequencer := SetupSequencerTest(t, sd, dp, log)
	miner.ActL1SetFeeRecipient(common.Address{'A'})
	sequencer.ActL2PipelineFull(t)
	client := miner.L1Client(t, sd.RollupCfg)
	verifEngine, verifier := SetupVerifier(t, sd, log, client, client.EthClient, miner.BlobStore(), miner.BeaconClient(), &sync.Config{})
	rollupSeqCl := sequencer.RollupClient()

	batcher := NewL2Batcher(log, sd.RollupCfg, BlobBatcherCfg(dp),
		rollupSeqCl, miner.EthClient(), seqEngine.EthClient(), seqEngine.EngineClient(t, sd.RollupCfg))
	return sd, dp, miner, sequencer, seqEngine, verifier, verifEngine, batcher
}
