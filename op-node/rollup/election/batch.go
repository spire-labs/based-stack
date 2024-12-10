package election

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	BatchTicketAccounting "github.com/ethereum-optimism/optimism/op-node/batch-contracts/bindings/BatchTicketAccounting"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

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

// Helper function to format the data into the type the rpc expects
func toBatchCallMsg(from common.Address, batchCreationCode string) map[string]interface{} {
	callMsg := map[string]interface{}{
		"from": from.Hex(),
		"to":   nil,
		"data": batchCreationCode,
	}

	return callMsg
}
