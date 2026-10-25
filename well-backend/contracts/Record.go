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
type TypesRecordMetadata struct {
	RecordType uint8
	CreatedAt  *big.Int
	UpdatedAt  *big.Int
}

// RecordMetaData contains all meta data concerning the Record contract.
var RecordMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_uploader\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_metadata\",\"type\":\"tuple\",\"internalType\":\"structTypes.RecordMetadata\",\"components\":[{\"name\":\"recordType\",\"type\":\"uint8\",\"internalType\":\"enumTypes.RecordType\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_contentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_storageReference\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_policyVerifier\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"POLICY_PROOF_TTL\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contentHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPolicy\",\"inputs\":[],\"outputs\":[{\"name\":\"deptCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"instCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecord\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structTypes.RecordMetadata\",\"components\":[{\"name\":\"recordType\",\"type\":\"uint8\",\"internalType\":\"enumTypes.RecordType\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"_contentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_storageReference\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_uploader\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantAccess\",\"inputs\":[{\"name\":\"doctor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"grantAccessWithProof\",\"inputs\":[{\"name\":\"pA\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pB\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pC\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"deptCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"instCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"doctorCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"evaluatedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"history\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAuthorized\",\"inputs\":[{\"name\":\"doctor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"metadata\",\"inputs\":[],\"outputs\":[{\"name\":\"recordType\",\"type\":\"uint8\",\"internalType\":\"enumTypes.RecordType\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"updatedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permissions\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractRecordPermissions\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policyVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractPolicyVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requiredDeptCommitment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requiredInstCommitment\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revokeAccess\",\"inputs\":[{\"name\":\"doctor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPolicy\",\"inputs\":[{\"name\":\"_requiredDeptCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_requiredInstCommitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storageReference\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateIntegrityHash\",\"inputs\":[{\"name\":\"newHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uploader\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyIntegrity\",\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AccessGranted\",\"inputs\":[{\"name\":\"doctor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"grantedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"grantedAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AccessRevoked\",\"inputs\":[{\"name\":\"doctor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"revokedAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IntegrityHashUpdated\",\"inputs\":[{\"name\":\"newHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PolicyUpdated\",\"inputs\":[{\"name\":\"requiredDeptCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"requiredInstCommitment\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PolicyMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaleProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// RecordABI is the input ABI used to generate the binding from.
// Deprecated: Use RecordMetaData.ABI instead.
var RecordABI = RecordMetaData.ABI

// Record is an auto generated Go binding around an Ethereum contract.
type Record struct {
	RecordCaller     // Read-only binding to the contract
	RecordTransactor // Write-only binding to the contract
	RecordFilterer   // Log filterer for contract events
}

