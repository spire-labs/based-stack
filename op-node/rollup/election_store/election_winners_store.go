package election_store

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionWinnersStore struct {
	electionWinnersMapByTime       map[uint64]*eth.ElectionWinner
	electionWinnersMapByParentSlot map[uint64]*eth.ElectionWinner
	log                            log.Logger
	latestWinner                   *eth.ElectionWinner
}

func NewElectionWinnersStore(log log.Logger) *ElectionWinnersStore {
	return &ElectionWinnersStore{
		electionWinnersMapByTime:       make(map[uint64]*eth.ElectionWinner),
		electionWinnersMapByParentSlot: make(map[uint64]*eth.ElectionWinner),
		log:                            log,
	}
}

func (e *ElectionWinnersStore) GetElectionWinnerByTime(timestamp uint64) *eth.ElectionWinner {
	return e.electionWinnersMapByTime[timestamp]
}

func (e *ElectionWinnersStore) GetElectionWinnerByParentSlot(timestamp uint64) *eth.ElectionWinner {
	return e.electionWinnersMapByParentSlot[timestamp]
}

func (e *ElectionWinnersStore) GetLatestElectionWinner() *eth.ElectionWinner {
	return e.latestWinner
}

// StoreElectionWinners stores election winners in the electionWinnersMap
//
// Parameters:
// - winners: election winners to store (sorted by timestamp asc)
func (e *ElectionWinnersStore) StoreElectionWinners(winners []*eth.ElectionWinner) {
	if len(winners) == 0 {
		e.log.Warn("No election winners to store")
		return
	}

	for _, winner := range winners {
		e.electionWinnersMapByTime[winner.Time] = winner
		e.electionWinnersMapByParentSlot[winner.ParentSlot] = winner
	}

	newLatest := winners[len(winners)-1]
	if e.latestWinner == nil || newLatest.Time > e.latestWinner.Time {
		e.latestWinner = newLatest
	}
}

func (e *ElectionWinnersStore) RemoveOutdatedElectionWinners(timestamp uint64) {
	e.log.Debug("Removing outdated election winners", "time", timestamp, "store", e.WinnersLength())

	for k := range e.electionWinnersMapByTime {
		if k < timestamp {
			delete(e.electionWinnersMapByTime, k)
		}
	}

	for k, v := range e.electionWinnersMapByParentSlot {
		if v.Time < timestamp {
			delete(e.electionWinnersMapByParentSlot, k)
		}
	}
	e.log.Debug("Removed outdated election winners", "map", e.WinnersLength())
}

func (e *ElectionWinnersStore) WinnersLength() string {
	return fmt.Sprintf("time: %d, parentSlot: %d", len(e.electionWinnersMapByTime), len(e.electionWinnersMapByParentSlot))
}
