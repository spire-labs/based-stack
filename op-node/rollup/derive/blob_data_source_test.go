package derive

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
	"github.com/ethereum-optimism/optimism/packages/contracts-bedrock/snapshots"
	"github.com/ethereum/go-ethereum/log"
)

type MockElectionWinnersProvider struct {
	electionWinner common.Address
}

func (m *MockElectionWinnersProvider) GetElectionWinner(_timestamp uint64) eth.ElectionWinner {
	return eth.ElectionWinner{
		Address: m.electionWinner,
		Time:    0x499602D2,
	}
}

func TestDataAndHashesFromTxs(t *testing.T) {
	// test setup
	rng := rand.New(rand.NewSource(12345))
	privateKey := testutils.InsecureRandomKey(rng)
	batchInboxAddr := testutils.RandomAddress(rng)
	electionWinnerAddr := testutils.RandomAddress(rng)
	logger := testlog.Logger(t, log.LvlInfo)

	batchInboxAbi := snapshots.LoadBatchInboxABI()
	batchSubmittedEventTopic := batchInboxAbi.Events["BatchSubmitted"].ID

	chainId := new(big.Int).SetUint64(rng.Uint64())
	signer := types.NewCancunSigner(chainId)
	config := DataSourceConfig{
		l1Signer:          signer,
		batchInboxAddress: batchInboxAddr,
	}

	mockElectionProvider := &MockElectionWinnersProvider{electionWinner: electionWinnerAddr}

	// create an instance of the blob data source for testing w/o calling a function. Just create the struct
	ds := BlobDataSource{
		ref:            eth.L1BlockRef{Time: 0x499602D2},
		dsCfg:          config,
		fetcher:        nil,
		log:            logger,
		batcherAddr:    batchInboxAddr,
		blobsFetcher:   nil,
		electionClient: mockElectionProvider,
	}

	// create a valid non-blob batcher transaction and make sure it's picked up

	calldata := submitCalldataTxData(t, testutils.RandomData(rng, rng.Intn(1000)))
	calldataTxData := &types.LegacyTx{
		Nonce:    rng.Uint64(),
		GasPrice: new(big.Int).SetUint64(rng.Uint64()),
		Gas:      2_000_000,
		To:       &batchInboxAddr,
		Value:    big.NewInt(10),
		Data:     calldata,
	}
	calldataTx, _ := types.SignNewTx(privateKey, signer, calldataTxData)
	calldataReceipt := &types.Receipt{
		Type: types.LegacyTxType,
		Logs: []*types.Log{{
			Address: batchInboxAddr,
			Topics:  []common.Hash{batchSubmittedEventTopic, padAddress(electionWinnerAddr)},
		}},
	}
	txs := []TxWithReceipt{{tx: calldataTx, receipt: calldataReceipt}}
	data, blobHashes := ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 1, len(data))
	require.Equal(t, 0, len(blobHashes))

	// create a valid blob batcher tx and make sure it's picked up
	blobHash := testutils.RandomHash(rng)
	blobTxData := &types.BlobTx{
		Nonce:      rng.Uint64(),
		Gas:        2_000_000,
		To:         batchInboxAddr,
		BlobHashes: []common.Hash{blobHash},
	}
	blobTx, _ := types.SignNewTx(privateKey, signer, blobTxData)
	blobReceipt := &types.Receipt{
		Logs: []*types.Log{{
			Address: batchInboxAddr,
			Topics:  []common.Hash{batchSubmittedEventTopic, padAddress(electionWinnerAddr)},
		}},
	}
	txs = []TxWithReceipt{{tx: blobTx, receipt: blobReceipt}}
	data, blobHashes = ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 1, len(data))
	require.Equal(t, 1, len(blobHashes))
	require.Nil(t, data[0].calldata)

	// try again with both the blob & calldata transactions and make sure both are picked up
	txs = []TxWithReceipt{{tx: blobTx, receipt: blobReceipt}, {tx: calldataTx, receipt: calldataReceipt}}
	data, blobHashes = ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 2, len(data))
	require.Equal(t, 1, len(blobHashes))
	require.NotNil(t, data[1].calldata)

	// make sure calldata tx to the batch inbox is ignored if election winner is invalid
	randomAddr := testutils.RandomAddress(rng)
	calldataTxData.Data = testutils.RandomData(rng, rng.Intn(1000))
	calldataTx, _ = types.SignNewTx(privateKey, signer, calldataTxData)
	txs = []TxWithReceipt{{tx: calldataTx, receipt: &types.Receipt{
		Logs: []*types.Log{{
			Address: batchInboxAddr,
			Topics:  []common.Hash{batchSubmittedEventTopic, padAddress(randomAddr)},
		}},
	}}}
	data, blobHashes = ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 0, len(data))
	require.Equal(t, 0, len(blobHashes))

	// make sure blob tx to the batch inbox is ignored if not calling the submit fn
	blobTxData.Data = testutils.RandomData(rng, rng.Intn(1000))
	blobTx, _ = types.SignNewTx(privateKey, signer, blobTxData)
	txs = []TxWithReceipt{{tx: blobTx, receipt: &types.Receipt{}}}
	data, blobHashes = ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 0, len(data))
	require.Equal(t, 0, len(blobHashes))

	// make sure blob tx ignored if the tx isn't going to the batch inbox addr, even if the
	// signature is valid.
	blobTxData.To = testutils.RandomAddress(rng)
	blobTx, _ = types.SignNewTx(privateKey, signer, blobTxData)
	txs = []TxWithReceipt{{tx: blobTx, receipt: &types.Receipt{}}}
	data, blobHashes = ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 0, len(data))
	require.Equal(t, 0, len(blobHashes))
}

