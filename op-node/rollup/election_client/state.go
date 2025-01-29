package election_client

import (
	"fmt"
	"sync"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionStore struct {
	electionWinnersMap map[uint64]*eth.ElectionWinner
	log                log.Logger
	latestWinner       *eth.ElectionWinner

	mu sync.Mutex
}

func NewElectionStore(log log.Logger) *ElectionStore {
	return &ElectionStore{
		electionWinnersMap: make(map[uint64]*eth.ElectionWinner),
		log:                log,
	}
}

func (e *ElectionStore) OnEvent(ev event.Event) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch x := ev.(type) {
	case rollup.ElectionWinnerEvent:
		e.StoreElectionWinners(x.ElectionWinners)
	case rollup.ElectionWinnerOutdatedEvent:
		// remove all election winners with a timestamp less than the outdated timestamp
		e.RemoveOutdatedElectionWinners(x.Time)
	default:
		return false
	}

	return true
}

func (e *ElectionStore) GetElectionWinner(timestamp uint64) eth.ElectionWinner {
	out := e.electionWinnersMap[timestamp]
	if out == nil {
		return eth.ElectionWinner{}
	}
	return *out
}

func (e *ElectionStore) GetLastWinnerInCurrentEpoch() eth.ElectionWinner {
	out := e.latestWinner
	if out == nil {
		return eth.ElectionWinner{}
	}
	return *out
}

// StoreElectionWinners stores election winners in the electionWinnersMap
//
// Parameters:
// - winners: election winners to store (sorted by timestamp asc)
func (e *ElectionStore) StoreElectionWinners(winners []*eth.ElectionWinner) {
	if len(winners) == 0 {
		e.log.Warn("No election winners to store")
		return
	}

	for _, winner := range winners {
		e.electionWinnersMap[winner.Time] = winner
	}

	newLatest := winners[len(winners)-1]
	if e.latestWinner == nil || newLatest.Time > e.latestWinner.Time {
		e.latestWinner = newLatest
	}
}

func (e *ElectionStore) RemoveOutdatedElectionWinners(timestamp uint64) {
	e.log.Debug("Removing outdated election winners", "time", timestamp, "store", e.winnersLength())

	for k := range e.electionWinnersMap {
		if k < timestamp {
			delete(e.electionWinnersMap, k)
		}
	}

	e.log.Debug("Removed outdated election winners", "map", e.winnersLength())
}

func (e *ElectionStore) winnersLength() string {
	return fmt.Sprintf("time: %d", len(e.electionWinnersMap))
}
