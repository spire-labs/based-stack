package election

import (
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
)

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
