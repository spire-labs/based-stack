package election

import (
	"github.com/ethereum-optimism/optimism/op-service/sources"
)

type Election struct {
	client sources.BeaconClient
}

func NewElection(client sources.BeaconClient) *Election {
	return &Election{
		client: client,
	}
}
