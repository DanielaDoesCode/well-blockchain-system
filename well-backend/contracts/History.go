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


// HistoryMetaData contains all meta data concerning the History contract.
var HistoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_accessControl\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_policyVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"accessControl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createRecord\",\"inputs\":[{\"name\":\"_uploader\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_metadata\",\"type\":\"tuple\",\"internalType\":\"structTypes.RecordMetadata\",\"components\":[{\"name\":\"recordType\",\"type\":\"uint8\",\"internalType\":\"enumTypes.RecordType\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_contentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_storageReference\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractRecord\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllRecords\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"contractRecord[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecord\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"internalType\":\"contractRecord\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTypes.RecordMetadata\",\"components\":[{\"name\":\"recordType\",\"type\":\"uint8\",\"internalType\":\"enumTypes.RecordType\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecordCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRecordAccess\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"internalType\":\"contractRecord\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasAccessToRecord\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"internalType\":\"contractRecord\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policyVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revokeRecordAccess\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"internalType\":\"contractRecord\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyRecordIntegrity\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"internalType\":\"contractRecord\"},{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AccessGranted\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AccessRevoked\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RecordCreated\",\"inputs\":[{\"name\":\"record\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"uploader\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recordType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumTypes.RecordType\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"RecordNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// HistoryABI is the input ABI used to generate the binding from.
// Deprecated: Use HistoryMetaData.ABI instead.
var HistoryABI = HistoryMetaData.ABI

// History is an auto generated Go binding around an Ethereum contract.
type History struct {
	HistoryCaller     // Read-only binding to the contract
	HistoryTransactor // Write-only binding to the contract
	HistoryFilterer   // Log filterer for contract events
}

// HistoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type HistoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HistoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HistoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HistorySession struct {
	Contract     *History          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HistoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HistoryCallerSession struct {
	Contract *HistoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// HistoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HistoryTransactorSession struct {
	Contract     *HistoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// HistoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type HistoryRaw struct {
	Contract *History // Generic contract binding to access the raw methods on
}

// HistoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HistoryCallerRaw struct {
	Contract *HistoryCaller // Generic read-only contract binding to access the raw methods on
}

// HistoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HistoryTransactorRaw struct {
	Contract *HistoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHistory creates a new instance of History, bound to a specific deployed contract.
func NewHistory(address common.Address, backend bind.ContractBackend) (*History, error) {
	contract, err := bindHistory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &History{HistoryCaller: HistoryCaller{contract: contract}, HistoryTransactor: HistoryTransactor{contract: contract}, HistoryFilterer: HistoryFilterer{contract: contract}}, nil
}

// NewHistoryCaller creates a new read-only instance of History, bound to a specific deployed contract.
func NewHistoryCaller(address common.Address, caller bind.ContractCaller) (*HistoryCaller, error) {
	contract, err := bindHistory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HistoryCaller{contract: contract}, nil
}

// NewHistoryTransactor creates a new write-only instance of History, bound to a specific deployed contract.
func NewHistoryTransactor(address common.Address, transactor bind.ContractTransactor) (*HistoryTransactor, error) {
	contract, err := bindHistory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HistoryTransactor{contract: contract}, nil
}

// NewHistoryFilterer creates a new log filterer instance of History, bound to a specific deployed contract.
func NewHistoryFilterer(address common.Address, filterer bind.ContractFilterer) (*HistoryFilterer, error) {
	contract, err := bindHistory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HistoryFilterer{contract: contract}, nil
}

// bindHistory binds a generic wrapper to an already deployed contract.
func bindHistory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HistoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_History *HistoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _History.Contract.HistoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_History *HistoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _History.Contract.HistoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_History *HistoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _History.Contract.HistoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_History *HistoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _History.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_History *HistoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _History.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_History *HistoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _History.Contract.contract.Transact(opts, method, params...)
}

// AccessControl is a free data retrieval call binding the contract method 0x13007d55.
//
// Solidity: function accessControl() view returns(address)
func (_History *HistoryCaller) AccessControl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "accessControl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessControl is a free data retrieval call binding the contract method 0x13007d55.
//
// Solidity: function accessControl() view returns(address)
func (_History *HistorySession) AccessControl() (common.Address, error) {
	return _History.Contract.AccessControl(&_History.CallOpts)
}

// AccessControl is a free data retrieval call binding the contract method 0x13007d55.
//
// Solidity: function accessControl() view returns(address)
func (_History *HistoryCallerSession) AccessControl() (common.Address, error) {
	return _History.Contract.AccessControl(&_History.CallOpts)
}

