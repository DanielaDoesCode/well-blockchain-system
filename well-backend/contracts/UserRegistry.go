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

// UserRegistryUser is an auto generated low-level Go binding around an user-defined struct.
type UserRegistryUser struct {
	Commitment  [32]byte
	InternalID  [32]byte
	Role        uint8
	UserAddress common.Address
}

// UserRegistryMetaData contains all meta data concerning the UserRegistry contract.
var UserRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addressToCommitment\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentToUser\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"},{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentByInternalID\",\"inputs\":[{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRole\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUser\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structUserRegistry.User\",\"components\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"},{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserByAddress\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structUserRegistry.User\",\"components\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"},{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"internalIDToCommitment\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isDoctor\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isFirstResponder\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isHospital\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"internalType\":\"enumTypes.Roles\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revoke\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Registered\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"role\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumTypes.Roles\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Revoked\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"internalID\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InternalIDAlreadyUsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[]}]",
}

// UserRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use UserRegistryMetaData.ABI instead.
var UserRegistryABI = UserRegistryMetaData.ABI

// UserRegistry is an auto generated Go binding around an Ethereum contract.
type UserRegistry struct {
	UserRegistryCaller     // Read-only binding to the contract
	UserRegistryTransactor // Write-only binding to the contract
	UserRegistryFilterer   // Log filterer for contract events
}

// UserRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type UserRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UserRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UserRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UserRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UserRegistrySession struct {
	Contract     *UserRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UserRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UserRegistryCallerSession struct {
	Contract *UserRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// UserRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UserRegistryTransactorSession struct {
	Contract     *UserRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// UserRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type UserRegistryRaw struct {
	Contract *UserRegistry // Generic contract binding to access the raw methods on
}

// UserRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UserRegistryCallerRaw struct {
	Contract *UserRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// UserRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UserRegistryTransactorRaw struct {
	Contract *UserRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUserRegistry creates a new instance of UserRegistry, bound to a specific deployed contract.
func NewUserRegistry(address common.Address, backend bind.ContractBackend) (*UserRegistry, error) {
	contract, err := bindUserRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UserRegistry{UserRegistryCaller: UserRegistryCaller{contract: contract}, UserRegistryTransactor: UserRegistryTransactor{contract: contract}, UserRegistryFilterer: UserRegistryFilterer{contract: contract}}, nil
}

// NewUserRegistryCaller creates a new read-only instance of UserRegistry, bound to a specific deployed contract.
func NewUserRegistryCaller(address common.Address, caller bind.ContractCaller) (*UserRegistryCaller, error) {
	contract, err := bindUserRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UserRegistryCaller{contract: contract}, nil
}

// NewUserRegistryTransactor creates a new write-only instance of UserRegistry, bound to a specific deployed contract.
func NewUserRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*UserRegistryTransactor, error) {
	contract, err := bindUserRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UserRegistryTransactor{contract: contract}, nil
}

// NewUserRegistryFilterer creates a new log filterer instance of UserRegistry, bound to a specific deployed contract.
func NewUserRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*UserRegistryFilterer, error) {
	contract, err := bindUserRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UserRegistryFilterer{contract: contract}, nil
}

// bindUserRegistry binds a generic wrapper to an already deployed contract.
func bindUserRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UserRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UserRegistry *UserRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UserRegistry.Contract.UserRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UserRegistry *UserRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UserRegistry.Contract.UserRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UserRegistry *UserRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UserRegistry.Contract.UserRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UserRegistry *UserRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UserRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UserRegistry *UserRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UserRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UserRegistry *UserRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UserRegistry.Contract.contract.Transact(opts, method, params...)
}

