package election

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-node/rollup/status"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionDeriver struct {
	client   ElectionClient
	election *Election
	log      log.Logger
	emitter  event.Emitter
	ctx      context.Context

	lastEpoch uint64

	mu sync.Mutex
}

func NewElectionDeriver(ctx context.Context, client ElectionClient, election *Election, log log.Logger) *ElectionDeriver {
	lastEpoch := uint64(0)

	return &ElectionDeriver{
		client:    client,
		election:  election,
		log:       log,
		ctx:       ctx,
		lastEpoch: lastEpoch,
	}
}

func (ed *ElectionDeriver) AttachEmitter(emitter event.Emitter) {
	ed.emitter = emitter
}

func (ed *ElectionDeriver) OnEvent(ev event.Event) bool {
	switch x := ev.(type) {
	// Do we want to do l1unsafe or l1safe here?
	case status.L1UnsafeEvent:
		ed.ProcessNewL1Block(x.L1Unsafe)

	default:
		return false
	}

	return true
}

func (ed *ElectionDeriver) ProcessNewL1Block(l1Head eth.L1BlockRef) {
	epoch, err := ed.client.GetEpochNumber(ed.ctx, l1Head.Time)

	if err != nil {
		log.Warn("Failed to get epoch number", "err", err)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
		return
	}

	// We dont need to recalculate the winners as we already did it for this epoch
	// If they are equal and its zero, then its the genesis epoch
	if epoch < ed.lastEpoch || (epoch == ed.lastEpoch && epoch != 0) {
		err := fmt.Errorf("epoch %d is not greater than the last epoch saved which was %d", epoch, ed.lastEpoch)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
		return
	}

	electionWinners, err := ed.election.GetWinnersAtEpoch(ed.ctx, epoch)

	for _, winner := range electionWinners {
		address := &winner.Address
		log.Info("This winners address is", "address", address)
	}

	if err != nil {
		log.Warn("Failed to get election winner", "err", err)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
	} else {
		ed.mu.Lock()
		defer ed.mu.Unlock()

		log.Info("Election winners", "electionWinners", electionWinners)
		ed.emitter.Emit(rollup.ElectionWinnerEvent{ElectionWinners: electionWinners})

		// Only update the last epoch if we got a valid winner
		ed.lastEpoch = epoch
	}
}
