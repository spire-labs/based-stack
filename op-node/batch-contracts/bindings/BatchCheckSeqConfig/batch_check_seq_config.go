// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BatchCheckSeqConfig

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

// BatchCheckSeqConfigMetaData contains all meta data concerning the BatchCheckSeqConfig contract.
var BatchCheckSeqConfigMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_sysConfig\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_addrs\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"nonpayable\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161054a38038061054a8339818101604052810190610032919061034f565b6000815167ffffffffffffffff81111561004f5761004e61020c565b5b60405190808252806020026020018201604052801561007d5781602001602082028036833780820191505090505b50905060005b8251811015610154578373ffffffffffffffffffffffffffffffffffffffff16630e61d0408483815181106100bb576100ba6103ab565b5b60200260200101516040518263ffffffff1660e01b81526004016100df91906103e9565b6020604051808303816000875af11580156100fe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610122919061043c565b828281518110610135576101346103ab565b5b6020026020010190151590811515815250508080600101915050610083565b506000816040516020016101689190610527565b6040516020818303038152906040529050602081018059038082f35b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101c382610198565b9050919050565b6101d3816101b8565b81146101de57600080fd5b50565b6000815190506101f0816101ca565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610244826101fb565b810181811067ffffffffffffffff821117156102635761026261020c565b5b80604052505050565b6000610276610184565b9050610282828261023b565b919050565b600067ffffffffffffffff8211156102a2576102a161020c565b5b602082029050602081019050919050565b600080fd5b60006102cb6102c684610287565b61026c565b905080838252602082019050602084028301858111156102ee576102ed6102b3565b5b835b81811015610317578061030388826101e1565b8452602084019350506020810190506102f0565b5050509392505050565b600082601f830112610336576103356101f6565b5b81516103468482602086016102b8565b91505092915050565b600080604083850312156103665761036561018e565b5b6000610374858286016101e1565b925050602083015167ffffffffffffffff81111561039557610394610193565b5b6103a185828601610321565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6103e3816101b8565b82525050565b60006020820190506103fe60008301846103da565b92915050565b60008115159050919050565b61041981610404565b811461042457600080fd5b50565b60008151905061043681610410565b92915050565b6000602082840312156104525761045161018e565b5b600061046084828501610427565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b61049e81610404565b82525050565b60006104b08383610495565b60208301905092915050565b6000602082019050919050565b60006104d482610469565b6104de8185610474565b93506104e983610485565b8060005b8381101561051a57815161050188826104a4565b975061050c836104bc565b9250506001810190506104ed565b5085935050505092915050565b6000602082019050818103600083015261054181846104c9565b90509291505056fe",
}

// BatchCheckSeqConfigABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchCheckSeqConfigMetaData.ABI instead.
var BatchCheckSeqConfigABI = BatchCheckSeqConfigMetaData.ABI

// BatchCheckSeqConfigBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchCheckSeqConfigMetaData.Bin instead.
var BatchCheckSeqConfigBin = BatchCheckSeqConfigMetaData.Bin

// DeployBatchCheckSeqConfig deploys a new Ethereum contract, binding an instance of BatchCheckSeqConfig to it.
func DeployBatchCheckSeqConfig(auth *bind.TransactOpts, backend bind.ContractBackend, _sysConfig common.Address, _addrs []common.Address) (common.Address, *types.Transaction, *BatchCheckSeqConfig, error) {
	parsed, err := BatchCheckSeqConfigMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchCheckSeqConfigBin), backend, _sysConfig, _addrs)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchCheckSeqConfig{BatchCheckSeqConfigCaller: BatchCheckSeqConfigCaller{contract: contract}, BatchCheckSeqConfigTransactor: BatchCheckSeqConfigTransactor{contract: contract}, BatchCheckSeqConfigFilterer: BatchCheckSeqConfigFilterer{contract: contract}}, nil
}

