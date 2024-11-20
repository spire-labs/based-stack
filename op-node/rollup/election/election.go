package election

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	// is there a better place to put this? making it its own package is difficult because of go modules
	// and us being a private fork
	BatchTicketAccounting "github.com/ethereum-optimism/optimism/op-node/batch-contracts/bindings"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type BeaconClient interface {
	GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error)
	GetEpochNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetSlotNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetTimeFromSlot(ctx context.Context, slot uint64) (uint64, error)
}

type Election struct {
	bc BeaconClient
	l2 *sources.EthClient

	log log.Logger

	cfg *rollup.Config
}

func NewElection(bc BeaconClient, l2 *sources.EthClient, log log.Logger, cfg *rollup.Config) *Election {
	return &Election{
		bc:  bc,
		l2:  l2,
		log: log,
		cfg: cfg,
	}
}

// l2PendingSafeBlock is passed in as a hexadecimal string
func (e *Election) GetWinnersAtEpoch(ctx context.Context, epoch uint64, l2PendingSafeBlock string, unsafeParentSlotTime uint64) ([]*eth.ElectionWinner, error) {
	var electionWinners []*eth.ElectionWinner
	var lookaheadAddresses []common.Address

	resp, err := e.bc.GetLookahead(ctx, epoch)

	if err != nil {
		return []*eth.ElectionWinner{}, err
	}

	validators := resp.Data

	for _, validator := range validators {
		address := common.BytesToAddress(validator.Pubkey[:20])
		lookaheadAddresses = append(lookaheadAddresses, address)
	}

	ticketCountPerValidator, err := e.GetBatchTicketAccounting(ctx, lookaheadAddresses, l2PendingSafeBlock)

	if err != nil {
		log.Error("Failed to get ticket count per validator", "err", err)
		return []*eth.ElectionWinner{}, err
	}

	for i, validator := range validators {
		address := common.BytesToAddress(validator.Pubkey[:20])
		time, err := e.bc.GetTimeFromSlot(ctx, uint64(validator.Slot))
		if err != nil {
			return []*eth.ElectionWinner{}, err
		}

		// At this stage if the ticket count is 0 then the validator is not a winner
		// We mark this as a nil address for now
		winner := eth.ElectionWinner{
			Address: func() common.Address {
				if ticketCountPerValidator[i].Cmp(big.NewInt(0)) > 0 {
					return address
				}
				return common.Address{}
			}(),
			Time: time,
			ParentSlot: func() uint64 {
				// TODO(spire): This breaks if L2 block time is not == L1 block time
				// Should also not be hardcoded but gotten from the rollup config
				if time == unsafeParentSlotTime+e.cfg.BlockTime {
					return unsafeParentSlotTime
				}

				// Sanity check, if first check fails at first index we need to return 0
				// To avoid out of bounds error
				if i == 0 {
					return 0
				}

				parent := electionWinners[i-1].ParentSlot

				// Keep chaining zeroes, we dont want to increment
				if parent == 0 {
					return 0
				}

				return parent + e.cfg.BlockTime
			}(),
		}
		electionWinners = append(electionWinners, &winner)
	}

	// For now winner will always return the current slots proposer, so we can just return the validators as is as each slot is its winner
	return electionWinners, nil
}

func (e *Election) GetBatchTicketAccounting(ctx context.Context, lookaheadAddresses []common.Address, blockNumber string) ([]*big.Int, error) {
	bin := BatchTicketAccounting.BatchTicketAccountingMetaData.Bin
	abiJson := BatchTicketAccounting.BatchTicketAccountingMetaData.ABI

	// Parse the ABI
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(abiJson)))

	if err != nil {
		// If we fail to parse the ABI, we should return an error
		// Because we cant continue without it
		return []*big.Int{}, err
	}

	constructorArgs, err := parsedABI.Pack("", lookaheadAddresses)

	if err != nil {
		e.log.Error("Failed to pack constructor args", "err", err)
		return []*big.Int{}, err
	}

	creationCode := "0x" + hex.EncodeToString(append(common.FromHex(bin), constructorArgs...))

	// NOTE: Should we be using latest here?
	encodedReturnData, err := e.l2.Call(ctx, toBatchCallMsg(common.Address{}, creationCode), blockNumber)

	if err != nil {
		// If we cant determine ticket accounting, we cant determine the winners
		return []*big.Int{}, err
	}

	encodedReturnData = strings.TrimPrefix(encodedReturnData, "0x")
	encodedReturnDataAsBytes, err := hex.DecodeString(encodedReturnData)
	if err != nil {
		return []*big.Int{}, err
	}

	// Hardcoded type as the abi is not aware of our custom constructor return type
	uint256ArrayType, err := abi.NewType("uint256[]", "", nil)

	if err != nil {
		return []*big.Int{}, err
	}

	args := abi.Arguments{
		{Type: uint256ArrayType},
	}

	// Parallel array to the lookahead for each validators ticket count
	decoded, err := args.Unpack(encodedReturnDataAsBytes)

	if err != nil {
		return []*big.Int{}, err
	}

	ticketCountPerValidator, ok := decoded[0].([]*big.Int)

	if !ok {
		err = fmt.Errorf("failed to convert raw data to []*big.Int")
		return []*big.Int{}, err
	}

	e.log.Info("BatchTicketAccounting contract called successfully, results were", "ticketCountPerValidator", ticketCountPerValidator)

	return ticketCountPerValidator, nil
}

// / Helper function to format the data into the type the rpc expects
func toBatchCallMsg(from common.Address, batchCreationCode string) map[string]interface{} {
	callMsg := map[string]interface{}{
		"from": from.Hex(),
		"to":   nil,
		"data": batchCreationCode,
	}

	return callMsg
}
