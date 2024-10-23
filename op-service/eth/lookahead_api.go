package eth

type Validator struct {
	Pubkey         Bytes48      `json:"pubkey"`
	ValidatorIndex Uint64String `json:"validator_index"`
	Slot           Uint64String `json:"slot"`
}

type APIGetLookaheadResponse struct {
	Data []*Validator `json:"data"`
}
