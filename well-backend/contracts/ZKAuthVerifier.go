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

// ZKAuthVerifierMetaData contains all meta data concerning the ZKAuthVerifier contract.
var ZKAuthVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_verifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_sessionDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authenticate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pA\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pB\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pC\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"nullifier\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"actionId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getSession\",\"inputs\":[{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"revokeSession\",\"inputs\":[{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeSessionByCommitter\",\"inputs\":[{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sessionDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"sessions\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usedNullifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractGroth16Verifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifySession\",\"inputs\":[{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Authenticated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SessionRevoked\",\"inputs\":[{\"name\":\"sessionToken\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSession\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NullifierAlreadyUsed\",\"inputs\":[{\"name\":\"nullifier\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SessionExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// ZKAuthVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ZKAuthVerifierMetaData.ABI instead.
var ZKAuthVerifierABI = ZKAuthVerifierMetaData.ABI

// ZKAuthVerifier is an auto generated Go binding around an Ethereum contract.
type ZKAuthVerifier struct {
	ZKAuthVerifierCaller     // Read-only binding to the contract
	ZKAuthVerifierTransactor // Write-only binding to the contract
	ZKAuthVerifierFilterer   // Log filterer for contract events
}

// ZKAuthVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZKAuthVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKAuthVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZKAuthVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKAuthVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZKAuthVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKAuthVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZKAuthVerifierSession struct {
	Contract     *ZKAuthVerifier   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZKAuthVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZKAuthVerifierCallerSession struct {
	Contract *ZKAuthVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ZKAuthVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZKAuthVerifierTransactorSession struct {
	Contract     *ZKAuthVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ZKAuthVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZKAuthVerifierRaw struct {
	Contract *ZKAuthVerifier // Generic contract binding to access the raw methods on
}

// ZKAuthVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZKAuthVerifierCallerRaw struct {
	Contract *ZKAuthVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ZKAuthVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZKAuthVerifierTransactorRaw struct {
	Contract *ZKAuthVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZKAuthVerifier creates a new instance of ZKAuthVerifier, bound to a specific deployed contract.
func NewZKAuthVerifier(address common.Address, backend bind.ContractBackend) (*ZKAuthVerifier, error) {
	contract, err := bindZKAuthVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZKAuthVerifier{ZKAuthVerifierCaller: ZKAuthVerifierCaller{contract: contract}, ZKAuthVerifierTransactor: ZKAuthVerifierTransactor{contract: contract}, ZKAuthVerifierFilterer: ZKAuthVerifierFilterer{contract: contract}}, nil
}

// NewZKAuthVerifierCaller creates a new read-only instance of ZKAuthVerifier, bound to a specific deployed contract.
func NewZKAuthVerifierCaller(address common.Address, caller bind.ContractCaller) (*ZKAuthVerifierCaller, error) {
	contract, err := bindZKAuthVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZKAuthVerifierCaller{contract: contract}, nil
}

// NewZKAuthVerifierTransactor creates a new write-only instance of ZKAuthVerifier, bound to a specific deployed contract.
func NewZKAuthVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ZKAuthVerifierTransactor, error) {
	contract, err := bindZKAuthVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZKAuthVerifierTransactor{contract: contract}, nil
}

// NewZKAuthVerifierFilterer creates a new log filterer instance of ZKAuthVerifier, bound to a specific deployed contract.
func NewZKAuthVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ZKAuthVerifierFilterer, error) {
	contract, err := bindZKAuthVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZKAuthVerifierFilterer{contract: contract}, nil
}

// bindZKAuthVerifier binds a generic wrapper to an already deployed contract.
func bindZKAuthVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ZKAuthVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZKAuthVerifier *ZKAuthVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZKAuthVerifier.Contract.ZKAuthVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZKAuthVerifier *ZKAuthVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.ZKAuthVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZKAuthVerifier *ZKAuthVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.ZKAuthVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZKAuthVerifier *ZKAuthVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZKAuthVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZKAuthVerifier *ZKAuthVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZKAuthVerifier *ZKAuthVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.contract.Transact(opts, method, params...)
}

