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

// EmergencySummaryFactoryMetaData contains all meta data concerning the EmergencySummaryFactory contract.
var EmergencySummaryFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_userRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"commitmentToEmergencySummary\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createEmergencySummary\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"patientAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getEmergencySummaryByAddress\",\"inputs\":[{\"name\":\"patientAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEmergencySummaryByCommitment\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"patientToEmergencySummary\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"riskAssessment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRiskAssessment\",\"inputs\":[{\"name\":\"_riskAssessment\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWellRegistry\",\"inputs\":[{\"name\":\"_wellRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractUserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"wellRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EmergencySummaryCreated\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"patient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"emergencySummary\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// EmergencySummaryFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use EmergencySummaryFactoryMetaData.ABI instead.
var EmergencySummaryFactoryABI = EmergencySummaryFactoryMetaData.ABI

// EmergencySummaryFactory is an auto generated Go binding around an Ethereum contract.
type EmergencySummaryFactory struct {
	EmergencySummaryFactoryCaller     // Read-only binding to the contract
	EmergencySummaryFactoryTransactor // Write-only binding to the contract
	EmergencySummaryFactoryFilterer   // Log filterer for contract events
}

// EmergencySummaryFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type EmergencySummaryFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmergencySummaryFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EmergencySummaryFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmergencySummaryFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EmergencySummaryFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmergencySummaryFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EmergencySummaryFactorySession struct {
	Contract     *EmergencySummaryFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// EmergencySummaryFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EmergencySummaryFactoryCallerSession struct {
	Contract *EmergencySummaryFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// EmergencySummaryFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EmergencySummaryFactoryTransactorSession struct {
	Contract     *EmergencySummaryFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// EmergencySummaryFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type EmergencySummaryFactoryRaw struct {
	Contract *EmergencySummaryFactory // Generic contract binding to access the raw methods on
}

// EmergencySummaryFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EmergencySummaryFactoryCallerRaw struct {
	Contract *EmergencySummaryFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// EmergencySummaryFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EmergencySummaryFactoryTransactorRaw struct {
	Contract *EmergencySummaryFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEmergencySummaryFactory creates a new instance of EmergencySummaryFactory, bound to a specific deployed contract.
func NewEmergencySummaryFactory(address common.Address, backend bind.ContractBackend) (*EmergencySummaryFactory, error) {
	contract, err := bindEmergencySummaryFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryFactory{EmergencySummaryFactoryCaller: EmergencySummaryFactoryCaller{contract: contract}, EmergencySummaryFactoryTransactor: EmergencySummaryFactoryTransactor{contract: contract}, EmergencySummaryFactoryFilterer: EmergencySummaryFactoryFilterer{contract: contract}}, nil
}

// NewEmergencySummaryFactoryCaller creates a new read-only instance of EmergencySummaryFactory, bound to a specific deployed contract.
func NewEmergencySummaryFactoryCaller(address common.Address, caller bind.ContractCaller) (*EmergencySummaryFactoryCaller, error) {
	contract, err := bindEmergencySummaryFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryFactoryCaller{contract: contract}, nil
}

// NewEmergencySummaryFactoryTransactor creates a new write-only instance of EmergencySummaryFactory, bound to a specific deployed contract.
func NewEmergencySummaryFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*EmergencySummaryFactoryTransactor, error) {
	contract, err := bindEmergencySummaryFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryFactoryTransactor{contract: contract}, nil
}

// NewEmergencySummaryFactoryFilterer creates a new log filterer instance of EmergencySummaryFactory, bound to a specific deployed contract.
func NewEmergencySummaryFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*EmergencySummaryFactoryFilterer, error) {
	contract, err := bindEmergencySummaryFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryFactoryFilterer{contract: contract}, nil
}

// bindEmergencySummaryFactory binds a generic wrapper to an already deployed contract.
func bindEmergencySummaryFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EmergencySummaryFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmergencySummaryFactory *EmergencySummaryFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmergencySummaryFactory.Contract.EmergencySummaryFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmergencySummaryFactory *EmergencySummaryFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.EmergencySummaryFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmergencySummaryFactory *EmergencySummaryFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.EmergencySummaryFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmergencySummaryFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.contract.Transact(opts, method, params...)
}