// AddressToCommitment is a free data retrieval call binding the contract method 0xed99a71a.
//
// Solidity: function addressToCommitment(address ) view returns(bytes32)
func (_UserRegistry *UserRegistryCaller) AddressToCommitment(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "addressToCommitment", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AddressToCommitment is a free data retrieval call binding the contract method 0xed99a71a.
//
// Solidity: function addressToCommitment(address ) view returns(bytes32)
func (_UserRegistry *UserRegistrySession) AddressToCommitment(arg0 common.Address) ([32]byte, error) {
	return _UserRegistry.Contract.AddressToCommitment(&_UserRegistry.CallOpts, arg0)
}

// AddressToCommitment is a free data retrieval call binding the contract method 0xed99a71a.
//
// Solidity: function addressToCommitment(address ) view returns(bytes32)
func (_UserRegistry *UserRegistryCallerSession) AddressToCommitment(arg0 common.Address) ([32]byte, error) {
	return _UserRegistry.Contract.AddressToCommitment(&_UserRegistry.CallOpts, arg0)
}

// CommitmentToUser is a free data retrieval call binding the contract method 0x5ed9bc29.
//
// Solidity: function commitmentToUser(bytes32 ) view returns(bytes32 commitment, bytes32 internalID, uint8 role, address userAddress)
func (_UserRegistry *UserRegistryCaller) CommitmentToUser(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Commitment  [32]byte
	InternalID  [32]byte
	Role        uint8
	UserAddress common.Address
}, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "commitmentToUser", arg0)

	outstruct := new(struct {
		Commitment  [32]byte
		InternalID  [32]byte
		Role        uint8
		UserAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Commitment = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.InternalID = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Role = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.UserAddress = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// CommitmentToUser is a free data retrieval call binding the contract method 0x5ed9bc29.
//
// Solidity: function commitmentToUser(bytes32 ) view returns(bytes32 commitment, bytes32 internalID, uint8 role, address userAddress)
func (_UserRegistry *UserRegistrySession) CommitmentToUser(arg0 [32]byte) (struct {
	Commitment  [32]byte
	InternalID  [32]byte
	Role        uint8
	UserAddress common.Address
}, error) {
	return _UserRegistry.Contract.CommitmentToUser(&_UserRegistry.CallOpts, arg0)
}

// CommitmentToUser is a free data retrieval call binding the contract method 0x5ed9bc29.
//
// Solidity: function commitmentToUser(bytes32 ) view returns(bytes32 commitment, bytes32 internalID, uint8 role, address userAddress)
func (_UserRegistry *UserRegistryCallerSession) CommitmentToUser(arg0 [32]byte) (struct {
	Commitment  [32]byte
	InternalID  [32]byte
	Role        uint8
	UserAddress common.Address
}, error) {
	return _UserRegistry.Contract.CommitmentToUser(&_UserRegistry.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool)
func (_UserRegistry *UserRegistryCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool)
func (_UserRegistry *UserRegistrySession) Commitments(arg0 [32]byte) (bool, error) {
	return _UserRegistry.Contract.Commitments(&_UserRegistry.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool)
func (_UserRegistry *UserRegistryCallerSession) Commitments(arg0 [32]byte) (bool, error) {
	return _UserRegistry.Contract.Commitments(&_UserRegistry.CallOpts, arg0)
}

// GetCommitmentByAddress is a free data retrieval call binding the contract method 0xaf0e9c8d.
//
// Solidity: function getCommitmentByAddress(address userAddress) view returns(bytes32 commitment)
func (_UserRegistry *UserRegistryCaller) GetCommitmentByAddress(opts *bind.CallOpts, userAddress common.Address) ([32]byte, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "getCommitmentByAddress", userAddress)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentByAddress is a free data retrieval call binding the contract method 0xaf0e9c8d.
//
// Solidity: function getCommitmentByAddress(address userAddress) view returns(bytes32 commitment)
func (_UserRegistry *UserRegistrySession) GetCommitmentByAddress(userAddress common.Address) ([32]byte, error) {
	return _UserRegistry.Contract.GetCommitmentByAddress(&_UserRegistry.CallOpts, userAddress)
}

// GetCommitmentByAddress is a free data retrieval call binding the contract method 0xaf0e9c8d.
//
// Solidity: function getCommitmentByAddress(address userAddress) view returns(bytes32 commitment)
func (_UserRegistry *UserRegistryCallerSession) GetCommitmentByAddress(userAddress common.Address) ([32]byte, error) {
	return _UserRegistry.Contract.GetCommitmentByAddress(&_UserRegistry.CallOpts, userAddress)
}

// GetCommitmentByInternalID is a free data retrieval call binding the contract method 0x632849ed.
//
// Solidity: function getCommitmentByInternalID(bytes32 internalID) view returns(bytes32 commitment)
func (_UserRegistry *UserRegistryCaller) GetCommitmentByInternalID(opts *bind.CallOpts, internalID [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "getCommitmentByInternalID", internalID)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentByInternalID is a free data retrieval call binding the contract method 0x632849ed.
//
// Solidity: function getCommitmentByInternalID(bytes32 internalID) view returns(bytes32 commitment)
func (_UserRegistry *UserRegistrySession) GetCommitmentByInternalID(internalID [32]byte) ([32]byte, error) {
	return _UserRegistry.Contract.GetCommitmentByInternalID(&_UserRegistry.CallOpts, internalID)
}

// GetCommitmentByInternalID is a free data retrieval call binding the contract method 0x632849ed.
//
// Solidity: function getCommitmentByInternalID(bytes32 internalID) view returns(bytes32 commitment)
func (_UserRegistry *UserRegistryCallerSession) GetCommitmentByInternalID(internalID [32]byte) ([32]byte, error) {
	return _UserRegistry.Contract.GetCommitmentByInternalID(&_UserRegistry.CallOpts, internalID)
}

// GetRole is a free data retrieval call binding the contract method 0x6e8eb7bc.
//
// Solidity: function getRole(bytes32 commitment) view returns(uint8)
func (_UserRegistry *UserRegistryCaller) GetRole(opts *bind.CallOpts, commitment [32]byte) (uint8, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "getRole", commitment)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRole is a free data retrieval call binding the contract method 0x6e8eb7bc.
//
// Solidity: function getRole(bytes32 commitment) view returns(uint8)
func (_UserRegistry *UserRegistrySession) GetRole(commitment [32]byte) (uint8, error) {
	return _UserRegistry.Contract.GetRole(&_UserRegistry.CallOpts, commitment)
}

// GetRole is a free data retrieval call binding the contract method 0x6e8eb7bc.
//
// Solidity: function getRole(bytes32 commitment) view returns(uint8)
func (_UserRegistry *UserRegistryCallerSession) GetRole(commitment [32]byte) (uint8, error) {
	return _UserRegistry.Contract.GetRole(&_UserRegistry.CallOpts, commitment)
}

// GetRoleByAddress is a free data retrieval call binding the contract method 0x7bff77d4.
//
// Solidity: function getRoleByAddress(address userAddress) view returns(uint8)
func (_UserRegistry *UserRegistryCaller) GetRoleByAddress(opts *bind.CallOpts, userAddress common.Address) (uint8, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "getRoleByAddress", userAddress)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetRoleByAddress is a free data retrieval call binding the contract method 0x7bff77d4.
//
// Solidity: function getRoleByAddress(address userAddress) view returns(uint8)
func (_UserRegistry *UserRegistrySession) GetRoleByAddress(userAddress common.Address) (uint8, error) {
	return _UserRegistry.Contract.GetRoleByAddress(&_UserRegistry.CallOpts, userAddress)
}

// GetRoleByAddress is a free data retrieval call binding the contract method 0x7bff77d4.
//
// Solidity: function getRoleByAddress(address userAddress) view returns(uint8)
func (_UserRegistry *UserRegistryCallerSession) GetRoleByAddress(userAddress common.Address) (uint8, error) {
	return _UserRegistry.Contract.GetRoleByAddress(&_UserRegistry.CallOpts, userAddress)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 commitment) view returns((bytes32,bytes32,uint8,address))
func (_UserRegistry *UserRegistryCaller) GetUser(opts *bind.CallOpts, commitment [32]byte) (UserRegistryUser, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "getUser", commitment)

	if err != nil {
		return *new(UserRegistryUser), err
	}

	out0 := *abi.ConvertType(out[0], new(UserRegistryUser)).(*UserRegistryUser)

	return out0, err

}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 commitment) view returns((bytes32,bytes32,uint8,address))
func (_UserRegistry *UserRegistrySession) GetUser(commitment [32]byte) (UserRegistryUser, error) {
	return _UserRegistry.Contract.GetUser(&_UserRegistry.CallOpts, commitment)
}

// GetUser is a free data retrieval call binding the contract method 0x6517579c.
//
// Solidity: function getUser(bytes32 commitment) view returns((bytes32,bytes32,uint8,address))
func (_UserRegistry *UserRegistryCallerSession) GetUser(commitment [32]byte) (UserRegistryUser, error) {
	return _UserRegistry.Contract.GetUser(&_UserRegistry.CallOpts, commitment)
}

// GetUserByAddress is a free data retrieval call binding the contract method 0x69c212f6.
//
// Solidity: function getUserByAddress(address userAddress) view returns((bytes32,bytes32,uint8,address))
func (_UserRegistry *UserRegistryCaller) GetUserByAddress(opts *bind.CallOpts, userAddress common.Address) (UserRegistryUser, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "getUserByAddress", userAddress)

	if err != nil {
		return *new(UserRegistryUser), err
	}

	out0 := *abi.ConvertType(out[0], new(UserRegistryUser)).(*UserRegistryUser)

	return out0, err

}

// GetUserByAddress is a free data retrieval call binding the contract method 0x69c212f6.
//
// Solidity: function getUserByAddress(address userAddress) view returns((bytes32,bytes32,uint8,address))
func (_UserRegistry *UserRegistrySession) GetUserByAddress(userAddress common.Address) (UserRegistryUser, error) {
	return _UserRegistry.Contract.GetUserByAddress(&_UserRegistry.CallOpts, userAddress)
}

// GetUserByAddress is a free data retrieval call binding the contract method 0x69c212f6.
//
// Solidity: function getUserByAddress(address userAddress) view returns((bytes32,bytes32,uint8,address))
func (_UserRegistry *UserRegistryCallerSession) GetUserByAddress(userAddress common.Address) (UserRegistryUser, error) {
	return _UserRegistry.Contract.GetUserByAddress(&_UserRegistry.CallOpts, userAddress)
}

// InternalIDToCommitment is a free data retrieval call binding the contract method 0x359635a2.
//
// Solidity: function internalIDToCommitment(bytes32 ) view returns(bytes32)
func (_UserRegistry *UserRegistryCaller) InternalIDToCommitment(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "internalIDToCommitment", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// InternalIDToCommitment is a free data retrieval call binding the contract method 0x359635a2.
//
// Solidity: function internalIDToCommitment(bytes32 ) view returns(bytes32)
func (_UserRegistry *UserRegistrySession) InternalIDToCommitment(arg0 [32]byte) ([32]byte, error) {
	return _UserRegistry.Contract.InternalIDToCommitment(&_UserRegistry.CallOpts, arg0)
}

// InternalIDToCommitment is a free data retrieval call binding the contract method 0x359635a2.
//
// Solidity: function internalIDToCommitment(bytes32 ) view returns(bytes32)
func (_UserRegistry *UserRegistryCallerSession) InternalIDToCommitment(arg0 [32]byte) ([32]byte, error) {
	return _UserRegistry.Contract.InternalIDToCommitment(&_UserRegistry.CallOpts, arg0)
}

// IsDoctor is a free data retrieval call binding the contract method 0x996440c6.
//
// Solidity: function isDoctor(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistryCaller) IsDoctor(opts *bind.CallOpts, userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "isDoctor", userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDoctor is a free data retrieval call binding the contract method 0x996440c6.
//
// Solidity: function isDoctor(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistrySession) IsDoctor(userAddress common.Address) (bool, error) {
	return _UserRegistry.Contract.IsDoctor(&_UserRegistry.CallOpts, userAddress)
}

// IsDoctor is a free data retrieval call binding the contract method 0x996440c6.
//
// Solidity: function isDoctor(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistryCallerSession) IsDoctor(userAddress common.Address) (bool, error) {
	return _UserRegistry.Contract.IsDoctor(&_UserRegistry.CallOpts, userAddress)
}

// IsFirstResponder is a free data retrieval call binding the contract method 0x19e4959d.
//
// Solidity: function isFirstResponder(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistryCaller) IsFirstResponder(opts *bind.CallOpts, userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "isFirstResponder", userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFirstResponder is a free data retrieval call binding the contract method 0x19e4959d.
//
// Solidity: function isFirstResponder(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistrySession) IsFirstResponder(userAddress common.Address) (bool, error) {
	return _UserRegistry.Contract.IsFirstResponder(&_UserRegistry.CallOpts, userAddress)
}

