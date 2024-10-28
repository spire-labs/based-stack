package election

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
)

type ElectionClient interface {
	GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error)
	GetEpochNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetSlotNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetTimeFromSlot(ctx context.Context, slot uint64) (uint64, error)
}

type Election struct {
	client ElectionClient
}

func NewElection(client ElectionClient) *Election {
	return &Election{
		client: client}
}

func (e *Election) GetWinnersAtEpoch(ctx context.Context, epoch uint64) ([]*eth.ElectionWinner, error) {
	resp, err := e.client.GetLookahead(ctx, epoch)
	var electionWinners []*eth.ElectionWinner

	if err != nil {
		return []*eth.ElectionWinner{}, err
	}

	validators := resp.Data

	for _, validator := range validators {
		address := common.BytesToAddress(validator.Pubkey[:20])
		time, err := e.client.GetTimeFromSlot(ctx, uint64(validator.Slot))
		if err != nil {
			return []*eth.ElectionWinner{}, err
		}
		winner := eth.ElectionWinner{
			Address: address,
			Time:    time,
		}

		electionWinners = append(electionWinners, &winner)
	}

	// For now winner will always return the current slots proposer, so we can just return the validators as is as each slot is its winner
	return electionWinners, nil
}
