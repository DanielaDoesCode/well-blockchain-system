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

// RiskAssessmentMetaData contains all meta data concerning the RiskAssessment contract.
var RiskAssessmentMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_incidentManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_userRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_emergencySummaryFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ACCESS_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"emergencyGrant\",\"inputs\":[{\"name\":\"emergencySummary\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"locationHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"patientCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emergencySummaryFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractEmergencySummaryFactory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIncidentPatient\",\"inputs\":[{\"name\":\"incidentId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"incidentManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIncidentManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractUserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EmergencyGranted\",\"inputs\":[{\"name\":\"responder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"incidentId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"expiresAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IncidentNotActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPatient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LocationMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAssignedResponder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotFirstResponder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PatientMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// RiskAssessmentABI is the input ABI used to generate the binding from.
// Deprecated: Use RiskAssessmentMetaData.ABI instead.
var RiskAssessmentABI = RiskAssessmentMetaData.ABI

// RiskAssessment is an auto generated Go binding around an Ethereum contract.
type RiskAssessment struct {
	RiskAssessmentCaller     // Read-only binding to the contract
	RiskAssessmentTransactor // Write-only binding to the contract
	RiskAssessmentFilterer   // Log filterer for contract events
}

// RiskAssessmentCaller is an auto generated read-only Go binding around an Ethereum contract.
type RiskAssessmentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RiskAssessmentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RiskAssessmentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RiskAssessmentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RiskAssessmentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RiskAssessmentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RiskAssessmentSession struct {
	Contract     *RiskAssessment   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RiskAssessmentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RiskAssessmentCallerSession struct {
	Contract *RiskAssessmentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// RiskAssessmentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RiskAssessmentTransactorSession struct {
	Contract     *RiskAssessmentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// RiskAssessmentRaw is an auto generated low-level Go binding around an Ethereum contract.
type RiskAssessmentRaw struct {
	Contract *RiskAssessment // Generic contract binding to access the raw methods on
}

// RiskAssessmentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RiskAssessmentCallerRaw struct {
	Contract *RiskAssessmentCaller // Generic read-only contract binding to access the raw methods on
}

// RiskAssessmentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RiskAssessmentTransactorRaw struct {
	Contract *RiskAssessmentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRiskAssessment creates a new instance of RiskAssessment, bound to a specific deployed contract.
func NewRiskAssessment(address common.Address, backend bind.ContractBackend) (*RiskAssessment, error) {
	contract, err := bindRiskAssessment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RiskAssessment{RiskAssessmentCaller: RiskAssessmentCaller{contract: contract}, RiskAssessmentTransactor: RiskAssessmentTransactor{contract: contract}, RiskAssessmentFilterer: RiskAssessmentFilterer{contract: contract}}, nil
}

// NewRiskAssessmentCaller creates a new read-only instance of RiskAssessment, bound to a specific deployed contract.
func NewRiskAssessmentCaller(address common.Address, caller bind.ContractCaller) (*RiskAssessmentCaller, error) {
	contract, err := bindRiskAssessment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RiskAssessmentCaller{contract: contract}, nil
}

// NewRiskAssessmentTransactor creates a new write-only instance of RiskAssessment, bound to a specific deployed contract.
func NewRiskAssessmentTransactor(address common.Address, transactor bind.ContractTransactor) (*RiskAssessmentTransactor, error) {
	contract, err := bindRiskAssessment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RiskAssessmentTransactor{contract: contract}, nil
}

// NewRiskAssessmentFilterer creates a new log filterer instance of RiskAssessment, bound to a specific deployed contract.
func NewRiskAssessmentFilterer(address common.Address, filterer bind.ContractFilterer) (*RiskAssessmentFilterer, error) {
	contract, err := bindRiskAssessment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RiskAssessmentFilterer{contract: contract}, nil
}

// bindRiskAssessment binds a generic wrapper to an already deployed contract.
func bindRiskAssessment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RiskAssessmentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RiskAssessment *RiskAssessmentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RiskAssessment.Contract.RiskAssessmentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RiskAssessment *RiskAssessmentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RiskAssessment.Contract.RiskAssessmentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RiskAssessment *RiskAssessmentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RiskAssessment.Contract.RiskAssessmentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RiskAssessment *RiskAssessmentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RiskAssessment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RiskAssessment *RiskAssessmentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RiskAssessment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RiskAssessment *RiskAssessmentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RiskAssessment.Contract.contract.Transact(opts, method, params...)
}

