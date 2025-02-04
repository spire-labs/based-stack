package election

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"

	// is there a better place to put this? making it its own package is difficult because of go modules
	// and us being a private fork

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	SYSTEM_CONFIG_ABI = `[
    {
        "inputs": [],
        "name": "electionFallbackList",
        "outputs": [
            {
                "internalType": "enum ElectionFallback[]",
                "name": "electionFallbackList_",
                "type": "uint8[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    }
	]`
)

type BeaconClient interface {
	GetLookahead(ctx context.Context, epoch uint64) (eth.APIGetLookaheadResponse, error)
	GetEpochNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetSlotNumber(ctx context.Context, timestamp uint64) (uint64, error)
	GetTimeFromSlot(ctx context.Context, slot uint64) (uint64, error)
}

type RpcClient interface {
	Call(ctx context.Context, callMsg map[string]interface{}, blockNumber string) (string, error)
}

type Election struct {
	bc BeaconClient
	l2 RpcClient
	l1 RpcClient

	log log.Logger

	cfg *rollup.Config
}

func NewElection(bc BeaconClient, l2 RpcClient, l1 RpcClient, log log.Logger, cfg *rollup.Config) *Election {
	return &Election{
		bc:  bc,
		l2:  l2,
		l1:  l1,
		log: log,
		cfg: cfg,
	}
}

// l2UnsafeBlock is passed in as a hexadecimal string
func (e *Election) GetWinnersAtEpoch(ctx context.Context, epoch uint64, l2UnsafeBlock string, l2UnsafeParentTime uint64, l1UnsafeBlock string) ([]*eth.ElectionWinner, error) {
	var operatorAddresses []common.Address

	resp, err := e.bc.GetLookahead(ctx, epoch)

	if err != nil {
		return []*eth.ElectionWinner{}, err
	}

	validators := resp.Data

	// TODO(spire): get operation addresses for each validator somehow
	for _, validator := range validators {
		address := common.BytesToAddress(validator.Pubkey[:20])
		operatorAddresses = append(operatorAddresses, address)
	}

	e.log.Info("Checking ticket count per validator at L2 unsafe block", "l2UnsafeBlock", l2UnsafeBlock)
	ticketCountPerValidator, err := e.GetBatchTicketAccounting(ctx, operatorAddresses, l2UnsafeBlock)

	tickets := make(map[common.Address]*big.Int)

	for i, operatorAddress := range operatorAddresses {
		// Already set
		if _, ok := tickets[operatorAddress]; ok {
			continue
		}

		tickets[operatorAddress] = ticketCountPerValidator[i]
	}

	e.log.Info("Ticket count per validator operator", "ticketCountPerValidator", tickets)

	if err != nil {
		log.Error("Failed to get ticket count per validator", "err", err)
		return []*eth.ElectionWinner{}, err
	}

	fallbacklist, err := e.GetElectionFallbackList(ctx, l1UnsafeBlock)
	if err != nil {
		log.Crit("Failed to get fallback list", "err", err)
		return []*eth.ElectionWinner{}, err
	}

	e.log.Info("Fallback list", "fallbacklist", fallbacklist)

	var electionWinners []*eth.ElectionWinner

	// Initialize election winners
	for _, validator := range validators {
		time, err := e.bc.GetTimeFromSlot(ctx, uint64(validator.Slot))
		if err != nil {
			return []*eth.ElectionWinner{}, err
		}

		winner := eth.ElectionWinner{
			// All winners initialized to zero address
			Address: func() common.Address {
				return common.Address{}
			}(),
			Time: time,
		}
		electionWinners = append(electionWinners, &winner)
	}

	return e.HandleInstructions(ctx, fallbacklist, electionWinners, operatorAddresses, tickets, l2UnsafeBlock, l1UnsafeBlock)
}

func (e *Election) GetElectionFallbackList(ctx context.Context, blockNumber string) ([]uint8, error) {
	sysConfig := e.cfg.L1SystemConfigAddress

	parsedSysConfigABI, err := abi.JSON(strings.NewReader(SYSTEM_CONFIG_ABI))
	if err != nil {
		return nil, err
	}

	calldataBytes, err := parsedSysConfigABI.Pack("electionFallbackList")

	if err != nil {
		return nil, err
	}

	calldata := "0x" + hex.EncodeToString(calldataBytes)

	encodedReturnData, err := e.l1.Call(ctx, toCallMsg(sysConfig, calldata), blockNumber)
	if err != nil {
		return nil, err
	}

	encodedReturnData = strings.TrimPrefix(encodedReturnData, "0x") // Remove the prefix
	retdata, err := hex.DecodeString(encodedReturnData)

	if err != nil {
		return nil, err
	}

	res, err := parsedSysConfigABI.Unpack("electionFallbackList", retdata)
	if err != nil {
		return nil, err
	}

	fallbacklist := res[0].([]uint8)

	return fallbacklist, nil
}

// Helper function to format the data into the type the rpc expects
func toCallMsg(to common.Address, data string) map[string]interface{} {
	callMsg := map[string]interface{}{
		"from": common.Address{}.Hex(),
		"to":   to.Hex(),
		"data": data,
	}

	return callMsg
}
