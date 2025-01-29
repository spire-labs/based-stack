package election_client

import (
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type ElectionClient struct {
	store *ElectionStore
}

func NewElectionClient(store *ElectionStore) *ElectionClient {
	return &ElectionClient{
		store: store,
	}
}

func (e *ElectionClient) GetElectionWinner(timestamp uint64) eth.ElectionWinner {
	return e.store.GetElectionWinner(timestamp)
}

func (e *ElectionClient) GetLastWinnerInCurrentEpoch() eth.ElectionWinner {
	return e.store.GetLastWinnerInCurrentEpoch()
}
