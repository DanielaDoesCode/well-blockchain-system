// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// HistoryFactoryMetaData contains all meta data concerning the HistoryFactory contract.
var HistoryFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_userRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitmentToHistory\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createHistory\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"patientAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getHistoryContract\",\"inputs\":[{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHistoryContractByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policyVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setPolicyVerifier\",\"inputs\":[{\"name\":\"_policyVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWellRegistry\",\"inputs\":[{\"name\":\"_wellRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractUserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"wellRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"HistoryCreated\",\"inputs\":[{\"name\":\"patient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"historyContract\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// HistoryFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use HistoryFactoryMetaData.ABI instead.
var HistoryFactoryABI = HistoryFactoryMetaData.ABI

// HistoryFactory is an auto generated Go binding around an Ethereum contract.
type HistoryFactory struct {
	HistoryFactoryCaller     // Read-only binding to the contract
	HistoryFactoryTransactor // Write-only binding to the contract
	HistoryFactoryFilterer   // Log filterer for contract events
}

// HistoryFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type HistoryFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoryFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HistoryFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoryFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HistoryFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoryFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HistoryFactorySession struct {
	Contract     *HistoryFactory   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HistoryFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HistoryFactoryCallerSession struct {
	Contract *HistoryFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// HistoryFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HistoryFactoryTransactorSession struct {
	Contract     *HistoryFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// HistoryFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type HistoryFactoryRaw struct {
	Contract *HistoryFactory // Generic contract binding to access the raw methods on
}

// HistoryFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HistoryFactoryCallerRaw struct {
	Contract *HistoryFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// HistoryFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HistoryFactoryTransactorRaw struct {
	Contract *HistoryFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHistoryFactory creates a new instance of HistoryFactory, bound to a specific deployed contract.
func NewHistoryFactory(address common.Address, backend bind.ContractBackend) (*HistoryFactory, error) {
	contract, err := bindHistoryFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HistoryFactory{HistoryFactoryCaller: HistoryFactoryCaller{contract: contract}, HistoryFactoryTransactor: HistoryFactoryTransactor{contract: contract}, HistoryFactoryFilterer: HistoryFactoryFilterer{contract: contract}}, nil
}

// NewHistoryFactoryCaller creates a new read-only instance of HistoryFactory, bound to a specific deployed contract.
func NewHistoryFactoryCaller(address common.Address, caller bind.ContractCaller) (*HistoryFactoryCaller, error) {
	contract, err := bindHistoryFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HistoryFactoryCaller{contract: contract}, nil
}

// NewHistoryFactoryTransactor creates a new write-only instance of HistoryFactory, bound to a specific deployed contract.
func NewHistoryFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*HistoryFactoryTransactor, error) {
	contract, err := bindHistoryFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HistoryFactoryTransactor{contract: contract}, nil
}

// NewHistoryFactoryFilterer creates a new log filterer instance of HistoryFactory, bound to a specific deployed contract.
func NewHistoryFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*HistoryFactoryFilterer, error) {
	contract, err := bindHistoryFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HistoryFactoryFilterer{contract: contract}, nil
}

// bindHistoryFactory binds a generic wrapper to an already deployed contract.
func bindHistoryFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HistoryFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HistoryFactory *HistoryFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HistoryFactory.Contract.HistoryFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HistoryFactory *HistoryFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HistoryFactory.Contract.HistoryFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HistoryFactory *HistoryFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HistoryFactory.Contract.HistoryFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HistoryFactory *HistoryFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HistoryFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HistoryFactory *HistoryFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HistoryFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HistoryFactory *HistoryFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HistoryFactory.Contract.contract.Transact(opts, method, params...)
}

