package election

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-service/eth"
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

func (e *Election) GetWinnersAtEpoch(ctx context.Context, epoch uint64) ([]*eth.Validator, error) {
	validators := []*eth.Validator{}

	resp, err := e.client.GetLookahead(ctx, epoch)

	if err != nil {
		return validators, err
	}

	validators = resp.Data

	// For now winner will always return the current slots proposer, so we can just return the validators as is as each slot is its winner
	return validators, nil
}
