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

// EmergencySummaryMetaData contains all meta data concerning the EmergencySummary contract.
var EmergencySummaryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_registry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_riskAssessement\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ACCESS_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"REQUEST_EXPIRY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approveAccessRequest\",\"inputs\":[{\"name\":\"_requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dataHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"dataVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"denyAccessRequest\",\"inputs\":[{\"name\":\"_requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"encryptedData\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccessInfo\",\"inputs\":[{\"name\":\"_user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"isGranted\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"grantedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiresAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccessLog\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccessRequest\",\"inputs\":[{\"name\":\"_requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"responder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"incidentDescription\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumEmergencySummary.AccessStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAuthorizedUsers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getData\",\"inputs\":[],\"outputs\":[{\"name\":\"_encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_version\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEncryptedData\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEncryptedKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMetadata\",\"inputs\":[],\"outputs\":[{\"name\":\"version\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updated\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"dataSize\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecoveryKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantAccess\",\"inputs\":[{\"name\":\"_user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_duration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasAccess\",\"inputs\":[{\"name\":\"_user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastUpdated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoveryKey\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recoveryVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerPublicKey\",\"inputs\":[{\"name\":\"_publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestAccess\",\"inputs\":[{\"name\":\"_publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_location\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_incidentDescription\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_patientIdentifier\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"riskAssessement\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"storeRecoveryKey\",\"inputs\":[{\"name\":\"_recoveryKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateData\",\"inputs\":[{\"name\":\"_encryptedData\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_dataHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKeys\",\"inputs\":[{\"name\":\"_users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_encryptedKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AccessGranted\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"expiresAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AccessRequested\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"responder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AccessRevoked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DataUpdated\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"dataHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PublicKeyRegistered\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RecoveryKeyStored\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoAccessKey\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PublicKeyRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequestExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequestNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// EmergencySummaryABI is the input ABI used to generate the binding from.
// Deprecated: Use EmergencySummaryMetaData.ABI instead.
var EmergencySummaryABI = EmergencySummaryMetaData.ABI

// EmergencySummary is an auto generated Go binding around an Ethereum contract.
type EmergencySummary struct {
	EmergencySummaryCaller     // Read-only binding to the contract
	EmergencySummaryTransactor // Write-only binding to the contract
	EmergencySummaryFilterer   // Log filterer for contract events
}

// EmergencySummaryCaller is an auto generated read-only Go binding around an Ethereum contract.
type EmergencySummaryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmergencySummaryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EmergencySummaryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmergencySummaryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EmergencySummaryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EmergencySummarySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EmergencySummarySession struct {
	Contract     *EmergencySummary // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EmergencySummaryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EmergencySummaryCallerSession struct {
	Contract *EmergencySummaryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// EmergencySummaryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EmergencySummaryTransactorSession struct {
	Contract     *EmergencySummaryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// EmergencySummaryRaw is an auto generated low-level Go binding around an Ethereum contract.
type EmergencySummaryRaw struct {
	Contract *EmergencySummary // Generic contract binding to access the raw methods on
}

// EmergencySummaryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EmergencySummaryCallerRaw struct {
	Contract *EmergencySummaryCaller // Generic read-only contract binding to access the raw methods on
}

// EmergencySummaryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EmergencySummaryTransactorRaw struct {
	Contract *EmergencySummaryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEmergencySummary creates a new instance of EmergencySummary, bound to a specific deployed contract.
