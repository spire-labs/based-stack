package election_store

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

func (e *ElectionClient) GetElectionWinnerByTime(timestamp uint64) *eth.ElectionWinner {
	return e.store.GetElectionWinnerByTime(timestamp)
}

func (e *ElectionClient) GetElectionWinnerByParentSlot(timestamp uint64) *eth.ElectionWinner {
	return e.store.GetElectionWinnerByParentSlot(timestamp)
}

func (e *ElectionClient) GetLatestElectionWinner() *eth.ElectionWinner {
	return e.store.GetLatestElectionWinner()
}
