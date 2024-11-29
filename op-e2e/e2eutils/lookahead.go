package e2eutils

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/core"
)

type FakeLookahead struct {
	genesisTimestamp uint64
	blockTime        uint64
	slotsPerEpoch    uint64
}

func NewBeaconClient(genesis *core.Genesis) FakeLookahead {
	return FakeLookahead{
		slotsPerEpoch:    8,
		blockTime:        12,
		genesisTimestamp: genesis.Timestamp}
}

func (l FakeLookahead) GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error) {
	out := eth.APIGetLookaheadResponse{}
	address, err := hex.DecodeString(strings.TrimPrefix("0x7c60541eB6f54f0F3c8B34D0a00De9045d2f5534", "0x"))
	if err != nil {
		return out, err
	}

	var pubkey eth.Bytes48
	copy(pubkey[:], address[:])

	for i := 0; i < int(l.slotsPerEpoch); i++ {
		if err != nil {
			return eth.APIGetLookaheadResponse{}, err
		}
		// TODO: fix parent slot
		out.Data = append(out.Data, &eth.Validator{Pubkey: pubkey, Slot: eth.Uint64String(i)})
	}

	return out, nil
}

func (l FakeLookahead) GetEpochNumber(ctx context.Context, timestamp uint64) (uint64, error) {
	elapsed := timestamp - l.genesisTimestamp
	if elapsed%l.blockTime != 0 {
		return 0, fmt.Errorf("timestamp invalid: does not match any slot. timestamp: %d", timestamp)
	}

	epoch := (elapsed / l.blockTime) / l.slotsPerEpoch
	return epoch, nil
}

// TODO(spire): is it even used?
func (l FakeLookahead) GetSlotNumber(ctx context.Context, timestamp uint64) (uint64, error) {
	return 1, nil
}

func (l FakeLookahead) GetTimeFromSlot(ctx context.Context, slot uint64) (uint64, error) {
	return l.genesisTimestamp + slot*l.blockTime, nil
}