// BatchCheckSeqConfig is an auto generated Go binding around an Ethereum contract.
type BatchCheckSeqConfig struct {
	BatchCheckSeqConfigCaller     // Read-only binding to the contract
	BatchCheckSeqConfigTransactor // Write-only binding to the contract
	BatchCheckSeqConfigFilterer   // Log filterer for contract events
}

// BatchCheckSeqConfigCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchCheckSeqConfigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchCheckSeqConfigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchCheckSeqConfigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchCheckSeqConfigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchCheckSeqConfigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchCheckSeqConfigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchCheckSeqConfigSession struct {
	Contract     *BatchCheckSeqConfig // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// BatchCheckSeqConfigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchCheckSeqConfigCallerSession struct {
	Contract *BatchCheckSeqConfigCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// BatchCheckSeqConfigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchCheckSeqConfigTransactorSession struct {
	Contract     *BatchCheckSeqConfigTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// BatchCheckSeqConfigRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchCheckSeqConfigRaw struct {
	Contract *BatchCheckSeqConfig // Generic contract binding to access the raw methods on
}

// BatchCheckSeqConfigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchCheckSeqConfigCallerRaw struct {
	Contract *BatchCheckSeqConfigCaller // Generic read-only contract binding to access the raw methods on
}

// BatchCheckSeqConfigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchCheckSeqConfigTransactorRaw struct {
	Contract *BatchCheckSeqConfigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchCheckSeqConfig creates a new instance of BatchCheckSeqConfig, bound to a specific deployed contract.
func NewBatchCheckSeqConfig(address common.Address, backend bind.ContractBackend) (*BatchCheckSeqConfig, error) {
	contract, err := bindBatchCheckSeqConfig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchCheckSeqConfig{BatchCheckSeqConfigCaller: BatchCheckSeqConfigCaller{contract: contract}, BatchCheckSeqConfigTransactor: BatchCheckSeqConfigTransactor{contract: contract}, BatchCheckSeqConfigFilterer: BatchCheckSeqConfigFilterer{contract: contract}}, nil
}

// NewBatchCheckSeqConfigCaller creates a new read-only instance of BatchCheckSeqConfig, bound to a specific deployed contract.
func NewBatchCheckSeqConfigCaller(address common.Address, caller bind.ContractCaller) (*BatchCheckSeqConfigCaller, error) {
	contract, err := bindBatchCheckSeqConfig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchCheckSeqConfigCaller{contract: contract}, nil
}

// NewBatchCheckSeqConfigTransactor creates a new write-only instance of BatchCheckSeqConfig, bound to a specific deployed contract.
func NewBatchCheckSeqConfigTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchCheckSeqConfigTransactor, error) {
	contract, err := bindBatchCheckSeqConfig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchCheckSeqConfigTransactor{contract: contract}, nil
}

// NewBatchCheckSeqConfigFilterer creates a new log filterer instance of BatchCheckSeqConfig, bound to a specific deployed contract.
func NewBatchCheckSeqConfigFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchCheckSeqConfigFilterer, error) {
	contract, err := bindBatchCheckSeqConfig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchCheckSeqConfigFilterer{contract: contract}, nil
}

// bindBatchCheckSeqConfig binds a generic wrapper to an already deployed contract.
func bindBatchCheckSeqConfig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchCheckSeqConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchCheckSeqConfig *BatchCheckSeqConfigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchCheckSeqConfig.Contract.BatchCheckSeqConfigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchCheckSeqConfig *BatchCheckSeqConfigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchCheckSeqConfig.Contract.BatchCheckSeqConfigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchCheckSeqConfig *BatchCheckSeqConfigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchCheckSeqConfig.Contract.BatchCheckSeqConfigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchCheckSeqConfig *BatchCheckSeqConfigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchCheckSeqConfig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchCheckSeqConfig *BatchCheckSeqConfigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchCheckSeqConfig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchCheckSeqConfig *BatchCheckSeqConfigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchCheckSeqConfig.Contract.contract.Transact(opts, method, params...)
}