// CommitmentToEmergencySummary is a free data retrieval call binding the contract method 0x4b6867e9.
//
// Solidity: function commitmentToEmergencySummary(bytes32 ) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) CommitmentToEmergencySummary(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "commitmentToEmergencySummary", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CommitmentToEmergencySummary is a free data retrieval call binding the contract method 0x4b6867e9.
//
// Solidity: function commitmentToEmergencySummary(bytes32 ) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) CommitmentToEmergencySummary(arg0 [32]byte) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.CommitmentToEmergencySummary(&_EmergencySummaryFactory.CallOpts, arg0)
}

// CommitmentToEmergencySummary is a free data retrieval call binding the contract method 0x4b6867e9.
//
// Solidity: function commitmentToEmergencySummary(bytes32 ) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) CommitmentToEmergencySummary(arg0 [32]byte) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.CommitmentToEmergencySummary(&_EmergencySummaryFactory.CallOpts, arg0)
}

// GetEmergencySummaryByAddress is a free data retrieval call binding the contract method 0x0e49adec.
//
// Solidity: function getEmergencySummaryByAddress(address patientAddress) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) GetEmergencySummaryByAddress(opts *bind.CallOpts, patientAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "getEmergencySummaryByAddress", patientAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEmergencySummaryByAddress is a free data retrieval call binding the contract method 0x0e49adec.
//
// Solidity: function getEmergencySummaryByAddress(address patientAddress) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) GetEmergencySummaryByAddress(patientAddress common.Address) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.GetEmergencySummaryByAddress(&_EmergencySummaryFactory.CallOpts, patientAddress)
}

// GetEmergencySummaryByAddress is a free data retrieval call binding the contract method 0x0e49adec.
//
// Solidity: function getEmergencySummaryByAddress(address patientAddress) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) GetEmergencySummaryByAddress(patientAddress common.Address) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.GetEmergencySummaryByAddress(&_EmergencySummaryFactory.CallOpts, patientAddress)
}

// GetEmergencySummaryByCommitment is a free data retrieval call binding the contract method 0x1089359f.
//
// Solidity: function getEmergencySummaryByCommitment(bytes32 commitment) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) GetEmergencySummaryByCommitment(opts *bind.CallOpts, commitment [32]byte) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "getEmergencySummaryByCommitment", commitment)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEmergencySummaryByCommitment is a free data retrieval call binding the contract method 0x1089359f.
//
// Solidity: function getEmergencySummaryByCommitment(bytes32 commitment) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) GetEmergencySummaryByCommitment(commitment [32]byte) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.GetEmergencySummaryByCommitment(&_EmergencySummaryFactory.CallOpts, commitment)
}

// GetEmergencySummaryByCommitment is a free data retrieval call binding the contract method 0x1089359f.
//
// Solidity: function getEmergencySummaryByCommitment(bytes32 commitment) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) GetEmergencySummaryByCommitment(commitment [32]byte) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.GetEmergencySummaryByCommitment(&_EmergencySummaryFactory.CallOpts, commitment)
}

// PatientToEmergencySummary is a free data retrieval call binding the contract method 0x5e01fbfa.
//
// Solidity: function patientToEmergencySummary(address ) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) PatientToEmergencySummary(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "patientToEmergencySummary", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PatientToEmergencySummary is a free data retrieval call binding the contract method 0x5e01fbfa.
//
// Solidity: function patientToEmergencySummary(address ) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) PatientToEmergencySummary(arg0 common.Address) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.PatientToEmergencySummary(&_EmergencySummaryFactory.CallOpts, arg0)
}

// PatientToEmergencySummary is a free data retrieval call binding the contract method 0x5e01fbfa.
//
// Solidity: function patientToEmergencySummary(address ) view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) PatientToEmergencySummary(arg0 common.Address) (common.Address, error) {
	return _EmergencySummaryFactory.Contract.PatientToEmergencySummary(&_EmergencySummaryFactory.CallOpts, arg0)
}

// RiskAssessment is a free data retrieval call binding the contract method 0x70a1a7d5.
//
// Solidity: function riskAssessment() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) RiskAssessment(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "riskAssessment")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RiskAssessment is a free data retrieval call binding the contract method 0x70a1a7d5.
//
// Solidity: function riskAssessment() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) RiskAssessment() (common.Address, error) {
	return _EmergencySummaryFactory.Contract.RiskAssessment(&_EmergencySummaryFactory.CallOpts)
}