func NewEmergencySummary(address common.Address, backend bind.ContractBackend) (*EmergencySummary, error) {
	contract, err := bindEmergencySummary(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EmergencySummary{EmergencySummaryCaller: EmergencySummaryCaller{contract: contract}, EmergencySummaryTransactor: EmergencySummaryTransactor{contract: contract}, EmergencySummaryFilterer: EmergencySummaryFilterer{contract: contract}}, nil
}

// NewEmergencySummaryCaller creates a new read-only instance of EmergencySummary, bound to a specific deployed contract.
func NewEmergencySummaryCaller(address common.Address, caller bind.ContractCaller) (*EmergencySummaryCaller, error) {
	contract, err := bindEmergencySummary(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryCaller{contract: contract}, nil
}

// NewEmergencySummaryTransactor creates a new write-only instance of EmergencySummary, bound to a specific deployed contract.
func NewEmergencySummaryTransactor(address common.Address, transactor bind.ContractTransactor) (*EmergencySummaryTransactor, error) {
	contract, err := bindEmergencySummary(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryTransactor{contract: contract}, nil
}

// NewEmergencySummaryFilterer creates a new log filterer instance of EmergencySummary, bound to a specific deployed contract.
func NewEmergencySummaryFilterer(address common.Address, filterer bind.ContractFilterer) (*EmergencySummaryFilterer, error) {
	contract, err := bindEmergencySummary(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryFilterer{contract: contract}, nil
}

// bindEmergencySummary binds a generic wrapper to an already deployed contract.
func bindEmergencySummary(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EmergencySummaryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmergencySummary *EmergencySummaryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmergencySummary.Contract.EmergencySummaryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmergencySummary *EmergencySummaryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmergencySummary.Contract.EmergencySummaryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmergencySummary *EmergencySummaryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmergencySummary.Contract.EmergencySummaryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EmergencySummary *EmergencySummaryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EmergencySummary.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EmergencySummary *EmergencySummaryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EmergencySummary.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EmergencySummary *EmergencySummaryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EmergencySummary.Contract.contract.Transact(opts, method, params...)
}

// ACCESSDURATION is a free data retrieval call binding the contract method 0x63a0cce6.
//
// Solidity: function ACCESS_DURATION() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCaller) ACCESSDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "ACCESS_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ACCESSDURATION is a free data retrieval call binding the contract method 0x63a0cce6.
//
// Solidity: function ACCESS_DURATION() view returns(uint256)
func (_EmergencySummary *EmergencySummarySession) ACCESSDURATION() (*big.Int, error) {
	return _EmergencySummary.Contract.ACCESSDURATION(&_EmergencySummary.CallOpts)
}

// ACCESSDURATION is a free data retrieval call binding the contract method 0x63a0cce6.
//
// Solidity: function ACCESS_DURATION() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCallerSession) ACCESSDURATION() (*big.Int, error) {
	return _EmergencySummary.Contract.ACCESSDURATION(&_EmergencySummary.CallOpts)
}

// REQUESTEXPIRY is a free data retrieval call binding the contract method 0x22f195b5.
//
// Solidity: function REQUEST_EXPIRY() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCaller) REQUESTEXPIRY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "REQUEST_EXPIRY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// REQUESTEXPIRY is a free data retrieval call binding the contract method 0x22f195b5.
//
// Solidity: function REQUEST_EXPIRY() view returns(uint256)
func (_EmergencySummary *EmergencySummarySession) REQUESTEXPIRY() (*big.Int, error) {
	return _EmergencySummary.Contract.REQUESTEXPIRY(&_EmergencySummary.CallOpts)
}

// REQUESTEXPIRY is a free data retrieval call binding the contract method 0x22f195b5.
//
// Solidity: function REQUEST_EXPIRY() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCallerSession) REQUESTEXPIRY() (*big.Int, error) {
	return _EmergencySummary.Contract.REQUESTEXPIRY(&_EmergencySummary.CallOpts)
}

// DataHash is a free data retrieval call binding the contract method 0x1b3012a3.
//
// Solidity: function dataHash() view returns(bytes32)
func (_EmergencySummary *EmergencySummaryCaller) DataHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "dataHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DataHash is a free data retrieval call binding the contract method 0x1b3012a3.
//
// Solidity: function dataHash() view returns(bytes32)
func (_EmergencySummary *EmergencySummarySession) DataHash() ([32]byte, error) {
	return _EmergencySummary.Contract.DataHash(&_EmergencySummary.CallOpts)
}

// DataHash is a free data retrieval call binding the contract method 0x1b3012a3.
//
// Solidity: function dataHash() view returns(bytes32)
func (_EmergencySummary *EmergencySummaryCallerSession) DataHash() ([32]byte, error) {
	return _EmergencySummary.Contract.DataHash(&_EmergencySummary.CallOpts)
}

// DataVersion is a free data retrieval call binding the contract method 0xcb5f1b09.
//
// Solidity: function dataVersion() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCaller) DataVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "dataVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DataVersion is a free data retrieval call binding the contract method 0xcb5f1b09.
//
// Solidity: function dataVersion() view returns(uint256)
func (_EmergencySummary *EmergencySummarySession) DataVersion() (*big.Int, error) {
	return _EmergencySummary.Contract.DataVersion(&_EmergencySummary.CallOpts)
}

// DataVersion is a free data retrieval call binding the contract method 0xcb5f1b09.
//
// Solidity: function dataVersion() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCallerSession) DataVersion() (*big.Int, error) {
	return _EmergencySummary.Contract.DataVersion(&_EmergencySummary.CallOpts)
}

// EncryptedData is a free data retrieval call binding the contract method 0xe9fc715a.
//
// Solidity: function encryptedData() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCaller) EncryptedData(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "encryptedData")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncryptedData is a free data retrieval call binding the contract method 0xe9fc715a.
//
// Solidity: function encryptedData() view returns(bytes)
func (_EmergencySummary *EmergencySummarySession) EncryptedData() ([]byte, error) {
	return _EmergencySummary.Contract.EncryptedData(&_EmergencySummary.CallOpts)
}

// EncryptedData is a free data retrieval call binding the contract method 0xe9fc715a.
//
// Solidity: function encryptedData() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCallerSession) EncryptedData() ([]byte, error) {
	return _EmergencySummary.Contract.EncryptedData(&_EmergencySummary.CallOpts)
}

