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

// TypesIncident is an auto generated low-level Go binding around an user-defined struct.
type TypesIncident struct {
	Id           [32]byte
	Hospital     common.Address
	LocationHash [32]byte
	CreatedAt    *big.Int
	ExpiresAt    *big.Int
	Active       bool
}

// IncidentManagerMetaData contains all meta data concerning the IncidentManager contract.
var IncidentManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_userRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_INCIDENT_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addResponders\",\"inputs\":[{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"responders\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createIncident\",\"inputs\":[{\"name\":\"locationHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deactivateIncident\",\"inputs\":[{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getIncident\",\"inputs\":[{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTypes.Incident\",\"components\":[{\"name\":\"id\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"hospital\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"locationHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"createdAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIncidentCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isIncidentActive\",\"inputs\":[{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isResponderInIncident\",\"inputs\":[{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"responder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractUserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"IncidentCreated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"hospital\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IncidentDeactivated\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RespondersAdded\",\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"responders\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IncidentExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncidentNotActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidIncident\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// IncidentManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use IncidentManagerMetaData.ABI instead.
var IncidentManagerABI = IncidentManagerMetaData.ABI

// IncidentManager is an auto generated Go binding around an Ethereum contract.
type IncidentManager struct {
	IncidentManagerCaller     // Read-only binding to the contract
	IncidentManagerTransactor // Write-only binding to the contract
	IncidentManagerFilterer   // Log filterer for contract events
}

// IncidentManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IncidentManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IncidentManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IncidentManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IncidentManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IncidentManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IncidentManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IncidentManagerSession struct {
	Contract     *IncidentManager  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IncidentManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IncidentManagerCallerSession struct {
	Contract *IncidentManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IncidentManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IncidentManagerTransactorSession struct {
	Contract     *IncidentManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IncidentManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IncidentManagerRaw struct {
	Contract *IncidentManager // Generic contract binding to access the raw methods on
}

// IncidentManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IncidentManagerCallerRaw struct {
	Contract *IncidentManagerCaller // Generic read-only contract binding to access the raw methods on
}

// IncidentManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IncidentManagerTransactorRaw struct {
	Contract *IncidentManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIncidentManager creates a new instance of IncidentManager, bound to a specific deployed contract.
func NewIncidentManager(address common.Address, backend bind.ContractBackend) (*IncidentManager, error) {
	contract, err := bindIncidentManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IncidentManager{IncidentManagerCaller: IncidentManagerCaller{contract: contract}, IncidentManagerTransactor: IncidentManagerTransactor{contract: contract}, IncidentManagerFilterer: IncidentManagerFilterer{contract: contract}}, nil
}

// NewIncidentManagerCaller creates a new read-only instance of IncidentManager, bound to a specific deployed contract.
func NewIncidentManagerCaller(address common.Address, caller bind.ContractCaller) (*IncidentManagerCaller, error) {
	contract, err := bindIncidentManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IncidentManagerCaller{contract: contract}, nil
}

// NewIncidentManagerTransactor creates a new write-only instance of IncidentManager, bound to a specific deployed contract.
func NewIncidentManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*IncidentManagerTransactor, error) {
	contract, err := bindIncidentManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IncidentManagerTransactor{contract: contract}, nil
}

// NewIncidentManagerFilterer creates a new log filterer instance of IncidentManager, bound to a specific deployed contract.
func NewIncidentManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*IncidentManagerFilterer, error) {
	contract, err := bindIncidentManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IncidentManagerFilterer{contract: contract}, nil
}

// bindIncidentManager binds a generic wrapper to an already deployed contract.
func bindIncidentManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IncidentManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IncidentManager *IncidentManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IncidentManager.Contract.IncidentManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IncidentManager *IncidentManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IncidentManager.Contract.IncidentManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IncidentManager *IncidentManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IncidentManager.Contract.IncidentManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IncidentManager *IncidentManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IncidentManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IncidentManager *IncidentManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IncidentManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IncidentManager *IncidentManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IncidentManager.Contract.contract.Transact(opts, method, params...)
}

