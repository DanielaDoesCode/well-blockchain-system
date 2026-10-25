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

// ResearchVerifierMetaData contains all meta data concerning the ResearchVerifier contract.
var ResearchVerifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"verifyProof\",\"inputs\":[{\"name\":\"_pA\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"_pB\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"_pC\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"_pubSignals\",\"type\":\"uint256[3]\",\"internalType\":\"uint256[3]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"}]",
}

// ResearchVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use ResearchVerifierMetaData.ABI instead.
var ResearchVerifierABI = ResearchVerifierMetaData.ABI

// ResearchVerifier is an auto generated Go binding around an Ethereum contract.
type ResearchVerifier struct {
	ResearchVerifierCaller     // Read-only binding to the contract
	ResearchVerifierTransactor // Write-only binding to the contract
	ResearchVerifierFilterer   // Log filterer for contract events
}

// ResearchVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type ResearchVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResearchVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ResearchVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResearchVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ResearchVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResearchVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ResearchVerifierSession struct {
	Contract     *ResearchVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ResearchVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ResearchVerifierCallerSession struct {
	Contract *ResearchVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ResearchVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ResearchVerifierTransactorSession struct {
	Contract     *ResearchVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ResearchVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type ResearchVerifierRaw struct {
	Contract *ResearchVerifier // Generic contract binding to access the raw methods on
}

// ResearchVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ResearchVerifierCallerRaw struct {
	Contract *ResearchVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// ResearchVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ResearchVerifierTransactorRaw struct {
	Contract *ResearchVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewResearchVerifier creates a new instance of ResearchVerifier, bound to a specific deployed contract.
func NewResearchVerifier(address common.Address, backend bind.ContractBackend) (*ResearchVerifier, error) {
	contract, err := bindResearchVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ResearchVerifier{ResearchVerifierCaller: ResearchVerifierCaller{contract: contract}, ResearchVerifierTransactor: ResearchVerifierTransactor{contract: contract}, ResearchVerifierFilterer: ResearchVerifierFilterer{contract: contract}}, nil
}

// NewResearchVerifierCaller creates a new read-only instance of ResearchVerifier, bound to a specific deployed contract.
func NewResearchVerifierCaller(address common.Address, caller bind.ContractCaller) (*ResearchVerifierCaller, error) {
	contract, err := bindResearchVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ResearchVerifierCaller{contract: contract}, nil
}

// NewResearchVerifierTransactor creates a new write-only instance of ResearchVerifier, bound to a specific deployed contract.
func NewResearchVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*ResearchVerifierTransactor, error) {
	contract, err := bindResearchVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ResearchVerifierTransactor{contract: contract}, nil
}

// NewResearchVerifierFilterer creates a new log filterer instance of ResearchVerifier, bound to a specific deployed contract.
func NewResearchVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*ResearchVerifierFilterer, error) {
	contract, err := bindResearchVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ResearchVerifierFilterer{contract: contract}, nil
}

// bindResearchVerifier binds a generic wrapper to an already deployed contract.
func bindResearchVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ResearchVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ResearchVerifier *ResearchVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ResearchVerifier.Contract.ResearchVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ResearchVerifier *ResearchVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ResearchVerifier.Contract.ResearchVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ResearchVerifier *ResearchVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ResearchVerifier.Contract.ResearchVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ResearchVerifier *ResearchVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ResearchVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ResearchVerifier *ResearchVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ResearchVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ResearchVerifier *ResearchVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ResearchVerifier.Contract.contract.Transact(opts, method, params...)
}

// VerifyProof is a free data retrieval call binding the contract method 0x11479fea.
//
// Solidity: function verifyProof(uint256[2] _pA, uint256[2][2] _pB, uint256[2] _pC, uint256[3] _pubSignals) view returns(bool)
func (_ResearchVerifier *ResearchVerifierCaller) VerifyProof(opts *bind.CallOpts, _pA [2]*big.Int, _pB [2][2]*big.Int, _pC [2]*big.Int, _pubSignals [3]*big.Int) (bool, error) {
	var out []interface{}
	err := _ResearchVerifier.contract.Call(opts, &out, "verifyProof", _pA, _pB, _pC, _pubSignals)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyProof is a free data retrieval call binding the contract method 0x11479fea.
//
// Solidity: function verifyProof(uint256[2] _pA, uint256[2][2] _pB, uint256[2] _pC, uint256[3] _pubSignals) view returns(bool)
func (_ResearchVerifier *ResearchVerifierSession) VerifyProof(_pA [2]*big.Int, _pB [2][2]*big.Int, _pC [2]*big.Int, _pubSignals [3]*big.Int) (bool, error) {
	return _ResearchVerifier.Contract.VerifyProof(&_ResearchVerifier.CallOpts, _pA, _pB, _pC, _pubSignals)
}

// VerifyProof is a free data retrieval call binding the contract method 0x11479fea.
//
// Solidity: function verifyProof(uint256[2] _pA, uint256[2][2] _pB, uint256[2] _pC, uint256[3] _pubSignals) view returns(bool)
func (_ResearchVerifier *ResearchVerifierCallerSession) VerifyProof(_pA [2]*big.Int, _pB [2][2]*big.Int, _pC [2]*big.Int, _pubSignals [3]*big.Int) (bool, error) {
	return _ResearchVerifier.Contract.VerifyProof(&_ResearchVerifier.CallOpts, _pA, _pB, _pC, _pubSignals)
}