// GetAccessInfo is a free data retrieval call binding the contract method 0x4e07292f.
//
// Solidity: function getAccessInfo(address _user) view returns(bool isGranted, uint256 grantedAt, uint256 expiresAt)
func (_EmergencySummary *EmergencySummaryCaller) GetAccessInfo(opts *bind.CallOpts, _user common.Address) (struct {
	IsGranted bool
	GrantedAt *big.Int
	ExpiresAt *big.Int
}, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getAccessInfo", _user)

	outstruct := new(struct {
		IsGranted bool
		GrantedAt *big.Int
		ExpiresAt *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsGranted = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.GrantedAt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ExpiresAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetAccessInfo is a free data retrieval call binding the contract method 0x4e07292f.
//
// Solidity: function getAccessInfo(address _user) view returns(bool isGranted, uint256 grantedAt, uint256 expiresAt)
func (_EmergencySummary *EmergencySummarySession) GetAccessInfo(_user common.Address) (struct {
	IsGranted bool
	GrantedAt *big.Int
	ExpiresAt *big.Int
}, error) {
	return _EmergencySummary.Contract.GetAccessInfo(&_EmergencySummary.CallOpts, _user)
}

// GetAccessInfo is a free data retrieval call binding the contract method 0x4e07292f.
//
// Solidity: function getAccessInfo(address _user) view returns(bool isGranted, uint256 grantedAt, uint256 expiresAt)
func (_EmergencySummary *EmergencySummaryCallerSession) GetAccessInfo(_user common.Address) (struct {
	IsGranted bool
	GrantedAt *big.Int
	ExpiresAt *big.Int
}, error) {
	return _EmergencySummary.Contract.GetAccessInfo(&_EmergencySummary.CallOpts, _user)
}

// GetAccessLog is a free data retrieval call binding the contract method 0x096d8009.
//
// Solidity: function getAccessLog() view returns(uint256[])
func (_EmergencySummary *EmergencySummaryCaller) GetAccessLog(opts *bind.CallOpts) ([]*big.Int, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getAccessLog")

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAccessLog is a free data retrieval call binding the contract method 0x096d8009.
//
// Solidity: function getAccessLog() view returns(uint256[])
func (_EmergencySummary *EmergencySummarySession) GetAccessLog() ([]*big.Int, error) {
	return _EmergencySummary.Contract.GetAccessLog(&_EmergencySummary.CallOpts)
}

// GetAccessLog is a free data retrieval call binding the contract method 0x096d8009.
//
// Solidity: function getAccessLog() view returns(uint256[])
func (_EmergencySummary *EmergencySummaryCallerSession) GetAccessLog() ([]*big.Int, error) {
	return _EmergencySummary.Contract.GetAccessLog(&_EmergencySummary.CallOpts)
}

// GetAccessRequest is a free data retrieval call binding the contract method 0xee12bc71.
//
// Solidity: function getAccessRequest(uint256 _requestId) view returns(address responder, string location, string incidentDescription, uint256 timestamp, uint8 status)
func (_EmergencySummary *EmergencySummaryCaller) GetAccessRequest(opts *bind.CallOpts, _requestId *big.Int) (struct {
	Responder           common.Address
	Location            string
	IncidentDescription string
	Timestamp           *big.Int
	Status              uint8
}, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getAccessRequest", _requestId)

	outstruct := new(struct {
		Responder           common.Address
		Location            string
		IncidentDescription string
		Timestamp           *big.Int
		Status              uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Responder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Location = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.IncidentDescription = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Timestamp = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetAccessRequest is a free data retrieval call binding the contract method 0xee12bc71.
//
// Solidity: function getAccessRequest(uint256 _requestId) view returns(address responder, string location, string incidentDescription, uint256 timestamp, uint8 status)
func (_EmergencySummary *EmergencySummarySession) GetAccessRequest(_requestId *big.Int) (struct {
	Responder           common.Address
	Location            string
	IncidentDescription string
	Timestamp           *big.Int
	Status              uint8
}, error) {
	return _EmergencySummary.Contract.GetAccessRequest(&_EmergencySummary.CallOpts, _requestId)
}

// GetAccessRequest is a free data retrieval call binding the contract method 0xee12bc71.
//
// Solidity: function getAccessRequest(uint256 _requestId) view returns(address responder, string location, string incidentDescription, uint256 timestamp, uint8 status)
func (_EmergencySummary *EmergencySummaryCallerSession) GetAccessRequest(_requestId *big.Int) (struct {
	Responder           common.Address
	Location            string
	IncidentDescription string
	Timestamp           *big.Int
	Status              uint8
}, error) {
	return _EmergencySummary.Contract.GetAccessRequest(&_EmergencySummary.CallOpts, _requestId)
}

// GetAuthorizedUsers is a free data retrieval call binding the contract method 0x1ba7214c.
//
// Solidity: function getAuthorizedUsers() view returns(address[])
func (_EmergencySummary *EmergencySummaryCaller) GetAuthorizedUsers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getAuthorizedUsers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAuthorizedUsers is a free data retrieval call binding the contract method 0x1ba7214c.
//
// Solidity: function getAuthorizedUsers() view returns(address[])
func (_EmergencySummary *EmergencySummarySession) GetAuthorizedUsers() ([]common.Address, error) {
	return _EmergencySummary.Contract.GetAuthorizedUsers(&_EmergencySummary.CallOpts)
}

// GetAuthorizedUsers is a free data retrieval call binding the contract method 0x1ba7214c.
//
// Solidity: function getAuthorizedUsers() view returns(address[])
func (_EmergencySummary *EmergencySummaryCallerSession) GetAuthorizedUsers() ([]common.Address, error) {
	return _EmergencySummary.Contract.GetAuthorizedUsers(&_EmergencySummary.CallOpts)
}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(bytes _encryptedData, bytes _encryptedKey, uint256 _version, bytes32 _hash)
func (_EmergencySummary *EmergencySummaryCaller) GetData(opts *bind.CallOpts) (struct {
	EncryptedData []byte
	EncryptedKey  []byte
	Version       *big.Int
	Hash          [32]byte
}, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getData")

	outstruct := new(struct {
		EncryptedData []byte
		EncryptedKey  []byte
		Version       *big.Int
		Hash          [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.EncryptedData = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.EncryptedKey = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Version = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Hash = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(bytes _encryptedData, bytes _encryptedKey, uint256 _version, bytes32 _hash)
func (_EmergencySummary *EmergencySummarySession) GetData() (struct {
	EncryptedData []byte
	EncryptedKey  []byte
	Version       *big.Int
	Hash          [32]byte
}, error) {
	return _EmergencySummary.Contract.GetData(&_EmergencySummary.CallOpts)
}

// GetData is a free data retrieval call binding the contract method 0x3bc5de30.
//
// Solidity: function getData() view returns(bytes _encryptedData, bytes _encryptedKey, uint256 _version, bytes32 _hash)
func (_EmergencySummary *EmergencySummaryCallerSession) GetData() (struct {
	EncryptedData []byte
	EncryptedKey  []byte
	Version       *big.Int
	Hash          [32]byte
}, error) {
	return _EmergencySummary.Contract.GetData(&_EmergencySummary.CallOpts)
}

// GetEncryptedData is a free data retrieval call binding the contract method 0x7ef508bb.
//
// Solidity: function getEncryptedData() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCaller) GetEncryptedData(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getEncryptedData")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetEncryptedData is a free data retrieval call binding the contract method 0x7ef508bb.
//
// Solidity: function getEncryptedData() view returns(bytes)
func (_EmergencySummary *EmergencySummarySession) GetEncryptedData() ([]byte, error) {
	return _EmergencySummary.Contract.GetEncryptedData(&_EmergencySummary.CallOpts)
}

// GetEncryptedData is a free data retrieval call binding the contract method 0x7ef508bb.
//
// Solidity: function getEncryptedData() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCallerSession) GetEncryptedData() ([]byte, error) {
	return _EmergencySummary.Contract.GetEncryptedData(&_EmergencySummary.CallOpts)
}

// GetEncryptedKey is a free data retrieval call binding the contract method 0xa1bc5a83.
//
// Solidity: function getEncryptedKey() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCaller) GetEncryptedKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getEncryptedKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetEncryptedKey is a free data retrieval call binding the contract method 0xa1bc5a83.
//
// Solidity: function getEncryptedKey() view returns(bytes)
func (_EmergencySummary *EmergencySummarySession) GetEncryptedKey() ([]byte, error) {
	return _EmergencySummary.Contract.GetEncryptedKey(&_EmergencySummary.CallOpts)
}

// GetEncryptedKey is a free data retrieval call binding the contract method 0xa1bc5a83.
//
// Solidity: function getEncryptedKey() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCallerSession) GetEncryptedKey() ([]byte, error) {
	return _EmergencySummary.Contract.GetEncryptedKey(&_EmergencySummary.CallOpts)
}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns(uint256 version, uint256 updated, uint256 dataSize)
func (_EmergencySummary *EmergencySummaryCaller) GetMetadata(opts *bind.CallOpts) (struct {
	Version  *big.Int
	Updated  *big.Int
	DataSize *big.Int
}, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getMetadata")

	outstruct := new(struct {
		Version  *big.Int
		Updated  *big.Int
		DataSize *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Updated = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DataSize = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns(uint256 version, uint256 updated, uint256 dataSize)
func (_EmergencySummary *EmergencySummarySession) GetMetadata() (struct {
	Version  *big.Int
	Updated  *big.Int
	DataSize *big.Int
}, error) {
	return _EmergencySummary.Contract.GetMetadata(&_EmergencySummary.CallOpts)
}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns(uint256 version, uint256 updated, uint256 dataSize)
func (_EmergencySummary *EmergencySummaryCallerSession) GetMetadata() (struct {
	Version  *big.Int
	Updated  *big.Int
	DataSize *big.Int
}, error) {
	return _EmergencySummary.Contract.GetMetadata(&_EmergencySummary.CallOpts)
}

// GetRecoveryKey is a free data retrieval call binding the contract method 0xfc4a7fca.
//
// Solidity: function getRecoveryKey() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCaller) GetRecoveryKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "getRecoveryKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRecoveryKey is a free data retrieval call binding the contract method 0xfc4a7fca.
//
// Solidity: function getRecoveryKey() view returns(bytes)
func (_EmergencySummary *EmergencySummarySession) GetRecoveryKey() ([]byte, error) {
	return _EmergencySummary.Contract.GetRecoveryKey(&_EmergencySummary.CallOpts)
}

// GetRecoveryKey is a free data retrieval call binding the contract method 0xfc4a7fca.
//
// Solidity: function getRecoveryKey() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCallerSession) GetRecoveryKey() ([]byte, error) {
	return _EmergencySummary.Contract.GetRecoveryKey(&_EmergencySummary.CallOpts)
}

// HasAccess is a free data retrieval call binding the contract method 0x95a078e8.
//
// Solidity: function hasAccess(address _user) view returns(bool)
func (_EmergencySummary *EmergencySummaryCaller) HasAccess(opts *bind.CallOpts, _user common.Address) (bool, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "hasAccess", _user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasAccess is a free data retrieval call binding the contract method 0x95a078e8.
//
// Solidity: function hasAccess(address _user) view returns(bool)
func (_EmergencySummary *EmergencySummarySession) HasAccess(_user common.Address) (bool, error) {
	return _EmergencySummary.Contract.HasAccess(&_EmergencySummary.CallOpts, _user)
}

// HasAccess is a free data retrieval call binding the contract method 0x95a078e8.
//
// Solidity: function hasAccess(address _user) view returns(bool)
func (_EmergencySummary *EmergencySummaryCallerSession) HasAccess(_user common.Address) (bool, error) {
	return _EmergencySummary.Contract.HasAccess(&_EmergencySummary.CallOpts, _user)
}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCaller) LastUpdated(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "lastUpdated")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns(uint256)
func (_EmergencySummary *EmergencySummarySession) LastUpdated() (*big.Int, error) {
	return _EmergencySummary.Contract.LastUpdated(&_EmergencySummary.CallOpts)
}

// LastUpdated is a free data retrieval call binding the contract method 0xd0b06f5d.
//
// Solidity: function lastUpdated() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCallerSession) LastUpdated() (*big.Int, error) {
	return _EmergencySummary.Contract.LastUpdated(&_EmergencySummary.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EmergencySummary *EmergencySummaryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EmergencySummary *EmergencySummarySession) Owner() (common.Address, error) {
	return _EmergencySummary.Contract.Owner(&_EmergencySummary.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EmergencySummary *EmergencySummaryCallerSession) Owner() (common.Address, error) {
	return _EmergencySummary.Contract.Owner(&_EmergencySummary.CallOpts)
}

// RecoveryKey is a free data retrieval call binding the contract method 0x97905a88.
//
// Solidity: function recoveryKey() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCaller) RecoveryKey(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "recoveryKey")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RecoveryKey is a free data retrieval call binding the contract method 0x97905a88.
//
// Solidity: function recoveryKey() view returns(bytes)
func (_EmergencySummary *EmergencySummarySession) RecoveryKey() ([]byte, error) {
	return _EmergencySummary.Contract.RecoveryKey(&_EmergencySummary.CallOpts)
}

// RecoveryKey is a free data retrieval call binding the contract method 0x97905a88.
//
// Solidity: function recoveryKey() view returns(bytes)
func (_EmergencySummary *EmergencySummaryCallerSession) RecoveryKey() ([]byte, error) {
	return _EmergencySummary.Contract.RecoveryKey(&_EmergencySummary.CallOpts)
}

// RecoveryVersion is a free data retrieval call binding the contract method 0x4ae1ec73.
//
// Solidity: function recoveryVersion() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCaller) RecoveryVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "recoveryVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RecoveryVersion is a free data retrieval call binding the contract method 0x4ae1ec73.
//
// Solidity: function recoveryVersion() view returns(uint256)
func (_EmergencySummary *EmergencySummarySession) RecoveryVersion() (*big.Int, error) {
	return _EmergencySummary.Contract.RecoveryVersion(&_EmergencySummary.CallOpts)
}

// RecoveryVersion is a free data retrieval call binding the contract method 0x4ae1ec73.
//
// Solidity: function recoveryVersion() view returns(uint256)
func (_EmergencySummary *EmergencySummaryCallerSession) RecoveryVersion() (*big.Int, error) {
	return _EmergencySummary.Contract.RecoveryVersion(&_EmergencySummary.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_EmergencySummary *EmergencySummaryCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_EmergencySummary *EmergencySummarySession) Registry() (common.Address, error) {
	return _EmergencySummary.Contract.Registry(&_EmergencySummary.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_EmergencySummary *EmergencySummaryCallerSession) Registry() (common.Address, error) {
	return _EmergencySummary.Contract.Registry(&_EmergencySummary.CallOpts)
}

// RiskAssessement is a free data retrieval call binding the contract method 0x79b2fc7a.
//
// Solidity: function riskAssessement() view returns(address)
func (_EmergencySummary *EmergencySummaryCaller) RiskAssessement(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EmergencySummary.contract.Call(opts, &out, "riskAssessement")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RiskAssessement is a free data retrieval call binding the contract method 0x79b2fc7a.
//
// Solidity: function riskAssessement() view returns(address)
func (_EmergencySummary *EmergencySummarySession) RiskAssessement() (common.Address, error) {
	return _EmergencySummary.Contract.RiskAssessement(&_EmergencySummary.CallOpts)
}

// RiskAssessement is a free data retrieval call binding the contract method 0x79b2fc7a.
//
// Solidity: function riskAssessement() view returns(address)
func (_EmergencySummary *EmergencySummaryCallerSession) RiskAssessement() (common.Address, error) {
	return _EmergencySummary.Contract.RiskAssessement(&_EmergencySummary.CallOpts)
}

// ApproveAccessRequest is a paid mutator transaction binding the contract method 0x84ab34e9.
//
// Solidity: function approveAccessRequest(uint256 _requestId, bytes _encryptedKey) returns()
func (_EmergencySummary *EmergencySummaryTransactor) ApproveAccessRequest(opts *bind.TransactOpts, _requestId *big.Int, _encryptedKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "approveAccessRequest", _requestId, _encryptedKey)
}

// ApproveAccessRequest is a paid mutator transaction binding the contract method 0x84ab34e9.
//
// Solidity: function approveAccessRequest(uint256 _requestId, bytes _encryptedKey) returns()
func (_EmergencySummary *EmergencySummarySession) ApproveAccessRequest(_requestId *big.Int, _encryptedKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.ApproveAccessRequest(&_EmergencySummary.TransactOpts, _requestId, _encryptedKey)
}

// ApproveAccessRequest is a paid mutator transaction binding the contract method 0x84ab34e9.
//
// Solidity: function approveAccessRequest(uint256 _requestId, bytes _encryptedKey) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) ApproveAccessRequest(_requestId *big.Int, _encryptedKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.ApproveAccessRequest(&_EmergencySummary.TransactOpts, _requestId, _encryptedKey)
}

// DenyAccessRequest is a paid mutator transaction binding the contract method 0x66f697b8.
//
// Solidity: function denyAccessRequest(uint256 _requestId) returns()
func (_EmergencySummary *EmergencySummaryTransactor) DenyAccessRequest(opts *bind.TransactOpts, _requestId *big.Int) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "denyAccessRequest", _requestId)
}

// DenyAccessRequest is a paid mutator transaction binding the contract method 0x66f697b8.
//
// Solidity: function denyAccessRequest(uint256 _requestId) returns()
func (_EmergencySummary *EmergencySummarySession) DenyAccessRequest(_requestId *big.Int) (*types.Transaction, error) {
	return _EmergencySummary.Contract.DenyAccessRequest(&_EmergencySummary.TransactOpts, _requestId)
}

// DenyAccessRequest is a paid mutator transaction binding the contract method 0x66f697b8.
//
// Solidity: function denyAccessRequest(uint256 _requestId) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) DenyAccessRequest(_requestId *big.Int) (*types.Transaction, error) {
	return _EmergencySummary.Contract.DenyAccessRequest(&_EmergencySummary.TransactOpts, _requestId)
}

// GrantAccess is a paid mutator transaction binding the contract method 0xf72bcde6.
//
// Solidity: function grantAccess(address _user, bytes _encryptedKey, uint256 _duration) returns()
func (_EmergencySummary *EmergencySummaryTransactor) GrantAccess(opts *bind.TransactOpts, _user common.Address, _encryptedKey []byte, _duration *big.Int) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "grantAccess", _user, _encryptedKey, _duration)
}

// GrantAccess is a paid mutator transaction binding the contract method 0xf72bcde6.
//
// Solidity: function grantAccess(address _user, bytes _encryptedKey, uint256 _duration) returns()
func (_EmergencySummary *EmergencySummarySession) GrantAccess(_user common.Address, _encryptedKey []byte, _duration *big.Int) (*types.Transaction, error) {
	return _EmergencySummary.Contract.GrantAccess(&_EmergencySummary.TransactOpts, _user, _encryptedKey, _duration)
}

// GrantAccess is a paid mutator transaction binding the contract method 0xf72bcde6.
//
// Solidity: function grantAccess(address _user, bytes _encryptedKey, uint256 _duration) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) GrantAccess(_user common.Address, _encryptedKey []byte, _duration *big.Int) (*types.Transaction, error) {
	return _EmergencySummary.Contract.GrantAccess(&_EmergencySummary.TransactOpts, _user, _encryptedKey, _duration)
}

// RegisterPublicKey is a paid mutator transaction binding the contract method 0x85623594.
//
// Solidity: function registerPublicKey(bytes _publicKey) returns()
func (_EmergencySummary *EmergencySummaryTransactor) RegisterPublicKey(opts *bind.TransactOpts, _publicKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "registerPublicKey", _publicKey)
}

// RegisterPublicKey is a paid mutator transaction binding the contract method 0x85623594.
//
// Solidity: function registerPublicKey(bytes _publicKey) returns()
func (_EmergencySummary *EmergencySummarySession) RegisterPublicKey(_publicKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.RegisterPublicKey(&_EmergencySummary.TransactOpts, _publicKey)
}

// RegisterPublicKey is a paid mutator transaction binding the contract method 0x85623594.
//
// Solidity: function registerPublicKey(bytes _publicKey) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) RegisterPublicKey(_publicKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.RegisterPublicKey(&_EmergencySummary.TransactOpts, _publicKey)
}

// RequestAccess is a paid mutator transaction binding the contract method 0x4f17e037.
//
// Solidity: function requestAccess(bytes _publicKey, string _location, string _incidentDescription, string _patientIdentifier) returns(uint256 requestId)
func (_EmergencySummary *EmergencySummaryTransactor) RequestAccess(opts *bind.TransactOpts, _publicKey []byte, _location string, _incidentDescription string, _patientIdentifier string) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "requestAccess", _publicKey, _location, _incidentDescription, _patientIdentifier)
}

// RequestAccess is a paid mutator transaction binding the contract method 0x4f17e037.
//
// Solidity: function requestAccess(bytes _publicKey, string _location, string _incidentDescription, string _patientIdentifier) returns(uint256 requestId)
func (_EmergencySummary *EmergencySummarySession) RequestAccess(_publicKey []byte, _location string, _incidentDescription string, _patientIdentifier string) (*types.Transaction, error) {
	return _EmergencySummary.Contract.RequestAccess(&_EmergencySummary.TransactOpts, _publicKey, _location, _incidentDescription, _patientIdentifier)
}

// RequestAccess is a paid mutator transaction binding the contract method 0x4f17e037.
//
// Solidity: function requestAccess(bytes _publicKey, string _location, string _incidentDescription, string _patientIdentifier) returns(uint256 requestId)
func (_EmergencySummary *EmergencySummaryTransactorSession) RequestAccess(_publicKey []byte, _location string, _incidentDescription string, _patientIdentifier string) (*types.Transaction, error) {
	return _EmergencySummary.Contract.RequestAccess(&_EmergencySummary.TransactOpts, _publicKey, _location, _incidentDescription, _patientIdentifier)
}

// StoreRecoveryKey is a paid mutator transaction binding the contract method 0x0a78d59c.
//
// Solidity: function storeRecoveryKey(bytes _recoveryKey) returns()
func (_EmergencySummary *EmergencySummaryTransactor) StoreRecoveryKey(opts *bind.TransactOpts, _recoveryKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "storeRecoveryKey", _recoveryKey)
}

// StoreRecoveryKey is a paid mutator transaction binding the contract method 0x0a78d59c.
//
// Solidity: function storeRecoveryKey(bytes _recoveryKey) returns()
func (_EmergencySummary *EmergencySummarySession) StoreRecoveryKey(_recoveryKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.StoreRecoveryKey(&_EmergencySummary.TransactOpts, _recoveryKey)
}

// StoreRecoveryKey is a paid mutator transaction binding the contract method 0x0a78d59c.
//
// Solidity: function storeRecoveryKey(bytes _recoveryKey) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) StoreRecoveryKey(_recoveryKey []byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.StoreRecoveryKey(&_EmergencySummary.TransactOpts, _recoveryKey)
}

// UpdateData is a paid mutator transaction binding the contract method 0x0f07e485.
//
// Solidity: function updateData(bytes _encryptedData, bytes32 _dataHash) returns()
func (_EmergencySummary *EmergencySummaryTransactor) UpdateData(opts *bind.TransactOpts, _encryptedData []byte, _dataHash [32]byte) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "updateData", _encryptedData, _dataHash)
}

