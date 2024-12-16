// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package BatchRandomTicketInstruction

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

// BatchRandomTicketInstructionMetaData contains all meta data concerning the BatchRandomTicketInstruction contract.
var BatchRandomTicketInstructionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"timestamps\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balances\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"seenWinners\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usedTickets\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610b21380380610b218339818101604052810190610032919061073a565b6000734200000000000000000000000000000000000028905060008173ffffffffffffffffffffffffffffffffffffffff166317d70f7c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610098573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100bc9190610783565b90506000835167ffffffffffffffff8111156100db576100da6105c1565b5b60405190808252806020026020018201604052801561011457816020015b610101610567565b8152602001906001900390816100f95790505b50905060005b845181101561053757600060018487848151811061013b5761013a6107b0565b5b602002602001015144604051602001610155929190610800565b6040516020818303038152906040528051906020012060001c610178919061085b565b61018291906108bb565b905060008573ffffffffffffffffffffffffffffffffffffffff16636352211e836040518263ffffffff1660e01b81526004016101bf91906108fe565b602060405180830381865afa1580156101dc573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102009190610977565b9050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036102a0576040518060400160405280888581518110610254576102536107b0565b5b602002602001015181526020018273ffffffffffffffffffffffffffffffffffffffff1681525084848151811061028e5761028d6107b0565b5b6020026020010181905250505061052a565b60008673ffffffffffffffffffffffffffffffffffffffff166370a08231836040518263ffffffff1660e01b81526004016102db91906109b3565b602060405180830381865afa1580156102f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031c9190610783565b90506000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1690508080156103ba575060008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b156103c457600092505b80610474576001826103d691906109ce565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060018060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506104c2565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000815460019003919050819055505b60405180604001604052808a87815181106104e0576104df6107b0565b5b602002602001015181526020018473ffffffffffffffffffffffffffffffffffffffff1681525086868151811061051a576105196107b0565b5b6020026020010181905250505050505b808060010191505061011a565b5060008160405160200161054b9190610afe565b6040516020818303038152906040529050602081018059038082f35b604051806040016040528060008152602001600073ffffffffffffffffffffffffffffffffffffffff1681525090565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105f9826105b0565b810181811067ffffffffffffffff82111715610618576106176105c1565b5b80604052505050565b600061062b610597565b905061063782826105f0565b919050565b600067ffffffffffffffff821115610657576106566105c1565b5b602082029050602081019050919050565b600080fd5b6000819050919050565b6106808161066d565b811461068b57600080fd5b50565b60008151905061069d81610677565b92915050565b60006106b66106b18461063c565b610621565b905080838252602082019050602084028301858111156106d9576106d8610668565b5b835b8181101561070257806106ee888261068e565b8452602084019350506020810190506106db565b5050509392505050565b600082601f830112610721576107206105ab565b5b81516107318482602086016106a3565b91505092915050565b6000602082840312156107505761074f6105a1565b5b600082015167ffffffffffffffff81111561076e5761076d6105a6565b5b61077a8482850161070c565b91505092915050565b600060208284031215610799576107986105a1565b5b60006107a78482850161068e565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b6107fa6107f58261066d565b6107df565b82525050565b600061080c82856107e9565b60208201915061081c82846107e9565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006108668261066d565b91506108718361066d565b9250826108815761088061082c565b5b828206905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006108c68261066d565b91506108d18361066d565b92508282019050808211156108e9576108e861088c565b5b92915050565b6108f88161066d565b82525050565b600060208201905061091360008301846108ef565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061094482610919565b9050919050565b61095481610939565b811461095f57600080fd5b50565b6000815190506109718161094b565b92915050565b60006020828403121561098d5761098c6105a1565b5b600061099b84828501610962565b91505092915050565b6109ad81610939565b82525050565b60006020820190506109c860008301846109a4565b92915050565b60006109d98261066d565b91506109e48361066d565b92508282039050818111156109fc576109fb61088c565b5b92915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b610a378161066d565b82525050565b610a4681610939565b82525050565b604082016000820151610a626000850182610a2e565b506020820151610a756020850182610a3d565b50505050565b6000610a878383610a4c565b60408301905092915050565b6000602082019050919050565b6000610aab82610a02565b610ab58185610a0d565b9350610ac083610a1e565b8060005b83811015610af1578151610ad88882610a7b565b9750610ae383610a93565b925050600181019050610ac4565b5085935050505092915050565b60006020820190508181036000830152610b188184610aa0565b90509291505056fe",
}

// BatchRandomTicketInstructionABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchRandomTicketInstructionMetaData.ABI instead.
var BatchRandomTicketInstructionABI = BatchRandomTicketInstructionMetaData.ABI

// BatchRandomTicketInstructionBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchRandomTicketInstructionMetaData.Bin instead.
var BatchRandomTicketInstructionBin = BatchRandomTicketInstructionMetaData.Bin

// DeployBatchRandomTicketInstruction deploys a new Ethereum contract, binding an instance of BatchRandomTicketInstruction to it.
func DeployBatchRandomTicketInstruction(auth *bind.TransactOpts, backend bind.ContractBackend, timestamps []*big.Int) (common.Address, *types.Transaction, *BatchRandomTicketInstruction, error) {
	parsed, err := BatchRandomTicketInstructionMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchRandomTicketInstructionBin), backend, timestamps)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchRandomTicketInstruction{BatchRandomTicketInstructionCaller: BatchRandomTicketInstructionCaller{contract: contract}, BatchRandomTicketInstructionTransactor: BatchRandomTicketInstructionTransactor{contract: contract}, BatchRandomTicketInstructionFilterer: BatchRandomTicketInstructionFilterer{contract: contract}}, nil
}

// BatchRandomTicketInstruction is an auto generated Go binding around an Ethereum contract.
type BatchRandomTicketInstruction struct {
	BatchRandomTicketInstructionCaller     // Read-only binding to the contract
	BatchRandomTicketInstructionTransactor // Write-only binding to the contract
	BatchRandomTicketInstructionFilterer   // Log filterer for contract events
}

// BatchRandomTicketInstructionCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchRandomTicketInstructionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchRandomTicketInstructionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchRandomTicketInstructionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchRandomTicketInstructionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchRandomTicketInstructionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchRandomTicketInstructionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchRandomTicketInstructionSession struct {
	Contract     *BatchRandomTicketInstruction // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// BatchRandomTicketInstructionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchRandomTicketInstructionCallerSession struct {
	Contract *BatchRandomTicketInstructionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// BatchRandomTicketInstructionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchRandomTicketInstructionTransactorSession struct {
	Contract     *BatchRandomTicketInstructionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// BatchRandomTicketInstructionRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchRandomTicketInstructionRaw struct {
	Contract *BatchRandomTicketInstruction // Generic contract binding to access the raw methods on
}

// BatchRandomTicketInstructionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchRandomTicketInstructionCallerRaw struct {
	Contract *BatchRandomTicketInstructionCaller // Generic read-only contract binding to access the raw methods on
}

// BatchRandomTicketInstructionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchRandomTicketInstructionTransactorRaw struct {
	Contract *BatchRandomTicketInstructionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchRandomTicketInstruction creates a new instance of BatchRandomTicketInstruction, bound to a specific deployed contract.
func NewBatchRandomTicketInstruction(address common.Address, backend bind.ContractBackend) (*BatchRandomTicketInstruction, error) {
	contract, err := bindBatchRandomTicketInstruction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchRandomTicketInstruction{BatchRandomTicketInstructionCaller: BatchRandomTicketInstructionCaller{contract: contract}, BatchRandomTicketInstructionTransactor: BatchRandomTicketInstructionTransactor{contract: contract}, BatchRandomTicketInstructionFilterer: BatchRandomTicketInstructionFilterer{contract: contract}}, nil
}

