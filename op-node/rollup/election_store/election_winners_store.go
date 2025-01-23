package election_store

import (
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionWinnersStore struct {
	electionWinnersMap map[uint64]*eth.ElectionWinner
	log                log.Logger
}

func NewElectionWinnersStore(log log.Logger) *ElectionWinnersStore {
	return &ElectionWinnersStore{
		electionWinnersMap: make(map[uint64]*eth.ElectionWinner),
		log:                log,
	}
}

func (e *ElectionWinnersStore) GetElectionWinner(timestamp uint64) *eth.ElectionWinner {
	return e.electionWinnersMap[timestamp]
}

func (e *ElectionWinnersStore) StoreElectionWinners(winners []*eth.ElectionWinner) {
	for _, winner := range winners {
		e.electionWinnersMap[winner.Time] = winner
	}
}

func (e *ElectionWinnersStore) RemoveOutdatedElectionWinners(timestamp uint64) {
	for k := range e.electionWinnersMap {
		if k < timestamp {
			delete(e.electionWinnersMap, k)
		}
	}
}

func (e *ElectionWinnersStore) WinnersLength() int {
	return len(e.electionWinnersMap)
}