// GetAllRecords is a free data retrieval call binding the contract method 0xa7f9fe72.
//
// Solidity: function getAllRecords() view returns(address[])
func (_History *HistoryCaller) GetAllRecords(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "getAllRecords")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAllRecords is a free data retrieval call binding the contract method 0xa7f9fe72.
//
// Solidity: function getAllRecords() view returns(address[])
func (_History *HistorySession) GetAllRecords() ([]common.Address, error) {
	return _History.Contract.GetAllRecords(&_History.CallOpts)
}

// GetAllRecords is a free data retrieval call binding the contract method 0xa7f9fe72.
//
// Solidity: function getAllRecords() view returns(address[])
func (_History *HistoryCallerSession) GetAllRecords() ([]common.Address, error) {
	return _History.Contract.GetAllRecords(&_History.CallOpts)
}

// GetRecord is a free data retrieval call binding the contract method 0x617fba04.
//
// Solidity: function getRecord(address record) view returns((uint8,uint256,uint256), bytes32, string, address, address)
func (_History *HistoryCaller) GetRecord(opts *bind.CallOpts, record common.Address) (TypesRecordMetadata, [32]byte, string, common.Address, common.Address, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "getRecord", record)

	if err != nil {
		return *new(TypesRecordMetadata), *new([32]byte), *new(string), *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesRecordMetadata)).(*TypesRecordMetadata)
	out1 := *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	out4 := *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return out0, out1, out2, out3, out4, err

}

// GetRecord is a free data retrieval call binding the contract method 0x617fba04.
//
// Solidity: function getRecord(address record) view returns((uint8,uint256,uint256), bytes32, string, address, address)
func (_History *HistorySession) GetRecord(record common.Address) (TypesRecordMetadata, [32]byte, string, common.Address, common.Address, error) {
	return _History.Contract.GetRecord(&_History.CallOpts, record)
}

// GetRecord is a free data retrieval call binding the contract method 0x617fba04.
//
// Solidity: function getRecord(address record) view returns((uint8,uint256,uint256), bytes32, string, address, address)
func (_History *HistoryCallerSession) GetRecord(record common.Address) (TypesRecordMetadata, [32]byte, string, common.Address, common.Address, error) {
	return _History.Contract.GetRecord(&_History.CallOpts, record)
}

// GetRecordCount is a free data retrieval call binding the contract method 0xca267f28.
//
// Solidity: function getRecordCount() view returns(uint256)
func (_History *HistoryCaller) GetRecordCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "getRecordCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRecordCount is a free data retrieval call binding the contract method 0xca267f28.
//
// Solidity: function getRecordCount() view returns(uint256)
func (_History *HistorySession) GetRecordCount() (*big.Int, error) {
	return _History.Contract.GetRecordCount(&_History.CallOpts)
}

// GetRecordCount is a free data retrieval call binding the contract method 0xca267f28.
//
// Solidity: function getRecordCount() view returns(uint256)
func (_History *HistoryCallerSession) GetRecordCount() (*big.Int, error) {
	return _History.Contract.GetRecordCount(&_History.CallOpts)
}

// HasAccessToRecord is a free data retrieval call binding the contract method 0x513fb46f.
//
// Solidity: function hasAccessToRecord(address record, address user) view returns(bool)
func (_History *HistoryCaller) HasAccessToRecord(opts *bind.CallOpts, record common.Address, user common.Address) (bool, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "hasAccessToRecord", record, user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccessToRecord is a free data retrieval call binding the contract method 0x513fb46f.
//
// Solidity: function hasAccessToRecord(address record, address user) view returns(bool)
func (_History *HistorySession) HasAccessToRecord(record common.Address, user common.Address) (bool, error) {
	return _History.Contract.HasAccessToRecord(&_History.CallOpts, record, user)
}

// HasAccessToRecord is a free data retrieval call binding the contract method 0x513fb46f.
//
// Solidity: function hasAccessToRecord(address record, address user) view returns(bool)
func (_History *HistoryCallerSession) HasAccessToRecord(record common.Address, user common.Address) (bool, error) {
	return _History.Contract.HasAccessToRecord(&_History.CallOpts, record, user)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_History *HistoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_History *HistorySession) Owner() (common.Address, error) {
	return _History.Contract.Owner(&_History.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_History *HistoryCallerSession) Owner() (common.Address, error) {
	return _History.Contract.Owner(&_History.CallOpts)
}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_History *HistoryCaller) PolicyVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "policyVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_History *HistorySession) PolicyVerifier() (common.Address, error) {
	return _History.Contract.PolicyVerifier(&_History.CallOpts)
}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_History *HistoryCallerSession) PolicyVerifier() (common.Address, error) {
	return _History.Contract.PolicyVerifier(&_History.CallOpts)
}

