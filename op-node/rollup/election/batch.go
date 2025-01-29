package election

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum-optimism/optimism/op-node/batch-contracts/bindings/BatchCheckSeqConfig"
	"github.com/ethereum-optimism/optimism/op-node/batch-contracts/bindings/BatchRandomTicketInstruction"
	BatchTicketAccounting "github.com/ethereum-optimism/optimism/op-node/batch-contracts/bindings/BatchTicketAccounting"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type RandomTicketInstructionRetdata struct {
	Timestamp uint64
	Winner    common.Address
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

func (e *Election) GetBatchRandomTicketInstruction(ctx context.Context, timestamps []uint64, blockNumber string) ([]RandomTicketInstructionRetdata, error) {
	bin := BatchRandomTicketInstruction.BatchRandomTicketInstructionMetaData.Bin
	abiJson := BatchRandomTicketInstruction.BatchRandomTicketInstructionMetaData.ABI

	// Parse the ABI
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(abiJson)))

	if err != nil {
		// If we fail to parse the ABI, we should return an error
		// Because we cant continue without it
		return []RandomTicketInstructionRetdata{}, err
	}

	timestampsAsBigInt := convertUint64ArrayToBigIntArray(timestamps)

	constructorArgs, err := parsedABI.Pack("", timestampsAsBigInt)

	if err != nil {
		return []RandomTicketInstructionRetdata{}, err
	}

	creationCode := "0x" + hex.EncodeToString(append(common.FromHex(bin), constructorArgs...))

	encodedReturnData, err := e.l2.Call(ctx, toBatchCallMsg(common.Address{}, creationCode), blockNumber)

	if err != nil {
		return []RandomTicketInstructionRetdata{}, err
	}

	encodedReturnData = strings.TrimPrefix(encodedReturnData, "0x")
	encodedReturnDataAsBytes, err := hex.DecodeString(encodedReturnData)

	if err != nil {
		return []RandomTicketInstructionRetdata{}, err
	}

	components := []abi.ArgumentMarshaling{
		{Name: "timestamp", Type: "uint256"}, // Corresponds to uint256
		{Name: "winner", Type: "address"},    // Corresponds to address
	}

	// Define the tuple array type
	retdataArrayType, err := abi.NewType("tuple[]", "", components)

	if err != nil {
		return []RandomTicketInstructionRetdata{}, err
	}

	args := abi.Arguments{
		{Type: retdataArrayType},
	}

	decoded, err := args.Unpack(encodedReturnDataAsBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack return data: %w", err)
	}

	// Doing the other decoding method was not working for some reason, so we use reflection to process decoded[0]
	// Can look more into this later if it becomes a problem
	decodedValue := reflect.ValueOf(decoded[0])
	if decodedValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("unexpected data type for decoded[0]: %T", decoded[0])
	}

	var results []RandomTicketInstructionRetdata
	for i := 0; i < decodedValue.Len(); i++ {
		// Access individual element
		element := decodedValue.Index(i)

		// Extract fields using reflection
		timestampField := element.FieldByName("Timestamp")
		winnerField := element.FieldByName("Winner")

		if !timestampField.IsValid() || !winnerField.IsValid() {
			return nil, fmt.Errorf("missing fields in tuple at index %d", i)
		}

		// Convert fields to their expected types
		timestamp, ok := timestampField.Interface().(*big.Int)
		if !ok {
			return nil, fmt.Errorf("invalid type for Timestamp at index %d: %T", i, timestampField.Interface())
		}

		winner, ok := winnerField.Interface().(common.Address)
		if !ok {
			return nil, fmt.Errorf("invalid type for Winner at index %d: %T", i, winnerField.Interface())
		}

		results = append(results, RandomTicketInstructionRetdata{
			Timestamp: timestamp.Uint64(),
			Winner:    winner,
		})
	}

	return results, nil

}

func (e *Election) GetBatchCheckSeqConfig(ctx context.Context, potentialWinners []common.Address, blockNumber string) ([]bool, error) {
	bin := BatchCheckSeqConfig.BatchCheckSeqConfigMetaData.Bin
	abiJson := BatchCheckSeqConfig.BatchCheckSeqConfigMetaData.ABI

	// Parse the ABI
	parsedABI, err := abi.JSON(bytes.NewReader([]byte(abiJson)))

	if err != nil {
		// If we fail to parse the ABI, we should return an error
		// Because we cant continue without it
		return []bool{}, err
	}

	constructorArgs, err := parsedABI.Pack("", e.cfg.L1SystemConfigAddress, potentialWinners)

	e.log.Info("Potential winners", "winners", potentialWinners)
	e.log.Info("Sys config address", "address", e.cfg.L1SystemConfigAddress)
	e.log.Info("Address zero", "address", common.Address{})

	if err != nil {
		return []bool{}, err
	}

	creationCode := "0x" + hex.EncodeToString(append(common.FromHex(bin), constructorArgs...))

	// NOTE: Should we be using latest here?
	encodedReturnData, err := e.l1.Call(ctx, toBatchCallMsg(common.Address{}, creationCode), blockNumber)

	if err != nil {
		e.log.Error("Failed to get config results", "err", err)
		// If we cant determien the config results, we cant determine the winners
		return []bool{}, err
	}

	encodedReturnData = strings.TrimPrefix(encodedReturnData, "0x")
	encodedReturnDataAsBytes, err := hex.DecodeString(encodedReturnData)
	if err != nil {
		return []bool{}, err
	}

	// Hardcoded type as the abi is not aware of our custom constructor return type
	boolArrayType, err := abi.NewType("bool[]", "", nil)

	if err != nil {
		return []bool{}, err
	}

	args := abi.Arguments{
		{Type: boolArrayType},
	}

	decoded, err := args.Unpack(encodedReturnDataAsBytes)

	e.log.Info("Decoded results", "results", decoded)

	if err != nil {
		return []bool{}, err
	}

	result, ok := decoded[0].([]bool)

	if !ok {
		err = fmt.Errorf("failed to convert raw data to []bool")
		return []bool{}, err
	}

	e.log.Info("BatchCheckSeqConfig called successfully, results were", "result", result)

	return result, nil
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

func convertUint64ArrayToBigIntArray(input []uint64) []*big.Int {
	result := make([]*big.Int, len(input))
	for i, v := range input {
		result[i] = new(big.Int).SetUint64(v)
	}
	return result
}