// IsFirstResponder is a free data retrieval call binding the contract method 0x19e4959d.
//
// Solidity: function isFirstResponder(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistryCallerSession) IsFirstResponder(userAddress common.Address) (bool, error) {
	return _UserRegistry.Contract.IsFirstResponder(&_UserRegistry.CallOpts, userAddress)
}

// IsHospital is a free data retrieval call binding the contract method 0xbdf6c673.
//
// Solidity: function isHospital(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistryCaller) IsHospital(opts *bind.CallOpts, userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "isHospital", userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsHospital is a free data retrieval call binding the contract method 0xbdf6c673.
//
// Solidity: function isHospital(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistrySession) IsHospital(userAddress common.Address) (bool, error) {
	return _UserRegistry.Contract.IsHospital(&_UserRegistry.CallOpts, userAddress)
}

// IsHospital is a free data retrieval call binding the contract method 0xbdf6c673.
//
// Solidity: function isHospital(address userAddress) view returns(bool)
func (_UserRegistry *UserRegistryCallerSession) IsHospital(userAddress common.Address) (bool, error) {
	return _UserRegistry.Contract.IsHospital(&_UserRegistry.CallOpts, userAddress)
}

// IsRegistered is a free data retrieval call binding the contract method 0x27258b22.
//
// Solidity: function isRegistered(bytes32 commitment) view returns(bool)
func (_UserRegistry *UserRegistryCaller) IsRegistered(opts *bind.CallOpts, commitment [32]byte) (bool, error) {
	var out []interface{}
	err := _UserRegistry.contract.Call(opts, &out, "isRegistered", commitment)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0x27258b22.
//
// Solidity: function isRegistered(bytes32 commitment) view returns(bool)
func (_UserRegistry *UserRegistrySession) IsRegistered(commitment [32]byte) (bool, error) {
	return _UserRegistry.Contract.IsRegistered(&_UserRegistry.CallOpts, commitment)
}

// IsRegistered is a free data retrieval call binding the contract method 0x27258b22.
//
// Solidity: function isRegistered(bytes32 commitment) view returns(bool)
func (_UserRegistry *UserRegistryCallerSession) IsRegistered(commitment [32]byte) (bool, error) {
	return _UserRegistry.Contract.IsRegistered(&_UserRegistry.CallOpts, commitment)
}

// Register is a paid mutator transaction binding the contract method 0x7793ce54.
//
// Solidity: function register(address userAddress, bytes32 commitment, bytes32 internalID, uint8 role) returns()
func (_UserRegistry *UserRegistryTransactor) Register(opts *bind.TransactOpts, userAddress common.Address, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _UserRegistry.contract.Transact(opts, "register", userAddress, commitment, internalID, role)
}

// Register is a paid mutator transaction binding the contract method 0x7793ce54.
//
// Solidity: function register(address userAddress, bytes32 commitment, bytes32 internalID, uint8 role) returns()
func (_UserRegistry *UserRegistrySession) Register(userAddress common.Address, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _UserRegistry.Contract.Register(&_UserRegistry.TransactOpts, userAddress, commitment, internalID, role)
}

// Register is a paid mutator transaction binding the contract method 0x7793ce54.
//
// Solidity: function register(address userAddress, bytes32 commitment, bytes32 internalID, uint8 role) returns()
func (_UserRegistry *UserRegistryTransactorSession) Register(userAddress common.Address, commitment [32]byte, internalID [32]byte, role uint8) (*types.Transaction, error) {
	return _UserRegistry.Contract.Register(&_UserRegistry.TransactOpts, userAddress, commitment, internalID, role)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 commitment) returns()
func (_UserRegistry *UserRegistryTransactor) Revoke(opts *bind.TransactOpts, commitment [32]byte) (*types.Transaction, error) {
	return _UserRegistry.contract.Transact(opts, "revoke", commitment)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 commitment) returns()
func (_UserRegistry *UserRegistrySession) Revoke(commitment [32]byte) (*types.Transaction, error) {
	return _UserRegistry.Contract.Revoke(&_UserRegistry.TransactOpts, commitment)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(bytes32 commitment) returns()
func (_UserRegistry *UserRegistryTransactorSession) Revoke(commitment [32]byte) (*types.Transaction, error) {
	return _UserRegistry.Contract.Revoke(&_UserRegistry.TransactOpts, commitment)
}

// UserRegistryRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the UserRegistry contract.
type UserRegistryRegisteredIterator struct {
	Event *UserRegistryRegistered // Event containing the contract specifics and raw log

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
func (it *UserRegistryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UserRegistryRegistered)
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
		it.Event = new(UserRegistryRegistered)
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
func (it *UserRegistryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UserRegistryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UserRegistryRegistered represents a Registered event raised by the UserRegistry contract.
type UserRegistryRegistered struct {
	Commitment [32]byte
	InternalID [32]byte
	Role       uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0xaa91fe5b2b0cb66eda7841ceac62bc8ce941ed13ee6966ba7295608c8744e6b9.
//
// Solidity: event Registered(bytes32 indexed commitment, bytes32 indexed internalID, uint8 role)
func (_UserRegistry *UserRegistryFilterer) FilterRegistered(opts *bind.FilterOpts, commitment [][32]byte, internalID [][32]byte) (*UserRegistryRegisteredIterator, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var internalIDRule []interface{}
	for _, internalIDItem := range internalID {
		internalIDRule = append(internalIDRule, internalIDItem)
	}

	logs, sub, err := _UserRegistry.contract.FilterLogs(opts, "Registered", commitmentRule, internalIDRule)
	if err != nil {
		return nil, err
	}
	return &UserRegistryRegisteredIterator{contract: _UserRegistry.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0xaa91fe5b2b0cb66eda7841ceac62bc8ce941ed13ee6966ba7295608c8744e6b9.
//
// Solidity: event Registered(bytes32 indexed commitment, bytes32 indexed internalID, uint8 role)
func (_UserRegistry *UserRegistryFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *UserRegistryRegistered, commitment [][32]byte, internalID [][32]byte) (event.Subscription, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var internalIDRule []interface{}
	for _, internalIDItem := range internalID {
		internalIDRule = append(internalIDRule, internalIDItem)
	}

	logs, sub, err := _UserRegistry.contract.WatchLogs(opts, "Registered", commitmentRule, internalIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UserRegistryRegistered)
				if err := _UserRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0xaa91fe5b2b0cb66eda7841ceac62bc8ce941ed13ee6966ba7295608c8744e6b9.
//
// Solidity: event Registered(bytes32 indexed commitment, bytes32 indexed internalID, uint8 role)
func (_UserRegistry *UserRegistryFilterer) ParseRegistered(log types.Log) (*UserRegistryRegistered, error) {
	event := new(UserRegistryRegistered)
	if err := _UserRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UserRegistryRevokedIterator is returned from FilterRevoked and is used to iterate over the raw logs and unpacked data for Revoked events raised by the UserRegistry contract.
type UserRegistryRevokedIterator struct {
	Event *UserRegistryRevoked // Event containing the contract specifics and raw log

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
func (it *UserRegistryRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UserRegistryRevoked)
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
		it.Event = new(UserRegistryRevoked)
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
func (it *UserRegistryRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UserRegistryRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UserRegistryRevoked represents a Revoked event raised by the UserRegistry contract.
type UserRegistryRevoked struct {
	Commitment [32]byte
	InternalID [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRevoked is a free log retrieval operation binding the contract event 0xa75610f599cea10d1a6bb6332878af5fffa58f57b227cba0d0229fa90a63cc51.
//
// Solidity: event Revoked(bytes32 indexed commitment, bytes32 indexed internalID)
func (_UserRegistry *UserRegistryFilterer) FilterRevoked(opts *bind.FilterOpts, commitment [][32]byte, internalID [][32]byte) (*UserRegistryRevokedIterator, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var internalIDRule []interface{}
	for _, internalIDItem := range internalID {
		internalIDRule = append(internalIDRule, internalIDItem)
	}

	logs, sub, err := _UserRegistry.contract.FilterLogs(opts, "Revoked", commitmentRule, internalIDRule)
	if err != nil {
		return nil, err
	}
	return &UserRegistryRevokedIterator{contract: _UserRegistry.contract, event: "Revoked", logs: logs, sub: sub}, nil
}

// WatchRevoked is a free log subscription operation binding the contract event 0xa75610f599cea10d1a6bb6332878af5fffa58f57b227cba0d0229fa90a63cc51.
//
// Solidity: event Revoked(bytes32 indexed commitment, bytes32 indexed internalID)
func (_UserRegistry *UserRegistryFilterer) WatchRevoked(opts *bind.WatchOpts, sink chan<- *UserRegistryRevoked, commitment [][32]byte, internalID [][32]byte) (event.Subscription, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var internalIDRule []interface{}
	for _, internalIDItem := range internalID {
		internalIDRule = append(internalIDRule, internalIDItem)
	}

	logs, sub, err := _UserRegistry.contract.WatchLogs(opts, "Revoked", commitmentRule, internalIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UserRegistryRevoked)
				if err := _UserRegistry.contract.UnpackLog(event, "Revoked", log); err != nil {
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

// ParseRevoked is a log parse operation binding the contract event 0xa75610f599cea10d1a6bb6332878af5fffa58f57b227cba0d0229fa90a63cc51.
//
// Solidity: event Revoked(bytes32 indexed commitment, bytes32 indexed internalID)
func (_UserRegistry *UserRegistryFilterer) ParseRevoked(log types.Log) (*UserRegistryRevoked, error) {
	event := new(UserRegistryRevoked)
	if err := _UserRegistry.contract.UnpackLog(event, "Revoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
