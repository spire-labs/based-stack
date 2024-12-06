package election

import (
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestProcessCurrentProposerInstruction(t *testing.T) {
	// Mock input data
	electionWinners := []*eth.ElectionWinner{
		{Address: common.Address{}},
		{Address: common.Address{}},
	}
	operatorAddresses := []common.Address{
		common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
	}
	tickets := map[common.Address]*big.Int{
		operatorAddresses[0]: big.NewInt(2),
		operatorAddresses[1]: big.NewInt(0),
	}

	e := &Election{}
	updatedWinners, err := e.ProcessCurrentProposerInstruction(electionWinners, operatorAddresses, tickets)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, operatorAddresses[0], updatedWinners[0].Address)
	assert.Equal(t, common.Address{}, updatedWinners[1].Address)
	assert.Equal(t, big.NewInt(1), tickets[operatorAddresses[0]])
	assert.Equal(t, big.NewInt(0), tickets[operatorAddresses[1]])
}

func TestProcessNextProposerInstruction(t *testing.T) {
	// Mock input data
	electionWinners := []*eth.ElectionWinner{
		{Address: common.Address{}},
		{Address: common.Address{}},
		{Address: common.Address{}},
	}
	operatorAddresses := []common.Address{
		common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		common.HexToAddress("0xabcdef1234567890abcdef1234567890abcdef12"),
		common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
	}
	tickets := map[common.Address]*big.Int{
		operatorAddresses[0]: big.NewInt(0),
		operatorAddresses[1]: big.NewInt(0),
		operatorAddresses[2]: big.NewInt(3),
	}

	e := &Election{}
	updatedWinners, err := e.ProcessNextProposerInstruction(electionWinners, operatorAddresses, tickets)

	// Assertions
	assert.NoError(t, err)
	// Next proposer instruction should not win its own slot
	assert.Equal(t, updatedWinners[0].Address, operatorAddresses[2])
	assert.Equal(t, updatedWinners[1].Address, operatorAddresses[2])
	assert.Equal(t, updatedWinners[2].Address, common.Address{})
	// Should decrement two for two slots won
	assert.Equal(t, big.NewInt(1), tickets[operatorAddresses[2]])
}
