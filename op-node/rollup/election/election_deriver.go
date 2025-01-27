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

	l2Finalized eth.L2BlockRef
	l2Unsafe    eth.L2BlockRef
	l1Unsafe    eth.L1BlockRef

	electionWinners []ElectionWinners

	// The timestamp of last slot in current epoch
	lastSlotTime uint64

	mu sync.Mutex
}

// TODO(spire): add this to network config or remove entirely.
const L1BlockTime = 12

func NewElectionDeriver(ctx context.Context, client BeaconClient, election *Election, log log.Logger) *ElectionDeriver {

	return &ElectionDeriver{
		client:   client,
		election: election,
		log:      log,
		ctx:      ctx,
	}
}

func (ed *ElectionDeriver) AttachEmitter(emitter event.Emitter) {
	ed.emitter = emitter
}

func (ed *ElectionDeriver) OnEvent(ev event.Event) bool {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	switch x := ev.(type) {
	case status.L1UnsafeEvent:
		ed.l1Unsafe = x.L1Unsafe
		ed.ProcessNewBlock()
	case engine.PendingSafeUpdateEvent:
		ed.l2Unsafe = x.Unsafe
		ed.ProcessNewBlock()
	case engine.ForkchoiceUpdateEvent:
		// optimization: only emit event if the finalized L2 head has changed
		if x.FinalizedL2Head.Number > ed.l2Finalized.Number {
			ed.l2Finalized = x.FinalizedL2Head
			ed.emitter.Emit(rollup.ElectionWinnerOutdatedEvent{Time: x.FinalizedL2Head.Time})
		}
	default:
		return false
	}

	return true
}

func (ed *ElectionDeriver) ProcessNewBlock() {
	if ed.l2Unsafe == (eth.L2BlockRef{}) {
		ed.log.Debug("Empty L2 block")
		return
	}

	if ed.l1Unsafe == (eth.L1BlockRef{}) {
		ed.log.Debug("Empty L1 block")
		return
	}

	// If epoch hasn't changed do nothing
	if ed.l1Unsafe.Time < ed.lastSlotTime {
		ed.log.Debug("Waiting for the last slot of the epoch", "current_time", ed.l1Unsafe.Time, "lastSlotTime", ed.lastSlotTime)
		return
	}

	if ed.lastSlotTime != 0 && ed.l2Unsafe.Time < ed.lastSlotTime {
		ed.log.Debug("L2 slot times mismatch", "l2Unsafe", ed.l2Unsafe, "time", ed.l2Unsafe.Time, "l1Unsafe", ed.l1Unsafe, "time", ed.l1Unsafe.Time, "lastSlotTime", ed.lastSlotTime)
		return
	}

	// Sanity check
	if ed.lastSlotTime != 0 && ed.l2Unsafe.Time != ed.lastSlotTime {
		ed.log.Error("l2Unsafe slot timestamp mismatch", "l2Unsafe", ed.l2Unsafe, "time", ed.l2Unsafe.Time, "l1Unsafe", ed.l1Unsafe, "time", ed.l1Unsafe.Time)
		return
	}

	ed.log.Debug("Processing", "l1Unsafe", ed.l1Unsafe.Number, "time", ed.l1Unsafe.Time, "l2Unsafe", ed.l2Unsafe.Number, "time", ed.l2Unsafe.Time, "lastSlotTime", ed.lastSlotTime)

	lastBlockNumberInEpoch := ed.l1Unsafe.Number
	nextEpochTime := ed.l1Unsafe.Time + L1BlockTime
	if ed.lastSlotTime != 0 && ed.l1Unsafe.Time != ed.lastSlotTime {
		ed.log.Warn("Slot times mismatch, L1 missed slot detected", "ed.l1Unsafe", ed.l1Unsafe, "time", ed.l1Unsafe.Time, "lastSlotTime", ed.lastSlotTime)
		lastBlockNumberInEpoch = ed.l1Unsafe.Number - 1
		nextEpochTime = ed.l1Unsafe.Time
	}

	newEpoch, err := ed.client.GetEpochNumber(ed.ctx, nextEpochTime)
	if err != nil {
		ed.log.Warn("Failed to get epoch number", "err", err)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
		return
	}

	electionWinners, err := ed.election.GetWinnersAtEpoch(ed.ctx, newEpoch, fmt.Sprintf("0x%x", ed.l2Unsafe.Number), ed.l2Unsafe.Time, fmt.Sprintf("0x%x", lastBlockNumberInEpoch))

	if err != nil {
		ed.log.Error("Failed to get election winner", "epoch", newEpoch, "err", err)
		ed.emitter.Emit(rollup.ElectionErrorEvent{Err: err})
	} else {
		ed.log.Info("Election winners", "epoch", newEpoch, "electionWinners", electionWinners)
		ed.emitter.Emit(rollup.ElectionWinnerEvent{ElectionWinners: electionWinners})

		ed.electionWinners = append(ed.electionWinners, ElectionWinners{winners: electionWinners, epoch: newEpoch})

		// Clear old election winners
		start := 0
		for i, electionWinners := range ed.electionWinners {
			if electionWinners.epoch < newEpoch {
				start = i + 1
			}
		}
		if start > 0 {
			ed.electionWinners = ed.electionWinners[start:]
		}

		// Update last slot in this epoch
		ed.lastSlotTime = electionWinners[len(electionWinners)-1].Time
	}
}

func (ed *ElectionDeriver) GetElectionWinners(ctx context.Context, epoch uint64) ([]eth.ElectionWinner, error) {
	ed.mu.Lock()
	defer ed.mu.Unlock()

	var winners []*eth.ElectionWinner
	storedEpochs := make([]uint64, len(ed.electionWinners))
	for i, stored := range ed.electionWinners {
		storedEpochs[i] = stored.epoch
		if stored.epoch == epoch {
			winners = stored.winners
		}
	}

	if winners == nil {
		return []eth.ElectionWinner{}, fmt.Errorf("no stored election winners for requested epoch, requested epoch: %d, stored_epochs: %v", epoch, storedEpochs)
	}

	out := make([]eth.ElectionWinner, len(winners))
	for i, winner := range winners {
		out[i] = *winner
	}

	return out, nil
}