// RiskAssessment is a free data retrieval call binding the contract method 0x70a1a7d5.
//
// Solidity: function riskAssessment() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) RiskAssessment() (common.Address, error) {
	return _EmergencySummaryFactory.Contract.RiskAssessment(&_EmergencySummaryFactory.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) UserRegistry() (common.Address, error) {
	return _EmergencySummaryFactory.Contract.UserRegistry(&_EmergencySummaryFactory.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) UserRegistry() (common.Address, error) {
	return _EmergencySummaryFactory.Contract.UserRegistry(&_EmergencySummaryFactory.CallOpts)
}

// WellRegistry is a free data retrieval call binding the contract method 0xc6583069.
//
// Solidity: function wellRegistry() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCaller) WellRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummaryFactory.contract.Call(opts, &out, "wellRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WellRegistry is a free data retrieval call binding the contract method 0xc6583069.
//
// Solidity: function wellRegistry() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) WellRegistry() (common.Address, error) {
	return _EmergencySummaryFactory.Contract.WellRegistry(&_EmergencySummaryFactory.CallOpts)
}

// WellRegistry is a free data retrieval call binding the contract method 0xc6583069.
//
// Solidity: function wellRegistry() view returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryCallerSession) WellRegistry() (common.Address, error) {
	return _EmergencySummaryFactory.Contract.WellRegistry(&_EmergencySummaryFactory.CallOpts)
}

// CreateEmergencySummary is a paid mutator transaction binding the contract method 0xaf88b417.
//
// Solidity: function createEmergencySummary(bytes32 commitment, address patientAddress) returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactor) CreateEmergencySummary(opts *bind.TransactOpts, commitment [32]byte, patientAddress common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.contract.Transact(opts, "createEmergencySummary", commitment, patientAddress)
}

// CreateEmergencySummary is a paid mutator transaction binding the contract method 0xaf88b417.
//
// Solidity: function createEmergencySummary(bytes32 commitment, address patientAddress) returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) CreateEmergencySummary(commitment [32]byte, patientAddress common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.CreateEmergencySummary(&_EmergencySummaryFactory.TransactOpts, commitment, patientAddress)
}

// CreateEmergencySummary is a paid mutator transaction binding the contract method 0xaf88b417.
//
// Solidity: function createEmergencySummary(bytes32 commitment, address patientAddress) returns(address)
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactorSession) CreateEmergencySummary(commitment [32]byte, patientAddress common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.CreateEmergencySummary(&_EmergencySummaryFactory.TransactOpts, commitment, patientAddress)
}

// SetRiskAssessment is a paid mutator transaction binding the contract method 0x76d2feee.
//
// Solidity: function setRiskAssessment(address _riskAssessment) returns()
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactor) SetRiskAssessment(opts *bind.TransactOpts, _riskAssessment common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.contract.Transact(opts, "setRiskAssessment", _riskAssessment)
}

// SetRiskAssessment is a paid mutator transaction binding the contract method 0x76d2feee.
//
// Solidity: function setRiskAssessment(address _riskAssessment) returns()
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) SetRiskAssessment(_riskAssessment common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.SetRiskAssessment(&_EmergencySummaryFactory.TransactOpts, _riskAssessment)
}

// SetRiskAssessment is a paid mutator transaction binding the contract method 0x76d2feee.
//
// Solidity: function setRiskAssessment(address _riskAssessment) returns()
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactorSession) SetRiskAssessment(_riskAssessment common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.SetRiskAssessment(&_EmergencySummaryFactory.TransactOpts, _riskAssessment)
}

// SetWellRegistry is a paid mutator transaction binding the contract method 0x16c03508.
//
// Solidity: function setWellRegistry(address _wellRegistry) returns()
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactor) SetWellRegistry(opts *bind.TransactOpts, _wellRegistry common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.contract.Transact(opts, "setWellRegistry", _wellRegistry)
}

// SetWellRegistry is a paid mutator transaction binding the contract method 0x16c03508.
//
// Solidity: function setWellRegistry(address _wellRegistry) returns()
func (_EmergencySummaryFactory *EmergencySummaryFactorySession) SetWellRegistry(_wellRegistry common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.SetWellRegistry(&_EmergencySummaryFactory.TransactOpts, _wellRegistry)
}

