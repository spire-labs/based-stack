package derive

import (
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testlog"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
	"github.com/ethereum-optimism/optimism/packages/contracts-bedrock/snapshots"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

type testTx struct {
	to      *common.Address
	dataLen int
	author  *ecdsa.PrivateKey
	good    bool
	value   int
	rawData []byte
}

func (tx *testTx) Create(t *testing.T, signer types.Signer, rng *rand.Rand) TxWithReceipt {
	t.Helper()
	tx.rawData = testutils.RandomData(rng, tx.dataLen)
	calldata := submitCalldataTxData(t, tx.rawData)
	outTx, err := types.SignNewTx(tx.author, signer, &types.DynamicFeeTx{
		ChainID:   signer.ChainID(),
		Nonce:     0,
		GasTipCap: big.NewInt(2 * params.GWei),
		GasFeeCap: big.NewInt(30 * params.GWei),
		Gas:       100_000,
		To:        tx.to,
		Value:     big.NewInt(int64(tx.value)),
		Data:      calldata,
	})
	require.NoError(t, err)

	authorAddr := crypto.PubkeyToAddress(tx.author.PublicKey)
	batchInboxAbi := snapshots.LoadBatchInboxABI()
	batchSubmittedEventTopic := batchInboxAbi.Events["BatchSubmitted"].ID

	outReceipt := &types.Receipt{
		Logs: []*types.Log{{
			Address: *tx.to,
			Topics:  []common.Hash{batchSubmittedEventTopic, padAddress(authorAddr)},
		}},
	}

	return TxWithReceipt{
		tx:      outTx,
		receipt: outReceipt,
	}
}

type calldataTest struct {
	name string
	txs  []testTx
}

// TestDataFromEVMTransactions creates some transactions from a specified template and asserts
// that DataFromEVMTransactions properly filters and returns the data from the authorized transactions
// inside the transaction set.
func TestDataFromEVMTransactions(t *testing.T) {
	inboxPriv := testutils.RandomKey()
	electionWinnerPriv := testutils.RandomKey()
	cfg := &rollup.Config{
		L1ChainID:                 big.NewInt(100),
		BatchInboxContractAddress: crypto.PubkeyToAddress(inboxPriv.PublicKey),
	}
	electionWinnerAddr := crypto.PubkeyToAddress(electionWinnerPriv.PublicKey)

	altInbox := testutils.RandomAddress(rand.New(rand.NewSource(1234)))
	altAuthor := testutils.RandomKey()

	testCases := []calldataTest{
		{
			name: "simple",
			txs:  []testTx{{to: &cfg.BatchInboxContractAddress, dataLen: 1234, author: electionWinnerPriv, good: true}},
		},
		{
			name: "other inbox",
			txs:  []testTx{{to: &altInbox, dataLen: 1234, author: electionWinnerPriv, good: false}}},
		{
			name: "other author",
			txs:  []testTx{{to: &cfg.BatchInboxContractAddress, dataLen: 1234, author: altAuthor, good: false}}},
		{
			name: "inbox is author",
			txs:  []testTx{{to: &cfg.BatchInboxContractAddress, dataLen: 1234, author: inboxPriv, good: false}}},
		{
			name: "author is inbox",
			txs:  []testTx{{to: &electionWinnerAddr, dataLen: 1234, author: electionWinnerPriv, good: false}}},
		{
			name: "unrelated",
			txs:  []testTx{{to: &altInbox, dataLen: 1234, author: altAuthor, good: false}}},
		{
			name: "empty tx",
			txs:  []testTx{{to: &cfg.BatchInboxContractAddress, dataLen: 0, author: electionWinnerPriv, good: true}}},
		{
			name: "value tx",
			txs:  []testTx{{to: &cfg.BatchInboxContractAddress, dataLen: 1234, value: 42, author: electionWinnerPriv, good: true}}},
		{
			name: "empty block", txs: []testTx{},
		},
		{
			name: "mixed txs",
			txs: []testTx{
				{to: &cfg.BatchInboxContractAddress, dataLen: 1234, value: 42, author: electionWinnerPriv, good: true},
				{to: &cfg.BatchInboxContractAddress, dataLen: 3333, value: 32, author: altAuthor, good: false},
				{to: &cfg.BatchInboxContractAddress, dataLen: 2000, value: 22, author: electionWinnerPriv, good: true},
				{to: &altInbox, dataLen: 2020, value: 12, author: electionWinnerPriv, good: false},
			},
		},
		// TODO: test with different batcher key, i.e. when it's changed from initial config value by L1 contract
	}

	for i, tc := range testCases {
		rng := rand.New(rand.NewSource(int64(i)))
		signer := cfg.L1Signer()

		var expectedData []eth.Data
		var txs []TxWithReceipt
		for _, tx := range tc.txs {
			txs = append(txs, tx.Create(t, signer, rng))
			if tx.good {
				expectedData = append(expectedData, tx.rawData)
			}
		}

		out := DataFromEVMTransactions(DataSourceConfig{cfg.L1Signer(), cfg.BatchInboxContractAddress, false}, electionWinnerAddr, txs, testlog.Logger(t, log.LevelCrit))
		require.ElementsMatch(t, expectedData, out)
	}

}

func submitCalldataTxData(t require.TestingT, payload []byte) []byte {
	batchInboxAbi := snapshots.LoadBatchInboxABI()
	submitCalldataMethod, ok := batchInboxAbi.Methods["submitCalldata"]
	require.True(t, ok)

	inputs, err := submitCalldataMethod.Inputs.Pack(new(big.Int).SetUint64(0), payload)
	require.NoError(t, err)

	submitSel := submitCalldataMethod.ID
	return append(submitSel[:], inputs...)
}
