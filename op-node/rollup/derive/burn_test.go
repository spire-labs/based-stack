package derive

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestSourceHash(t *testing.T) {
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	burnSource := BurnSource{
		Address: address,
	}

	sourceHash := burnSource.SourceHash()

	assert.NotNil(t, sourceHash, "SourceHash should not be nil")
	assert.Equal(t, 32, len(sourceHash.Bytes()), "SourceHash should be 32 bytes")
}

func TestMarshalBinaryBurnTx(t *testing.T) {
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	data, err := marshsalBinaryBurnTx(address)

	assert.NoError(t, err, "marshsalBinaryBurnTx should not return an error")
	assert.Equal(t, BurnLen, len(data), "marshsalBinaryBurnTx should return data of correct length")
}

func TestBuildBurnTx(t *testing.T) {
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	tx, err := BuildBurnTx(address)

	assert.NoError(t, err, "BuildBurnTx should not return an error")
	assert.NotNil(t, tx, "BuildBurnTx should return a valid transaction")

	assert.Equal(t, BurnDepositorAddress, tx.From, "Transaction 'From' field should match BurnDepositorAddress")
	assert.Equal(t, &ElectionTickets, tx.To, "Transaction 'To' field should match ElectionTickets")
	assert.Equal(t, big.NewInt(0), tx.Value, "Transaction 'Value' should be zero")
	assert.False(t, tx.IsSystemTransaction, "IsSystemTransaction should be false")
	assert.NotNil(t, tx.Data, "Transaction 'Data' should not be nil")
}

func TestBurnTxBytes(t *testing.T) {
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	data, err := BurnTxBytes(address)

	assert.NoError(t, err, "BurnTxBytes should not return an error")
	assert.NotNil(t, data, "BurnTxBytes should return valid byte array")
}

func TestBuildAndMarshalBurnTx(t *testing.T) {
	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")

	// Build transaction
	tx, err := BuildBurnTx(address)
	assert.NoError(t, err, "BuildBurnTx should not return an error")
	assert.NotNil(t, tx, "BuildBurnTx should return a valid transaction")

	// Marshal the transaction
	marshaledTx, err := BurnTxBytes(address)
	assert.NoError(t, err, "BurnTxBytes should not return an error")
	assert.NotNil(t, marshaledTx, "BurnTxBytes should return valid byte array")

	// Check if the marshaled transaction is non-empty
	assert.NotEqual(t, 0, len(marshaledTx), "Marshalled transaction should not be empty")
}

func TestBurnSelector(t *testing.T) {
	expectedSelector := crypto.Keccak256([]byte(BurnSignature))[:4]
	assert.Equal(t, BurnSelector, expectedSelector, "BurnSelector should match the Keccak256 hash of the signature")
}
