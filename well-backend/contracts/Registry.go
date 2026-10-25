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

// TypesRecordMetadata is an auto generated low-level Go binding around an user-defined struct.


// WellRegistryUser is an auto generated low-level Go binding around an user-defined struct.
type WellRegistryUser struct {
	Commitment      [32]byte
	InternalID      [32]byte
	Role            uint8
	UserAddress     common.Address
	HistoryContract common.Address
}

// RegistryMetaData contains all meta data concerning the Registry contract.
var RegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_userRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_historyFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createRecord\",\"inputs\":[{\"name\":\"uploader\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"metadata\",\"type\":\"tuple\",\"internalType\":\"structTypes.RecordMetadata\",\"components\":[{\"name\":\"recordType\",\"type\":\"uint8\",\"internalType\":\"enumTypes.RecordType\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"contentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"storageReference\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractRecord\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emergencySummaryFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractEmergencySummaryFactory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"firstAdminRegistered\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentByInternalID\",\"inputs\":[{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEmergencySummaryByAddress\",\"inputs\":[{\"name\":\"patientAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHistoryContract\",\"inputs\":[{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHistoryContractByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIncidentManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRiskAssessment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRole\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUser\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structWellRegistry.User\",\"components\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"},{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"historyContract\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structWellRegistry.User\",\"components\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"},{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"historyContract\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"historyFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractHistoryFactory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"incidentManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"outputs\":[{\"name\":\"historyContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerByAdmin\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"outputs\":[{\"name\":\"historyContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revoke\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"riskAssessment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setEmergencySummaryFactory\",\"inputs\":[{\"name\":\"_emergencySummaryFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIncidentManager\",\"inputs\":[{\"name\":\"_incidentManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRiskAssessment\",\"inputs\":[{\"name\":\"_riskAssessment\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractUserRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"AdminAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoSelfRegistrationDoctor\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// RegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistryMetaData.ABI instead.
var RegistryABI = RegistryMetaData.ABI

// Registry is an auto generated Go binding around an Ethereum contract.
type Registry struct {
	RegistryCaller     // Read-only binding to the contract
	RegistryTransactor // Write-only binding to the contract
	RegistryFilterer   // Log filterer for contract events
}

// RegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistrySession struct {
	Contract     *Registry         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistryCallerSession struct {
	Contract *RegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistryTransactorSession struct {
	Contract     *RegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistryRaw struct {
	Contract *Registry // Generic contract binding to access the raw methods on
}

// RegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistryCallerRaw struct {
	Contract *RegistryCaller // Generic read-only contract binding to access the raw methods on
}

// RegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistryTransactorRaw struct {
	Contract *RegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistry creates a new instance of Registry, bound to a specific deployed contract.
func NewRegistry(address common.Address, backend bind.ContractBackend) (*Registry, error) {
	contract, err := bindRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registry{RegistryCaller: RegistryCaller{contract: contract}, RegistryTransactor: RegistryTransactor{contract: contract}, RegistryFilterer: RegistryFilterer{contract: contract}}, nil
}

// NewRegistryCaller creates a new read-only instance of Registry, bound to a specific deployed contract.
func NewRegistryCaller(address common.Address, caller bind.ContractCaller) (*RegistryCaller, error) {
	contract, err := bindRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryCaller{contract: contract}, nil
}

// NewRegistryTransactor creates a new write-only instance of Registry, bound to a specific deployed contract.
func NewRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistryTransactor, error) {
	contract, err := bindRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistryTransactor{contract: contract}, nil
}

// NewRegistryFilterer creates a new log filterer instance of Registry, bound to a specific deployed contract.
func NewRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistryFilterer, error) {
	contract, err := bindRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistryFilterer{contract: contract}, nil
}

// bindRegistry binds a generic wrapper to an already deployed contract.
func bindRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.RegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.RegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registry *RegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registry *RegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registry *RegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registry.Contract.contract.Transact(opts, method, params...)
}

