package derive

import (
	"fmt"
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

type MockElectionWinnersProvider struct{}

func (m *MockElectionWinnersProvider) GetElectionWinners() []*eth.ElectionWinner {
	return []*eth.ElectionWinner{
		{
			Address: common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			Time:    0x499602D2,
		},
		{
			Address: common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Time:    9876543210,
		},
	}
}

func TestDataAndHashesFromTxs(t *testing.T) {
	// test setup
	rng := rand.New(rand.NewSource(12345))
	privateKey := testutils.InsecureRandomKey(rng)
	batchInboxAddr := testutils.RandomAddress(rng)
	logger := testlog.Logger(t, log.LvlInfo)

	batchInboxAbi := snapshots.LoadBatchInboxABI()
	batchSubmittedEventTopic := batchInboxAbi.Events["BatchSubmitted"].ID

	chainId := new(big.Int).SetUint64(rng.Uint64())
	signer := types.NewCancunSigner(chainId)
	config := DataSourceConfig{
		l1Signer:          signer,
		batchInboxAddress: batchInboxAddr,
	}

	mockElectionProvider := &MockElectionWinnersProvider{}

	// create an instance of the blob data source for testing w/o calling a function. Just create the struct
	ds := BlobDataSource{
		ref:              eth.L1BlockRef{Time: 0x499602D2},
		dsCfg:            config,
		fetcher:          nil,
		log:              logger,
		batcherAddr:      batchInboxAddr,
		blobsFetcher:     nil,
		electionProvider: mockElectionProvider,
	}

	// TODO(miszke): enable other DA sources
	// create a valid non-blob batcher transaction and make sure it's picked up
	// txData := &types.LegacyTx{
	// 	Nonce:    rng.Uint64(),
	// 	GasPrice: new(big.Int).SetUint64(rng.Uint64()),
	// 	Gas:      2_000_000,
	// 	To:       &batchInboxAddr,
	// 	Value:    big.NewInt(10),
	// 	Data:     testutils.RandomData(rng, rng.Intn(1000)),
	// }
	// calldataTx, _ := types.SignNewTx(privateKey, signer, txData)
	// txs := types.Transactions{calldataTx}
	// data, blobHashes := dataAndHashesFromTxs(txs, &config, batcherAddr, logger)
	// require.Equal(t, 1, len(data))
	// require.Equal(t, 0, len(blobHashes))

	electionWinnerAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	// create a valid blob batcher tx and make sure it's picked up
	blobHash := testutils.RandomHash(rng)
	blobTxData := &types.BlobTx{
		Nonce:      rng.Uint64(),
		Gas:        2_000_000,
		To:         electionWinnerAddress,
		BlobHashes: []common.Hash{blobHash},
	}
	blobTx, _ := types.SignNewTx(privateKey, signer, blobTxData)
	receipt := &types.Receipt{
		Type: types.BlobTxType,
		Logs: []*types.Log{{
			Address: electionWinnerAddress,
			Topics:  []common.Hash{batchSubmittedEventTopic},
		}},
	}
	txs := []TxWithReceipt{{tx: blobTx, receipt: receipt}}
	data, blobHashes := ds.dataAndHashesFromTxs(txs, &config, logger)
	require.Equal(t, 1, len(data))
	require.Equal(t, 1, len(blobHashes))
	require.Nil(t, data[0].calldata)

	// TODO(miszke): enable other DA sources
	// try again with both the blob & calldata transactions and make sure both are picked up
	// txs = types.Transactions{blobTx, calldataTx}
	// data, blobHashes = dataAndHashesFromTxs(txs, &config, batcherAddr, logger)
	// require.Equal(t, 2, len(data))
	// require.Equal(t, 1, len(blobHashes))
	// require.NotNil(t, data[1].calldata)

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
	batcherAddr := testutils.RandomAddress(rng)
	logger := testlog.Logger(t, log.LvlInfo)

	batchInboxAbi := snapshots.LoadBatchInboxABI()
	batchSubmittedEventTopic := batchInboxAbi.Events["BatchSubmitted"].ID

	chainId := new(big.Int).SetUint64(rng.Uint64())
	signer := types.NewCancunSigner(chainId)

	ds := BlobDataSource{
		log: logger,
	}

	t.Run("Valid blob batch transaction", func(t *testing.T) {
		electionWinnerAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
		blobTxData := &types.BlobTx{
			Nonce:      rng.Uint64(),
			Gas:        2_000_000,
			To:         electionWinnerAddress,
			BlobHashes: []common.Hash{testutils.RandomHash(rng)},
		}
		blobTx, _ := types.SignNewTx(privateKey, signer, blobTxData)
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{{
				Address: batcherAddr,
				Topics:  []common.Hash{batchSubmittedEventTopic},
			}},
		}
		fmt.Println(blobTx)
		valid := ds.isValidBatchTx(receipt, batcherAddr, logger)
		require.True(t, valid, "Expected transaction to be valid")
	})

	t.Run("Invalid receipt type", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.LegacyTxType,
			Logs: []*types.Log{{
				Address: batcherAddr,
				Topics:  []common.Hash{batchSubmittedEventTopic},
			}},
		}

		valid := ds.isValidBatchTx(receipt, batcherAddr, logger)
		require.False(t, valid, "Expected transaction to be invalid due to receipt type")
	})

	t.Run("Log address does not match batcher address", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{{
				Address: testutils.RandomAddress(rng),
				Topics:  []common.Hash{batchSubmittedEventTopic},
			}},
		}

		valid := ds.isValidBatchTx(receipt, batcherAddr, logger)
		require.False(t, valid, "Expected transaction to be invalid due to log address mismatch")
	})

	t.Run("Log topic does not match BatchSubmitted event", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{{
				Address: batcherAddr,
				Topics:  []common.Hash{testutils.RandomHash(rng)},
			}},
		}

		valid := ds.isValidBatchTx(receipt, batcherAddr, logger)
		require.False(t, valid, "Expected transaction to be invalid due to incorrect log topic")
	})

	t.Run("Receipt has no logs", func(t *testing.T) {
		receipt := &types.Receipt{
			Type: types.BlobTxType,
			Logs: []*types.Log{},
		}

		valid := ds.isValidBatchTx(receipt, batcherAddr, logger)
		require.False(t, valid, "Expected transaction to be invalid due to missing logs")
	})
}