// CommitmentToHistory is a free data retrieval call binding the contract method 0xaaabf72a.
//
// Solidity: function commitmentToHistory(bytes32 ) view returns(address)
func (_HistoryFactory *HistoryFactoryCaller) CommitmentToHistory(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _HistoryFactory.contract.Call(opts, &out, "commitmentToHistory", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CommitmentToHistory is a free data retrieval call binding the contract method 0xaaabf72a.
//
// Solidity: function commitmentToHistory(bytes32 ) view returns(address)
func (_HistoryFactory *HistoryFactorySession) CommitmentToHistory(arg0 [32]byte) (common.Address, error) {
	return _HistoryFactory.Contract.CommitmentToHistory(&_HistoryFactory.CallOpts, arg0)
}

// CommitmentToHistory is a free data retrieval call binding the contract method 0xaaabf72a.
//
// Solidity: function commitmentToHistory(bytes32 ) view returns(address)
func (_HistoryFactory *HistoryFactoryCallerSession) CommitmentToHistory(arg0 [32]byte) (common.Address, error) {
	return _HistoryFactory.Contract.CommitmentToHistory(&_HistoryFactory.CallOpts, arg0)
}

// GetHistoryContract is a free data retrieval call binding the contract method 0xafc3b238.
//
// Solidity: function getHistoryContract(bytes32 internalID) view returns(address)
func (_HistoryFactory *HistoryFactoryCaller) GetHistoryContract(opts *bind.CallOpts, internalID [32]byte) (common.Address, error) {
	var out []interface{}
	err := _HistoryFactory.contract.Call(opts, &out, "getHistoryContract", internalID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetHistoryContract is a free data retrieval call binding the contract method 0xafc3b238.
//
// Solidity: function getHistoryContract(bytes32 internalID) view returns(address)
func (_HistoryFactory *HistoryFactorySession) GetHistoryContract(internalID [32]byte) (common.Address, error) {
	return _HistoryFactory.Contract.GetHistoryContract(&_HistoryFactory.CallOpts, internalID)
}

// GetHistoryContract is a free data retrieval call binding the contract method 0xafc3b238.
//
// Solidity: function getHistoryContract(bytes32 internalID) view returns(address)
func (_HistoryFactory *HistoryFactoryCallerSession) GetHistoryContract(internalID [32]byte) (common.Address, error) {
	return _HistoryFactory.Contract.GetHistoryContract(&_HistoryFactory.CallOpts, internalID)
}

// GetHistoryContractByAddress is a free data retrieval call binding the contract method 0x145f6705.
//
// Solidity: function getHistoryContractByAddress(address userAddress) view returns(address)
func (_HistoryFactory *HistoryFactoryCaller) GetHistoryContractByAddress(opts *bind.CallOpts, userAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _HistoryFactory.contract.Call(opts, &out, "getHistoryContractByAddress", userAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetHistoryContractByAddress is a free data retrieval call binding the contract method 0x145f6705.
//
// Solidity: function getHistoryContractByAddress(address userAddress) view returns(address)
func (_HistoryFactory *HistoryFactorySession) GetHistoryContractByAddress(userAddress common.Address) (common.Address, error) {
	return _HistoryFactory.Contract.GetHistoryContractByAddress(&_HistoryFactory.CallOpts, userAddress)
}

// GetHistoryContractByAddress is a free data retrieval call binding the contract method 0x145f6705.
//
// Solidity: function getHistoryContractByAddress(address userAddress) view returns(address)
func (_HistoryFactory *HistoryFactoryCallerSession) GetHistoryContractByAddress(userAddress common.Address) (common.Address, error) {
	return _HistoryFactory.Contract.GetHistoryContractByAddress(&_HistoryFactory.CallOpts, userAddress)
}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_HistoryFactory *HistoryFactoryCaller) PolicyVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HistoryFactory.contract.Call(opts, &out, "policyVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_HistoryFactory *HistoryFactorySession) PolicyVerifier() (common.Address, error) {
	return _HistoryFactory.Contract.PolicyVerifier(&_HistoryFactory.CallOpts)
}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_HistoryFactory *HistoryFactoryCallerSession) PolicyVerifier() (common.Address, error) {
	return _HistoryFactory.Contract.PolicyVerifier(&_HistoryFactory.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_HistoryFactory *HistoryFactoryCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HistoryFactory.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_HistoryFactory *HistoryFactorySession) UserRegistry() (common.Address, error) {
	return _HistoryFactory.Contract.UserRegistry(&_HistoryFactory.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_HistoryFactory *HistoryFactoryCallerSession) UserRegistry() (common.Address, error) {
	return _HistoryFactory.Contract.UserRegistry(&_HistoryFactory.CallOpts)
}

// WellRegistry is a free data retrieval call binding the contract method 0xc6583069.
//
// Solidity: function wellRegistry() view returns(address)
func (_HistoryFactory *HistoryFactoryCaller) WellRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HistoryFactory.contract.Call(opts, &out, "wellRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WellRegistry is a free data retrieval call binding the contract method 0xc6583069.
//
// Solidity: function wellRegistry() view returns(address)
func (_HistoryFactory *HistoryFactorySession) WellRegistry() (common.Address, error) {
	return _HistoryFactory.Contract.WellRegistry(&_HistoryFactory.CallOpts)
}

// WellRegistry is a free data retrieval call binding the contract method 0xc6583069.
//
// Solidity: function wellRegistry() view returns(address)
func (_HistoryFactory *HistoryFactoryCallerSession) WellRegistry() (common.Address, error) {
	return _HistoryFactory.Contract.WellRegistry(&_HistoryFactory.CallOpts)
}

// CreateHistory is a paid mutator transaction binding the contract method 0x6fb5d4a0.
//
// Solidity: function createHistory(bytes32 commitment, address patientAddress) returns(address)
func (_HistoryFactory *HistoryFactoryTransactor) CreateHistory(opts *bind.TransactOpts, commitment [32]byte, patientAddress common.Address) (*types.Transaction, error) {
	return _HistoryFactory.contract.Transact(opts, "createHistory", commitment, patientAddress)
}

// CreateHistory is a paid mutator transaction binding the contract method 0x6fb5d4a0.
//
// Solidity: function createHistory(bytes32 commitment, address patientAddress) returns(address)
func (_HistoryFactory *HistoryFactorySession) CreateHistory(commitment [32]byte, patientAddress common.Address) (*types.Transaction, error) {
	return _HistoryFactory.Contract.CreateHistory(&_HistoryFactory.TransactOpts, commitment, patientAddress)
}

// CreateHistory is a paid mutator transaction binding the contract method 0x6fb5d4a0.
//
// Solidity: function createHistory(bytes32 commitment, address patientAddress) returns(address)
func (_HistoryFactory *HistoryFactoryTransactorSession) CreateHistory(commitment [32]byte, patientAddress common.Address) (*types.Transaction, error) {
	return _HistoryFactory.Contract.CreateHistory(&_HistoryFactory.TransactOpts, commitment, patientAddress)
}

// SetPolicyVerifier is a paid mutator transaction binding the contract method 0xdb59f31f.
//
// Solidity: function setPolicyVerifier(address _policyVerifier) returns()
func (_HistoryFactory *HistoryFactoryTransactor) SetPolicyVerifier(opts *bind.TransactOpts, _policyVerifier common.Address) (*types.Transaction, error) {
	return _HistoryFactory.contract.Transact(opts, "setPolicyVerifier", _policyVerifier)
}

// SetPolicyVerifier is a paid mutator transaction binding the contract method 0xdb59f31f.
//
// Solidity: function setPolicyVerifier(address _policyVerifier) returns()
func (_HistoryFactory *HistoryFactorySession) SetPolicyVerifier(_policyVerifier common.Address) (*types.Transaction, error) {
	return _HistoryFactory.Contract.SetPolicyVerifier(&_HistoryFactory.TransactOpts, _policyVerifier)
}

// SetPolicyVerifier is a paid mutator transaction binding the contract method 0xdb59f31f.
//
// Solidity: function setPolicyVerifier(address _policyVerifier) returns()
func (_HistoryFactory *HistoryFactoryTransactorSession) SetPolicyVerifier(_policyVerifier common.Address) (*types.Transaction, error) {
	return _HistoryFactory.Contract.SetPolicyVerifier(&_HistoryFactory.TransactOpts, _policyVerifier)
}

// SetWellRegistry is a paid mutator transaction binding the contract method 0x16c03508.
//
// Solidity: function setWellRegistry(address _wellRegistry) returns()
func (_HistoryFactory *HistoryFactoryTransactor) SetWellRegistry(opts *bind.TransactOpts, _wellRegistry common.Address) (*types.Transaction, error) {
	return _HistoryFactory.contract.Transact(opts, "setWellRegistry", _wellRegistry)
}

// SetWellRegistry is a paid mutator transaction binding the contract method 0x16c03508.
//
// Solidity: function setWellRegistry(address _wellRegistry) returns()
func (_HistoryFactory *HistoryFactorySession) SetWellRegistry(_wellRegistry common.Address) (*types.Transaction, error) {
	return _HistoryFactory.Contract.SetWellRegistry(&_HistoryFactory.TransactOpts, _wellRegistry)
}

// SetWellRegistry is a paid mutator transaction binding the contract method 0x16c03508.
//
// Solidity: function setWellRegistry(address _wellRegistry) returns()
func (_HistoryFactory *HistoryFactoryTransactorSession) SetWellRegistry(_wellRegistry common.Address) (*types.Transaction, error) {
	return _HistoryFactory.Contract.SetWellRegistry(&_HistoryFactory.TransactOpts, _wellRegistry)
}

// HistoryFactoryHistoryCreatedIterator is returned from FilterHistoryCreated and is used to iterate over the raw logs and unpacked data for HistoryCreated events raised by the HistoryFactory contract.
type HistoryFactoryHistoryCreatedIterator struct {
	Event *HistoryFactoryHistoryCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HistoryFactoryHistoryCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoryFactoryHistoryCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HistoryFactoryHistoryCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HistoryFactoryHistoryCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoryFactoryHistoryCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoryFactoryHistoryCreated represents a HistoryCreated event raised by the HistoryFactory contract.
type HistoryFactoryHistoryCreated struct {
	Patient         common.Address
	HistoryContract common.Address
	Commitment      [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterHistoryCreated is a free log retrieval operation binding the contract event 0x6eaa6b2c062081a6d337d0c7edd3a597f1e3e735ec4a4c159b116ba0166ad8bc.
//
// Solidity: event HistoryCreated(address indexed patient, address indexed historyContract, bytes32 indexed commitment)
func (_HistoryFactory *HistoryFactoryFilterer) FilterHistoryCreated(opts *bind.FilterOpts, patient []common.Address, historyContract []common.Address, commitment [][32]byte) (*HistoryFactoryHistoryCreatedIterator, error) {

	var patientRule []interface{}
	for _, patientItem := range patient {
		patientRule = append(patientRule, patientItem)
	}
	var historyContractRule []interface{}
	for _, historyContractItem := range historyContract {
		historyContractRule = append(historyContractRule, historyContractItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}

	logs, sub, err := _HistoryFactory.contract.FilterLogs(opts, "HistoryCreated", patientRule, historyContractRule, commitmentRule)
	if err != nil {
		return nil, err
	}
	return &HistoryFactoryHistoryCreatedIterator{contract: _HistoryFactory.contract, event: "HistoryCreated", logs: logs, sub: sub}, nil
}

// WatchHistoryCreated is a free log subscription operation binding the contract event 0x6eaa6b2c062081a6d337d0c7edd3a597f1e3e735ec4a4c159b116ba0166ad8bc.
//
// Solidity: event HistoryCreated(address indexed patient, address indexed historyContract, bytes32 indexed commitment)
func (_HistoryFactory *HistoryFactoryFilterer) WatchHistoryCreated(opts *bind.WatchOpts, sink chan<- *HistoryFactoryHistoryCreated, patient []common.Address, historyContract []common.Address, commitment [][32]byte) (event.Subscription, error) {

	var patientRule []interface{}
	for _, patientItem := range patient {
		patientRule = append(patientRule, patientItem)
	}
	var historyContractRule []interface{}
	for _, historyContractItem := range historyContract {
		historyContractRule = append(historyContractRule, historyContractItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}

	logs, sub, err := _HistoryFactory.contract.WatchLogs(opts, "HistoryCreated", patientRule, historyContractRule, commitmentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoryFactoryHistoryCreated)
				if err := _HistoryFactory.contract.UnpackLog(event, "HistoryCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseHistoryCreated is a log parse operation binding the contract event 0x6eaa6b2c062081a6d337d0c7edd3a597f1e3e735ec4a4c159b116ba0166ad8bc.
//
// Solidity: event HistoryCreated(address indexed patient, address indexed historyContract, bytes32 indexed commitment)
func (_HistoryFactory *HistoryFactoryFilterer) ParseHistoryCreated(log types.Log) (*HistoryFactoryHistoryCreated, error) {
	event := new(HistoryFactoryHistoryCreated)
	if err := _HistoryFactory.contract.UnpackLog(event, "HistoryCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