// EmergencySummaryFactory is a free data retrieval call binding the contract method 0x8229f127.
//
// Solidity: function emergencySummaryFactory() view returns(address)
func (_Registry *RegistryCaller) EmergencySummaryFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "emergencySummaryFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EmergencySummaryFactory is a free data retrieval call binding the contract method 0x8229f127.
//
// Solidity: function emergencySummaryFactory() view returns(address)
func (_Registry *RegistrySession) EmergencySummaryFactory() (common.Address, error) {
	return _Registry.Contract.EmergencySummaryFactory(&_Registry.CallOpts)
}

// EmergencySummaryFactory is a free data retrieval call binding the contract method 0x8229f127.
//
// Solidity: function emergencySummaryFactory() view returns(address)
func (_Registry *RegistryCallerSession) EmergencySummaryFactory() (common.Address, error) {
	return _Registry.Contract.EmergencySummaryFactory(&_Registry.CallOpts)
}

// FirstAdminRegistered is a free data retrieval call binding the contract method 0xb9ff5576.
//
// Solidity: function firstAdminRegistered() view returns(bool)
func (_Registry *RegistryCaller) FirstAdminRegistered(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "firstAdminRegistered")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FirstAdminRegistered is a free data retrieval call binding the contract method 0xb9ff5576.
//
// Solidity: function firstAdminRegistered() view returns(bool)
func (_Registry *RegistrySession) FirstAdminRegistered() (bool, error) {
	return _Registry.Contract.FirstAdminRegistered(&_Registry.CallOpts)
}

// FirstAdminRegistered is a free data retrieval call binding the contract method 0xb9ff5576.
//
// Solidity: function firstAdminRegistered() view returns(bool)
func (_Registry *RegistryCallerSession) FirstAdminRegistered() (bool, error) {
	return _Registry.Contract.FirstAdminRegistered(&_Registry.CallOpts)
}