// MAXINCIDENTDURATION is a free data retrieval call binding the contract method 0xa4859da8.
//
// Solidity: function MAX_INCIDENT_DURATION() view returns(uint256)
func (_IncidentManager *IncidentManagerCaller) MAXINCIDENTDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IncidentManager.contract.Call(opts, &out, "MAX_INCIDENT_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXINCIDENTDURATION is a free data retrieval call binding the contract method 0xa4859da8.
//
// Solidity: function MAX_INCIDENT_DURATION() view returns(uint256)
func (_IncidentManager *IncidentManagerSession) MAXINCIDENTDURATION() (*big.Int, error) {
	return _IncidentManager.Contract.MAXINCIDENTDURATION(&_IncidentManager.CallOpts)
}

// MAXINCIDENTDURATION is a free data retrieval call binding the contract method 0xa4859da8.
//
// Solidity: function MAX_INCIDENT_DURATION() view returns(uint256)
func (_IncidentManager *IncidentManagerCallerSession) MAXINCIDENTDURATION() (*big.Int, error) {
	return _IncidentManager.Contract.MAXINCIDENTDURATION(&_IncidentManager.CallOpts)
}

// GetIncident is a free data retrieval call binding the contract method 0xb8f03828.
//
// Solidity: function getIncident(bytes32 incidentId) view returns((bytes32,address,bytes32,uint48,uint48,bool))
func (_IncidentManager *IncidentManagerCaller) GetIncident(opts *bind.CallOpts, incidentId [32]byte) (TypesIncident, error) {
	var out []interface{}
	err := _IncidentManager.contract.Call(opts, &out, "getIncident", incidentId)

	if err != nil {
		return *new(TypesIncident), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesIncident)).(*TypesIncident)

	return out0, err

}

// GetIncident is a free data retrieval call binding the contract method 0xb8f03828.
//
// Solidity: function getIncident(bytes32 incidentId) view returns((bytes32,address,bytes32,uint48,uint48,bool))
func (_IncidentManager *IncidentManagerSession) GetIncident(incidentId [32]byte) (TypesIncident, error) {
	return _IncidentManager.Contract.GetIncident(&_IncidentManager.CallOpts, incidentId)
}

// GetIncident is a free data retrieval call binding the contract method 0xb8f03828.
//
// Solidity: function getIncident(bytes32 incidentId) view returns((bytes32,address,bytes32,uint48,uint48,bool))
func (_IncidentManager *IncidentManagerCallerSession) GetIncident(incidentId [32]byte) (TypesIncident, error) {
	return _IncidentManager.Contract.GetIncident(&_IncidentManager.CallOpts, incidentId)
}

