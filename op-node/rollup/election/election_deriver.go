package election

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/engine"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-node/rollup/status"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionDeriver struct {
	client   BeaconClient
	election *Election
	log      log.Logger
	emitter  event.Emitter
	ctx      context.Context

	nextEpochToUse uint64

	l2Unsafe      eth.L2BlockRef
	l2PendingSafe eth.L2BlockRef
	l1Unsafe      eth.L1BlockRef

	mu sync.Mutex
}

func NewElectionDeriver(ctx context.Context, client BeaconClient, election *Election, log log.Logger) *ElectionDeriver {
	nextEpochToUse := uint64(0)

	return &ElectionDeriver{
		client:         client,
		election:       election,
		log:            log,
		ctx:            ctx,
		nextEpochToUse: nextEpochToUse,
	}
}

func (ed *ElectionDeriver) AttachEmitter(emitter event.Emitter) {
	ed.emitter = emitter
}

func (ed *ElectionDeriver) OnEvent(ev event.Event) bool {
	switch x := ev.(type) {
	// Do we want to do l1unsafe or l1safe here?
	case status.L1UnsafeEvent:
		ed.l1Unsafe = x.L1Unsafe
		ed.ProcessNewL1Block(x.L1Unsafe)
	case engine.PendingSafeUpdateEvent:
		ed.l2PendingSafe = x.PendingSafe
		ed.l2Unsafe = x.Unsafe

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

	// If epoch hasnt changed, do nothing
	if epoch != ed.nextEpochToUse {
		// This is not an error or worth a log so we just return
		return
	}

	// We use unsafe because even if there is a reorg, the time should still be the same
	electionWinners, err := ed.election.GetWinnersAtEpoch(ed.ctx, epoch, "0x"+fmt.Sprintf("%x", ed.l2PendingSafe.Number), ed.l2Unsafe.Time)

	// TODO(spire): Adjust this to handle potential state changes from the previous election
	nextElectionWinners, _ := ed.election.GetWinnersAtEpoch(ed.ctx, epoch+1, "0x"+fmt.Sprintf("%x", ed.l2PendingSafe.Number), electionWinners[len(electionWinners)-1].Time)

	for _, winner := range electionWinners {
		address := &winner.Address
		log.Debug("This winners address is", "address", address)
	}

	if err != nil {
		log.Warn("Failed to get election winner", "err", err)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
	} else {
		ed.mu.Lock()
		defer ed.mu.Unlock()

		log.Info("Election winners", "epoch", epoch, "electionWinners", electionWinners)
		ed.emitter.Emit(rollup.ElectionWinnerEvent{ElectionWinners: electionWinners})
		ed.emitter.Emit(rollup.NextElectionWinnerEvent{ElectionWinners: nextElectionWinners})

		// Update the next epoch to use
		ed.nextEpochToUse = epoch + 1
	}
}
