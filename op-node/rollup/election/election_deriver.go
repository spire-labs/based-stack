package election

import (
	"context"
	"sync"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-node/rollup/status"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionDeriver struct {
	client   ElectionClient
	election *Election
	log      log.Logger
	emitter  event.Emitter

	mu sync.Mutex
}

func NewElectionDeriver(client ElectionClient, election *Election, log log.Logger) *ElectionDeriver {
	return &ElectionDeriver{
		client:   client,
		election: election,
		log:      log,
	}
}

func (ed *ElectionDeriver) AttachEmitter(emitter event.Emitter) {
	ed.emitter = emitter
}

func (ed *ElectionDeriver) OnEvent(ev event.Event) bool {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	switch x := ev.(type) {
	// Do we want to do l1unsafe or l1safe here?
	case status.L1UnsafeEvent:
		// how do we get the correct ctx variable here?
		ctx := context.Background()
		validator, err := ed.election.GetWinner(ctx, x.L1Unsafe.Time, ed.log)
		if err != nil {
			log.Warn("Failed to get election winner", "err", err)
			ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
		} else {
			log.Info("Election winner", "validator", validator)
			ed.emitter.Emit(rollup.ElectionWinnerEvent{Validator: validator})
		}
	default:
		return false
	}

	return true
}
