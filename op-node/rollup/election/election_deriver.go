package election

import (
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
)

type ElectionDeriver struct {
	client   ElectionClient
	election *Election
	emitter  event.Emitter
}

func NewElectionDeriver(client ElectionClient, election *Election) *ElectionDeriver {
	return &ElectionDeriver{
		client:   client,
		election: election,
	}
}

func (ed *ElectionDeriver) AttachEmitter(emitter event.Emitter) {
	ed.emitter = emitter
}

func (ed *ElectionDeriver) OnEvent(ev event.Event) bool {
	// Unimplemented, this is a placeholder
	return false
}
