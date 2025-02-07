package derive

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/packages/contracts-bedrock/snapshots"
)

// CalldataSource is a fault tolerant approach to fetching data.
// The constructor will never fail & it will instead re-attempt the fetcher
// at a later point.
type CalldataSource struct {
	// Internal state + data
	data []eth.Data
	// Required to re-attempt fetching
	ref     eth.L1BlockRef
	dsCfg   DataSourceConfig
	fetcher L1TransactionFetcher
	log     log.Logger

	electionClient ElectionWinnersProvider
}

// NewCalldataSource creates a new calldata source. It suppresses errors in fetching the L1 block if they occur.
// If there is an error, it will attempt to fetch the result on the next call to `Next`.
func NewCalldataSource(ctx context.Context, log log.Logger, dsCfg DataSourceConfig, fetcher L1TransactionFetcher, ref eth.L1BlockRef, electionClient ElectionWinnersProvider) DataIter {
	return &CalldataSource{
		ref:            ref,
		dsCfg:          dsCfg,
		fetcher:        fetcher,
		log:            log.New("origin", ref),
		electionClient: electionClient,
	}
}

// Next returns the next piece of data if it has it. If the constructor failed, this
// will attempt to reinitialize itself. If it cannot find the block it returns a ResetError
// otherwise it returns a temporary error if fetching the block returns an error.
func (ds *CalldataSource) Next(ctx context.Context) (eth.Data, error) {
	electionWinner := ds.electionClient.GetElectionWinner(ds.ref.Time)
	if electionWinner.Address == (common.Address{}) {
		ds.log.Warn("No election winner found for block", "time", ds.ref.Time)
		return nil, io.EOF
	}

	_, txs, err := ds.fetcher.InfoAndTxsByHash(ctx, ds.ref.Hash)
	if err != nil && errors.Is(err, ethereum.NotFound) {
		return nil, NewResetError(fmt.Errorf("failed to open calldata source: %w", err))
	} else if err != nil {
		return nil, NewTemporaryError(fmt.Errorf("failed to open calldata source: %w", err))
	}

	_, receipts, err := ds.fetcher.FetchReceipts(ctx, ds.ref.Hash)
	if err != nil {
		return nil, NewTemporaryError(fmt.Errorf("failed to fetch tx receipts: %w", err))
	}

	txsWithReceipts := make([]TxWithReceipt, len(txs))
	for i, tx := range txs {
		txsWithReceipts[i] = TxWithReceipt{tx: tx, receipt: receipts[i]}
	}

	ds.data = DataFromEVMTransactions(ds.dsCfg, electionWinner.Address, txsWithReceipts, ds.log)

	if len(ds.data) == 0 {
		return nil, io.EOF
	} else {
		data := ds.data[0]
		ds.data = ds.data[1:]
		return data, nil
	}
}

// DataFromEVMTransactions filters all of the transactions and returns the calldata from transactions
// that are sent to the batch inbox address from the batch sender address.
// This will return an empty array if no valid transactions are found.
func DataFromEVMTransactions(dsCfg DataSourceConfig, electionWinner common.Address, txs []TxWithReceipt, log log.Logger) []eth.Data {
	out := []eth.Data{}

	for _, tx := range txs {
		if tx.tx.Type() == types.BlobTxType {
			continue
		}
		if isValidBatchTx(tx.receipt, electionWinner, &dsCfg, log) {
			payload, err := ExtractPayload(tx.tx)
			if err != nil {
				log.Error("Could not extract payload", "err", err, "tx", tx.tx)
				continue
			}
			out = append(out, payload)
		}
	}
	return out
}

func ExtractPayload(tx *types.Transaction) (eth.Data, error) {
	data := tx.Data()

	// We need to skip past the selector and first argument (uint256 _targetTimestamp)
	batchInboxAbi := snapshots.LoadBatchInboxABI()
	method := batchInboxAbi.Methods["submitCalldata"]

	// Unpack the method inputs
	inputs, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		return nil, err
	}

	return inputs[1].([]byte), nil
}
