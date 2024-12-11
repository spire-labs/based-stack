package sources

import (
	"context"
	"math/rand"
	"net/url"
	"path"
	"strconv"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
)

type BeaconHTTPClientFakeLookahead struct {
	BeaconHTTPClient
	validatorPubkeys []eth.Bytes48
}

func NewBeaconHTTPClientFakeLookahead(cl client.HTTP, validatorAddresses []common.Address) *BeaconHTTPClientFakeLookahead {

	pubkeys := []eth.Bytes48{}

	for _, validator := range validatorAddresses {
		pubkeys = append(pubkeys, fakePubkeyFromAddress(validator))
	}

	return &BeaconHTTPClientFakeLookahead{BeaconHTTPClient: *NewBeaconHTTPClient(cl), validatorPubkeys: pubkeys}
}

func fakePubkeyFromAddress(address common.Address) eth.Bytes48 {
	var pubkey eth.Bytes48
	copy(pubkey[:], address.Bytes()[:])
	return pubkey
}

func (bc *BeaconHTTPClientFakeLookahead) GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error) {
	reqPath := path.Join(lookaheadMethodPrefix, strconv.FormatUint(epoch, 10))

	var reqQuery url.Values
	var resp eth.APIGetLookaheadResponse

	if err := bc.apiReq(ctx, &resp, reqPath, reqQuery); err != nil {
		return eth.APIGetLookaheadResponse{}, err
	}
	rng := rand.New(rand.NewSource(int64(epoch)))
	// Change the addresses randomly
	for _, validator := range resp.Data {
		validator.Pubkey = bc.validatorPubkeys[rng.Intn(len(bc.validatorPubkeys))]
	}

	return resp, nil
}