// SetWellRegistry is a paid mutator transaction binding the contract method 0x16c03508.
//
// Solidity: function setWellRegistry(address _wellRegistry) returns()
func (_EmergencySummaryFactory *EmergencySummaryFactoryTransactorSession) SetWellRegistry(_wellRegistry common.Address) (*types.Transaction, error) {
	return _EmergencySummaryFactory.Contract.SetWellRegistry(&_EmergencySummaryFactory.TransactOpts, _wellRegistry)
}

// EmergencySummaryFactoryEmergencySummaryCreatedIterator is returned from FilterEmergencySummaryCreated and is used to iterate over the raw logs and unpacked data for EmergencySummaryCreated events raised by the EmergencySummaryFactory contract.
type EmergencySummaryFactoryEmergencySummaryCreatedIterator struct {
	Event *EmergencySummaryFactoryEmergencySummaryCreated // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryFactoryEmergencySummaryCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryFactoryEmergencySummaryCreated)
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
		it.Event = new(EmergencySummaryFactoryEmergencySummaryCreated)
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
func (it *EmergencySummaryFactoryEmergencySummaryCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryFactoryEmergencySummaryCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryFactoryEmergencySummaryCreated represents a EmergencySummaryCreated event raised by the EmergencySummaryFactory contract.
type EmergencySummaryFactoryEmergencySummaryCreated struct {
	Commitment       [32]byte
	Patient          common.Address
	EmergencySummary common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterEmergencySummaryCreated is a free log retrieval operation binding the contract event 0x657755e6b091b06b07b9e7963dd8bcd27a68cc4b4e5fb1922c4edf5d5053913a.
//
// Solidity: event EmergencySummaryCreated(bytes32 indexed commitment, address indexed patient, address emergencySummary)
func (_EmergencySummaryFactory *EmergencySummaryFactoryFilterer) FilterEmergencySummaryCreated(opts *bind.FilterOpts, commitment [][32]byte, patient []common.Address) (*EmergencySummaryFactoryEmergencySummaryCreatedIterator, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var patientRule []interface{}
	for _, patientItem := range patient {
		patientRule = append(patientRule, patientItem)
	}

	logs, sub, err := _EmergencySummaryFactory.contract.FilterLogs(opts, "EmergencySummaryCreated", commitmentRule, patientRule)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryFactoryEmergencySummaryCreatedIterator{contract: _EmergencySummaryFactory.contract, event: "EmergencySummaryCreated", logs: logs, sub: sub}, nil
}

// WatchEmergencySummaryCreated is a free log subscription operation binding the contract event 0x657755e6b091b06b07b9e7963dd8bcd27a68cc4b4e5fb1922c4edf5d5053913a.
//
// Solidity: event EmergencySummaryCreated(bytes32 indexed commitment, address indexed patient, address emergencySummary)
func (_EmergencySummaryFactory *EmergencySummaryFactoryFilterer) WatchEmergencySummaryCreated(opts *bind.WatchOpts, sink chan<- *EmergencySummaryFactoryEmergencySummaryCreated, commitment [][32]byte, patient []common.Address) (event.Subscription, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var patientRule []interface{}
	for _, patientItem := range patient {
		patientRule = append(patientRule, patientItem)
	}

	logs, sub, err := _EmergencySummaryFactory.contract.WatchLogs(opts, "EmergencySummaryCreated", commitmentRule, patientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryFactoryEmergencySummaryCreated)
				if err := _EmergencySummaryFactory.contract.UnpackLog(event, "EmergencySummaryCreated", log); err != nil {
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

// ParseEmergencySummaryCreated is a log parse operation binding the contract event 0x657755e6b091b06b07b9e7963dd8bcd27a68cc4b4e5fb1922c4edf5d5053913a.
//
// Solidity: event EmergencySummaryCreated(bytes32 indexed commitment, address indexed patient, address emergencySummary)
func (_EmergencySummaryFactory *EmergencySummaryFactoryFilterer) ParseEmergencySummaryCreated(log types.Log) (*EmergencySummaryFactoryEmergencySummaryCreated, error) {
	event := new(EmergencySummaryFactoryEmergencySummaryCreated)
	if err := _EmergencySummaryFactory.contract.UnpackLog(event, "EmergencySummaryCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
