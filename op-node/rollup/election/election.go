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

// This is the representation of the order of precedence, each value represents a "opcode" for how to process the election winners, based on the system configuration
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

	e.log.Info("Checking ticket count per validator at L2 safe block", "l2SafeBlock", l2PendingSafeBlock)
	ticketCountPerValidator, err := e.GetBatchTicketAccounting(ctx, lookaheadAddresses, l2PendingSafeBlock)

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
			electionWinners, err = e.ProcessCurrentProposerInstruction(electionWinners, lookaheadAddresses, ticketCountPerValidator)
			if err != nil {
				return []*eth.ElectionWinner{}, err
			}
			e.log.Info("Ticket count per validator", "ticketCountPerValidator", ticketCountPerValidator)
			continue
		case CURRENT_PROPOSER_WITH_CONFIG:
			// TODO(spire): This is not implemented yet
			continue
		case NEXT_PROPOSER:
			// TODO(spire): This is not implemented yet
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

// All 3 arrays are parallel arrays
func (e *Election) ProcessCurrentProposerInstruction(electionWinners []*eth.ElectionWinner, validatorOperatorAddresses []common.Address, validatorTickets []*big.Int) ([]*eth.ElectionWinner, error) {
	if len(electionWinners) != len(validatorOperatorAddresses) || len(electionWinners) != len(validatorTickets) {
		return []*eth.ElectionWinner{}, fmt.Errorf("invalid input lengths for this instruction")
	}

	addressZero := common.Address{}

	for i, winner := range electionWinners {
		// Keep invariant that we dont replace any slot that already has a winner
		if winner.Address != addressZero {
			continue
		}

		if validatorTickets[i].Cmp(big.NewInt(0)) > 0 {
			winner.Address = validatorOperatorAddresses[i]
			validatorTickets[i] = validatorTickets[i].Sub(validatorTickets[i], big.NewInt(1))
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