// UpdateData is a paid mutator transaction binding the contract method 0x0f07e485.
//
// Solidity: function updateData(bytes _encryptedData, bytes32 _dataHash) returns()
func (_EmergencySummary *EmergencySummarySession) UpdateData(_encryptedData []byte, _dataHash [32]byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.UpdateData(&_EmergencySummary.TransactOpts, _encryptedData, _dataHash)
}

// UpdateData is a paid mutator transaction binding the contract method 0x0f07e485.
//
// Solidity: function updateData(bytes _encryptedData, bytes32 _dataHash) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) UpdateData(_encryptedData []byte, _dataHash [32]byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.UpdateData(&_EmergencySummary.TransactOpts, _encryptedData, _dataHash)
}

// UpdateKeys is a paid mutator transaction binding the contract method 0x57e1f3d6.
//
// Solidity: function updateKeys(address[] _users, bytes[] _encryptedKeys) returns()
func (_EmergencySummary *EmergencySummaryTransactor) UpdateKeys(opts *bind.TransactOpts, _users []common.Address, _encryptedKeys [][]byte) (*types.Transaction, error) {
	return _EmergencySummary.contract.Transact(opts, "updateKeys", _users, _encryptedKeys)
}

// UpdateKeys is a paid mutator transaction binding the contract method 0x57e1f3d6.
//
// Solidity: function updateKeys(address[] _users, bytes[] _encryptedKeys) returns()
func (_EmergencySummary *EmergencySummarySession) UpdateKeys(_users []common.Address, _encryptedKeys [][]byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.UpdateKeys(&_EmergencySummary.TransactOpts, _users, _encryptedKeys)
}