// VerifyRecordIntegrity is a free data retrieval call binding the contract method 0x41510560.
//
// Solidity: function verifyRecordIntegrity(address record, bytes32 hash) view returns(bool)
func (_History *HistoryCaller) VerifyRecordIntegrity(opts *bind.CallOpts, record common.Address, hash [32]byte) (bool, error) {
	var out []interface{}
	err := _History.contract.Call(opts, &out, "verifyRecordIntegrity", record, hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyRecordIntegrity is a free data retrieval call binding the contract method 0x41510560.
//
// Solidity: function verifyRecordIntegrity(address record, bytes32 hash) view returns(bool)
func (_History *HistorySession) VerifyRecordIntegrity(record common.Address, hash [32]byte) (bool, error) {
	return _History.Contract.VerifyRecordIntegrity(&_History.CallOpts, record, hash)
}

// VerifyRecordIntegrity is a free data retrieval call binding the contract method 0x41510560.
//
// Solidity: function verifyRecordIntegrity(address record, bytes32 hash) view returns(bool)
func (_History *HistoryCallerSession) VerifyRecordIntegrity(record common.Address, hash [32]byte) (bool, error) {
	return _History.Contract.VerifyRecordIntegrity(&_History.CallOpts, record, hash)
}

// CreateRecord is a paid mutator transaction binding the contract method 0x3e51e13c.
//
// Solidity: function createRecord(address _uploader, (uint8,uint256,uint256) _metadata, bytes32 _contentHash, string _storageReference) returns(address)
func (_History *HistoryTransactor) CreateRecord(opts *bind.TransactOpts, _uploader common.Address, _metadata TypesRecordMetadata, _contentHash [32]byte, _storageReference string) (*types.Transaction, error) {
	return _History.contract.Transact(opts, "createRecord", _uploader, _metadata, _contentHash, _storageReference)
}

// CreateRecord is a paid mutator transaction binding the contract method 0x3e51e13c.
//
// Solidity: function createRecord(address _uploader, (uint8,uint256,uint256) _metadata, bytes32 _contentHash, string _storageReference) returns(address)
func (_History *HistorySession) CreateRecord(_uploader common.Address, _metadata TypesRecordMetadata, _contentHash [32]byte, _storageReference string) (*types.Transaction, error) {
	return _History.Contract.CreateRecord(&_History.TransactOpts, _uploader, _metadata, _contentHash, _storageReference)
}

// CreateRecord is a paid mutator transaction binding the contract method 0x3e51e13c.
//
// Solidity: function createRecord(address _uploader, (uint8,uint256,uint256) _metadata, bytes32 _contentHash, string _storageReference) returns(address)
func (_History *HistoryTransactorSession) CreateRecord(_uploader common.Address, _metadata TypesRecordMetadata, _contentHash [32]byte, _storageReference string) (*types.Transaction, error) {
	return _History.Contract.CreateRecord(&_History.TransactOpts, _uploader, _metadata, _contentHash, _storageReference)
}

// GrantRecordAccess is a paid mutator transaction binding the contract method 0xec2a8d82.
//
// Solidity: function grantRecordAccess(address record, address user) returns()
func (_History *HistoryTransactor) GrantRecordAccess(opts *bind.TransactOpts, record common.Address, user common.Address) (*types.Transaction, error) {
	return _History.contract.Transact(opts, "grantRecordAccess", record, user)
}

// GrantRecordAccess is a paid mutator transaction binding the contract method 0xec2a8d82.
//
// Solidity: function grantRecordAccess(address record, address user) returns()
func (_History *HistorySession) GrantRecordAccess(record common.Address, user common.Address) (*types.Transaction, error) {
	return _History.Contract.GrantRecordAccess(&_History.TransactOpts, record, user)
}

// GrantRecordAccess is a paid mutator transaction binding the contract method 0xec2a8d82.
//
// Solidity: function grantRecordAccess(address record, address user) returns()
func (_History *HistoryTransactorSession) GrantRecordAccess(record common.Address, user common.Address) (*types.Transaction, error) {
	return _History.Contract.GrantRecordAccess(&_History.TransactOpts, record, user)
}

// RevokeRecordAccess is a paid mutator transaction binding the contract method 0xd73d01ee.
//
// Solidity: function revokeRecordAccess(address record, address user) returns()
func (_History *HistoryTransactor) RevokeRecordAccess(opts *bind.TransactOpts, record common.Address, user common.Address) (*types.Transaction, error) {
	return _History.contract.Transact(opts, "revokeRecordAccess", record, user)
}

// RevokeRecordAccess is a paid mutator transaction binding the contract method 0xd73d01ee.
//
// Solidity: function revokeRecordAccess(address record, address user) returns()
func (_History *HistorySession) RevokeRecordAccess(record common.Address, user common.Address) (*types.Transaction, error) {
	return _History.Contract.RevokeRecordAccess(&_History.TransactOpts, record, user)
}

// RevokeRecordAccess is a paid mutator transaction binding the contract method 0xd73d01ee.
//
// Solidity: function revokeRecordAccess(address record, address user) returns()
func (_History *HistoryTransactorSession) RevokeRecordAccess(record common.Address, user common.Address) (*types.Transaction, error) {
	return _History.Contract.RevokeRecordAccess(&_History.TransactOpts, record, user)
}

// HistoryAccessGrantedIterator is returned from FilterAccessGranted and is used to iterate over the raw logs and unpacked data for AccessGranted events raised by the History contract.
type HistoryAccessGrantedIterator struct {
	Event *HistoryAccessGranted // Event containing the contract specifics and raw log

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
func (it *HistoryAccessGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoryAccessGranted)
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
		it.Event = new(HistoryAccessGranted)
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
func (it *HistoryAccessGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoryAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoryAccessGranted represents a AccessGranted event raised by the History contract.
type HistoryAccessGranted struct {
	Record common.Address
	User   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAccessGranted is a free log retrieval operation binding the contract event 0x9f520b3bfc56d06f7065d10b4683b3f57ac8895d5026cd176dee95ce1454cb8d.
//
// Solidity: event AccessGranted(address indexed record, address indexed user)
func (_History *HistoryFilterer) FilterAccessGranted(opts *bind.FilterOpts, record []common.Address, user []common.Address) (*HistoryAccessGrantedIterator, error) {

	var recordRule []interface{}
	for _, recordItem := range record {
		recordRule = append(recordRule, recordItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _History.contract.FilterLogs(opts, "AccessGranted", recordRule, userRule)
	if err != nil {
		return nil, err
	}
	return &HistoryAccessGrantedIterator{contract: _History.contract, event: "AccessGranted", logs: logs, sub: sub}, nil
}

// WatchAccessGranted is a free log subscription operation binding the contract event 0x9f520b3bfc56d06f7065d10b4683b3f57ac8895d5026cd176dee95ce1454cb8d.
//
// Solidity: event AccessGranted(address indexed record, address indexed user)
func (_History *HistoryFilterer) WatchAccessGranted(opts *bind.WatchOpts, sink chan<- *HistoryAccessGranted, record []common.Address, user []common.Address) (event.Subscription, error) {

	var recordRule []interface{}
	for _, recordItem := range record {
		recordRule = append(recordRule, recordItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _History.contract.WatchLogs(opts, "AccessGranted", recordRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoryAccessGranted)
				if err := _History.contract.UnpackLog(event, "AccessGranted", log); err != nil {
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

// ParseAccessGranted is a log parse operation binding the contract event 0x9f520b3bfc56d06f7065d10b4683b3f57ac8895d5026cd176dee95ce1454cb8d.
//
// Solidity: event AccessGranted(address indexed record, address indexed user)
func (_History *HistoryFilterer) ParseAccessGranted(log types.Log) (*HistoryAccessGranted, error) {
	event := new(HistoryAccessGranted)
	if err := _History.contract.UnpackLog(event, "AccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HistoryAccessRevokedIterator is returned from FilterAccessRevoked and is used to iterate over the raw logs and unpacked data for AccessRevoked events raised by the History contract.
type HistoryAccessRevokedIterator struct {
	Event *HistoryAccessRevoked // Event containing the contract specifics and raw log

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
func (it *HistoryAccessRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoryAccessRevoked)
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
		it.Event = new(HistoryAccessRevoked)
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
func (it *HistoryAccessRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoryAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoryAccessRevoked represents a AccessRevoked event raised by the History contract.
type HistoryAccessRevoked struct {
	Record common.Address
	User   common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAccessRevoked is a free log retrieval operation binding the contract event 0x825c8be24eb0df19500f63e86c29e7d0d951e73056b889b891d85e40938d9b6e.
//
// Solidity: event AccessRevoked(address indexed record, address indexed user)
func (_History *HistoryFilterer) FilterAccessRevoked(opts *bind.FilterOpts, record []common.Address, user []common.Address) (*HistoryAccessRevokedIterator, error) {

	var recordRule []interface{}
	for _, recordItem := range record {
		recordRule = append(recordRule, recordItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _History.contract.FilterLogs(opts, "AccessRevoked", recordRule, userRule)
	if err != nil {
		return nil, err
	}
	return &HistoryAccessRevokedIterator{contract: _History.contract, event: "AccessRevoked", logs: logs, sub: sub}, nil
}

// WatchAccessRevoked is a free log subscription operation binding the contract event 0x825c8be24eb0df19500f63e86c29e7d0d951e73056b889b891d85e40938d9b6e.
//
// Solidity: event AccessRevoked(address indexed record, address indexed user)
func (_History *HistoryFilterer) WatchAccessRevoked(opts *bind.WatchOpts, sink chan<- *HistoryAccessRevoked, record []common.Address, user []common.Address) (event.Subscription, error) {

	var recordRule []interface{}
	for _, recordItem := range record {
		recordRule = append(recordRule, recordItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _History.contract.WatchLogs(opts, "AccessRevoked", recordRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoryAccessRevoked)
				if err := _History.contract.UnpackLog(event, "AccessRevoked", log); err != nil {
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

// ParseAccessRevoked is a log parse operation binding the contract event 0x825c8be24eb0df19500f63e86c29e7d0d951e73056b889b891d85e40938d9b6e.
//
// Solidity: event AccessRevoked(address indexed record, address indexed user)
func (_History *HistoryFilterer) ParseAccessRevoked(log types.Log) (*HistoryAccessRevoked, error) {
	event := new(HistoryAccessRevoked)
	if err := _History.contract.UnpackLog(event, "AccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HistoryRecordCreatedIterator is returned from FilterRecordCreated and is used to iterate over the raw logs and unpacked data for RecordCreated events raised by the History contract.
type HistoryRecordCreatedIterator struct {
	Event *HistoryRecordCreated // Event containing the contract specifics and raw log

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
func (it *HistoryRecordCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoryRecordCreated)
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
		it.Event = new(HistoryRecordCreated)
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
func (it *HistoryRecordCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoryRecordCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoryRecordCreated represents a RecordCreated event raised by the History contract.
type HistoryRecordCreated struct {
	Record     common.Address
	Uploader   common.Address
	RecordType uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRecordCreated is a free log retrieval operation binding the contract event 0x6680dcae98d75849ef75b956cfbbc634aa68b616a6123f52022bb72ce4b775ac.
//
// Solidity: event RecordCreated(address indexed record, address indexed uploader, uint8 recordType)
func (_History *HistoryFilterer) FilterRecordCreated(opts *bind.FilterOpts, record []common.Address, uploader []common.Address) (*HistoryRecordCreatedIterator, error) {

	var recordRule []interface{}
	for _, recordItem := range record {
		recordRule = append(recordRule, recordItem)
	}
	var uploaderRule []interface{}
	for _, uploaderItem := range uploader {
		uploaderRule = append(uploaderRule, uploaderItem)
	}

	logs, sub, err := _History.contract.FilterLogs(opts, "RecordCreated", recordRule, uploaderRule)
	if err != nil {
		return nil, err
	}
	return &HistoryRecordCreatedIterator{contract: _History.contract, event: "RecordCreated", logs: logs, sub: sub}, nil
}

// WatchRecordCreated is a free log subscription operation binding the contract event 0x6680dcae98d75849ef75b956cfbbc634aa68b616a6123f52022bb72ce4b775ac.
//
// Solidity: event RecordCreated(address indexed record, address indexed uploader, uint8 recordType)
func (_History *HistoryFilterer) WatchRecordCreated(opts *bind.WatchOpts, sink chan<- *HistoryRecordCreated, record []common.Address, uploader []common.Address) (event.Subscription, error) {

	var recordRule []interface{}
	for _, recordItem := range record {
		recordRule = append(recordRule, recordItem)
	}
	var uploaderRule []interface{}
	for _, uploaderItem := range uploader {
		uploaderRule = append(uploaderRule, uploaderItem)
	}

	logs, sub, err := _History.contract.WatchLogs(opts, "RecordCreated", recordRule, uploaderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoryRecordCreated)
				if err := _History.contract.UnpackLog(event, "RecordCreated", log); err != nil {
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

// ParseRecordCreated is a log parse operation binding the contract event 0x6680dcae98d75849ef75b956cfbbc634aa68b616a6123f52022bb72ce4b775ac.
//
// Solidity: event RecordCreated(address indexed record, address indexed uploader, uint8 recordType)
func (_History *HistoryFilterer) ParseRecordCreated(log types.Log) (*HistoryRecordCreated, error) {
	event := new(HistoryRecordCreated)
	if err := _History.contract.UnpackLog(event, "RecordCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