// NewBatchRandomTicketInstructionCaller creates a new read-only instance of BatchRandomTicketInstruction, bound to a specific deployed contract.
func NewBatchRandomTicketInstructionCaller(address common.Address, caller bind.ContractCaller) (*BatchRandomTicketInstructionCaller, error) {
	contract, err := bindBatchRandomTicketInstruction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchRandomTicketInstructionCaller{contract: contract}, nil
}

// NewBatchRandomTicketInstructionTransactor creates a new write-only instance of BatchRandomTicketInstruction, bound to a specific deployed contract.
func NewBatchRandomTicketInstructionTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchRandomTicketInstructionTransactor, error) {
	contract, err := bindBatchRandomTicketInstruction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchRandomTicketInstructionTransactor{contract: contract}, nil
}

// NewBatchRandomTicketInstructionFilterer creates a new log filterer instance of BatchRandomTicketInstruction, bound to a specific deployed contract.
func NewBatchRandomTicketInstructionFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchRandomTicketInstructionFilterer, error) {
	contract, err := bindBatchRandomTicketInstruction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchRandomTicketInstructionFilterer{contract: contract}, nil
}

// bindBatchRandomTicketInstruction binds a generic wrapper to an already deployed contract.
func bindBatchRandomTicketInstruction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchRandomTicketInstructionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchRandomTicketInstruction.Contract.BatchRandomTicketInstructionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchRandomTicketInstruction.Contract.BatchRandomTicketInstructionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchRandomTicketInstruction.Contract.BatchRandomTicketInstructionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchRandomTicketInstruction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchRandomTicketInstruction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchRandomTicketInstruction.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BatchRandomTicketInstruction.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _BatchRandomTicketInstruction.Contract.Balances(&_BatchRandomTicketInstruction.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _BatchRandomTicketInstruction.Contract.Balances(&_BatchRandomTicketInstruction.CallOpts, arg0)
}

// SeenWinners is a free data retrieval call binding the contract method 0x44c603ab.
//
// Solidity: function seenWinners(address ) view returns(bool)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCaller) SeenWinners(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _BatchRandomTicketInstruction.contract.Call(opts, &out, "seenWinners", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SeenWinners is a free data retrieval call binding the contract method 0x44c603ab.
//
// Solidity: function seenWinners(address ) view returns(bool)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionSession) SeenWinners(arg0 common.Address) (bool, error) {
	return _BatchRandomTicketInstruction.Contract.SeenWinners(&_BatchRandomTicketInstruction.CallOpts, arg0)
}

// SeenWinners is a free data retrieval call binding the contract method 0x44c603ab.
//
// Solidity: function seenWinners(address ) view returns(bool)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCallerSession) SeenWinners(arg0 common.Address) (bool, error) {
	return _BatchRandomTicketInstruction.Contract.SeenWinners(&_BatchRandomTicketInstruction.CallOpts, arg0)
}

// UsedTickets is a free data retrieval call binding the contract method 0x71efdc21.
//
// Solidity: function usedTickets(uint256 ) view returns(bool)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCaller) UsedTickets(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _BatchRandomTicketInstruction.contract.Call(opts, &out, "usedTickets", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedTickets is a free data retrieval call binding the contract method 0x71efdc21.
//
// Solidity: function usedTickets(uint256 ) view returns(bool)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionSession) UsedTickets(arg0 *big.Int) (bool, error) {
	return _BatchRandomTicketInstruction.Contract.UsedTickets(&_BatchRandomTicketInstruction.CallOpts, arg0)
}

// UsedTickets is a free data retrieval call binding the contract method 0x71efdc21.
//
// Solidity: function usedTickets(uint256 ) view returns(bool)
func (_BatchRandomTicketInstruction *BatchRandomTicketInstructionCallerSession) UsedTickets(arg0 *big.Int) (bool, error) {
	return _BatchRandomTicketInstruction.Contract.UsedTickets(&_BatchRandomTicketInstruction.CallOpts, arg0)
}
