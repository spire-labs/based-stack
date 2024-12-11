package e2eutils

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/core"
)

type FakeLookahead struct {
	genesisTimestamp uint64
	blockTime        uint64
	slotsPerEpoch    uint64
}

func NewBeaconClient(genesis *core.Genesis, cfg *genesis.DeployConfig) FakeLookahead {
	return FakeLookahead{
		// TODO(spire): mocking a long lookahead makes it easy to fix some tests.
		// we should test the epoch switch in e2e as well.
		// This field should be configurable.
		slotsPerEpoch:    90,
		blockTime:        cfg.L1BlockTime,
		genesisTimestamp: genesis.Timestamp}
}

func (l FakeLookahead) GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error) {
	out := eth.APIGetLookaheadResponse{}
	address, err := hex.DecodeString(strings.TrimPrefix("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "0x"))
	if err != nil {
		return out, err
	}

	var pubkey eth.Bytes48
	copy(pubkey[:], address[:])

	for i := 0; i < int(l.slotsPerEpoch); i++ {
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
