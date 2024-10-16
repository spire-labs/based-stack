package election

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionClient interface {
	GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error)
	GetEpochNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetSlotNumber(ctx context.Context, timestamp uint64) (uint64, error)
}

type Election struct {
	client ElectionClient
}

func NewElection(client ElectionClient) *Election {
	return &Election{
		client: client}
}

func (e *Election) GetWinner(ctx context.Context, timestamp uint64, log log.Logger) (eth.Validator, error) {
	var err error

	epoch, err := e.client.GetEpochNumber(ctx, timestamp)

	if err != nil {
		return eth.Validator{}, err
	}

	data, err := e.client.GetLookahead(ctx, epoch)

	if err != nil {
		return eth.Validator{}, err
	}

	slot, err := e.client.GetSlotNumber(ctx, timestamp)

	if err != nil {
		return eth.Validator{}, err
	}

	length := uint64(len(data.Data))

	log.Info("Choosing validator at slot", "slot", slot)
	// For now winner will always return the current slots proposer
	return *data.Data[slot%length], nil
}