// ACCESSDURATION is a free data retrieval call binding the contract method 0x63a0cce6.
//
// Solidity: function ACCESS_DURATION() view returns(uint256)
func (_RiskAssessment *RiskAssessmentCaller) ACCESSDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RiskAssessment.contract.Call(opts, &out, "ACCESS_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ACCESSDURATION is a free data retrieval call binding the contract method 0x63a0cce6.
//
// Solidity: function ACCESS_DURATION() view returns(uint256)
func (_RiskAssessment *RiskAssessmentSession) ACCESSDURATION() (*big.Int, error) {
	return _RiskAssessment.Contract.ACCESSDURATION(&_RiskAssessment.CallOpts)
}

// ACCESSDURATION is a free data retrieval call binding the contract method 0x63a0cce6.
//
// Solidity: function ACCESS_DURATION() view returns(uint256)
func (_RiskAssessment *RiskAssessmentCallerSession) ACCESSDURATION() (*big.Int, error) {
	return _RiskAssessment.Contract.ACCESSDURATION(&_RiskAssessment.CallOpts)
}

// EmergencySummaryFactory is a free data retrieval call binding the contract method 0x8229f127.
//
// Solidity: function emergencySummaryFactory() view returns(address)
func (_RiskAssessment *RiskAssessmentCaller) EmergencySummaryFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RiskAssessment.contract.Call(opts, &out, "emergencySummaryFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EmergencySummaryFactory is a free data retrieval call binding the contract method 0x8229f127.
//
// Solidity: function emergencySummaryFactory() view returns(address)
func (_RiskAssessment *RiskAssessmentSession) EmergencySummaryFactory() (common.Address, error) {
	return _RiskAssessment.Contract.EmergencySummaryFactory(&_RiskAssessment.CallOpts)
}

// EmergencySummaryFactory is a free data retrieval call binding the contract method 0x8229f127.
//
// Solidity: function emergencySummaryFactory() view returns(address)
func (_RiskAssessment *RiskAssessmentCallerSession) EmergencySummaryFactory() (common.Address, error) {
	return _RiskAssessment.Contract.EmergencySummaryFactory(&_RiskAssessment.CallOpts)
}

// GetIncidentPatient is a free data retrieval call binding the contract method 0x853e0bb1.
//
// Solidity: function getIncidentPatient(bytes32 incidentId) view returns(bytes32)
func (_RiskAssessment *RiskAssessmentCaller) GetIncidentPatient(opts *bind.CallOpts, incidentId [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _RiskAssessment.contract.Call(opts, &out, "getIncidentPatient", incidentId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetIncidentPatient is a free data retrieval call binding the contract method 0x853e0bb1.
//
// Solidity: function getIncidentPatient(bytes32 incidentId) view returns(bytes32)
func (_RiskAssessment *RiskAssessmentSession) GetIncidentPatient(incidentId [32]byte) ([32]byte, error) {
	return _RiskAssessment.Contract.GetIncidentPatient(&_RiskAssessment.CallOpts, incidentId)
}

// GetIncidentPatient is a free data retrieval call binding the contract method 0x853e0bb1.
//
// Solidity: function getIncidentPatient(bytes32 incidentId) view returns(bytes32)
func (_RiskAssessment *RiskAssessmentCallerSession) GetIncidentPatient(incidentId [32]byte) ([32]byte, error) {
	return _RiskAssessment.Contract.GetIncidentPatient(&_RiskAssessment.CallOpts, incidentId)
}

// IncidentManager is a free data retrieval call binding the contract method 0xe7030023.
//
// Solidity: function incidentManager() view returns(address)
func (_RiskAssessment *RiskAssessmentCaller) IncidentManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RiskAssessment.contract.Call(opts, &out, "incidentManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IncidentManager is a free data retrieval call binding the contract method 0xe7030023.
//
// Solidity: function incidentManager() view returns(address)
func (_RiskAssessment *RiskAssessmentSession) IncidentManager() (common.Address, error) {
	return _RiskAssessment.Contract.IncidentManager(&_RiskAssessment.CallOpts)
}

// IncidentManager is a free data retrieval call binding the contract method 0xe7030023.
//
// Solidity: function incidentManager() view returns(address)
func (_RiskAssessment *RiskAssessmentCallerSession) IncidentManager() (common.Address, error) {
	return _RiskAssessment.Contract.IncidentManager(&_RiskAssessment.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_RiskAssessment *RiskAssessmentCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RiskAssessment.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_RiskAssessment *RiskAssessmentSession) UserRegistry() (common.Address, error) {
	return _RiskAssessment.Contract.UserRegistry(&_RiskAssessment.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_RiskAssessment *RiskAssessmentCallerSession) UserRegistry() (common.Address, error) {
	return _RiskAssessment.Contract.UserRegistry(&_RiskAssessment.CallOpts)
}

// EmergencyGrant is a paid mutator transaction binding the contract method 0x62adf75d.
//
// Solidity: function emergencyGrant(address emergencySummary, bytes32 incidentId, bytes32 locationHash, bytes32 patientCommitment, bytes encryptedKey) returns()
func (_RiskAssessment *RiskAssessmentTransactor) EmergencyGrant(opts *bind.TransactOpts, emergencySummary common.Address, incidentId [32]byte, locationHash [32]byte, patientCommitment [32]byte, encryptedKey []byte) (*types.Transaction, error) {
	return _RiskAssessment.contract.Transact(opts, "emergencyGrant", emergencySummary, incidentId, locationHash, patientCommitment, encryptedKey)
}

// EmergencyGrant is a paid mutator transaction binding the contract method 0x62adf75d.
//
// Solidity: function emergencyGrant(address emergencySummary, bytes32 incidentId, bytes32 locationHash, bytes32 patientCommitment, bytes encryptedKey) returns()
func (_RiskAssessment *RiskAssessmentSession) EmergencyGrant(emergencySummary common.Address, incidentId [32]byte, locationHash [32]byte, patientCommitment [32]byte, encryptedKey []byte) (*types.Transaction, error) {
	return _RiskAssessment.Contract.EmergencyGrant(&_RiskAssessment.TransactOpts, emergencySummary, incidentId, locationHash, patientCommitment, encryptedKey)
}

// EmergencyGrant is a paid mutator transaction binding the contract method 0x62adf75d.
//
// Solidity: function emergencyGrant(address emergencySummary, bytes32 incidentId, bytes32 locationHash, bytes32 patientCommitment, bytes encryptedKey) returns()
func (_RiskAssessment *RiskAssessmentTransactorSession) EmergencyGrant(emergencySummary common.Address, incidentId [32]byte, locationHash [32]byte, patientCommitment [32]byte, encryptedKey []byte) (*types.Transaction, error) {
	return _RiskAssessment.Contract.EmergencyGrant(&_RiskAssessment.TransactOpts, emergencySummary, incidentId, locationHash, patientCommitment, encryptedKey)
}

// RiskAssessmentEmergencyGrantedIterator is returned from FilterEmergencyGranted and is used to iterate over the raw logs and unpacked data for EmergencyGranted events raised by the RiskAssessment contract.
type RiskAssessmentEmergencyGrantedIterator struct {
	Event *RiskAssessmentEmergencyGranted // Event containing the contract specifics and raw log

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
func (it *RiskAssessmentEmergencyGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RiskAssessmentEmergencyGranted)
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
		it.Event = new(RiskAssessmentEmergencyGranted)
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
func (it *RiskAssessmentEmergencyGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RiskAssessmentEmergencyGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RiskAssessmentEmergencyGranted represents a EmergencyGranted event raised by the RiskAssessment contract.
type RiskAssessmentEmergencyGranted struct {
	Responder  common.Address
	IncidentId [32]byte
	ExpiresAt  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEmergencyGranted is a free log retrieval operation binding the contract event 0x0152296501168c94de9cf6398871483d4beb834aea5803d67a9d6a0618dda379.
//
// Solidity: event EmergencyGranted(address indexed responder, bytes32 indexed incidentId, uint256 expiresAt)
func (_RiskAssessment *RiskAssessmentFilterer) FilterEmergencyGranted(opts *bind.FilterOpts, responder []common.Address, incidentId [][32]byte) (*RiskAssessmentEmergencyGrantedIterator, error) {

	var responderRule []interface{}
	for _, responderItem := range responder {
		responderRule = append(responderRule, responderItem)
	}
	var incidentIdRule []interface{}
	for _, incidentIdItem := range incidentId {
		incidentIdRule = append(incidentIdRule, incidentIdItem)
	}

	logs, sub, err := _RiskAssessment.contract.FilterLogs(opts, "EmergencyGranted", responderRule, incidentIdRule)
	if err != nil {
		return nil, err
	}
	return &RiskAssessmentEmergencyGrantedIterator{contract: _RiskAssessment.contract, event: "EmergencyGranted", logs: logs, sub: sub}, nil
}

// WatchEmergencyGranted is a free log subscription operation binding the contract event 0x0152296501168c94de9cf6398871483d4beb834aea5803d67a9d6a0618dda379.
//
// Solidity: event EmergencyGranted(address indexed responder, bytes32 indexed incidentId, uint256 expiresAt)
func (_RiskAssessment *RiskAssessmentFilterer) WatchEmergencyGranted(opts *bind.WatchOpts, sink chan<- *RiskAssessmentEmergencyGranted, responder []common.Address, incidentId [][32]byte) (event.Subscription, error) {

	var responderRule []interface{}
	for _, responderItem := range responder {
		responderRule = append(responderRule, responderItem)
	}
	var incidentIdRule []interface{}
	for _, incidentIdItem := range incidentId {
		incidentIdRule = append(incidentIdRule, incidentIdItem)
	}

	logs, sub, err := _RiskAssessment.contract.WatchLogs(opts, "EmergencyGranted", responderRule, incidentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RiskAssessmentEmergencyGranted)
				if err := _RiskAssessment.contract.UnpackLog(event, "EmergencyGranted", log); err != nil {
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

// ParseEmergencyGranted is a log parse operation binding the contract event 0x0152296501168c94de9cf6398871483d4beb834aea5803d67a9d6a0618dda379.
//
// Solidity: event EmergencyGranted(address indexed responder, bytes32 indexed incidentId, uint256 expiresAt)
func (_RiskAssessment *RiskAssessmentFilterer) ParseEmergencyGranted(log types.Log) (*RiskAssessmentEmergencyGranted, error) {
	event := new(RiskAssessmentEmergencyGranted)
	if err := _RiskAssessment.contract.UnpackLog(event, "EmergencyGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
