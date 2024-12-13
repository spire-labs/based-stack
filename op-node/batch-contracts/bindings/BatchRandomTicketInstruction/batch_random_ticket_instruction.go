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
	Bin: "0x608060405234801561001057600080fd5b50604051610aa9380380610aa9833981810160405281019061003291906106c2565b6000734200000000000000000000000000000000000028905060008173ffffffffffffffffffffffffffffffffffffffff166317d70f7c6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610098573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100bc919061070b565b90506000835167ffffffffffffffff8111156100db576100da610549565b5b60405190808252806020026020018201604052801561011457816020015b6101016104ef565b8152602001906001900390816100f95790505b50905060005b84518110156104bf5760006001846101329190610767565b8683815181106101455761014461079b565b5b60200260200101514460405160200161015f9291906107eb565b6040516020818303038152906040528051906020012060001c6101829190610846565b905060008573ffffffffffffffffffffffffffffffffffffffff16636352211e836040518263ffffffff1660e01b81526004016101bf9190610886565b602060405180830381865afa1580156101dc573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061020091906108ff565b905060008673ffffffffffffffffffffffffffffffffffffffff166370a08231836040518263ffffffff1660e01b815260040161023d919061093b565b602060405180830381865afa15801561025a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027e919061070b565b90506000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905080801561031c575060008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b1561032657600092505b806103d2576001826103389190610956565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555060018060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055505b60405180604001604052808a87815181106103f0576103ef61079b565b5b602002602001015181526020018973ffffffffffffffffffffffffffffffffffffffff16636352211e876040518263ffffffff1660e01b81526004016104369190610886565b602060405180830381865afa158015610453573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061047791906108ff565b73ffffffffffffffffffffffffffffffffffffffff168152508686815181106104a3576104a261079b565b5b602002602001018190525050505050808060010191505061011a565b506000816040516020016104d39190610a86565b6040516020818303038152906040529050602081018059038082f35b604051806040016040528060008152602001600073ffffffffffffffffffffffffffffffffffffffff1681525090565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61058182610538565b810181811067ffffffffffffffff821117156105a05761059f610549565b5b80604052505050565b60006105b361051f565b90506105bf8282610578565b919050565b600067ffffffffffffffff8211156105df576105de610549565b5b602082029050602081019050919050565b600080fd5b6000819050919050565b610608816105f5565b811461061357600080fd5b50565b600081519050610625816105ff565b92915050565b600061063e610639846105c4565b6105a9565b90508083825260208201905060208402830185811115610661576106606105f0565b5b835b8181101561068a57806106768882610616565b845260208401935050602081019050610663565b5050509392505050565b600082601f8301126106a9576106a8610533565b5b81516106b984826020860161062b565b91505092915050565b6000602082840312156106d8576106d7610529565b5b600082015167ffffffffffffffff8111156106f6576106f561052e565b5b61070284828501610694565b91505092915050565b60006020828403121561072157610720610529565b5b600061072f84828501610616565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610772826105f5565b915061077d836105f5565b925082820190508082111561079557610794610738565b5b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6000819050919050565b6107e56107e0826105f5565b6107ca565b82525050565b60006107f782856107d4565b60208201915061080782846107d4565b6020820191508190509392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000610851826105f5565b915061085c836105f5565b92508261086c5761086b610817565b5b828206905092915050565b610880816105f5565b82525050565b600060208201905061089b6000830184610877565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006108cc826108a1565b9050919050565b6108dc816108c1565b81146108e757600080fd5b50565b6000815190506108f9816108d3565b92915050565b60006020828403121561091557610914610529565b5b6000610923848285016108ea565b91505092915050565b610935816108c1565b82525050565b6000602082019050610950600083018461092c565b92915050565b6000610961826105f5565b915061096c836105f5565b925082820390508181111561098457610983610738565b5b92915050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b6109bf816105f5565b82525050565b6109ce816108c1565b82525050565b6040820160008201516109ea60008501826109b6565b5060208201516109fd60208501826109c5565b50505050565b6000610a0f83836109d4565b60408301905092915050565b6000602082019050919050565b6000610a338261098a565b610a3d8185610995565b9350610a48836109a6565b8060005b83811015610a79578151610a608882610a03565b9750610a6b83610a1b565b925050600181019050610a4c565b5085935050505092915050565b60006020820190508181036000830152610aa08184610a28565b90509291505056fe",
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