// GetIncidentCount is a free data retrieval call binding the contract method 0x95962692.
//
// Solidity: function getIncidentCount() view returns(uint256)
func (_IncidentManager *IncidentManagerCaller) GetIncidentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IncidentManager.contract.Call(opts, &out, "getIncidentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetIncidentCount is a free data retrieval call binding the contract method 0x95962692.
//
// Solidity: function getIncidentCount() view returns(uint256)
func (_IncidentManager *IncidentManagerSession) GetIncidentCount() (*big.Int, error) {
	return _IncidentManager.Contract.GetIncidentCount(&_IncidentManager.CallOpts)
}

// GetIncidentCount is a free data retrieval call binding the contract method 0x95962692.
//
// Solidity: function getIncidentCount() view returns(uint256)
func (_IncidentManager *IncidentManagerCallerSession) GetIncidentCount() (*big.Int, error) {
	return _IncidentManager.Contract.GetIncidentCount(&_IncidentManager.CallOpts)
}

// IsIncidentActive is a free data retrieval call binding the contract method 0xb7d2c631.
//
// Solidity: function isIncidentActive(bytes32 incidentId) view returns(bool)
func (_IncidentManager *IncidentManagerCaller) IsIncidentActive(opts *bind.CallOpts, incidentId [32]byte) (bool, error) {
	var out []interface{}
	err := _IncidentManager.contract.Call(opts, &out, "isIncidentActive", incidentId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsIncidentActive is a free data retrieval call binding the contract method 0xb7d2c631.
//
// Solidity: function isIncidentActive(bytes32 incidentId) view returns(bool)
func (_IncidentManager *IncidentManagerSession) IsIncidentActive(incidentId [32]byte) (bool, error) {
	return _IncidentManager.Contract.IsIncidentActive(&_IncidentManager.CallOpts, incidentId)
}

// IsIncidentActive is a free data retrieval call binding the contract method 0xb7d2c631.
//
// Solidity: function isIncidentActive(bytes32 incidentId) view returns(bool)
func (_IncidentManager *IncidentManagerCallerSession) IsIncidentActive(incidentId [32]byte) (bool, error) {
	return _IncidentManager.Contract.IsIncidentActive(&_IncidentManager.CallOpts, incidentId)
}

// IsResponderInIncident is a free data retrieval call binding the contract method 0x79812168.
//
// Solidity: function isResponderInIncident(bytes32 incidentId, address responder) view returns(bool)
func (_IncidentManager *IncidentManagerCaller) IsResponderInIncident(opts *bind.CallOpts, incidentId [32]byte, responder common.Address) (bool, error) {
	var out []interface{}
	err := _IncidentManager.contract.Call(opts, &out, "isResponderInIncident", incidentId, responder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsResponderInIncident is a free data retrieval call binding the contract method 0x79812168.
//
// Solidity: function isResponderInIncident(bytes32 incidentId, address responder) view returns(bool)
func (_IncidentManager *IncidentManagerSession) IsResponderInIncident(incidentId [32]byte, responder common.Address) (bool, error) {
	return _IncidentManager.Contract.IsResponderInIncident(&_IncidentManager.CallOpts, incidentId, responder)
}

// IsResponderInIncident is a free data retrieval call binding the contract method 0x79812168.
//
// Solidity: function isResponderInIncident(bytes32 incidentId, address responder) view returns(bool)
func (_IncidentManager *IncidentManagerCallerSession) IsResponderInIncident(incidentId [32]byte, responder common.Address) (bool, error) {
	return _IncidentManager.Contract.IsResponderInIncident(&_IncidentManager.CallOpts, incidentId, responder)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_IncidentManager *IncidentManagerCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IncidentManager.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_IncidentManager *IncidentManagerSession) UserRegistry() (common.Address, error) {
	return _IncidentManager.Contract.UserRegistry(&_IncidentManager.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_IncidentManager *IncidentManagerCallerSession) UserRegistry() (common.Address, error) {
	return _IncidentManager.Contract.UserRegistry(&_IncidentManager.CallOpts)
}

// AddResponders is a paid mutator transaction binding the contract method 0x998b2ab0.
//
// Solidity: function addResponders(bytes32 incidentId, address[] responders) returns()
func (_IncidentManager *IncidentManagerTransactor) AddResponders(opts *bind.TransactOpts, incidentId [32]byte, responders []common.Address) (*types.Transaction, error) {
	return _IncidentManager.contract.Transact(opts, "addResponders", incidentId, responders)
}

// AddResponders is a paid mutator transaction binding the contract method 0x998b2ab0.
//
// Solidity: function addResponders(bytes32 incidentId, address[] responders) returns()
func (_IncidentManager *IncidentManagerSession) AddResponders(incidentId [32]byte, responders []common.Address) (*types.Transaction, error) {
	return _IncidentManager.Contract.AddResponders(&_IncidentManager.TransactOpts, incidentId, responders)
}

// AddResponders is a paid mutator transaction binding the contract method 0x998b2ab0.
//
// Solidity: function addResponders(bytes32 incidentId, address[] responders) returns()
func (_IncidentManager *IncidentManagerTransactorSession) AddResponders(incidentId [32]byte, responders []common.Address) (*types.Transaction, error) {
	return _IncidentManager.Contract.AddResponders(&_IncidentManager.TransactOpts, incidentId, responders)
}

// CreateIncident is a paid mutator transaction binding the contract method 0x034da2b8.
//
// Solidity: function createIncident(bytes32 locationHash, uint48 expiresAt) returns(bytes32)
func (_IncidentManager *IncidentManagerTransactor) CreateIncident(opts *bind.TransactOpts, locationHash [32]byte, expiresAt *big.Int) (*types.Transaction, error) {
	return _IncidentManager.contract.Transact(opts, "createIncident", locationHash, expiresAt)
}

// CreateIncident is a paid mutator transaction binding the contract method 0x034da2b8.
//
// Solidity: function createIncident(bytes32 locationHash, uint48 expiresAt) returns(bytes32)
func (_IncidentManager *IncidentManagerSession) CreateIncident(locationHash [32]byte, expiresAt *big.Int) (*types.Transaction, error) {
	return _IncidentManager.Contract.CreateIncident(&_IncidentManager.TransactOpts, locationHash, expiresAt)
}

// CreateIncident is a paid mutator transaction binding the contract method 0x034da2b8.
//
// Solidity: function createIncident(bytes32 locationHash, uint48 expiresAt) returns(bytes32)
func (_IncidentManager *IncidentManagerTransactorSession) CreateIncident(locationHash [32]byte, expiresAt *big.Int) (*types.Transaction, error) {
	return _IncidentManager.Contract.CreateIncident(&_IncidentManager.TransactOpts, locationHash, expiresAt)
}

// DeactivateIncident is a paid mutator transaction binding the contract method 0x5c81d45f.
//
// Solidity: function deactivateIncident(bytes32 incidentId) returns()
func (_IncidentManager *IncidentManagerTransactor) DeactivateIncident(opts *bind.TransactOpts, incidentId [32]byte) (*types.Transaction, error) {
	return _IncidentManager.contract.Transact(opts, "deactivateIncident", incidentId)
}

// DeactivateIncident is a paid mutator transaction binding the contract method 0x5c81d45f.
//
// Solidity: function deactivateIncident(bytes32 incidentId) returns()
func (_IncidentManager *IncidentManagerSession) DeactivateIncident(incidentId [32]byte) (*types.Transaction, error) {
	return _IncidentManager.Contract.DeactivateIncident(&_IncidentManager.TransactOpts, incidentId)
}

// DeactivateIncident is a paid mutator transaction binding the contract method 0x5c81d45f.
//
// Solidity: function deactivateIncident(bytes32 incidentId) returns()
func (_IncidentManager *IncidentManagerTransactorSession) DeactivateIncident(incidentId [32]byte) (*types.Transaction, error) {
	return _IncidentManager.Contract.DeactivateIncident(&_IncidentManager.TransactOpts, incidentId)
}

// IncidentManagerIncidentCreatedIterator is returned from FilterIncidentCreated and is used to iterate over the raw logs and unpacked data for IncidentCreated events raised by the IncidentManager contract.
type IncidentManagerIncidentCreatedIterator struct {
	Event *IncidentManagerIncidentCreated // Event containing the contract specifics and raw log

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
func (it *IncidentManagerIncidentCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IncidentManagerIncidentCreated)
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
		it.Event = new(IncidentManagerIncidentCreated)
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
func (it *IncidentManagerIncidentCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IncidentManagerIncidentCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IncidentManagerIncidentCreated represents a IncidentCreated event raised by the IncidentManager contract.
type IncidentManagerIncidentCreated struct {
	Id        [32]byte
	Hospital  common.Address
	ExpiresAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIncidentCreated is a free log retrieval operation binding the contract event 0x53596902d4ac2bd51f676322e33c940130b9d3ce0646613995b3844cc3a20e34.
//
// Solidity: event IncidentCreated(bytes32 indexed id, address indexed hospital, uint48 expiresAt)
func (_IncidentManager *IncidentManagerFilterer) FilterIncidentCreated(opts *bind.FilterOpts, id [][32]byte, hospital []common.Address) (*IncidentManagerIncidentCreatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var hospitalRule []interface{}
	for _, hospitalItem := range hospital {
		hospitalRule = append(hospitalRule, hospitalItem)
	}

	logs, sub, err := _IncidentManager.contract.FilterLogs(opts, "IncidentCreated", idRule, hospitalRule)
	if err != nil {
		return nil, err
	}
	return &IncidentManagerIncidentCreatedIterator{contract: _IncidentManager.contract, event: "IncidentCreated", logs: logs, sub: sub}, nil
}

// WatchIncidentCreated is a free log subscription operation binding the contract event 0x53596902d4ac2bd51f676322e33c940130b9d3ce0646613995b3844cc3a20e34.
//
// Solidity: event IncidentCreated(bytes32 indexed id, address indexed hospital, uint48 expiresAt)
func (_IncidentManager *IncidentManagerFilterer) WatchIncidentCreated(opts *bind.WatchOpts, sink chan<- *IncidentManagerIncidentCreated, id [][32]byte, hospital []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var hospitalRule []interface{}
	for _, hospitalItem := range hospital {
		hospitalRule = append(hospitalRule, hospitalItem)
	}

	logs, sub, err := _IncidentManager.contract.WatchLogs(opts, "IncidentCreated", idRule, hospitalRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IncidentManagerIncidentCreated)
				if err := _IncidentManager.contract.UnpackLog(event, "IncidentCreated", log); err != nil {
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

// ParseIncidentCreated is a log parse operation binding the contract event 0x53596902d4ac2bd51f676322e33c940130b9d3ce0646613995b3844cc3a20e34.
//
// Solidity: event IncidentCreated(bytes32 indexed id, address indexed hospital, uint48 expiresAt)
func (_IncidentManager *IncidentManagerFilterer) ParseIncidentCreated(log types.Log) (*IncidentManagerIncidentCreated, error) {
	event := new(IncidentManagerIncidentCreated)
	if err := _IncidentManager.contract.UnpackLog(event, "IncidentCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IncidentManagerIncidentDeactivatedIterator is returned from FilterIncidentDeactivated and is used to iterate over the raw logs and unpacked data for IncidentDeactivated events raised by the IncidentManager contract.
type IncidentManagerIncidentDeactivatedIterator struct {
	Event *IncidentManagerIncidentDeactivated // Event containing the contract specifics and raw log

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
func (it *IncidentManagerIncidentDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IncidentManagerIncidentDeactivated)
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
		it.Event = new(IncidentManagerIncidentDeactivated)
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
func (it *IncidentManagerIncidentDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IncidentManagerIncidentDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IncidentManagerIncidentDeactivated represents a IncidentDeactivated event raised by the IncidentManager contract.
type IncidentManagerIncidentDeactivated struct {
	Id  [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterIncidentDeactivated is a free log retrieval operation binding the contract event 0xfdee85bbf887847652b152da3f3e1b7bf0ec162e10b20b8db1e3cdb71bedf29a.
//
// Solidity: event IncidentDeactivated(bytes32 indexed id)
func (_IncidentManager *IncidentManagerFilterer) FilterIncidentDeactivated(opts *bind.FilterOpts, id [][32]byte) (*IncidentManagerIncidentDeactivatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IncidentManager.contract.FilterLogs(opts, "IncidentDeactivated", idRule)
	if err != nil {
		return nil, err
	}
	return &IncidentManagerIncidentDeactivatedIterator{contract: _IncidentManager.contract, event: "IncidentDeactivated", logs: logs, sub: sub}, nil
}

// WatchIncidentDeactivated is a free log subscription operation binding the contract event 0xfdee85bbf887847652b152da3f3e1b7bf0ec162e10b20b8db1e3cdb71bedf29a.
//
// Solidity: event IncidentDeactivated(bytes32 indexed id)
func (_IncidentManager *IncidentManagerFilterer) WatchIncidentDeactivated(opts *bind.WatchOpts, sink chan<- *IncidentManagerIncidentDeactivated, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IncidentManager.contract.WatchLogs(opts, "IncidentDeactivated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IncidentManagerIncidentDeactivated)
				if err := _IncidentManager.contract.UnpackLog(event, "IncidentDeactivated", log); err != nil {
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

// ParseIncidentDeactivated is a log parse operation binding the contract event 0xfdee85bbf887847652b152da3f3e1b7bf0ec162e10b20b8db1e3cdb71bedf29a.
//
// Solidity: event IncidentDeactivated(bytes32 indexed id)
func (_IncidentManager *IncidentManagerFilterer) ParseIncidentDeactivated(log types.Log) (*IncidentManagerIncidentDeactivated, error) {
	event := new(IncidentManagerIncidentDeactivated)
	if err := _IncidentManager.contract.UnpackLog(event, "IncidentDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IncidentManagerRespondersAddedIterator is returned from FilterRespondersAdded and is used to iterate over the raw logs and unpacked data for RespondersAdded events raised by the IncidentManager contract.
type IncidentManagerRespondersAddedIterator struct {
	Event *IncidentManagerRespondersAdded // Event containing the contract specifics and raw log

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
func (it *IncidentManagerRespondersAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IncidentManagerRespondersAdded)
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
		it.Event = new(IncidentManagerRespondersAdded)
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
func (it *IncidentManagerRespondersAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IncidentManagerRespondersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IncidentManagerRespondersAdded represents a RespondersAdded event raised by the IncidentManager contract.
type IncidentManagerRespondersAdded struct {
	Id         [32]byte
	Responders []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRespondersAdded is a free log retrieval operation binding the contract event 0xf04417362c6ded240de171f1d2bdba69a886c3675520871b7d510cbe108e75e0.
//
// Solidity: event RespondersAdded(bytes32 indexed id, address[] responders)
func (_IncidentManager *IncidentManagerFilterer) FilterRespondersAdded(opts *bind.FilterOpts, id [][32]byte) (*IncidentManagerRespondersAddedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IncidentManager.contract.FilterLogs(opts, "RespondersAdded", idRule)
	if err != nil {
		return nil, err
	}
	return &IncidentManagerRespondersAddedIterator{contract: _IncidentManager.contract, event: "RespondersAdded", logs: logs, sub: sub}, nil
}

// WatchRespondersAdded is a free log subscription operation binding the contract event 0xf04417362c6ded240de171f1d2bdba69a886c3675520871b7d510cbe108e75e0.
//
// Solidity: event RespondersAdded(bytes32 indexed id, address[] responders)
func (_IncidentManager *IncidentManagerFilterer) WatchRespondersAdded(opts *bind.WatchOpts, sink chan<- *IncidentManagerRespondersAdded, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _IncidentManager.contract.WatchLogs(opts, "RespondersAdded", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IncidentManagerRespondersAdded)
				if err := _IncidentManager.contract.UnpackLog(event, "RespondersAdded", log); err != nil {
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

// ParseRespondersAdded is a log parse operation binding the contract event 0xf04417362c6ded240de171f1d2bdba69a886c3675520871b7d510cbe108e75e0.
//
// Solidity: event RespondersAdded(bytes32 indexed id, address[] responders)
func (_IncidentManager *IncidentManagerFilterer) ParseRespondersAdded(log types.Log) (*IncidentManagerRespondersAdded, error) {
	event := new(IncidentManagerRespondersAdded)
	if err := _IncidentManager.contract.UnpackLog(event, "RespondersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