// GetSession is a free data retrieval call binding the contract method 0x39b240bd.
//
// Solidity: function getSession(bytes32 sessionToken) view returns(address user, uint48 expiresAt, bytes32 commitment)
func (_ZKAuthVerifier *ZKAuthVerifierCaller) GetSession(opts *bind.CallOpts, sessionToken [32]byte) (struct {
	User       common.Address
	ExpiresAt  *big.Int
	Commitment [32]byte
}, error) {
	var out []interface{}
	err := _ZKAuthVerifier.contract.Call(opts, &out, "getSession", sessionToken)

	outstruct := new(struct {
		User       common.Address
		ExpiresAt  *big.Int
		Commitment [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ExpiresAt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Commitment = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// GetSession is a free data retrieval call binding the contract method 0x39b240bd.
//
// Solidity: function getSession(bytes32 sessionToken) view returns(address user, uint48 expiresAt, bytes32 commitment)
func (_ZKAuthVerifier *ZKAuthVerifierSession) GetSession(sessionToken [32]byte) (struct {
	User       common.Address
	ExpiresAt  *big.Int
	Commitment [32]byte
}, error) {
	return _ZKAuthVerifier.Contract.GetSession(&_ZKAuthVerifier.CallOpts, sessionToken)
}

// GetSession is a free data retrieval call binding the contract method 0x39b240bd.
//
// Solidity: function getSession(bytes32 sessionToken) view returns(address user, uint48 expiresAt, bytes32 commitment)
func (_ZKAuthVerifier *ZKAuthVerifierCallerSession) GetSession(sessionToken [32]byte) (struct {
	User       common.Address
	ExpiresAt  *big.Int
	Commitment [32]byte
}, error) {
	return _ZKAuthVerifier.Contract.GetSession(&_ZKAuthVerifier.CallOpts, sessionToken)
}

// SessionDuration is a free data retrieval call binding the contract method 0xece35d01.
//
// Solidity: function sessionDuration() view returns(uint256)
func (_ZKAuthVerifier *ZKAuthVerifierCaller) SessionDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZKAuthVerifier.contract.Call(opts, &out, "sessionDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SessionDuration is a free data retrieval call binding the contract method 0xece35d01.
//
// Solidity: function sessionDuration() view returns(uint256)
func (_ZKAuthVerifier *ZKAuthVerifierSession) SessionDuration() (*big.Int, error) {
	return _ZKAuthVerifier.Contract.SessionDuration(&_ZKAuthVerifier.CallOpts)
}

// SessionDuration is a free data retrieval call binding the contract method 0xece35d01.
//
// Solidity: function sessionDuration() view returns(uint256)
func (_ZKAuthVerifier *ZKAuthVerifierCallerSession) SessionDuration() (*big.Int, error) {
	return _ZKAuthVerifier.Contract.SessionDuration(&_ZKAuthVerifier.CallOpts)
}

// Sessions is a free data retrieval call binding the contract method 0x7dbd2832.
//
// Solidity: function sessions(bytes32 ) view returns(address user, uint48 expiresAt, bytes32 commitment)
func (_ZKAuthVerifier *ZKAuthVerifierCaller) Sessions(opts *bind.CallOpts, arg0 [32]byte) (struct {
	User       common.Address
	ExpiresAt  *big.Int
	Commitment [32]byte
}, error) {
	var out []interface{}
	err := _ZKAuthVerifier.contract.Call(opts, &out, "sessions", arg0)

	outstruct := new(struct {
		User       common.Address
		ExpiresAt  *big.Int
		Commitment [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.User = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ExpiresAt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Commitment = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Sessions is a free data retrieval call binding the contract method 0x7dbd2832.
//
// Solidity: function sessions(bytes32 ) view returns(address user, uint48 expiresAt, bytes32 commitment)
func (_ZKAuthVerifier *ZKAuthVerifierSession) Sessions(arg0 [32]byte) (struct {
	User       common.Address
	ExpiresAt  *big.Int
	Commitment [32]byte
}, error) {
	return _ZKAuthVerifier.Contract.Sessions(&_ZKAuthVerifier.CallOpts, arg0)
}

// Sessions is a free data retrieval call binding the contract method 0x7dbd2832.
//
// Solidity: function sessions(bytes32 ) view returns(address user, uint48 expiresAt, bytes32 commitment)
func (_ZKAuthVerifier *ZKAuthVerifierCallerSession) Sessions(arg0 [32]byte) (struct {
	User       common.Address
	ExpiresAt  *big.Int
	Commitment [32]byte
}, error) {
	return _ZKAuthVerifier.Contract.Sessions(&_ZKAuthVerifier.CallOpts, arg0)
}

// UsedNullifiers is a free data retrieval call binding the contract method 0x206137aa.
//
// Solidity: function usedNullifiers(bytes32 ) view returns(bool)
func (_ZKAuthVerifier *ZKAuthVerifierCaller) UsedNullifiers(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ZKAuthVerifier.contract.Call(opts, &out, "usedNullifiers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedNullifiers is a free data retrieval call binding the contract method 0x206137aa.
//
// Solidity: function usedNullifiers(bytes32 ) view returns(bool)
func (_ZKAuthVerifier *ZKAuthVerifierSession) UsedNullifiers(arg0 [32]byte) (bool, error) {
	return _ZKAuthVerifier.Contract.UsedNullifiers(&_ZKAuthVerifier.CallOpts, arg0)
}

// UsedNullifiers is a free data retrieval call binding the contract method 0x206137aa.
//
// Solidity: function usedNullifiers(bytes32 ) view returns(bool)
func (_ZKAuthVerifier *ZKAuthVerifierCallerSession) UsedNullifiers(arg0 [32]byte) (bool, error) {
	return _ZKAuthVerifier.Contract.UsedNullifiers(&_ZKAuthVerifier.CallOpts, arg0)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_ZKAuthVerifier *ZKAuthVerifierCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ZKAuthVerifier.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_ZKAuthVerifier *ZKAuthVerifierSession) Verifier() (common.Address, error) {
	return _ZKAuthVerifier.Contract.Verifier(&_ZKAuthVerifier.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_ZKAuthVerifier *ZKAuthVerifierCallerSession) Verifier() (common.Address, error) {
	return _ZKAuthVerifier.Contract.Verifier(&_ZKAuthVerifier.CallOpts)
}

// VerifySession is a free data retrieval call binding the contract method 0xab0a6bea.
//
// Solidity: function verifySession(bytes32 sessionToken) view returns(bool)
func (_ZKAuthVerifier *ZKAuthVerifierCaller) VerifySession(opts *bind.CallOpts, sessionToken [32]byte) (bool, error) {
	var out []interface{}
	err := _ZKAuthVerifier.contract.Call(opts, &out, "verifySession", sessionToken)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifySession is a free data retrieval call binding the contract method 0xab0a6bea.
//
// Solidity: function verifySession(bytes32 sessionToken) view returns(bool)
func (_ZKAuthVerifier *ZKAuthVerifierSession) VerifySession(sessionToken [32]byte) (bool, error) {
	return _ZKAuthVerifier.Contract.VerifySession(&_ZKAuthVerifier.CallOpts, sessionToken)
}

// VerifySession is a free data retrieval call binding the contract method 0xab0a6bea.
//
// Solidity: function verifySession(bytes32 sessionToken) view returns(bool)
func (_ZKAuthVerifier *ZKAuthVerifierCallerSession) VerifySession(sessionToken [32]byte) (bool, error) {
	return _ZKAuthVerifier.Contract.VerifySession(&_ZKAuthVerifier.CallOpts, sessionToken)
}

// Authenticate is a paid mutator transaction binding the contract method 0x52211bf5.
//
// Solidity: function authenticate(address user, bytes32 commitment, uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 nullifier, bytes32 actionId) returns(bytes32 sessionToken)
func (_ZKAuthVerifier *ZKAuthVerifierTransactor) Authenticate(opts *bind.TransactOpts, user common.Address, commitment [32]byte, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, nullifier [32]byte, actionId [32]byte) (*types.Transaction, error) {
	return _ZKAuthVerifier.contract.Transact(opts, "authenticate", user, commitment, pA, pB, pC, nullifier, actionId)
}

// Authenticate is a paid mutator transaction binding the contract method 0x52211bf5.
//
// Solidity: function authenticate(address user, bytes32 commitment, uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 nullifier, bytes32 actionId) returns(bytes32 sessionToken)
func (_ZKAuthVerifier *ZKAuthVerifierSession) Authenticate(user common.Address, commitment [32]byte, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, nullifier [32]byte, actionId [32]byte) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.Authenticate(&_ZKAuthVerifier.TransactOpts, user, commitment, pA, pB, pC, nullifier, actionId)
}

// Authenticate is a paid mutator transaction binding the contract method 0x52211bf5.
//
// Solidity: function authenticate(address user, bytes32 commitment, uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 nullifier, bytes32 actionId) returns(bytes32 sessionToken)
func (_ZKAuthVerifier *ZKAuthVerifierTransactorSession) Authenticate(user common.Address, commitment [32]byte, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, nullifier [32]byte, actionId [32]byte) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.Authenticate(&_ZKAuthVerifier.TransactOpts, user, commitment, pA, pB, pC, nullifier, actionId)
}

// RevokeSession is a paid mutator transaction binding the contract method 0xa7fed385.
//
// Solidity: function revokeSession(bytes32 sessionToken) returns()
func (_ZKAuthVerifier *ZKAuthVerifierTransactor) RevokeSession(opts *bind.TransactOpts, sessionToken [32]byte) (*types.Transaction, error) {
	return _ZKAuthVerifier.contract.Transact(opts, "revokeSession", sessionToken)
}

// RevokeSession is a paid mutator transaction binding the contract method 0xa7fed385.
//
// Solidity: function revokeSession(bytes32 sessionToken) returns()
func (_ZKAuthVerifier *ZKAuthVerifierSession) RevokeSession(sessionToken [32]byte) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.RevokeSession(&_ZKAuthVerifier.TransactOpts, sessionToken)
}

// RevokeSession is a paid mutator transaction binding the contract method 0xa7fed385.
//
// Solidity: function revokeSession(bytes32 sessionToken) returns()
func (_ZKAuthVerifier *ZKAuthVerifierTransactorSession) RevokeSession(sessionToken [32]byte) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.RevokeSession(&_ZKAuthVerifier.TransactOpts, sessionToken)
}

// RevokeSessionByCommitter is a paid mutator transaction binding the contract method 0x77dc44ae.
//
// Solidity: function revokeSessionByCommitter(bytes32 sessionToken, address committer) returns()
func (_ZKAuthVerifier *ZKAuthVerifierTransactor) RevokeSessionByCommitter(opts *bind.TransactOpts, sessionToken [32]byte, committer common.Address) (*types.Transaction, error) {
	return _ZKAuthVerifier.contract.Transact(opts, "revokeSessionByCommitter", sessionToken, committer)
}

// RevokeSessionByCommitter is a paid mutator transaction binding the contract method 0x77dc44ae.
//
// Solidity: function revokeSessionByCommitter(bytes32 sessionToken, address committer) returns()
func (_ZKAuthVerifier *ZKAuthVerifierSession) RevokeSessionByCommitter(sessionToken [32]byte, committer common.Address) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.RevokeSessionByCommitter(&_ZKAuthVerifier.TransactOpts, sessionToken, committer)
}

// RevokeSessionByCommitter is a paid mutator transaction binding the contract method 0x77dc44ae.
//
// Solidity: function revokeSessionByCommitter(bytes32 sessionToken, address committer) returns()
func (_ZKAuthVerifier *ZKAuthVerifierTransactorSession) RevokeSessionByCommitter(sessionToken [32]byte, committer common.Address) (*types.Transaction, error) {
	return _ZKAuthVerifier.Contract.RevokeSessionByCommitter(&_ZKAuthVerifier.TransactOpts, sessionToken, committer)
}

// ZKAuthVerifierAuthenticatedIterator is returned from FilterAuthenticated and is used to iterate over the raw logs and unpacked data for Authenticated events raised by the ZKAuthVerifier contract.
type ZKAuthVerifierAuthenticatedIterator struct {
	Event *ZKAuthVerifierAuthenticated // Event containing the contract specifics and raw log

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
func (it *ZKAuthVerifierAuthenticatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKAuthVerifierAuthenticated)
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
		it.Event = new(ZKAuthVerifierAuthenticated)
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
func (it *ZKAuthVerifierAuthenticatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKAuthVerifierAuthenticatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKAuthVerifierAuthenticated represents a Authenticated event raised by the ZKAuthVerifier contract.
type ZKAuthVerifierAuthenticated struct {
	User         common.Address
	Commitment   [32]byte
	SessionToken [32]byte
	ExpiresAt    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAuthenticated is a free log retrieval operation binding the contract event 0x33be472cd42e19c5fa37891c4fb00870c7021dfea8d369de13a822d6a87e06c7.
//
// Solidity: event Authenticated(address indexed user, bytes32 indexed commitment, bytes32 indexed sessionToken, uint48 expiresAt)
func (_ZKAuthVerifier *ZKAuthVerifierFilterer) FilterAuthenticated(opts *bind.FilterOpts, user []common.Address, commitment [][32]byte, sessionToken [][32]byte) (*ZKAuthVerifierAuthenticatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var sessionTokenRule []interface{}
	for _, sessionTokenItem := range sessionToken {
		sessionTokenRule = append(sessionTokenRule, sessionTokenItem)
	}

	logs, sub, err := _ZKAuthVerifier.contract.FilterLogs(opts, "Authenticated", userRule, commitmentRule, sessionTokenRule)
	if err != nil {
		return nil, err
	}
	return &ZKAuthVerifierAuthenticatedIterator{contract: _ZKAuthVerifier.contract, event: "Authenticated", logs: logs, sub: sub}, nil
}

// WatchAuthenticated is a free log subscription operation binding the contract event 0x33be472cd42e19c5fa37891c4fb00870c7021dfea8d369de13a822d6a87e06c7.
//
// Solidity: event Authenticated(address indexed user, bytes32 indexed commitment, bytes32 indexed sessionToken, uint48 expiresAt)
func (_ZKAuthVerifier *ZKAuthVerifierFilterer) WatchAuthenticated(opts *bind.WatchOpts, sink chan<- *ZKAuthVerifierAuthenticated, user []common.Address, commitment [][32]byte, sessionToken [][32]byte) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var sessionTokenRule []interface{}
	for _, sessionTokenItem := range sessionToken {
		sessionTokenRule = append(sessionTokenRule, sessionTokenItem)
	}

	logs, sub, err := _ZKAuthVerifier.contract.WatchLogs(opts, "Authenticated", userRule, commitmentRule, sessionTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKAuthVerifierAuthenticated)
				if err := _ZKAuthVerifier.contract.UnpackLog(event, "Authenticated", log); err != nil {
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

// ParseAuthenticated is a log parse operation binding the contract event 0x33be472cd42e19c5fa37891c4fb00870c7021dfea8d369de13a822d6a87e06c7.
//
// Solidity: event Authenticated(address indexed user, bytes32 indexed commitment, bytes32 indexed sessionToken, uint48 expiresAt)
func (_ZKAuthVerifier *ZKAuthVerifierFilterer) ParseAuthenticated(log types.Log) (*ZKAuthVerifierAuthenticated, error) {
	event := new(ZKAuthVerifierAuthenticated)
	if err := _ZKAuthVerifier.contract.UnpackLog(event, "Authenticated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKAuthVerifierSessionRevokedIterator is returned from FilterSessionRevoked and is used to iterate over the raw logs and unpacked data for SessionRevoked events raised by the ZKAuthVerifier contract.
type ZKAuthVerifierSessionRevokedIterator struct {
	Event *ZKAuthVerifierSessionRevoked // Event containing the contract specifics and raw log

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
func (it *ZKAuthVerifierSessionRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKAuthVerifierSessionRevoked)
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
		it.Event = new(ZKAuthVerifierSessionRevoked)
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
func (it *ZKAuthVerifierSessionRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKAuthVerifierSessionRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKAuthVerifierSessionRevoked represents a SessionRevoked event raised by the ZKAuthVerifier contract.
type ZKAuthVerifierSessionRevoked struct {
	SessionToken [32]byte
	User         common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSessionRevoked is a free log retrieval operation binding the contract event 0x53c7d74d36edc91dbf6b53bb0ea80bbdb2ca5894d54f6f62beacc537701822fd.
//
// Solidity: event SessionRevoked(bytes32 indexed sessionToken, address indexed user)
func (_ZKAuthVerifier *ZKAuthVerifierFilterer) FilterSessionRevoked(opts *bind.FilterOpts, sessionToken [][32]byte, user []common.Address) (*ZKAuthVerifierSessionRevokedIterator, error) {

	var sessionTokenRule []interface{}
	for _, sessionTokenItem := range sessionToken {
		sessionTokenRule = append(sessionTokenRule, sessionTokenItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ZKAuthVerifier.contract.FilterLogs(opts, "SessionRevoked", sessionTokenRule, userRule)
	if err != nil {
		return nil, err
	}
	return &ZKAuthVerifierSessionRevokedIterator{contract: _ZKAuthVerifier.contract, event: "SessionRevoked", logs: logs, sub: sub}, nil
}

// WatchSessionRevoked is a free log subscription operation binding the contract event 0x53c7d74d36edc91dbf6b53bb0ea80bbdb2ca5894d54f6f62beacc537701822fd.
//
// Solidity: event SessionRevoked(bytes32 indexed sessionToken, address indexed user)
func (_ZKAuthVerifier *ZKAuthVerifierFilterer) WatchSessionRevoked(opts *bind.WatchOpts, sink chan<- *ZKAuthVerifierSessionRevoked, sessionToken [][32]byte, user []common.Address) (event.Subscription, error) {

	var sessionTokenRule []interface{}
	for _, sessionTokenItem := range sessionToken {
		sessionTokenRule = append(sessionTokenRule, sessionTokenItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ZKAuthVerifier.contract.WatchLogs(opts, "SessionRevoked", sessionTokenRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKAuthVerifierSessionRevoked)
				if err := _ZKAuthVerifier.contract.UnpackLog(event, "SessionRevoked", log); err != nil {
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

// ParseSessionRevoked is a log parse operation binding the contract event 0x53c7d74d36edc91dbf6b53bb0ea80bbdb2ca5894d54f6f62beacc537701822fd.
//
// Solidity: event SessionRevoked(bytes32 indexed sessionToken, address indexed user)
func (_ZKAuthVerifier *ZKAuthVerifierFilterer) ParseSessionRevoked(log types.Log) (*ZKAuthVerifierSessionRevoked, error) {
	event := new(ZKAuthVerifierSessionRevoked)
	if err := _ZKAuthVerifier.contract.UnpackLog(event, "SessionRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
