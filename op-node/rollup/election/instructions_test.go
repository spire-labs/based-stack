package election

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

type MockElection struct {
	*Election
	GetBatchRandomTicketInstructionFunc func(ctx context.Context, timestamps []uint64, blockNumber string) ([]RandomTicketInstructionRetdata, error)
}

// Override the function to mock the batch
func (m *MockElection) GetBatchRandomTicketInstruction(ctx context.Context, timestamps []uint64, blockNumber string) ([]RandomTicketInstructionRetdata, error) {

	return m.GetBatchRandomTicketInstructionFunc(ctx, timestamps, blockNumber)
}

func createMockInputsForHandleInstructions() (context.Context, []*eth.ElectionWinner, []common.Address, map[common.Address]*big.Int, string) {
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
		operatorAddresses[1]: big.NewInt(3),
	}
	blockNumber := "0x1"
	ctx := context.Background()
	return ctx, electionWinners, operatorAddresses, tickets, blockNumber
}

func TestSingleCurrentProposerInstruction(t *testing.T) {
	ctx, electionWinners, operatorAddresses, tickets, blockNumber := createMockInputsForHandleInstructions()
	e := &Election{}

	instructions := []uint8{CURRENT_PROPOSER}
	result, err := e.HandleInstructions(ctx, instructions, electionWinners, operatorAddresses, tickets, blockNumber, blockNumber)

	assert.NoError(t, err)
	assert.Equal(t, electionWinners, result)
}

func TestSingleNextProposerInstruction(t *testing.T) {
	ctx, electionWinners, operatorAddresses, tickets, blockNumber := createMockInputsForHandleInstructions()
	e := &Election{}

	instructions := []uint8{NEXT_PROPOSER}
	result, err := e.HandleInstructions(ctx, instructions, electionWinners, operatorAddresses, tickets, blockNumber, blockNumber)

	assert.NoError(t, err)
	assert.Equal(t, electionWinners, result)
}

func TestMultipleInstructions(t *testing.T) {
	ctx, electionWinners, operatorAddresses, tickets, blockNumber := createMockInputsForHandleInstructions()
	e := &Election{}

	instructions := []uint8{CURRENT_PROPOSER, NEXT_PROPOSER}
	result, err := e.HandleInstructions(ctx, instructions, electionWinners, operatorAddresses, tickets, blockNumber, blockNumber)

	assert.NoError(t, err)
	assert.Equal(t, electionWinners, result)
}

func TestInvalidInstruction(t *testing.T) {
	ctx, electionWinners, operatorAddresses, tickets, blockNumber := createMockInputsForHandleInstructions()
	e := &Election{}

	instructions := []uint8{255} // Unknown instruction
	result, err := e.HandleInstructions(ctx, instructions, electionWinners, operatorAddresses, tickets, blockNumber, blockNumber)

	assert.Error(t, err)
	assert.EqualError(t, err, "unknown fallback instruction: 255")
	assert.Empty(t, result, "Result should be empty on error")
}

func TestNoInstructions(t *testing.T) {
	ctx, electionWinners, operatorAddresses, tickets, blockNumber := createMockInputsForHandleInstructions()
	e := &Election{}

	instructions := []uint8{}
	result, err := e.HandleInstructions(ctx, instructions, electionWinners, operatorAddresses, tickets, blockNumber, blockNumber)

	assert.NoError(t, err)
	assert.Equal(t, electionWinners, result)
}

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
