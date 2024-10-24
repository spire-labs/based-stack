package eth

import "github.com/ethereum/go-ethereum/common"

type Validator struct {
	Pubkey         Bytes48      `json:"pubkey"`
	ValidatorIndex Uint64String `json:"validator_index"`
	Slot           Uint64String `json:"slot"`
}

type ElectionWinner struct {
	Address common.Address
	Time    uint64
}

type APIGetLookaheadResponse struct {
	Data []*Validator `json:"data"`
}