// RecordCaller is an auto generated read-only Go binding around an Ethereum contract.
type RecordCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RecordTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RecordTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RecordFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RecordFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RecordSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RecordSession struct {
	Contract     *Record           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RecordCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RecordCallerSession struct {
	Contract *RecordCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RecordTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RecordTransactorSession struct {
	Contract     *RecordTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RecordRaw is an auto generated low-level Go binding around an Ethereum contract.
type RecordRaw struct {
	Contract *Record // Generic contract binding to access the raw methods on
}

// RecordCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RecordCallerRaw struct {
	Contract *RecordCaller // Generic read-only contract binding to access the raw methods on
}

// RecordTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RecordTransactorRaw struct {
	Contract *RecordTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRecord creates a new instance of Record, bound to a specific deployed contract.
func NewRecord(address common.Address, backend bind.ContractBackend) (*Record, error) {
	contract, err := bindRecord(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Record{RecordCaller: RecordCaller{contract: contract}, RecordTransactor: RecordTransactor{contract: contract}, RecordFilterer: RecordFilterer{contract: contract}}, nil
}

// NewRecordCaller creates a new read-only instance of Record, bound to a specific deployed contract.
func NewRecordCaller(address common.Address, caller bind.ContractCaller) (*RecordCaller, error) {
	contract, err := bindRecord(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RecordCaller{contract: contract}, nil
}

// NewRecordTransactor creates a new write-only instance of Record, bound to a specific deployed contract.
func NewRecordTransactor(address common.Address, transactor bind.ContractTransactor) (*RecordTransactor, error) {
	contract, err := bindRecord(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RecordTransactor{contract: contract}, nil
}

// NewRecordFilterer creates a new log filterer instance of Record, bound to a specific deployed contract.
func NewRecordFilterer(address common.Address, filterer bind.ContractFilterer) (*RecordFilterer, error) {
	contract, err := bindRecord(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RecordFilterer{contract: contract}, nil
}

// bindRecord binds a generic wrapper to an already deployed contract.
func bindRecord(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RecordMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Record *RecordRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Record.Contract.RecordCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Record *RecordRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Record.Contract.RecordTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Record *RecordRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Record.Contract.RecordTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Record *RecordCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Record.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Record *RecordTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Record.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Record *RecordTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Record.Contract.contract.Transact(opts, method, params...)
}

// POLICYPROOFTTL is a free data retrieval call binding the contract method 0xb420ddf6.
//
// Solidity: function POLICY_PROOF_TTL() view returns(uint256)
func (_Record *RecordCaller) POLICYPROOFTTL(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "POLICY_PROOF_TTL")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// POLICYPROOFTTL is a free data retrieval call binding the contract method 0xb420ddf6.
//
// Solidity: function POLICY_PROOF_TTL() view returns(uint256)
func (_Record *RecordSession) POLICYPROOFTTL() (*big.Int, error) {
	return _Record.Contract.POLICYPROOFTTL(&_Record.CallOpts)
}

// POLICYPROOFTTL is a free data retrieval call binding the contract method 0xb420ddf6.
//
// Solidity: function POLICY_PROOF_TTL() view returns(uint256)
func (_Record *RecordCallerSession) POLICYPROOFTTL() (*big.Int, error) {
	return _Record.Contract.POLICYPROOFTTL(&_Record.CallOpts)
}

// ContentHash is a free data retrieval call binding the contract method 0x646c2e33.
//
// Solidity: function contentHash() view returns(bytes32)
func (_Record *RecordCaller) ContentHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "contentHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ContentHash is a free data retrieval call binding the contract method 0x646c2e33.
//
// Solidity: function contentHash() view returns(bytes32)
func (_Record *RecordSession) ContentHash() ([32]byte, error) {
	return _Record.Contract.ContentHash(&_Record.CallOpts)
}

// ContentHash is a free data retrieval call binding the contract method 0x646c2e33.
//
// Solidity: function contentHash() view returns(bytes32)
func (_Record *RecordCallerSession) ContentHash() ([32]byte, error) {
	return _Record.Contract.ContentHash(&_Record.CallOpts)
}

// GetPolicy is a free data retrieval call binding the contract method 0xce1e4626.
//
// Solidity: function getPolicy() view returns(bytes32 deptCommitment, bytes32 instCommitment)
func (_Record *RecordCaller) GetPolicy(opts *bind.CallOpts) (struct {
	DeptCommitment [32]byte
	InstCommitment [32]byte
}, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "getPolicy")

	outstruct := new(struct {
		DeptCommitment [32]byte
		InstCommitment [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DeptCommitment = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.InstCommitment = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// GetPolicy is a free data retrieval call binding the contract method 0xce1e4626.
//
// Solidity: function getPolicy() view returns(bytes32 deptCommitment, bytes32 instCommitment)
func (_Record *RecordSession) GetPolicy() (struct {
	DeptCommitment [32]byte
	InstCommitment [32]byte
}, error) {
	return _Record.Contract.GetPolicy(&_Record.CallOpts)
}

// GetPolicy is a free data retrieval call binding the contract method 0xce1e4626.
//
// Solidity: function getPolicy() view returns(bytes32 deptCommitment, bytes32 instCommitment)
func (_Record *RecordCallerSession) GetPolicy() (struct {
	DeptCommitment [32]byte
	InstCommitment [32]byte
}, error) {
	return _Record.Contract.GetPolicy(&_Record.CallOpts)
}

// GetRecord is a free data retrieval call binding the contract method 0xf18ad6fb.
//
// Solidity: function getRecord() view returns((uint8,uint256,uint256), bytes32 _contentHash, string _storageReference, address _owner, address _uploader)
func (_Record *RecordCaller) GetRecord(opts *bind.CallOpts) (TypesRecordMetadata, [32]byte, string, common.Address, common.Address, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "getRecord")

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

// GetRecord is a free data retrieval call binding the contract method 0xf18ad6fb.
//
// Solidity: function getRecord() view returns((uint8,uint256,uint256), bytes32 _contentHash, string _storageReference, address _owner, address _uploader)
func (_Record *RecordSession) GetRecord() (TypesRecordMetadata, [32]byte, string, common.Address, common.Address, error) {
	return _Record.Contract.GetRecord(&_Record.CallOpts)
}

// GetRecord is a free data retrieval call binding the contract method 0xf18ad6fb.
//
// Solidity: function getRecord() view returns((uint8,uint256,uint256), bytes32 _contentHash, string _storageReference, address _owner, address _uploader)
func (_Record *RecordCallerSession) GetRecord() (TypesRecordMetadata, [32]byte, string, common.Address, common.Address, error) {
	return _Record.Contract.GetRecord(&_Record.CallOpts)
}

// History is a free data retrieval call binding the contract method 0x98daac83.
//
// Solidity: function history() view returns(address)
func (_Record *RecordCaller) History(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "history")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// History is a free data retrieval call binding the contract method 0x98daac83.
//
// Solidity: function history() view returns(address)
func (_Record *RecordSession) History() (common.Address, error) {
	return _Record.Contract.History(&_Record.CallOpts)
}

// History is a free data retrieval call binding the contract method 0x98daac83.
//
// Solidity: function history() view returns(address)
func (_Record *RecordCallerSession) History() (common.Address, error) {
	return _Record.Contract.History(&_Record.CallOpts)
}

// IsAuthorized is a free data retrieval call binding the contract method 0xfe9fbb80.
//
// Solidity: function isAuthorized(address doctor) view returns(bool)
func (_Record *RecordCaller) IsAuthorized(opts *bind.CallOpts, doctor common.Address) (bool, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "isAuthorized", doctor)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAuthorized is a free data retrieval call binding the contract method 0xfe9fbb80.
//
// Solidity: function isAuthorized(address doctor) view returns(bool)
func (_Record *RecordSession) IsAuthorized(doctor common.Address) (bool, error) {
	return _Record.Contract.IsAuthorized(&_Record.CallOpts, doctor)
}

// IsAuthorized is a free data retrieval call binding the contract method 0xfe9fbb80.
//
// Solidity: function isAuthorized(address doctor) view returns(bool)
func (_Record *RecordCallerSession) IsAuthorized(doctor common.Address) (bool, error) {
	return _Record.Contract.IsAuthorized(&_Record.CallOpts, doctor)
}

// Metadata is a free data retrieval call binding the contract method 0x392f37e9.
//
// Solidity: function metadata() view returns(uint8 recordType, uint256 createdAt, uint256 updatedAt)
func (_Record *RecordCaller) Metadata(opts *bind.CallOpts) (struct {
	RecordType uint8
	CreatedAt  *big.Int
	UpdatedAt  *big.Int
}, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "metadata")

	outstruct := new(struct {
		RecordType uint8
		CreatedAt  *big.Int
		UpdatedAt  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RecordType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.CreatedAt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.UpdatedAt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Metadata is a free data retrieval call binding the contract method 0x392f37e9.
//
// Solidity: function metadata() view returns(uint8 recordType, uint256 createdAt, uint256 updatedAt)
func (_Record *RecordSession) Metadata() (struct {
	RecordType uint8
	CreatedAt  *big.Int
	UpdatedAt  *big.Int
}, error) {
	return _Record.Contract.Metadata(&_Record.CallOpts)
}

// Metadata is a free data retrieval call binding the contract method 0x392f37e9.
//
// Solidity: function metadata() view returns(uint8 recordType, uint256 createdAt, uint256 updatedAt)
func (_Record *RecordCallerSession) Metadata() (struct {
	RecordType uint8
	CreatedAt  *big.Int
	UpdatedAt  *big.Int
}, error) {
	return _Record.Contract.Metadata(&_Record.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Record *RecordCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Record *RecordSession) Owner() (common.Address, error) {
	return _Record.Contract.Owner(&_Record.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Record *RecordCallerSession) Owner() (common.Address, error) {
	return _Record.Contract.Owner(&_Record.CallOpts)
}

// Permissions is a free data retrieval call binding the contract method 0xab8c71c0.
//
// Solidity: function permissions() view returns(address)
func (_Record *RecordCaller) Permissions(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "permissions")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Permissions is a free data retrieval call binding the contract method 0xab8c71c0.
//
// Solidity: function permissions() view returns(address)
func (_Record *RecordSession) Permissions() (common.Address, error) {
	return _Record.Contract.Permissions(&_Record.CallOpts)
}

// Permissions is a free data retrieval call binding the contract method 0xab8c71c0.
//
// Solidity: function permissions() view returns(address)
func (_Record *RecordCallerSession) Permissions() (common.Address, error) {
	return _Record.Contract.Permissions(&_Record.CallOpts)
}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_Record *RecordCaller) PolicyVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "policyVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_Record *RecordSession) PolicyVerifier() (common.Address, error) {
	return _Record.Contract.PolicyVerifier(&_Record.CallOpts)
}

// PolicyVerifier is a free data retrieval call binding the contract method 0xe4d16fbb.
//
// Solidity: function policyVerifier() view returns(address)
func (_Record *RecordCallerSession) PolicyVerifier() (common.Address, error) {
	return _Record.Contract.PolicyVerifier(&_Record.CallOpts)
}

// RequiredDeptCommitment is a free data retrieval call binding the contract method 0xfb5728a3.
//
// Solidity: function requiredDeptCommitment() view returns(bytes32)
func (_Record *RecordCaller) RequiredDeptCommitment(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "requiredDeptCommitment")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RequiredDeptCommitment is a free data retrieval call binding the contract method 0xfb5728a3.
//
// Solidity: function requiredDeptCommitment() view returns(bytes32)
func (_Record *RecordSession) RequiredDeptCommitment() ([32]byte, error) {
	return _Record.Contract.RequiredDeptCommitment(&_Record.CallOpts)
}

// RequiredDeptCommitment is a free data retrieval call binding the contract method 0xfb5728a3.
//
// Solidity: function requiredDeptCommitment() view returns(bytes32)
func (_Record *RecordCallerSession) RequiredDeptCommitment() ([32]byte, error) {
	return _Record.Contract.RequiredDeptCommitment(&_Record.CallOpts)
}

// RequiredInstCommitment is a free data retrieval call binding the contract method 0xc2cb8837.
//
// Solidity: function requiredInstCommitment() view returns(bytes32)
func (_Record *RecordCaller) RequiredInstCommitment(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "requiredInstCommitment")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RequiredInstCommitment is a free data retrieval call binding the contract method 0xc2cb8837.
//
// Solidity: function requiredInstCommitment() view returns(bytes32)
func (_Record *RecordSession) RequiredInstCommitment() ([32]byte, error) {
	return _Record.Contract.RequiredInstCommitment(&_Record.CallOpts)
}

// RequiredInstCommitment is a free data retrieval call binding the contract method 0xc2cb8837.
//
// Solidity: function requiredInstCommitment() view returns(bytes32)
func (_Record *RecordCallerSession) RequiredInstCommitment() ([32]byte, error) {
	return _Record.Contract.RequiredInstCommitment(&_Record.CallOpts)
}

// StorageReference is a free data retrieval call binding the contract method 0xcdb18091.
//
// Solidity: function storageReference() view returns(string)
func (_Record *RecordCaller) StorageReference(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "storageReference")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// StorageReference is a free data retrieval call binding the contract method 0xcdb18091.
//
// Solidity: function storageReference() view returns(string)
func (_Record *RecordSession) StorageReference() (string, error) {
	return _Record.Contract.StorageReference(&_Record.CallOpts)
}

// StorageReference is a free data retrieval call binding the contract method 0xcdb18091.
//
// Solidity: function storageReference() view returns(string)
func (_Record *RecordCallerSession) StorageReference() (string, error) {
	return _Record.Contract.StorageReference(&_Record.CallOpts)
}

// Uploader is a free data retrieval call binding the contract method 0x65fe2a0b.
//
// Solidity: function uploader() view returns(address)
func (_Record *RecordCaller) Uploader(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "uploader")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Uploader is a free data retrieval call binding the contract method 0x65fe2a0b.
//
// Solidity: function uploader() view returns(address)
func (_Record *RecordSession) Uploader() (common.Address, error) {
	return _Record.Contract.Uploader(&_Record.CallOpts)
}

// Uploader is a free data retrieval call binding the contract method 0x65fe2a0b.
//
// Solidity: function uploader() view returns(address)
func (_Record *RecordCallerSession) Uploader() (common.Address, error) {
	return _Record.Contract.Uploader(&_Record.CallOpts)
}

// VerifyIntegrity is a free data retrieval call binding the contract method 0x14f38116.
//
// Solidity: function verifyIntegrity(bytes32 hash) view returns(bool)
func (_Record *RecordCaller) VerifyIntegrity(opts *bind.CallOpts, hash [32]byte) (bool, error) {
	var out []interface{}
	err := _Record.contract.Call(opts, &out, "verifyIntegrity", hash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyIntegrity is a free data retrieval call binding the contract method 0x14f38116.
//
// Solidity: function verifyIntegrity(bytes32 hash) view returns(bool)
func (_Record *RecordSession) VerifyIntegrity(hash [32]byte) (bool, error) {
	return _Record.Contract.VerifyIntegrity(&_Record.CallOpts, hash)
}

// VerifyIntegrity is a free data retrieval call binding the contract method 0x14f38116.
//
// Solidity: function verifyIntegrity(bytes32 hash) view returns(bool)
func (_Record *RecordCallerSession) VerifyIntegrity(hash [32]byte) (bool, error) {
	return _Record.Contract.VerifyIntegrity(&_Record.CallOpts, hash)
}

// GrantAccess is a paid mutator transaction binding the contract method 0x0ae5e739.
//
// Solidity: function grantAccess(address doctor) returns()
func (_Record *RecordTransactor) GrantAccess(opts *bind.TransactOpts, doctor common.Address) (*types.Transaction, error) {
	return _Record.contract.Transact(opts, "grantAccess", doctor)
}

// GrantAccess is a paid mutator transaction binding the contract method 0x0ae5e739.
//
// Solidity: function grantAccess(address doctor) returns()
func (_Record *RecordSession) GrantAccess(doctor common.Address) (*types.Transaction, error) {
	return _Record.Contract.GrantAccess(&_Record.TransactOpts, doctor)
}

// GrantAccess is a paid mutator transaction binding the contract method 0x0ae5e739.
//
// Solidity: function grantAccess(address doctor) returns()
func (_Record *RecordTransactorSession) GrantAccess(doctor common.Address) (*types.Transaction, error) {
	return _Record.Contract.GrantAccess(&_Record.TransactOpts, doctor)
}

// GrantAccessWithProof is a paid mutator transaction binding the contract method 0x6286f3b7.
//
// Solidity: function grantAccessWithProof(uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 deptCommitment, bytes32 instCommitment, bytes32 doctorCommitment, uint48 evaluatedAt) returns()
func (_Record *RecordTransactor) GrantAccessWithProof(opts *bind.TransactOpts, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, deptCommitment [32]byte, instCommitment [32]byte, doctorCommitment [32]byte, evaluatedAt *big.Int) (*types.Transaction, error) {
	return _Record.contract.Transact(opts, "grantAccessWithProof", pA, pB, pC, deptCommitment, instCommitment, doctorCommitment, evaluatedAt)
}

// GrantAccessWithProof is a paid mutator transaction binding the contract method 0x6286f3b7.
//
// Solidity: function grantAccessWithProof(uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 deptCommitment, bytes32 instCommitment, bytes32 doctorCommitment, uint48 evaluatedAt) returns()
func (_Record *RecordSession) GrantAccessWithProof(pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, deptCommitment [32]byte, instCommitment [32]byte, doctorCommitment [32]byte, evaluatedAt *big.Int) (*types.Transaction, error) {
	return _Record.Contract.GrantAccessWithProof(&_Record.TransactOpts, pA, pB, pC, deptCommitment, instCommitment, doctorCommitment, evaluatedAt)
}

// GrantAccessWithProof is a paid mutator transaction binding the contract method 0x6286f3b7.
//
// Solidity: function grantAccessWithProof(uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 deptCommitment, bytes32 instCommitment, bytes32 doctorCommitment, uint48 evaluatedAt) returns()
func (_Record *RecordTransactorSession) GrantAccessWithProof(pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, deptCommitment [32]byte, instCommitment [32]byte, doctorCommitment [32]byte, evaluatedAt *big.Int) (*types.Transaction, error) {
	return _Record.Contract.GrantAccessWithProof(&_Record.TransactOpts, pA, pB, pC, deptCommitment, instCommitment, doctorCommitment, evaluatedAt)
}

// RevokeAccess is a paid mutator transaction binding the contract method 0x85e68531.
//
// Solidity: function revokeAccess(address doctor) returns()
func (_Record *RecordTransactor) RevokeAccess(opts *bind.TransactOpts, doctor common.Address) (*types.Transaction, error) {
	return _Record.contract.Transact(opts, "revokeAccess", doctor)
}

// RevokeAccess is a paid mutator transaction binding the contract method 0x85e68531.
//
// Solidity: function revokeAccess(address doctor) returns()
func (_Record *RecordSession) RevokeAccess(doctor common.Address) (*types.Transaction, error) {
	return _Record.Contract.RevokeAccess(&_Record.TransactOpts, doctor)
}

// RevokeAccess is a paid mutator transaction binding the contract method 0x85e68531.
//
// Solidity: function revokeAccess(address doctor) returns()
func (_Record *RecordTransactorSession) RevokeAccess(doctor common.Address) (*types.Transaction, error) {
	return _Record.Contract.RevokeAccess(&_Record.TransactOpts, doctor)
}

// SetPolicy is a paid mutator transaction binding the contract method 0x8135da93.
//
// Solidity: function setPolicy(bytes32 _requiredDeptCommitment, bytes32 _requiredInstCommitment) returns()
func (_Record *RecordTransactor) SetPolicy(opts *bind.TransactOpts, _requiredDeptCommitment [32]byte, _requiredInstCommitment [32]byte) (*types.Transaction, error) {
	return _Record.contract.Transact(opts, "setPolicy", _requiredDeptCommitment, _requiredInstCommitment)
}

// SetPolicy is a paid mutator transaction binding the contract method 0x8135da93.
//
// Solidity: function setPolicy(bytes32 _requiredDeptCommitment, bytes32 _requiredInstCommitment) returns()
func (_Record *RecordSession) SetPolicy(_requiredDeptCommitment [32]byte, _requiredInstCommitment [32]byte) (*types.Transaction, error) {
	return _Record.Contract.SetPolicy(&_Record.TransactOpts, _requiredDeptCommitment, _requiredInstCommitment)
}

// SetPolicy is a paid mutator transaction binding the contract method 0x8135da93.
//
// Solidity: function setPolicy(bytes32 _requiredDeptCommitment, bytes32 _requiredInstCommitment) returns()
func (_Record *RecordTransactorSession) SetPolicy(_requiredDeptCommitment [32]byte, _requiredInstCommitment [32]byte) (*types.Transaction, error) {
	return _Record.Contract.SetPolicy(&_Record.TransactOpts, _requiredDeptCommitment, _requiredInstCommitment)
}

// UpdateIntegrityHash is a paid mutator transaction binding the contract method 0x71cde89f.
//
// Solidity: function updateIntegrityHash(bytes32 newHash) returns()
func (_Record *RecordTransactor) UpdateIntegrityHash(opts *bind.TransactOpts, newHash [32]byte) (*types.Transaction, error) {
	return _Record.contract.Transact(opts, "updateIntegrityHash", newHash)
}

// UpdateIntegrityHash is a paid mutator transaction binding the contract method 0x71cde89f.
//
// Solidity: function updateIntegrityHash(bytes32 newHash) returns()
func (_Record *RecordSession) UpdateIntegrityHash(newHash [32]byte) (*types.Transaction, error) {
	return _Record.Contract.UpdateIntegrityHash(&_Record.TransactOpts, newHash)
}

// UpdateIntegrityHash is a paid mutator transaction binding the contract method 0x71cde89f.
//
// Solidity: function updateIntegrityHash(bytes32 newHash) returns()
func (_Record *RecordTransactorSession) UpdateIntegrityHash(newHash [32]byte) (*types.Transaction, error) {
	return _Record.Contract.UpdateIntegrityHash(&_Record.TransactOpts, newHash)
}

// RecordAccessGrantedIterator is returned from FilterAccessGranted and is used to iterate over the raw logs and unpacked data for AccessGranted events raised by the Record contract.
type RecordAccessGrantedIterator struct {
	Event *RecordAccessGranted // Event containing the contract specifics and raw log

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
func (it *RecordAccessGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RecordAccessGranted)
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
		it.Event = new(RecordAccessGranted)
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
func (it *RecordAccessGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RecordAccessGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RecordAccessGranted represents a AccessGranted event raised by the Record contract.
type RecordAccessGranted struct {
	Doctor    common.Address
	GrantedBy common.Address
	GrantedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAccessGranted is a free log retrieval operation binding the contract event 0x9a8bff21fc3398b6b251a2896c3af7f54b60c58f385e158d79ac75058471bfcb.
//
// Solidity: event AccessGranted(address indexed doctor, address indexed grantedBy, uint48 grantedAt)
func (_Record *RecordFilterer) FilterAccessGranted(opts *bind.FilterOpts, doctor []common.Address, grantedBy []common.Address) (*RecordAccessGrantedIterator, error) {

	var doctorRule []interface{}
	for _, doctorItem := range doctor {
		doctorRule = append(doctorRule, doctorItem)
	}
	var grantedByRule []interface{}
	for _, grantedByItem := range grantedBy {
		grantedByRule = append(grantedByRule, grantedByItem)
	}

	logs, sub, err := _Record.contract.FilterLogs(opts, "AccessGranted", doctorRule, grantedByRule)
	if err != nil {
		return nil, err
	}
	return &RecordAccessGrantedIterator{contract: _Record.contract, event: "AccessGranted", logs: logs, sub: sub}, nil
}

// WatchAccessGranted is a free log subscription operation binding the contract event 0x9a8bff21fc3398b6b251a2896c3af7f54b60c58f385e158d79ac75058471bfcb.
//
// Solidity: event AccessGranted(address indexed doctor, address indexed grantedBy, uint48 grantedAt)
func (_Record *RecordFilterer) WatchAccessGranted(opts *bind.WatchOpts, sink chan<- *RecordAccessGranted, doctor []common.Address, grantedBy []common.Address) (event.Subscription, error) {

	var doctorRule []interface{}
	for _, doctorItem := range doctor {
		doctorRule = append(doctorRule, doctorItem)
	}
	var grantedByRule []interface{}
	for _, grantedByItem := range grantedBy {
		grantedByRule = append(grantedByRule, grantedByItem)
	}

	logs, sub, err := _Record.contract.WatchLogs(opts, "AccessGranted", doctorRule, grantedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RecordAccessGranted)
				if err := _Record.contract.UnpackLog(event, "AccessGranted", log); err != nil {
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

// ParseAccessGranted is a log parse operation binding the contract event 0x9a8bff21fc3398b6b251a2896c3af7f54b60c58f385e158d79ac75058471bfcb.
//
// Solidity: event AccessGranted(address indexed doctor, address indexed grantedBy, uint48 grantedAt)
func (_Record *RecordFilterer) ParseAccessGranted(log types.Log) (*RecordAccessGranted, error) {
	event := new(RecordAccessGranted)
	if err := _Record.contract.UnpackLog(event, "AccessGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RecordAccessRevokedIterator is returned from FilterAccessRevoked and is used to iterate over the raw logs and unpacked data for AccessRevoked events raised by the Record contract.
type RecordAccessRevokedIterator struct {
	Event *RecordAccessRevoked // Event containing the contract specifics and raw log

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
func (it *RecordAccessRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RecordAccessRevoked)
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
		it.Event = new(RecordAccessRevoked)
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
func (it *RecordAccessRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RecordAccessRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RecordAccessRevoked represents a AccessRevoked event raised by the Record contract.
type RecordAccessRevoked struct {
	Doctor    common.Address
	RevokedAt *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAccessRevoked is a free log retrieval operation binding the contract event 0xe2a1495e1ebb017d1b24edb1cc9979e3f6e8ef90a9b9742c93e50ddfb2827992.
//
// Solidity: event AccessRevoked(address indexed doctor, uint48 revokedAt)
func (_Record *RecordFilterer) FilterAccessRevoked(opts *bind.FilterOpts, doctor []common.Address) (*RecordAccessRevokedIterator, error) {

	var doctorRule []interface{}
	for _, doctorItem := range doctor {
		doctorRule = append(doctorRule, doctorItem)
	}

	logs, sub, err := _Record.contract.FilterLogs(opts, "AccessRevoked", doctorRule)
	if err != nil {
		return nil, err
	}
	return &RecordAccessRevokedIterator{contract: _Record.contract, event: "AccessRevoked", logs: logs, sub: sub}, nil
}

// WatchAccessRevoked is a free log subscription operation binding the contract event 0xe2a1495e1ebb017d1b24edb1cc9979e3f6e8ef90a9b9742c93e50ddfb2827992.
//
// Solidity: event AccessRevoked(address indexed doctor, uint48 revokedAt)
func (_Record *RecordFilterer) WatchAccessRevoked(opts *bind.WatchOpts, sink chan<- *RecordAccessRevoked, doctor []common.Address) (event.Subscription, error) {

	var doctorRule []interface{}
	for _, doctorItem := range doctor {
		doctorRule = append(doctorRule, doctorItem)
	}

	logs, sub, err := _Record.contract.WatchLogs(opts, "AccessRevoked", doctorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RecordAccessRevoked)
				if err := _Record.contract.UnpackLog(event, "AccessRevoked", log); err != nil {
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

// ParseAccessRevoked is a log parse operation binding the contract event 0xe2a1495e1ebb017d1b24edb1cc9979e3f6e8ef90a9b9742c93e50ddfb2827992.
//
// Solidity: event AccessRevoked(address indexed doctor, uint48 revokedAt)
func (_Record *RecordFilterer) ParseAccessRevoked(log types.Log) (*RecordAccessRevoked, error) {
	event := new(RecordAccessRevoked)
	if err := _Record.contract.UnpackLog(event, "AccessRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RecordIntegrityHashUpdatedIterator is returned from FilterIntegrityHashUpdated and is used to iterate over the raw logs and unpacked data for IntegrityHashUpdated events raised by the Record contract.
type RecordIntegrityHashUpdatedIterator struct {
	Event *RecordIntegrityHashUpdated // Event containing the contract specifics and raw log

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
func (it *RecordIntegrityHashUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RecordIntegrityHashUpdated)
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
		it.Event = new(RecordIntegrityHashUpdated)
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
func (it *RecordIntegrityHashUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RecordIntegrityHashUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RecordIntegrityHashUpdated represents a IntegrityHashUpdated event raised by the Record contract.
type RecordIntegrityHashUpdated struct {
	NewHash [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterIntegrityHashUpdated is a free log retrieval operation binding the contract event 0xfb5a49bdba0ec3da0cb6850783a8250c8addeb3bdbd1bc555c688bb8e39175e7.
//
// Solidity: event IntegrityHashUpdated(bytes32 newHash)
func (_Record *RecordFilterer) FilterIntegrityHashUpdated(opts *bind.FilterOpts) (*RecordIntegrityHashUpdatedIterator, error) {

	logs, sub, err := _Record.contract.FilterLogs(opts, "IntegrityHashUpdated")
	if err != nil {
		return nil, err
	}
	return &RecordIntegrityHashUpdatedIterator{contract: _Record.contract, event: "IntegrityHashUpdated", logs: logs, sub: sub}, nil
}

// WatchIntegrityHashUpdated is a free log subscription operation binding the contract event 0xfb5a49bdba0ec3da0cb6850783a8250c8addeb3bdbd1bc555c688bb8e39175e7.
//
// Solidity: event IntegrityHashUpdated(bytes32 newHash)
func (_Record *RecordFilterer) WatchIntegrityHashUpdated(opts *bind.WatchOpts, sink chan<- *RecordIntegrityHashUpdated) (event.Subscription, error) {

	logs, sub, err := _Record.contract.WatchLogs(opts, "IntegrityHashUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RecordIntegrityHashUpdated)
				if err := _Record.contract.UnpackLog(event, "IntegrityHashUpdated", log); err != nil {
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

// ParseIntegrityHashUpdated is a log parse operation binding the contract event 0xfb5a49bdba0ec3da0cb6850783a8250c8addeb3bdbd1bc555c688bb8e39175e7.
//
// Solidity: event IntegrityHashUpdated(bytes32 newHash)
func (_Record *RecordFilterer) ParseIntegrityHashUpdated(log types.Log) (*RecordIntegrityHashUpdated, error) {
	event := new(RecordIntegrityHashUpdated)
	if err := _Record.contract.UnpackLog(event, "IntegrityHashUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RecordPolicyUpdatedIterator is returned from FilterPolicyUpdated and is used to iterate over the raw logs and unpacked data for PolicyUpdated events raised by the Record contract.
type RecordPolicyUpdatedIterator struct {
	Event *RecordPolicyUpdated // Event containing the contract specifics and raw log

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
func (it *RecordPolicyUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RecordPolicyUpdated)
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
		it.Event = new(RecordPolicyUpdated)
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
func (it *RecordPolicyUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RecordPolicyUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RecordPolicyUpdated represents a PolicyUpdated event raised by the Record contract.
type RecordPolicyUpdated struct {
	RequiredDeptCommitment [32]byte
	RequiredInstCommitment [32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterPolicyUpdated is a free log retrieval operation binding the contract event 0x66bcfecc06b590f2d3d4fea9dc6173f115cfa08e6a5ee6f06214e92e52efea2b.
//
// Solidity: event PolicyUpdated(bytes32 requiredDeptCommitment, bytes32 requiredInstCommitment)
func (_Record *RecordFilterer) FilterPolicyUpdated(opts *bind.FilterOpts) (*RecordPolicyUpdatedIterator, error) {

	logs, sub, err := _Record.contract.FilterLogs(opts, "PolicyUpdated")
	if err != nil {
		return nil, err
	}
	return &RecordPolicyUpdatedIterator{contract: _Record.contract, event: "PolicyUpdated", logs: logs, sub: sub}, nil
}

// WatchPolicyUpdated is a free log subscription operation binding the contract event 0x66bcfecc06b590f2d3d4fea9dc6173f115cfa08e6a5ee6f06214e92e52efea2b.
//
// Solidity: event PolicyUpdated(bytes32 requiredDeptCommitment, bytes32 requiredInstCommitment)
func (_Record *RecordFilterer) WatchPolicyUpdated(opts *bind.WatchOpts, sink chan<- *RecordPolicyUpdated) (event.Subscription, error) {

	logs, sub, err := _Record.contract.WatchLogs(opts, "PolicyUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RecordPolicyUpdated)
				if err := _Record.contract.UnpackLog(event, "PolicyUpdated", log); err != nil {
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

// ParsePolicyUpdated is a log parse operation binding the contract event 0x66bcfecc06b590f2d3d4fea9dc6173f115cfa08e6a5ee6f06214e92e52efea2b.
//
// Solidity: event PolicyUpdated(bytes32 requiredDeptCommitment, bytes32 requiredInstCommitment)
func (_Record *RecordFilterer) ParsePolicyUpdated(log types.Log) (*RecordPolicyUpdated, error) {
	event := new(RecordPolicyUpdated)
	if err := _Record.contract.UnpackLog(event, "PolicyUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
