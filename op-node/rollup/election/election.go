package election

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type ElectionClient interface {
	GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error)
	GetCurrentEpoch(ctx context.Context) (uint64, error)
}

type Election struct {
	client ElectionClient
}

func NewElection(client ElectionClient) *Election {
	return &Election{
		client: client,
	}
}

func (e *Election) GetWinner(ctx context.Context) (eth.Validator, error) {
	var err error

	epoch, err := e.client.GetCurrentEpoch(ctx)

	if err != nil {
		return eth.Validator{}, err
	}

	data, err := e.client.GetLookahead(ctx, epoch)

	if err != nil {
		return eth.Validator{}, err
	}

	// TODO: Adjust this to add the election rules
	// Just return the current proposer for now
	return *data.Data[0], nil
}
