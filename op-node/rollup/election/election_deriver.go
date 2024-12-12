package election

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/engine"
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-node/rollup/status"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/log"
)

type ElectionWinners struct {
	winners []*eth.ElectionWinner
	epoch   uint64
}

type ElectionDeriver struct {
	client   BeaconClient
	election *Election
	log      log.Logger
	emitter  event.Emitter
	ctx      context.Context

	nextEpochToUse uint64

	l2Unsafe eth.L2BlockRef
	l1Unsafe eth.L1BlockRef

	electionWinners []ElectionWinners

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
	case status.L1UnsafeEvent:
		ed.l1Unsafe = x.L1Unsafe
		log.Info("L1 Unsafe is at block", "block", x.L1Unsafe.Number, "time", x.L1Unsafe.Time)
		ed.ProcessNewL1Block(x.L1Unsafe)
	case engine.PendingSafeUpdateEvent:
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

	log.Info("From Deriver:", "l2Block", ed.l2Unsafe.Number)
	log.Info("From Deriver:", "l2ParentTimestamp", ed.l2Unsafe.Time)
	log.Info("From Deriver:", "l1Block", ed.l1Unsafe.Number)
	// We use unsafe because even if there is a reorg, the time should still be the same
	electionWinners, err := ed.election.GetWinnersAtEpoch(ed.ctx, epoch, fmt.Sprintf("0x%x", ed.l2Unsafe.Number), ed.l2Unsafe.Time, fmt.Sprintf("0x%x", ed.l1Unsafe.Number))

	if err != nil {
		log.Error("Failed to get election winner", "err", err)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
	} else {
		ed.mu.Lock()
		defer ed.mu.Unlock()

		log.Info("Election winners", "epoch", epoch, "electionWinners", electionWinners)
		ed.emitter.Emit(rollup.ElectionWinnerEvent{ElectionWinners: electionWinners})

		ed.electionWinners = append(ed.electionWinners, ElectionWinners{winners: electionWinners, epoch: epoch})

		// Clear old election winners
		start := 0
		for i, electionWinners := range ed.electionWinners {
			if electionWinners.epoch < epoch {
				start = i + 1
			}
		}
		if start > 0 {
			ed.electionWinners = ed.electionWinners[start:]
		}

		// Update the next epoch to use
		ed.nextEpochToUse = epoch + 1
	}
}

func (ed *ElectionDeriver) GetElectionWinners(ctx context.Context, epoch uint64) ([]eth.ElectionWinner, error) {
	var winners []*eth.ElectionWinner
	for _, stored := range ed.electionWinners {
		if stored.epoch == epoch {
			winners = stored.winners
		}
	}

	if winners == nil {
		return []eth.ElectionWinner{}, errors.New("invalid epoch")
	}

	out := make([]eth.ElectionWinner, len(winners))
	for i, winner := range winners {
		out[i] = *winner
	}

	return out, nil
}
