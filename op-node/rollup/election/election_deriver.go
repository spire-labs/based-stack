package election

import (
	"github.com/ethereum-optimism/optimism/op-node/rollup/event"
	"github.com/ethereum-optimism/optimism/op-service/sources"
)

type ElectionDeriver struct {
	client   sources.BeaconClient
	election *Election
	emitter  event.Emitter
}

func NewElectionDeriver(client sources.BeaconClient, election *Election) *ElectionDeriver {
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
