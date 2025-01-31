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
	Bin: "0x608060405234801561001057600080fd5b506040516105c53803806105c5833981810160405281019061003291906103ca565b6000815167ffffffffffffffff81111561004f5761004e610287565b5b60405190808252806020026020018201604052801561007d5781602001602082028036833780820191505090505b50905060005b82518110156101cf57600073ffffffffffffffffffffffffffffffffffffffff168382815181106100b7576100b6610426565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff16036101065760008282815181106100ef576100ee610426565b5b6020026020010190151590811515815250506101c2565b8373ffffffffffffffffffffffffffffffffffffffff16630e61d04084838151811061013557610134610426565b5b60200260200101516040518263ffffffff1660e01b81526004016101599190610464565b6020604051808303816000875af1158015610178573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061019c91906104b7565b8282815181106101af576101ae610426565b5b6020026020010190151590811515815250505b8080600101915050610083565b506000816040516020016101e391906105a2565b6040516020818303038152906040529050602081018059038082f35b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061023e82610213565b9050919050565b61024e81610233565b811461025957600080fd5b50565b60008151905061026b81610245565b92915050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6102bf82610276565b810181811067ffffffffffffffff821117156102de576102dd610287565b5b80604052505050565b60006102f16101ff565b90506102fd82826102b6565b919050565b600067ffffffffffffffff82111561031d5761031c610287565b5b602082029050602081019050919050565b600080fd5b600061034661034184610302565b6102e7565b905080838252602082019050602084028301858111156103695761036861032e565b5b835b81811015610392578061037e888261025c565b84526020840193505060208101905061036b565b5050509392505050565b600082601f8301126103b1576103b0610271565b5b81516103c1848260208601610333565b91505092915050565b600080604083850312156103e1576103e0610209565b5b60006103ef8582860161025c565b925050602083015167ffffffffffffffff8111156104105761040f61020e565b5b61041c8582860161039c565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b61045e81610233565b82525050565b60006020820190506104796000830184610455565b92915050565b60008115159050919050565b6104948161047f565b811461049f57600080fd5b50565b6000815190506104b18161048b565b92915050565b6000602082840312156104cd576104cc610209565b5b60006104db848285016104a2565b91505092915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6105198161047f565b82525050565b600061052b8383610510565b60208301905092915050565b6000602082019050919050565b600061054f826104e4565b61055981856104ef565b935061056483610500565b8060005b8381101561059557815161057c888261051f565b975061058783610537565b925050600181019050610568565b5085935050505092915050565b600060208201905081810360008301526105bc8184610544565b90509291505056fe",
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
