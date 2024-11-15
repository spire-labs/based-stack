// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BatchTicketAccounting

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BatchTicketAccountingMetaData contains all meta data concerning the BatchTicketAccounting contract.
var BatchTicketAccountingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161054638038061054683398181016040528101906100329190610360565b600073420000000000000000000000000000000000002890506000825167ffffffffffffffff811115610068576100676101bf565b5b6040519080825280602002602001820160405280156100965781602001602082028036833780820191505090505b50905060005b8351811015610165578273ffffffffffffffffffffffffffffffffffffffff166370ca5bbe8583815181106100d4576100d36103a9565b5b60200260200101516040518263ffffffff1660e01b81526004016100f891906103e7565b602060405180830381865afa158015610115573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101399190610438565b82828151811061014c5761014b6103a9565b5b602002602001018181525050808060010191505061009c565b506000816040516020016101799190610523565b6040516020818303038152906040529050602081018059038082f35b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6101f7826101ae565b810181811067ffffffffffffffff82111715610216576102156101bf565b5b80604052505050565b6000610229610195565b905061023582826101ee565b919050565b600067ffffffffffffffff821115610255576102546101bf565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006102968261026b565b9050919050565b6102a68161028b565b81146102b157600080fd5b50565b6000815190506102c38161029d565b92915050565b60006102dc6102d78461023a565b61021f565b905080838252602082019050602084028301858111156102ff576102fe610266565b5b835b81811015610328578061031488826102b4565b845260208401935050602081019050610301565b5050509392505050565b600082601f830112610347576103466101a9565b5b81516103578482602086016102c9565b91505092915050565b6000602082840312156103765761037561019f565b5b600082015167ffffffffffffffff811115610394576103936101a4565b5b6103a084828501610332565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6103e18161028b565b82525050565b60006020820190506103fc60008301846103d8565b92915050565b6000819050919050565b61041581610402565b811461042057600080fd5b50565b6000815190506104328161040c565b92915050565b60006020828403121561044e5761044d61019f565b5b600061045c84828501610423565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61049a81610402565b82525050565b60006104ac8383610491565b60208301905092915050565b6000602082019050919050565b60006104d082610465565b6104da8185610470565b93506104e583610481565b8060005b838110156105165781516104fd88826104a0565b9750610508836104b8565b9250506001810190506104e9565b5085935050505092915050565b6000602082019050818103600083015261053d81846104c5565b90509291505056fe",
}

// BatchTicketAccountingABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchTicketAccountingMetaData.ABI instead.
var BatchTicketAccountingABI = BatchTicketAccountingMetaData.ABI

// BatchTicketAccountingBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchTicketAccountingMetaData.Bin instead.
var BatchTicketAccountingBin = BatchTicketAccountingMetaData.Bin

// DeployBatchTicketAccounting deploys a new Ethereum contract, binding an instance of BatchTicketAccounting to it.
func DeployBatchTicketAccounting(auth *bind.TransactOpts, backend bind.ContractBackend, addresses []common.Address) (common.Address, *types.Transaction, *BatchTicketAccounting, error) {
	parsed, err := BatchTicketAccountingMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchTicketAccountingBin), backend, addresses)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchTicketAccounting{BatchTicketAccountingCaller: BatchTicketAccountingCaller{contract: contract}, BatchTicketAccountingTransactor: BatchTicketAccountingTransactor{contract: contract}, BatchTicketAccountingFilterer: BatchTicketAccountingFilterer{contract: contract}}, nil
}

// BatchTicketAccounting is an auto generated Go binding around an Ethereum contract.
type BatchTicketAccounting struct {
	BatchTicketAccountingCaller     // Read-only binding to the contract
	BatchTicketAccountingTransactor // Write-only binding to the contract
	BatchTicketAccountingFilterer   // Log filterer for contract events
}

// BatchTicketAccountingCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchTicketAccountingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchTicketAccountingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchTicketAccountingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchTicketAccountingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchTicketAccountingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchTicketAccountingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchTicketAccountingSession struct {
	Contract     *BatchTicketAccounting // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BatchTicketAccountingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchTicketAccountingCallerSession struct {
	Contract *BatchTicketAccountingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// BatchTicketAccountingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchTicketAccountingTransactorSession struct {
	Contract     *BatchTicketAccountingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// BatchTicketAccountingRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchTicketAccountingRaw struct {
	Contract *BatchTicketAccounting // Generic contract binding to access the raw methods on
}

// BatchTicketAccountingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchTicketAccountingCallerRaw struct {
	Contract *BatchTicketAccountingCaller // Generic read-only contract binding to access the raw methods on
}

// BatchTicketAccountingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchTicketAccountingTransactorRaw struct {
	Contract *BatchTicketAccountingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchTicketAccounting creates a new instance of BatchTicketAccounting, bound to a specific deployed contract.
func NewBatchTicketAccounting(address common.Address, backend bind.ContractBackend) (*BatchTicketAccounting, error) {
	contract, err := bindBatchTicketAccounting(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchTicketAccounting{BatchTicketAccountingCaller: BatchTicketAccountingCaller{contract: contract}, BatchTicketAccountingTransactor: BatchTicketAccountingTransactor{contract: contract}, BatchTicketAccountingFilterer: BatchTicketAccountingFilterer{contract: contract}}, nil
}

// NewBatchTicketAccountingCaller creates a new read-only instance of BatchTicketAccounting, bound to a specific deployed contract.
func NewBatchTicketAccountingCaller(address common.Address, caller bind.ContractCaller) (*BatchTicketAccountingCaller, error) {
	contract, err := bindBatchTicketAccounting(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchTicketAccountingCaller{contract: contract}, nil
}

// NewBatchTicketAccountingTransactor creates a new write-only instance of BatchTicketAccounting, bound to a specific deployed contract.
func NewBatchTicketAccountingTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchTicketAccountingTransactor, error) {
	contract, err := bindBatchTicketAccounting(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchTicketAccountingTransactor{contract: contract}, nil
}

// NewBatchTicketAccountingFilterer creates a new log filterer instance of BatchTicketAccounting, bound to a specific deployed contract.
func NewBatchTicketAccountingFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchTicketAccountingFilterer, error) {
	contract, err := bindBatchTicketAccounting(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchTicketAccountingFilterer{contract: contract}, nil
}

// bindBatchTicketAccounting binds a generic wrapper to an already deployed contract.
func bindBatchTicketAccounting(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchTicketAccountingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchTicketAccounting *BatchTicketAccountingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchTicketAccounting.Contract.BatchTicketAccountingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchTicketAccounting *BatchTicketAccountingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchTicketAccounting.Contract.BatchTicketAccountingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchTicketAccounting *BatchTicketAccountingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchTicketAccounting.Contract.BatchTicketAccountingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchTicketAccounting *BatchTicketAccountingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchTicketAccounting.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchTicketAccounting *BatchTicketAccountingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchTicketAccounting.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchTicketAccounting *BatchTicketAccountingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchTicketAccounting.Contract.contract.Transact(opts, method, params...)
}
