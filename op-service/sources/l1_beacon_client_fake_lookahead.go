package sources

import (
	"context"
	"encoding/hex"
	"math/rand"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type BeaconHTTPClientFakeLookahead struct {
	BeaconHTTPClient
	validatorPubkeys []eth.Bytes48
}

func NewBeaconHTTPClientFakeLookahead(cl client.HTTP) *BeaconHTTPClientFakeLookahead {
	pubkeys := []eth.Bytes48{fakePubkeyFromAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"),
		fakePubkeyFromAddress("0x976EA74026E726554dB657fA54763abd0C3a0aa9")}

	return &BeaconHTTPClientFakeLookahead{BeaconHTTPClient: *NewBeaconHTTPClient(cl), validatorPubkeys: pubkeys}
}

func fakePubkeyFromAddress(addressStr string) eth.Bytes48 {
	address, _ := hex.DecodeString(strings.TrimPrefix(addressStr, "0x"))
	var pubkey eth.Bytes48
	copy(pubkey[:], address[:])
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