// UpdateKeys is a paid mutator transaction binding the contract method 0x57e1f3d6.
//
// Solidity: function updateKeys(address[] _users, bytes[] _encryptedKeys) returns()
func (_EmergencySummary *EmergencySummaryTransactorSession) UpdateKeys(_users []common.Address, _encryptedKeys [][]byte) (*types.Transaction, error) {
	return _EmergencySummary.Contract.UpdateKeys(&_EmergencySummary.TransactOpts, _users, _encryptedKeys)
}

// EmergencySummaryAccessGrantedIterator is returned from FilterAccessGranted and is used to iterate over the raw logs and unpacked data for AccessGranted events raised by the EmergencySummary contract.
type EmergencySummaryAccessGrantedIterator struct {
	Event *EmergencySummaryAccessGranted // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryAccessGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryAccessGranted)
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
		it.Event = new(EmergencySummaryAccessGranted)
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
func (it *EmergencySummaryAccessGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryAccessGranted represents a AccessGranted event raised by the EmergencySummary contract.
type EmergencySummaryAccessGranted struct {
	User      common.Address
	ExpiresAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAccessGranted is a free log retrieval operation binding the contract event 0xb4c6779ceb4a20f448e76a0e11f39bd183cff9c9dbac53df6bfcc202e2eb32f1.
//
// Solidity: event AccessGranted(address indexed user, uint256 expiresAt)
func (_EmergencySummary *EmergencySummaryFilterer) FilterAccessGranted(opts *bind.FilterOpts, user []common.Address) (*EmergencySummaryAccessGrantedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EmergencySummary.contract.FilterLogs(opts, "AccessGranted", userRule)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryAccessGrantedIterator{contract: _EmergencySummary.contract, event: "AccessGranted", logs: logs, sub: sub}, nil
}

// WatchAccessGranted is a free log subscription operation binding the contract event 0xb4c6779ceb4a20f448e76a0e11f39bd183cff9c9dbac53df6bfcc202e2eb32f1.
//
// Solidity: event AccessGranted(address indexed user, uint256 expiresAt)
func (_EmergencySummary *EmergencySummaryFilterer) WatchAccessGranted(opts *bind.WatchOpts, sink chan<- *EmergencySummaryAccessGranted, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EmergencySummary.contract.WatchLogs(opts, "AccessGranted", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryAccessGranted)
				if err := _EmergencySummary.contract.UnpackLog(event, "AccessGranted", log); err != nil {
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

// ParseAccessGranted is a log parse operation binding the contract event 0xb4c6779ceb4a20f448e76a0e11f39bd183cff9c9dbac53df6bfcc202e2eb32f1.
//
// Solidity: event AccessGranted(address indexed user, uint256 expiresAt)
func (_EmergencySummary *EmergencySummaryFilterer) ParseAccessGranted(log types.Log) (*EmergencySummaryAccessGranted, error) {
	event := new(EmergencySummaryAccessGranted)
	if err := _EmergencySummary.contract.UnpackLog(event, "AccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmergencySummaryAccessRequestedIterator is returned from FilterAccessRequested and is used to iterate over the raw logs and unpacked data for AccessRequested events raised by the EmergencySummary contract.
type EmergencySummaryAccessRequestedIterator struct {
	Event *EmergencySummaryAccessRequested // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryAccessRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryAccessRequested)
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
		it.Event = new(EmergencySummaryAccessRequested)
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
func (it *EmergencySummaryAccessRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryAccessRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryAccessRequested represents a AccessRequested event raised by the EmergencySummary contract.
type EmergencySummaryAccessRequested struct {
	RequestId *big.Int
	Responder common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAccessRequested is a free log retrieval operation binding the contract event 0xfdb34af38bf2bba5bb086deda26262c9e7cd30b4a10a79d9297b60760f4a251f.
//
// Solidity: event AccessRequested(uint256 indexed requestId, address indexed responder, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) FilterAccessRequested(opts *bind.FilterOpts, requestId []*big.Int, responder []common.Address) (*EmergencySummaryAccessRequestedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var responderRule []interface{}
	for _, responderItem := range responder {
		responderRule = append(responderRule, responderItem)
	}

	logs, sub, err := _EmergencySummary.contract.FilterLogs(opts, "AccessRequested", requestIdRule, responderRule)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryAccessRequestedIterator{contract: _EmergencySummary.contract, event: "AccessRequested", logs: logs, sub: sub}, nil
}

// WatchAccessRequested is a free log subscription operation binding the contract event 0xfdb34af38bf2bba5bb086deda26262c9e7cd30b4a10a79d9297b60760f4a251f.
//
// Solidity: event AccessRequested(uint256 indexed requestId, address indexed responder, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) WatchAccessRequested(opts *bind.WatchOpts, sink chan<- *EmergencySummaryAccessRequested, requestId []*big.Int, responder []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var responderRule []interface{}
	for _, responderItem := range responder {
		responderRule = append(responderRule, responderItem)
	}

	logs, sub, err := _EmergencySummary.contract.WatchLogs(opts, "AccessRequested", requestIdRule, responderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryAccessRequested)
				if err := _EmergencySummary.contract.UnpackLog(event, "AccessRequested", log); err != nil {
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

// ParseAccessRequested is a log parse operation binding the contract event 0xfdb34af38bf2bba5bb086deda26262c9e7cd30b4a10a79d9297b60760f4a251f.
//
// Solidity: event AccessRequested(uint256 indexed requestId, address indexed responder, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) ParseAccessRequested(log types.Log) (*EmergencySummaryAccessRequested, error) {
	event := new(EmergencySummaryAccessRequested)
	if err := _EmergencySummary.contract.UnpackLog(event, "AccessRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmergencySummaryAccessRevokedIterator is returned from FilterAccessRevoked and is used to iterate over the raw logs and unpacked data for AccessRevoked events raised by the EmergencySummary contract.
type EmergencySummaryAccessRevokedIterator struct {
	Event *EmergencySummaryAccessRevoked // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryAccessRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryAccessRevoked)
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
		it.Event = new(EmergencySummaryAccessRevoked)
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
func (it *EmergencySummaryAccessRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryAccessRevoked represents a AccessRevoked event raised by the EmergencySummary contract.
type EmergencySummaryAccessRevoked struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAccessRevoked is a free log retrieval operation binding the contract event 0x1b9b72fde9da721e70e6aca3b0cf4cbe73e82765ef1f280157740376531bfdd8.
//
// Solidity: event AccessRevoked(address indexed user)
func (_EmergencySummary *EmergencySummaryFilterer) FilterAccessRevoked(opts *bind.FilterOpts, user []common.Address) (*EmergencySummaryAccessRevokedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EmergencySummary.contract.FilterLogs(opts, "AccessRevoked", userRule)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryAccessRevokedIterator{contract: _EmergencySummary.contract, event: "AccessRevoked", logs: logs, sub: sub}, nil
}

// WatchAccessRevoked is a free log subscription operation binding the contract event 0x1b9b72fde9da721e70e6aca3b0cf4cbe73e82765ef1f280157740376531bfdd8.
//
// Solidity: event AccessRevoked(address indexed user)
func (_EmergencySummary *EmergencySummaryFilterer) WatchAccessRevoked(opts *bind.WatchOpts, sink chan<- *EmergencySummaryAccessRevoked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EmergencySummary.contract.WatchLogs(opts, "AccessRevoked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryAccessRevoked)
				if err := _EmergencySummary.contract.UnpackLog(event, "AccessRevoked", log); err != nil {
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

// ParseAccessRevoked is a log parse operation binding the contract event 0x1b9b72fde9da721e70e6aca3b0cf4cbe73e82765ef1f280157740376531bfdd8.
//
// Solidity: event AccessRevoked(address indexed user)
func (_EmergencySummary *EmergencySummaryFilterer) ParseAccessRevoked(log types.Log) (*EmergencySummaryAccessRevoked, error) {
	event := new(EmergencySummaryAccessRevoked)
	if err := _EmergencySummary.contract.UnpackLog(event, "AccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmergencySummaryDataUpdatedIterator is returned from FilterDataUpdated and is used to iterate over the raw logs and unpacked data for DataUpdated events raised by the EmergencySummary contract.
type EmergencySummaryDataUpdatedIterator struct {
	Event *EmergencySummaryDataUpdated // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryDataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryDataUpdated)
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
		it.Event = new(EmergencySummaryDataUpdated)
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
func (it *EmergencySummaryDataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryDataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryDataUpdated represents a DataUpdated event raised by the EmergencySummary contract.
type EmergencySummaryDataUpdated struct {
	Version   *big.Int
	DataHash  [32]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDataUpdated is a free log retrieval operation binding the contract event 0xdcb8c1999099ea0e5ee2598d27790e847d62a30ab3916dc4b3302baba949b09f.
//
// Solidity: event DataUpdated(uint256 version, bytes32 dataHash, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) FilterDataUpdated(opts *bind.FilterOpts) (*EmergencySummaryDataUpdatedIterator, error) {

	logs, sub, err := _EmergencySummary.contract.FilterLogs(opts, "DataUpdated")
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryDataUpdatedIterator{contract: _EmergencySummary.contract, event: "DataUpdated", logs: logs, sub: sub}, nil
}

// WatchDataUpdated is a free log subscription operation binding the contract event 0xdcb8c1999099ea0e5ee2598d27790e847d62a30ab3916dc4b3302baba949b09f.
//
// Solidity: event DataUpdated(uint256 version, bytes32 dataHash, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) WatchDataUpdated(opts *bind.WatchOpts, sink chan<- *EmergencySummaryDataUpdated) (event.Subscription, error) {

	logs, sub, err := _EmergencySummary.contract.WatchLogs(opts, "DataUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryDataUpdated)
				if err := _EmergencySummary.contract.UnpackLog(event, "DataUpdated", log); err != nil {
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

// ParseDataUpdated is a log parse operation binding the contract event 0xdcb8c1999099ea0e5ee2598d27790e847d62a30ab3916dc4b3302baba949b09f.
//
// Solidity: event DataUpdated(uint256 version, bytes32 dataHash, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) ParseDataUpdated(log types.Log) (*EmergencySummaryDataUpdated, error) {
	event := new(EmergencySummaryDataUpdated)
	if err := _EmergencySummary.contract.UnpackLog(event, "DataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmergencySummaryPublicKeyRegisteredIterator is returned from FilterPublicKeyRegistered and is used to iterate over the raw logs and unpacked data for PublicKeyRegistered events raised by the EmergencySummary contract.
type EmergencySummaryPublicKeyRegisteredIterator struct {
	Event *EmergencySummaryPublicKeyRegistered // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryPublicKeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryPublicKeyRegistered)
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
		it.Event = new(EmergencySummaryPublicKeyRegistered)
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
func (it *EmergencySummaryPublicKeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryPublicKeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryPublicKeyRegistered represents a PublicKeyRegistered event raised by the EmergencySummary contract.
type EmergencySummaryPublicKeyRegistered struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPublicKeyRegistered is a free log retrieval operation binding the contract event 0xfa3256bc03b03ee751ee91c53b7318d6e4321ab4d09757c08a50a3f65f678300.
//
// Solidity: event PublicKeyRegistered(address indexed user)
func (_EmergencySummary *EmergencySummaryFilterer) FilterPublicKeyRegistered(opts *bind.FilterOpts, user []common.Address) (*EmergencySummaryPublicKeyRegisteredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EmergencySummary.contract.FilterLogs(opts, "PublicKeyRegistered", userRule)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryPublicKeyRegisteredIterator{contract: _EmergencySummary.contract, event: "PublicKeyRegistered", logs: logs, sub: sub}, nil
}

// WatchPublicKeyRegistered is a free log subscription operation binding the contract event 0xfa3256bc03b03ee751ee91c53b7318d6e4321ab4d09757c08a50a3f65f678300.
//
// Solidity: event PublicKeyRegistered(address indexed user)
func (_EmergencySummary *EmergencySummaryFilterer) WatchPublicKeyRegistered(opts *bind.WatchOpts, sink chan<- *EmergencySummaryPublicKeyRegistered, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _EmergencySummary.contract.WatchLogs(opts, "PublicKeyRegistered", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryPublicKeyRegistered)
				if err := _EmergencySummary.contract.UnpackLog(event, "PublicKeyRegistered", log); err != nil {
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

// ParsePublicKeyRegistered is a log parse operation binding the contract event 0xfa3256bc03b03ee751ee91c53b7318d6e4321ab4d09757c08a50a3f65f678300.
//
// Solidity: event PublicKeyRegistered(address indexed user)
func (_EmergencySummary *EmergencySummaryFilterer) ParsePublicKeyRegistered(log types.Log) (*EmergencySummaryPublicKeyRegistered, error) {
	event := new(EmergencySummaryPublicKeyRegistered)
	if err := _EmergencySummary.contract.UnpackLog(event, "PublicKeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EmergencySummaryRecoveryKeyStoredIterator is returned from FilterRecoveryKeyStored and is used to iterate over the raw logs and unpacked data for RecoveryKeyStored events raised by the EmergencySummary contract.
type EmergencySummaryRecoveryKeyStoredIterator struct {
	Event *EmergencySummaryRecoveryKeyStored // Event containing the contract specifics and raw log

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
func (it *EmergencySummaryRecoveryKeyStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EmergencySummaryRecoveryKeyStored)
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
		it.Event = new(EmergencySummaryRecoveryKeyStored)
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
func (it *EmergencySummaryRecoveryKeyStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EmergencySummaryRecoveryKeyStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EmergencySummaryRecoveryKeyStored represents a RecoveryKeyStored event raised by the EmergencySummary contract.
type EmergencySummaryRecoveryKeyStored struct {
	Version   *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRecoveryKeyStored is a free log retrieval operation binding the contract event 0xd4ef68535320f052f7e62f0505a9ad40fa87a5c6ad2f7cf705928257ba4a4184.
//
// Solidity: event RecoveryKeyStored(uint256 indexed version, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) FilterRecoveryKeyStored(opts *bind.FilterOpts, version []*big.Int) (*EmergencySummaryRecoveryKeyStoredIterator, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _EmergencySummary.contract.FilterLogs(opts, "RecoveryKeyStored", versionRule)
	if err != nil {
		return nil, err
	}
	return &EmergencySummaryRecoveryKeyStoredIterator{contract: _EmergencySummary.contract, event: "RecoveryKeyStored", logs: logs, sub: sub}, nil
}

// WatchRecoveryKeyStored is a free log subscription operation binding the contract event 0xd4ef68535320f052f7e62f0505a9ad40fa87a5c6ad2f7cf705928257ba4a4184.
//
// Solidity: event RecoveryKeyStored(uint256 indexed version, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) WatchRecoveryKeyStored(opts *bind.WatchOpts, sink chan<- *EmergencySummaryRecoveryKeyStored, version []*big.Int) (event.Subscription, error) {

	var versionRule []interface{}
	for _, versionItem := range version {
		versionRule = append(versionRule, versionItem)
	}

	logs, sub, err := _EmergencySummary.contract.WatchLogs(opts, "RecoveryKeyStored", versionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EmergencySummaryRecoveryKeyStored)
				if err := _EmergencySummary.contract.UnpackLog(event, "RecoveryKeyStored", log); err != nil {
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

// ParseRecoveryKeyStored is a log parse operation binding the contract event 0xd4ef68535320f052f7e62f0505a9ad40fa87a5c6ad2f7cf705928257ba4a4184.
//
// Solidity: event RecoveryKeyStored(uint256 indexed version, uint256 timestamp)
func (_EmergencySummary *EmergencySummaryFilterer) ParseRecoveryKeyStored(log types.Log) (*EmergencySummaryRecoveryKeyStored, error) {
	event := new(EmergencySummaryRecoveryKeyStored)
	if err := _EmergencySummary.contract.UnpackLog(event, "RecoveryKeyStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