func TestFillBlobPointers(t *testing.T) {
	blob := eth.Blob{}
	rng := rand.New(rand.NewSource(1234))
	calldata := eth.Data{}

	for i := 0; i < 100; i++ {
		// create a random length input data array w/ len = [0-10)
		dataLen := rng.Intn(10)
		data := make([]blobOrCalldata, dataLen)

		// pick some subset of those to be blobs, and the rest calldata
		blobLen := 0
		if dataLen != 0 {
			blobLen = rng.Intn(dataLen)
		}
		calldataLen := dataLen - blobLen

		// fill in the calldata entries at random indices
		for j := 0; j < calldataLen; j++ {
			randomIndex := rng.Intn(dataLen)
			for data[randomIndex].calldata != nil {
				randomIndex = (randomIndex + 1) % dataLen
			}
			data[randomIndex].calldata = &calldata
		}

		// create the input blobs array and call fillBlobPointers on it
		blobs := make([]*eth.Blob, blobLen)
		for j := 0; j < blobLen; j++ {
			blobs[j] = &blob
		}
		err := fillBlobPointers(data, blobs)
		require.NoError(t, err)

		// check that we get the expected number of calldata vs blobs results
		blobCount := 0
		calldataCount := 0
		for j := 0; j < dataLen; j++ {
			if data[j].calldata != nil {
				calldataCount++
			}
			if data[j].blob != nil {
				blobCount++
			}
		}
		require.Equal(t, blobLen, blobCount)
		require.Equal(t, calldataLen, calldataCount)
	}
}

func TestIsValidBatchTx(t *testing.T) {
	rng := rand.New(rand.NewSource(12345))
	privateKey := testutils.InsecureRandomKey(rng)
	batchInboxAddr := testutils.RandomAddress(rng)
	electionWinnerAddr := testutils.RandomAddress(rng)
	randomAddr := testutils.RandomAddress(rng)

	logger := testlog.Logger(t, log.LvlInfo)
	cfg := &DataSourceConfig{
		batchInboxAddress: batchInboxAddr,
	}

	electionWinnerTopic := padAddress(electionWinnerAddr)

	batchInboxAbi := snapshots.LoadBatchInboxABI()
	batchSubmittedEventTopic := batchInboxAbi.Events["BatchSubmitted"].ID

	chainId := new(big.Int).SetUint64(rng.Uint64())
	signer := types.NewCancunSigner(chainId)

	t.Run("Valid blob batch transaction with correct election winner", func(t *testing.T) {
		blobTxData := &types.BlobTx{
			Nonce:      rng.Uint64(),
			Gas:        2_000_000,
			To:         electionWinnerAddr,
			BlobHashes: []common.Hash{testutils.RandomHash(rng)},
		}
		_, _ = types.SignNewTx(privateKey, signer, blobTxData)
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{{
				Address: batchInboxAddr,
				Topics:  []common.Hash{batchSubmittedEventTopic, electionWinnerTopic},
			}},
		}
		valid := isValidBatchTx(receipt, electionWinnerAddr, cfg, logger)
		require.True(t, valid, "Expected transaction and winner to be valid")
	})

	t.Run("Valid blob batch transaction with incorrect election winner", func(t *testing.T) {
		blobTxData := &types.BlobTx{
			Nonce:      rng.Uint64(),
			Gas:        2_000_000,
			To:         batchInboxAddr,
			BlobHashes: []common.Hash{testutils.RandomHash(rng)},
		}
		_, _ = types.SignNewTx(privateKey, signer, blobTxData)
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{{
				Address: batchInboxAddr,
				Topics:  []common.Hash{batchSubmittedEventTopic, padAddress(randomAddr)},
			}},
		}
		valid := isValidBatchTx(receipt, batchInboxAddr, cfg, logger)
		require.False(t, valid, "Expected transaction to be invalid due to incorrect election winner")
	})

	t.Run("Invalid receipt type", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.LegacyTxType,
			Logs: []*types.Log{{
				Address: electionWinnerAddr,
				Topics:  []common.Hash{batchSubmittedEventTopic, electionWinnerTopic},
			}},
		}

		valid := isValidBatchTx(receipt, electionWinnerAddr, cfg, logger)
		require.False(t, valid, "Expected transaction to be invalid due to receipt type")
	})

	t.Run("Log topic does not match BatchSubmitted event", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{{
				Address: batchInboxAddr,
				Topics:  []common.Hash{testutils.RandomHash(rng)},
			}},
		}

		valid := isValidBatchTx(receipt, electionWinnerAddr, cfg, logger)
		require.False(t, valid, "Expected transaction to be invalid due to incorrect log topic")
	})

	t.Run("Receipt has no logs", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{},
		}

		valid := isValidBatchTx(receipt, electionWinnerAddr, cfg, logger)
		require.False(t, valid, "Expected transaction to be invalid due to missing logs")
	})
}

func padAddress(address common.Address) common.Hash {
	padded := make([]byte, 32)
	copy(padded[32-len(address.Bytes()):], address.Bytes())
	return common.BytesToHash(padded)
}
