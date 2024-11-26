package derive

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/packages/contracts-bedrock/snapshots"
)

type DataIter interface {
	Next(ctx context.Context) (eth.Data, error)
}

type L1TransactionFetcher interface {
	InfoAndTxsByHash(ctx context.Context, hash common.Hash) (eth.BlockInfo, types.Transactions, error)
	FetchReceipts(ctx context.Context, hast common.Hash) (eth.BlockInfo, types.Receipts, error)
}

type L1BlobsFetcher interface {
	// GetBlobs fetches blobs that were confirmed in the given L1 block with the given indexed hashes.
	GetBlobs(ctx context.Context, ref eth.L1BlockRef, hashes []eth.IndexedBlobHash) ([]*eth.Blob, error)
}

type AltDAInputFetcher interface {
	// GetInput fetches the input for the given commitment at the given block number from the DA storage service.
	GetInput(ctx context.Context, l1 altda.L1Fetcher, c altda.CommitmentData, blockId eth.L1BlockRef) (eth.Data, error)
	// AdvanceL1Origin advances the L1 origin to the given block number, syncing the DA challenge events.
	AdvanceL1Origin(ctx context.Context, l1 altda.L1Fetcher, blockId eth.BlockID) error
	// Reset the challenge origin in case of L1 reorg
	Reset(ctx context.Context, base eth.L1BlockRef, baseCfg eth.SystemConfig) error
}

type ElectionWinnersProvider interface {
	GetElectionWinners() []*eth.ElectionWinner
}

// DataSourceFactory reads raw transactions from a given block & then filters for
// batch submitter transactions.
// This is not a stage in the pipeline, but a wrapper for another stage in the pipeline
type DataSourceFactory struct {
	log             log.Logger
	dsCfg           DataSourceConfig
	fetcher         L1Fetcher
	blobsFetcher    L1BlobsFetcher
	altDAFetcher    AltDAInputFetcher
	ecotoneTime     *uint64
	emitter         event.Emitter
	electionWinners []*eth.ElectionWinner
}

func NewDataSourceFactory(log log.Logger, cfg *rollup.Config, fetcher L1Fetcher, blobsFetcher L1BlobsFetcher, altDAFetcher AltDAInputFetcher) *DataSourceFactory {
	config := DataSourceConfig{
		l1Signer:          cfg.L1Signer(),
		batchInboxAddress: cfg.BatchInboxContractAddress,
		altDAEnabled:      cfg.AltDAEnabled(),
	}
	return &DataSourceFactory{
		log:             log,
		dsCfg:           config,
		fetcher:         fetcher,
		blobsFetcher:    blobsFetcher,
		altDAFetcher:    altDAFetcher,
		ecotoneTime:     cfg.EcotoneTime,
		emitter:         nil,
		electionWinners: []*eth.ElectionWinner{},
	}
}

// OpenData returns the appropriate data source for the L1 block `ref`.
func (ds *DataSourceFactory) OpenData(ctx context.Context, ref eth.L1BlockRef, batcherAddr common.Address) (DataIter, error) {
	// Creates a data iterator from blob or calldata source so we can forward it to the altDA source
	// if enabled as it still requires an L1 data source for fetching input commmitments.
	var src DataIter
	if ds.ecotoneTime != nil && ref.Time >= *ds.ecotoneTime {
		if ds.blobsFetcher == nil {
			return nil, fmt.Errorf("ecotone upgrade active but beacon endpoint not configured")
		}
		src = NewBlobDataSource(ctx, ds.log, ds.dsCfg, ds.fetcher, ds.blobsFetcher, ref, batcherAddr, ds)
	} else {
		src = NewCalldataSource(ctx, ds.log, ds.dsCfg, ds.fetcher, ref, batcherAddr)
	}
	if ds.dsCfg.altDAEnabled {
		// altDA([calldata | blobdata](l1Ref)) -> data
		return NewAltDADataSource(ds.log, src, ds.fetcher, ds.altDAFetcher, ref), nil
	}
	return src, nil
}

// DataSourceConfig regroups the mandatory rollup.Config fields needed for DataFromEVMTransactions.
type DataSourceConfig struct {
	l1Signer          types.Signer
	batchInboxAddress common.Address
	altDAEnabled      bool
}

func isValidBatchTx(receipt *types.Receipt, batcherAddr common.Address, logger log.Logger) bool {
	if receipt.Type != types.BlobTxType {
		// TODO(spire): enable other DA sources
		logger.Warn("not a blob tx")
		return false
	}
	batchInboxAbi := snapshots.LoadBatchInboxABI()
	topic0 := batchInboxAbi.Events["BatchSubmitted"].ID
	for _, log := range receipt.Logs {
		if log.Address != batcherAddr {
			continue
		}
		if log.Topics[0] != topic0 {
			continue
		}
		return true
	}

	return false
}

func (ds *DataSourceFactory) GetElectionWinners() []*eth.ElectionWinner {
	return ds.electionWinners
}

func (ds *DataSourceFactory) AttachEmitter(em event.Emitter) {
	ds.emitter = em
}

func (ds *DataSourceFactory) OnEvent(ev event.Event) bool {
	switch x := ev.(type) {
	case rollup.ElectionWinnerEvent:
		ds.electionWinners = x.ElectionWinners
	default:
		return false
	}
	return true
}