// GetCommitmentByAddress is a free data retrieval call binding the contract method 0xaf0e9c8d.
//
// Solidity: function getCommitmentByAddress(address userAddress) view returns(bytes32 commitment)
func (_Registry *RegistryCaller) GetCommitmentByAddress(opts *bind.CallOpts, userAddress common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getCommitmentByAddress", userAddress)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentByAddress is a free data retrieval call binding the contract method 0xaf0e9c8d.
//
// Solidity: function getCommitmentByAddress(address userAddress) view returns(bytes32 commitment)
func (_Registry *RegistrySession) GetCommitmentByAddress(userAddress common.Address) ([32]byte, error) {
	return _Registry.Contract.GetCommitmentByAddress(&_Registry.CallOpts, userAddress)
}

// GetCommitmentByAddress is a free data retrieval call binding the contract method 0xaf0e9c8d.
//
// Solidity: function getCommitmentByAddress(address userAddress) view returns(bytes32 commitment)
func (_Registry *RegistryCallerSession) GetCommitmentByAddress(userAddress common.Address) ([32]byte, error) {
	return _Registry.Contract.GetCommitmentByAddress(&_Registry.CallOpts, userAddress)
}

// GetCommitmentByInternalID is a free data retrieval call binding the contract method 0x632849ed.
//
// Solidity: function getCommitmentByInternalID(bytes32 internalID) view returns(bytes32 commitment)
func (_Registry *RegistryCaller) GetCommitmentByInternalID(opts *bind.CallOpts, internalID [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getCommitmentByInternalID", internalID)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentByInternalID is a free data retrieval call binding the contract method 0x632849ed.
//
// Solidity: function getCommitmentByInternalID(bytes32 internalID) view returns(bytes32 commitment)
func (_Registry *RegistrySession) GetCommitmentByInternalID(internalID [32]byte) ([32]byte, error) {
	return _Registry.Contract.GetCommitmentByInternalID(&_Registry.CallOpts, internalID)
}

// GetCommitmentByInternalID is a free data retrieval call binding the contract method 0x632849ed.
//
// Solidity: function getCommitmentByInternalID(bytes32 internalID) view returns(bytes32 commitment)
func (_Registry *RegistryCallerSession) GetCommitmentByInternalID(internalID [32]byte) ([32]byte, error) {
	return _Registry.Contract.GetCommitmentByInternalID(&_Registry.CallOpts, internalID)
}

// GetEmergencySummaryByAddress is a free data retrieval call binding the contract method 0x0e49adec.
//
// Solidity: function getEmergencySummaryByAddress(address patientAddress) view returns(address)
func (_Registry *RegistryCaller) GetEmergencySummaryByAddress(opts *bind.CallOpts, patientAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getEmergencySummaryByAddress", patientAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEmergencySummaryByAddress is a free data retrieval call binding the contract method 0x0e49adec.
//
// Solidity: function getEmergencySummaryByAddress(address patientAddress) view returns(address)
func (_Registry *RegistrySession) GetEmergencySummaryByAddress(patientAddress common.Address) (common.Address, error) {
	return _Registry.Contract.GetEmergencySummaryByAddress(&_Registry.CallOpts, patientAddress)
}

// GetEmergencySummaryByAddress is a free data retrieval call binding the contract method 0x0e49adec.
//
// Solidity: function getEmergencySummaryByAddress(address patientAddress) view returns(address)
func (_Registry *RegistryCallerSession) GetEmergencySummaryByAddress(patientAddress common.Address) (common.Address, error) {
	return _Registry.Contract.GetEmergencySummaryByAddress(&_Registry.CallOpts, patientAddress)
}

// GetHistoryContract is a free data retrieval call binding the contract method 0xafc3b238.
//
// Solidity: function getHistoryContract(bytes32 internalID) view returns(address)
func (_Registry *RegistryCaller) GetHistoryContract(opts *bind.CallOpts, internalID [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getHistoryContract", internalID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetHistoryContract is a free data retrieval call binding the contract method 0xafc3b238.
//
// Solidity: function getHistoryContract(bytes32 internalID) view returns(address)
func (_Registry *RegistrySession) GetHistoryContract(internalID [32]byte) (common.Address, error) {
	return _Registry.Contract.GetHistoryContract(&_Registry.CallOpts, internalID)
}

// GetHistoryContract is a free data retrieval call binding the contract method 0xafc3b238.
//
// Solidity: function getHistoryContract(bytes32 internalID) view returns(address)
func (_Registry *RegistryCallerSession) GetHistoryContract(internalID [32]byte) (common.Address, error) {
	return _Registry.Contract.GetHistoryContract(&_Registry.CallOpts, internalID)
}

// GetHistoryContractByAddress is a free data retrieval call binding the contract method 0x145f6705.
//
// Solidity: function getHistoryContractByAddress(address userAddress) view returns(address)
func (_Registry *RegistryCaller) GetHistoryContractByAddress(opts *bind.CallOpts, userAddress common.Address) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getHistoryContractByAddress", userAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetHistoryContractByAddress is a free data retrieval call binding the contract method 0x145f6705.
//
// Solidity: function getHistoryContractByAddress(address userAddress) view returns(address)
func (_Registry *RegistrySession) GetHistoryContractByAddress(userAddress common.Address) (common.Address, error) {
	return _Registry.Contract.GetHistoryContractByAddress(&_Registry.CallOpts, userAddress)
}

// GetHistoryContractByAddress is a free data retrieval call binding the contract method 0x145f6705.
//
// Solidity: function getHistoryContractByAddress(address userAddress) view returns(address)
func (_Registry *RegistryCallerSession) GetHistoryContractByAddress(userAddress common.Address) (common.Address, error) {
	return _Registry.Contract.GetHistoryContractByAddress(&_Registry.CallOpts, userAddress)
}

// GetIncidentManager is a free data retrieval call binding the contract method 0x417048ae.
//
// Solidity: function getIncidentManager() view returns(address)
func (_Registry *RegistryCaller) GetIncidentManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getIncidentManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetIncidentManager is a free data retrieval call binding the contract method 0x417048ae.
//
// Solidity: function getIncidentManager() view returns(address)
func (_Registry *RegistrySession) GetIncidentManager() (common.Address, error) {
	return _Registry.Contract.GetIncidentManager(&_Registry.CallOpts)
}

// GetIncidentManager is a free data retrieval call binding the contract method 0x417048ae.
//
// Solidity: function getIncidentManager() view returns(address)
func (_Registry *RegistryCallerSession) GetIncidentManager() (common.Address, error) {
	return _Registry.Contract.GetIncidentManager(&_Registry.CallOpts)
}

// GetRiskAssessment is a free data retrieval call binding the contract method 0x3acd6d50.
//
// Solidity: function getRiskAssessment() view returns(address)
func (_Registry *RegistryCaller) GetRiskAssessment(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getRiskAssessment")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRiskAssessment is a free data retrieval call binding the contract method 0x3acd6d50.
//
// Solidity: function getRiskAssessment() view returns(address)
func (_Registry *RegistrySession) GetRiskAssessment() (common.Address, error) {
	return _Registry.Contract.GetRiskAssessment(&_Registry.CallOpts)
}

// GetRiskAssessment is a free data retrieval call binding the contract method 0x3acd6d50.
//
// Solidity: function getRiskAssessment() view returns(address)
func (_Registry *RegistryCallerSession) GetRiskAssessment() (common.Address, error) {
	return _Registry.Contract.GetRiskAssessment(&_Registry.CallOpts)
}

// GetRole is a free data retrieval call binding the contract method 0x6e8eb7bc.
//
// Solidity: function getRole(bytes32 commitment) view returns(uint8)
func (_Registry *RegistryCaller) GetRole(opts *bind.CallOpts, commitment [32]byte) (uint8, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getRole", commitment)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRole is a free data retrieval call binding the contract method 0x6e8eb7bc.
//
// Solidity: function getRole(bytes32 commitment) view returns(uint8)
func (_Registry *RegistrySession) GetRole(commitment [32]byte) (uint8, error) {
	return _Registry.Contract.GetRole(&_Registry.CallOpts, commitment)
}

// GetRole is a free data retrieval call binding the contract method 0x6e8eb7bc.
//
// Solidity: function getRole(bytes32 commitment) view returns(uint8)
func (_Registry *RegistryCallerSession) GetRole(commitment [32]byte) (uint8, error) {
	return _Registry.Contract.GetRole(&_Registry.CallOpts, commitment)
}

// GetRoleByAddress is a free data retrieval call binding the contract method 0x7bff77d4.
//
// Solidity: function getRoleByAddress(address userAddress) view returns(uint8)
func (_Registry *RegistryCaller) GetRoleByAddress(opts *bind.CallOpts, userAddress common.Address) (uint8, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getRoleByAddress", userAddress)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRoleByAddress is a free data retrieval call binding the contract method 0x7bff77d4.
//
// Solidity: function getRoleByAddress(address userAddress) view returns(uint8)
func (_Registry *RegistrySession) GetRoleByAddress(userAddress common.Address) (uint8, error) {
	return _Registry.Contract.GetRoleByAddress(&_Registry.CallOpts, userAddress)
}

// GetRoleByAddress is a free data retrieval call binding the contract method 0x7bff77d4.
//
// Solidity: function getRoleByAddress(address userAddress) view returns(uint8)
func (_Registry *RegistryCallerSession) GetRoleByAddress(userAddress common.Address) (uint8, error) {
	return _Registry.Contract.GetRoleByAddress(&_Registry.CallOpts, userAddress)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 commitment) view returns((bytes32,bytes32,uint8,address,address))
func (_Registry *RegistryCaller) GetUser(opts *bind.CallOpts, commitment [32]byte) (WellRegistryUser, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getUser", commitment)

	if err != nil {
		return *new(WellRegistryUser), err
	}

	out0 := *abi.ConvertType(out[0], new(WellRegistryUser)).(*WellRegistryUser)

	return out0, err

}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 commitment) view returns((bytes32,bytes32,uint8,address,address))
func (_Registry *RegistrySession) GetUser(commitment [32]byte) (WellRegistryUser, error) {
	return _Registry.Contract.GetUser(&_Registry.CallOpts, commitment)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 commitment) view returns((bytes32,bytes32,uint8,address,address))
func (_Registry *RegistryCallerSession) GetUser(commitment [32]byte) (WellRegistryUser, error) {
	return _Registry.Contract.GetUser(&_Registry.CallOpts, commitment)
}

// GetUserByAddress is a free data retrieval call binding the contract method 0x69c212f6.
//
// Solidity: function getUserByAddress(address userAddress) view returns((bytes32,bytes32,uint8,address,address))
func (_Registry *RegistryCaller) GetUserByAddress(opts *bind.CallOpts, userAddress common.Address) (WellRegistryUser, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "getUserByAddress", userAddress)

	if err != nil {
		return *new(WellRegistryUser), err
	}

	out0 := *abi.ConvertType(out[0], new(WellRegistryUser)).(*WellRegistryUser)

	return out0, err

}

// GetUserByAddress is a free data retrieval call binding the contract method 0x69c212f6.
//
// Solidity: function getUserByAddress(address userAddress) view returns((bytes32,bytes32,uint8,address,address))
func (_Registry *RegistrySession) GetUserByAddress(userAddress common.Address) (WellRegistryUser, error) {
	return _Registry.Contract.GetUserByAddress(&_Registry.CallOpts, userAddress)
}

// GetUserByAddress is a free data retrieval call binding the contract method 0x69c212f6.
//
// Solidity: function getUserByAddress(address userAddress) view returns((bytes32,bytes32,uint8,address,address))
func (_Registry *RegistryCallerSession) GetUserByAddress(userAddress common.Address) (WellRegistryUser, error) {
	return _Registry.Contract.GetUserByAddress(&_Registry.CallOpts, userAddress)
}

// HistoryFactory is a free data retrieval call binding the contract method 0xa1968579.
//
// Solidity: function historyFactory() view returns(address)
func (_Registry *RegistryCaller) HistoryFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "historyFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// HistoryFactory is a free data retrieval call binding the contract method 0xa1968579.
//
// Solidity: function historyFactory() view returns(address)
func (_Registry *RegistrySession) HistoryFactory() (common.Address, error) {
	return _Registry.Contract.HistoryFactory(&_Registry.CallOpts)
}

// HistoryFactory is a free data retrieval call binding the contract method 0xa1968579.
//
// Solidity: function historyFactory() view returns(address)
func (_Registry *RegistryCallerSession) HistoryFactory() (common.Address, error) {
	return _Registry.Contract.HistoryFactory(&_Registry.CallOpts)
}

// IncidentManager is a free data retrieval call binding the contract method 0xe7030023.
//
// Solidity: function incidentManager() view returns(address)
func (_Registry *RegistryCaller) IncidentManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "incidentManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// IncidentManager is a free data retrieval call binding the contract method 0xe7030023.
//
// Solidity: function incidentManager() view returns(address)
func (_Registry *RegistrySession) IncidentManager() (common.Address, error) {
	return _Registry.Contract.IncidentManager(&_Registry.CallOpts)
}

// IncidentManager is a free data retrieval call binding the contract method 0xe7030023.
//
// Solidity: function incidentManager() view returns(address)
func (_Registry *RegistryCallerSession) IncidentManager() (common.Address, error) {
	return _Registry.Contract.IncidentManager(&_Registry.CallOpts)
}

// IsRegistered is a free data retrieval call binding the contract method 0x27258b22.
//
// Solidity: function isRegistered(bytes32 commitment) view returns(bool)
func (_Registry *RegistryCaller) IsRegistered(opts *bind.CallOpts, commitment [32]byte) (bool, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "isRegistered", commitment)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0x27258b22.
//
// Solidity: function isRegistered(bytes32 commitment) view returns(bool)
func (_Registry *RegistrySession) IsRegistered(commitment [32]byte) (bool, error) {
	return _Registry.Contract.IsRegistered(&_Registry.CallOpts, commitment)
}

// IsRegistered is a free data retrieval call binding the contract method 0x27258b22.
//
// Solidity: function isRegistered(bytes32 commitment) view returns(bool)
func (_Registry *RegistryCallerSession) IsRegistered(commitment [32]byte) (bool, error) {
	return _Registry.Contract.IsRegistered(&_Registry.CallOpts, commitment)
}

// RiskAssessment is a free data retrieval call binding the contract method 0x70a1a7d5.
//
// Solidity: function riskAssessment() view returns(address)
func (_Registry *RegistryCaller) RiskAssessment(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "riskAssessment")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RiskAssessment is a free data retrieval call binding the contract method 0x70a1a7d5.
//
// Solidity: function riskAssessment() view returns(address)
func (_Registry *RegistrySession) RiskAssessment() (common.Address, error) {
	return _Registry.Contract.RiskAssessment(&_Registry.CallOpts)
}

// RiskAssessment is a free data retrieval call binding the contract method 0x70a1a7d5.
//
// Solidity: function riskAssessment() view returns(address)
func (_Registry *RegistryCallerSession) RiskAssessment() (common.Address, error) {
	return _Registry.Contract.RiskAssessment(&_Registry.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Registry *RegistryCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Registry.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Registry *RegistrySession) UserRegistry() (common.Address, error) {
	return _Registry.Contract.UserRegistry(&_Registry.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Registry *RegistryCallerSession) UserRegistry() (common.Address, error) {
	return _Registry.Contract.UserRegistry(&_Registry.CallOpts)
}

// CreateRecord is a paid mutator transaction binding the contract method 0x3e51e13c.
//
// Solidity: function createRecord(address uploader, (uint8,uint256,uint256) metadata, bytes32 contentHash, string storageReference) returns(address)
func (_Registry *RegistryTransactor) CreateRecord(opts *bind.TransactOpts, uploader common.Address, metadata TypesRecordMetadata, contentHash [32]byte, storageReference string) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "createRecord", uploader, metadata, contentHash, storageReference)
}

// CreateRecord is a paid mutator transaction binding the contract method 0x3e51e13c.
//
// Solidity: function createRecord(address uploader, (uint8,uint256,uint256) metadata, bytes32 contentHash, string storageReference) returns(address)
func (_Registry *RegistrySession) CreateRecord(uploader common.Address, metadata TypesRecordMetadata, contentHash [32]byte, storageReference string) (*types.Transaction, error) {
	return _Registry.Contract.CreateRecord(&_Registry.TransactOpts, uploader, metadata, contentHash, storageReference)
}

// CreateRecord is a paid mutator transaction binding the contract method 0x3e51e13c.
//
// Solidity: function createRecord(address uploader, (uint8,uint256,uint256) metadata, bytes32 contentHash, string storageReference) returns(address)
func (_Registry *RegistryTransactorSession) CreateRecord(uploader common.Address, metadata TypesRecordMetadata, contentHash [32]byte, storageReference string) (*types.Transaction, error) {
	return _Registry.Contract.CreateRecord(&_Registry.TransactOpts, uploader, metadata, contentHash, storageReference)
}

// Register is a paid mutator transaction binding the contract method 0x4fe183dd.
//
// Solidity: function register(bytes32 commitment, bytes32 internalID, uint8 role) returns(address historyContract)
func (_Registry *RegistryTransactor) Register(opts *bind.TransactOpts, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "register", commitment, internalID, role)
}

// Register is a paid mutator transaction binding the contract method 0x4fe183dd.
//
// Solidity: function register(bytes32 commitment, bytes32 internalID, uint8 role) returns(address historyContract)
func (_Registry *RegistrySession) Register(commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, commitment, internalID, role)
}

// Register is a paid mutator transaction binding the contract method 0x4fe183dd.
//
// Solidity: function register(bytes32 commitment, bytes32 internalID, uint8 role) returns(address historyContract)
func (_Registry *RegistryTransactorSession) Register(commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _Registry.Contract.Register(&_Registry.TransactOpts, commitment, internalID, role)
}

// RegisterByAdmin is a paid mutator transaction binding the contract method 0x78e61e5c.
//
// Solidity: function registerByAdmin(address userAddress, bytes32 commitment, bytes32 internalID, uint8 role) returns(address historyContract)
func (_Registry *RegistryTransactor) RegisterByAdmin(opts *bind.TransactOpts, userAddress common.Address, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "registerByAdmin", userAddress, commitment, internalID, role)
}

// RegisterByAdmin is a paid mutator transaction binding the contract method 0x78e61e5c.
//
// Solidity: function registerByAdmin(address userAddress, bytes32 commitment, bytes32 internalID, uint8 role) returns(address historyContract)
func (_Registry *RegistrySession) RegisterByAdmin(userAddress common.Address, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _Registry.Contract.RegisterByAdmin(&_Registry.TransactOpts, userAddress, commitment, internalID, role)
}

// RegisterByAdmin is a paid mutator transaction binding the contract method 0x78e61e5c.
//
// Solidity: function registerByAdmin(address userAddress, bytes32 commitment, bytes32 internalID, uint8 role) returns(address historyContract)
func (_Registry *RegistryTransactorSession) RegisterByAdmin(userAddress common.Address, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _Registry.Contract.RegisterByAdmin(&_Registry.TransactOpts, userAddress, commitment, internalID, role)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 commitment) returns()
func (_Registry *RegistryTransactor) Revoke(opts *bind.TransactOpts, commitment [32]byte) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "revoke", commitment)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 commitment) returns()
func (_Registry *RegistrySession) Revoke(commitment [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.Revoke(&_Registry.TransactOpts, commitment)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 commitment) returns()
func (_Registry *RegistryTransactorSession) Revoke(commitment [32]byte) (*types.Transaction, error) {
	return _Registry.Contract.Revoke(&_Registry.TransactOpts, commitment)
}

// SetEmergencySummaryFactory is a paid mutator transaction binding the contract method 0x6291a18f.
//
// Solidity: function setEmergencySummaryFactory(address _emergencySummaryFactory) returns()
func (_Registry *RegistryTransactor) SetEmergencySummaryFactory(opts *bind.TransactOpts, _emergencySummaryFactory common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setEmergencySummaryFactory", _emergencySummaryFactory)
}

// SetEmergencySummaryFactory is a paid mutator transaction binding the contract method 0x6291a18f.
//
// Solidity: function setEmergencySummaryFactory(address _emergencySummaryFactory) returns()
func (_Registry *RegistrySession) SetEmergencySummaryFactory(_emergencySummaryFactory common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetEmergencySummaryFactory(&_Registry.TransactOpts, _emergencySummaryFactory)
}

// SetEmergencySummaryFactory is a paid mutator transaction binding the contract method 0x6291a18f.
//
// Solidity: function setEmergencySummaryFactory(address _emergencySummaryFactory) returns()
func (_Registry *RegistryTransactorSession) SetEmergencySummaryFactory(_emergencySummaryFactory common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetEmergencySummaryFactory(&_Registry.TransactOpts, _emergencySummaryFactory)
}

// SetIncidentManager is a paid mutator transaction binding the contract method 0x4d281948.
//
// Solidity: function setIncidentManager(address _incidentManager) returns()
func (_Registry *RegistryTransactor) SetIncidentManager(opts *bind.TransactOpts, _incidentManager common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setIncidentManager", _incidentManager)
}

// SetIncidentManager is a paid mutator transaction binding the contract method 0x4d281948.
//
// Solidity: function setIncidentManager(address _incidentManager) returns()
func (_Registry *RegistrySession) SetIncidentManager(_incidentManager common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetIncidentManager(&_Registry.TransactOpts, _incidentManager)
}

// SetIncidentManager is a paid mutator transaction binding the contract method 0x4d281948.
//
// Solidity: function setIncidentManager(address _incidentManager) returns()
func (_Registry *RegistryTransactorSession) SetIncidentManager(_incidentManager common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetIncidentManager(&_Registry.TransactOpts, _incidentManager)
}

// SetRiskAssessment is a paid mutator transaction binding the contract method 0x76d2feee.
//
// Solidity: function setRiskAssessment(address _riskAssessment) returns()
func (_Registry *RegistryTransactor) SetRiskAssessment(opts *bind.TransactOpts, _riskAssessment common.Address) (*types.Transaction, error) {
	return _Registry.contract.Transact(opts, "setRiskAssessment", _riskAssessment)
}

// SetRiskAssessment is a paid mutator transaction binding the contract method 0x76d2feee.
//
// Solidity: function setRiskAssessment(address _riskAssessment) returns()
func (_Registry *RegistrySession) SetRiskAssessment(_riskAssessment common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetRiskAssessment(&_Registry.TransactOpts, _riskAssessment)
}

// SetRiskAssessment is a paid mutator transaction binding the contract method 0x76d2feee.
//
// Solidity: function setRiskAssessment(address _riskAssessment) returns()
func (_Registry *RegistryTransactorSession) SetRiskAssessment(_riskAssessment common.Address) (*types.Transaction, error) {
	return _Registry.Contract.SetRiskAssessment(&_Registry.TransactOpts, _riskAssessment)
}
