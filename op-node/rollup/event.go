package rollup

import (
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// L1TemporaryErrorEvent identifies a temporary issue with the L1 data.
type L1TemporaryErrorEvent struct {
	Err error
}

var _ event.Event = L1TemporaryErrorEvent{}

func (ev L1TemporaryErrorEvent) String() string {
	return "l1-temporary-error"
}

// EngineTemporaryErrorEvent identifies a temporary processing issue.
// It applies to both L1 and L2 data, often inter-related.
// This scope will be reduced over time, to only capture L2-engine specific temporary errors.
// See L1TemporaryErrorEvent for L1 related temporary errors.
type EngineTemporaryErrorEvent struct {
	Err error
}

var _ event.Event = EngineTemporaryErrorEvent{}

func (ev EngineTemporaryErrorEvent) String() string {
	return "engine-temporary-error"
}

type ResetEvent struct {
	Err error
}

var _ event.Event = ResetEvent{}

func (ev ResetEvent) String() string {
	return "reset-event"
}

// CriticalErrorEvent is an alias for event.CriticalErrorEvent
type CriticalErrorEvent = event.CriticalErrorEvent

type LivenessErrorEvent struct {
	Err error
}

func (ev LivenessErrorEvent) String() string {
	return "liveness-error"
}

// NOTE: We cant store this in election because derive needs to import it
// there is a circular dependency
// is there a better place to put this / better way to do this?
type ElectionWinnerEvent struct {
	ElectionWinners []*eth.ElectionWinner
}

func (ev ElectionWinnerEvent) String() string {
	return "election-winner"
}

type ElectionWinnerOutdatedEvent struct {
	Time uint64
}

func (ev ElectionWinnerOutdatedEvent) String() string {
	return "election-winner-outdated"
}

type ElectionErrorEvent struct {
	Err error
}

func (ev ElectionErrorEvent) String() string {
	return "election-error"
}
