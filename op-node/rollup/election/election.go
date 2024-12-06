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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// This is the representation of the order of precedence, each value represents an instruction for how to process the election winners, based on the system configuration
// NO_FALLBACK is how we know to stop processing the list on-chain, and is an instruction we should never receive here, we check it only for sanity
// TODO(spire): The current unimplemented codes are:
//   - CURRENT_PROPOSER_WITH_CONFIG
//   - NEXT_PROPOSER_WITH_CONFIG
//   - PERMISSIONLESS
const (
	NO_FALLBACK                  = 0x00
	CURRENT_PROPOSER             = 0x01
	CURRENT_PROPOSER_WITH_CONFIG = 0x02
	NEXT_PROPOSER                = 0x03
	NEXT_PROPOSER_WITH_CONFIG    = 0x04
	RANDOM_TICKET_HOLDER         = 0x05
	PERMISSIONLESS               = 0x06
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

// l2PendingSafeBlock is passed in as a hexadecimal string
func (e *Election) GetWinnersAtEpoch(ctx context.Context, epoch uint64, l2PendingSafeBlock string, unsafeParentSlotTime uint64) ([]*eth.ElectionWinner, error) {
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

	e.log.Info("Checking ticket count per validator at L2 safe block", "l2SafeBlock", l2PendingSafeBlock)
	ticketCountPerValidator, err := e.GetBatchTicketAccounting(ctx, operatorAddresses, l2PendingSafeBlock)

	tickets := make(map[common.Address]*big.Int)

	for i, operatorAddress := range operatorAddresses {
		// Already set
		if _, ok := tickets[operatorAddress]; ok {
			continue
		}

		tickets[operatorAddress] = ticketCountPerValidator[i]
	}

	e.log.Info("Operator addresses and tickets", "operatorAddresses", operatorAddresses, "ticketCountPerValidator", tickets)

	if err != nil {
		log.Error("Failed to get ticket count per validator", "err", err)
		return []*eth.ElectionWinner{}, err
	}

	fallbacklist, err := e.GetElectionFallbackList(ctx, l2PendingSafeBlock)
	if err != nil {
		log.Crit("Failed to get fallback list", "err", err)
		return []*eth.ElectionWinner{}, err
	}

	e.log.Info("Fallback list", "fallbacklist", fallbacklist)

	var electionWinners []*eth.ElectionWinner

	// Initialize election winners
	for i, validator := range validators {
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

	// Process instructions
	for _, instruction := range fallbacklist {
		switch instruction {
		case NO_FALLBACK:
			// We should never get here, but if we do something is wrong with the system config
			e.log.Crit("Fallback list contained NO_FALLBACK instruction")
			return []*eth.ElectionWinner{}, fmt.Errorf("fallback list contained NO_FALLBACK instruction")
		case CURRENT_PROPOSER:
			e.log.Info("Processing CURRENT_PROPOSER instruction")
			electionWinners, err = e.ProcessCurrentProposerInstruction(electionWinners, operatorAddresses, tickets)
			if err != nil {
				return []*eth.ElectionWinner{}, err
			}
			e.log.Info("Ticket count per validator", "ticketCountPerValidator", tickets)
			e.log.Info("Election winners after CURRENT_PROPOSER", "electionWinners", electionWinners)
			continue
		case CURRENT_PROPOSER_WITH_CONFIG:
			// TODO(spire): This is not implemented yet
			continue
		case NEXT_PROPOSER:
			e.log.Info("Processing NEXT_PROPOSER instruction")
			electionWinners, err = e.ProcessNextProposerInstruction(electionWinners, operatorAddresses, tickets)
			if err != nil {
				return []*eth.ElectionWinner{}, err
			}
			e.log.Info("Ticket count per validator", "tickets", tickets)
			e.log.Info("Election winners after NEXT_PROPOSER", "electionWinners", electionWinners)
			continue
		case NEXT_PROPOSER_WITH_CONFIG:
			// TODO(spire): This is not implemented yet
			continue
		case RANDOM_TICKET_HOLDER:
			// TODO(spire): This is not implemented yet
			continue
		case PERMISSIONLESS:
			// TODO(spire): This is not implemented yet
			continue
		default:
			e.log.Crit("Unknown fallback instruction", "instruction", instruction)
			return []*eth.ElectionWinner{}, fmt.Errorf("unknown fallback instruction: %d", instruction)
		}
	}

	return electionWinners, nil
}

// If the current proposer is the winner of a slot and holds a ticket he wins the electtion for that slot
func (e *Election) ProcessCurrentProposerInstruction(electionWinners []*eth.ElectionWinner, operatorAddresses []common.Address, tickets map[common.Address]*big.Int) ([]*eth.ElectionWinner, error) {
	if len(electionWinners) != len(operatorAddresses) {
		return []*eth.ElectionWinner{}, fmt.Errorf("invalid input lengths for this instruction")
	}

	addressZero := common.Address{}

	for i, winner := range electionWinners {
		// Keep invariant that we dont replace any slot that already has a winner
		if winner.Address != addressZero {
			continue
		}

		operator := operatorAddresses[i]

		operatorTickets, ok := tickets[operator]

		if !ok {
			return []*eth.ElectionWinner{}, fmt.Errorf("failed to find tickets for operator %s", winner.Address.Hex())
		}

		if operatorTickets.Cmp(big.NewInt(0)) > 0 {
			electionWinners[i].Address = operator
			tickets[operator] = operatorTickets.Sub(operatorTickets, big.NewInt(1))
		}
	}

	return electionWinners, nil
}

// If a proposer is the winner of a slot and does not hold a ticket, the next proposer in the lookahead who does is the winner
func (e *Election) ProcessNextProposerInstruction(electionWinners []*eth.ElectionWinner, operatorAddresses []common.Address, tickets map[common.Address]*big.Int) ([]*eth.ElectionWinner, error) {
	// must be parallel arrays
	if len(electionWinners) != len(operatorAddresses) {
		return []*eth.ElectionWinner{}, fmt.Errorf("invalid input lengths for this instruction")
	}

	addressZero := common.Address{}

	for i, winner := range electionWinners {
		// slot has a winner, skipping
		if winner.Address != addressZero {
			continue
		}

		operator := operatorAddresses[i]
		operatorTickets, ok := tickets[operator]

		if !ok {
			return []*eth.ElectionWinner{}, fmt.Errorf("failed to find tickets for operator %s", winner.Address.Hex())
		}

		// If this proposer does hold a ticket, it should be handled in a different instruction to add it to the electionWinners
		if operatorTickets.Cmp(big.NewInt(0)) == 0 {
			// proposer does not hold a ticket, finding the next proposer in the lookahead who does
			for j := i + 1; j < len(electionWinners); j++ {

				// Get the new values for the operator
				operator := operatorAddresses[j]
				operatorTickets, ok := tickets[operator]

				if !ok {
					return []*eth.ElectionWinner{}, fmt.Errorf("failed to find tickets for operator %s", winner.Address.Hex())
				}

				// When the next proposer has a ticket, he is the winner and we can finish searching
				if operatorTickets.Cmp(big.NewInt(0)) > 0 {
					electionWinners[i].Address = operatorAddresses[j]
					tickets[operator] = operatorTickets.Sub(operatorTickets, big.NewInt(1))
					break
				}
			}
		}
	}

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

// / Helper function to format the data into the type the rpc expects
func toBatchCallMsg(from common.Address, batchCreationCode string) map[string]interface{} {
	callMsg := map[string]interface{}{
		"from": from.Hex(),
		"to":   nil,
		"data": batchCreationCode,
	}

	return callMsg
}

func toCallMsg(to common.Address, data string) map[string]interface{} {
	callMsg := map[string]interface{}{
		"from": common.Address{}.Hex(),
		"to":   to.Hex(),
		"data": data,
	}

	return callMsg
}
