package derive

import (
	"bytes"
	"context"
	"encoding/binary"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/solabi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	BurnSignature = "burn(address)"
	BurnLen       = 4 + 32 // selector + padded address
)

var (
	BurnSelector         = crypto.Keccak256([]byte(BurnSignature))[:4]
	ElectionTickets      = common.HexToAddress("0x4200000000000000000000000000000000000028")
	BurnDepositerAddress = common.HexToAddress("0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001")
)

type BurnSource struct {
	Address common.Address
}

func (src *BurnSource) SourceHash() common.Hash {
	winnerHash := crypto.Keccak256Hash([]byte(src.Address.Hex()))

	var domainInput [32 * 2]byte
	binary.BigEndian.PutUint64(domainInput[32-8:32], BurnSourceDomain)
	copy(domainInput[32:], winnerHash[:])
	return crypto.Keccak256Hash(domainInput[:])

}

func BuildBurnTx(ctx context.Context, address common.Address) (*types.DepositTx, error) {
	source := BurnSource{
		Address: address,
	}

	data, err := marhsalBinaryBurnTx(address)
	if err != nil {
		return nil, err
	}

	// NOTE(spire): For some reason if system transaction is true, the tx fails
	// it passes with it being false though
	// i believe it is still forcefully included in the deriv pipeline anyway
	out := &types.DepositTx{
		SourceHash:          source.SourceHash(),
		From:                BurnDepositerAddress,
		IsSystemTransaction: false,
		To:                  &ElectionTickets,
		Mint:                nil,
		Value:               big.NewInt(0),
		Data:                data,
		Gas:                 1_000_000,
	}

	return out, nil
}

func BurnTxBytes(ctx context.Context, address common.Address) ([]byte, error) {
	tx, err := BuildBurnTx(ctx, address)
	if err != nil {
		return nil, err
	}

	l1Tx := types.NewTx(tx)
	opaqueL1Tx, err := l1Tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return opaqueL1Tx, nil

}

func marhsalBinaryBurnTx(address common.Address) ([]byte, error) {
	w := bytes.NewBuffer(make([]byte, 0, BurnLen))

	if err := solabi.WriteSignature(w, BurnSelector); err != nil {
		return nil, err
	}
	if err := solabi.WriteAddress(w, address); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
