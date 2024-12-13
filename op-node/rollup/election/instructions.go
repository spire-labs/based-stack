package election

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
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

func (e *Election) HandleInstructions(ctx context.Context, instructions []uint8, electionWinners []*eth.ElectionWinner, operatorAddresses []common.Address, tickets map[common.Address]*big.Int, l2UnsafeBlock string) ([]*eth.ElectionWinner, error) {
	var err error

	// Process instructions
	for _, instruction := range instructions {
		switch instruction {
		case NO_FALLBACK:
			// We should never get here, but if we do something is wrong with the system config
			return []*eth.ElectionWinner{}, fmt.Errorf("fallback list contained NO_FALLBACK instruction")
		case CURRENT_PROPOSER:
			electionWinners, err = e.ProcessCurrentProposerInstruction(electionWinners, operatorAddresses, tickets)
			if err != nil {
				return []*eth.ElectionWinner{}, err
			}
			continue
		case CURRENT_PROPOSER_WITH_CONFIG:
			// TODO(spire): This is not implemented yet
			continue
		case NEXT_PROPOSER:
			electionWinners, err = e.ProcessNextProposerInstruction(electionWinners, operatorAddresses, tickets)
			if err != nil {
				return []*eth.ElectionWinner{}, err
			}
			continue
		case NEXT_PROPOSER_WITH_CONFIG:
			// TODO(spire): This is not implemented yet
			continue
		case RANDOM_TICKET_HOLDER:
			electionWinners, err = e.ProcessRandomTicketInstruction(ctx, electionWinners, l2UnsafeBlock)
			if err != nil {
				return []*eth.ElectionWinner{}, err
			}
			continue
		case PERMISSIONLESS:
			// TODO(spire): This is not implemented yet
			continue
		default:
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

func (e *Election) ProcessRandomTicketInstruction(ctx context.Context, electionWinners []*eth.ElectionWinner, blockNumber string) ([]*eth.ElectionWinner, error) {
	var timestamps []uint64
	// Only get the timestamps that still need to be filled
	for _, winner := range electionWinners {
		if winner.Address == (common.Address{}) {
			timestamps = append(timestamps, winner.Time)
		}
	}

	newWinners, err := e.GetBatchRandomTicketInstruction(ctx, timestamps, blockNumber)
	if err != nil {
		return []*eth.ElectionWinner{}, err
	}

	for _, winner := range newWinners {
		// Missed slot
		if winner.Winner == (common.Address{}) {
			continue
		}

		for i, electionWinner := range electionWinners {
			// Search for matching
			if electionWinner.Time == winner.Timestamp {
				electionWinners[i].Address = winner.Winner
				break
			}
		}
	}

	return electionWinners, nil

}
