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

// ResearchContributerMetaData contains all meta data concerning the ResearchContributer contract.
var ResearchContributerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_researchVerifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_registry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_attestationManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"attestationManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAttestationManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contribute\",\"inputs\":[{\"name\":\"pA\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"pB\",\"type\":\"uint256[2][2]\",\"internalType\":\"uint256[2][2]\"},{\"name\":\"pC\",\"type\":\"uint256[2]\",\"internalType\":\"uint256[2]\"},{\"name\":\"commitment\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"nullifier\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"merkleRoot\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractWellRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"researchVerifier\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIResearchVerifier\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"usedNullifiers\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ContributionSubmitted\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"nullifier\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CommitmentNotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NullifierAlreadyUsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StaleOrInvalidRoot\",\"inputs\":[]}]",
}

// ResearchContributerABI is the input ABI used to generate the binding from.
// Deprecated: Use ResearchContributerMetaData.ABI instead.
var ResearchContributerABI = ResearchContributerMetaData.ABI

// ResearchContributer is an auto generated Go binding around an Ethereum contract.
type ResearchContributer struct {
	ResearchContributerCaller     // Read-only binding to the contract
	ResearchContributerTransactor // Write-only binding to the contract
	ResearchContributerFilterer   // Log filterer for contract events
}

// ResearchContributerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ResearchContributerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResearchContributerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ResearchContributerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResearchContributerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ResearchContributerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ResearchContributerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ResearchContributerSession struct {
	Contract     *ResearchContributer // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ResearchContributerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ResearchContributerCallerSession struct {
	Contract *ResearchContributerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ResearchContributerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ResearchContributerTransactorSession struct {
	Contract     *ResearchContributerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ResearchContributerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ResearchContributerRaw struct {
	Contract *ResearchContributer // Generic contract binding to access the raw methods on
}

// ResearchContributerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ResearchContributerCallerRaw struct {
	Contract *ResearchContributerCaller // Generic read-only contract binding to access the raw methods on
}

// ResearchContributerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ResearchContributerTransactorRaw struct {
	Contract *ResearchContributerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewResearchContributer creates a new instance of ResearchContributer, bound to a specific deployed contract.
func NewResearchContributer(address common.Address, backend bind.ContractBackend) (*ResearchContributer, error) {
	contract, err := bindResearchContributer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ResearchContributer{ResearchContributerCaller: ResearchContributerCaller{contract: contract}, ResearchContributerTransactor: ResearchContributerTransactor{contract: contract}, ResearchContributerFilterer: ResearchContributerFilterer{contract: contract}}, nil
}

// NewResearchContributerCaller creates a new read-only instance of ResearchContributer, bound to a specific deployed contract.
func NewResearchContributerCaller(address common.Address, caller bind.ContractCaller) (*ResearchContributerCaller, error) {
	contract, err := bindResearchContributer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ResearchContributerCaller{contract: contract}, nil
}

// NewResearchContributerTransactor creates a new write-only instance of ResearchContributer, bound to a specific deployed contract.
func NewResearchContributerTransactor(address common.Address, transactor bind.ContractTransactor) (*ResearchContributerTransactor, error) {
	contract, err := bindResearchContributer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ResearchContributerTransactor{contract: contract}, nil
}

// NewResearchContributerFilterer creates a new log filterer instance of ResearchContributer, bound to a specific deployed contract.
func NewResearchContributerFilterer(address common.Address, filterer bind.ContractFilterer) (*ResearchContributerFilterer, error) {
	contract, err := bindResearchContributer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ResearchContributerFilterer{contract: contract}, nil
}

// bindResearchContributer binds a generic wrapper to an already deployed contract.
func bindResearchContributer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ResearchContributerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ResearchContributer *ResearchContributerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ResearchContributer.Contract.ResearchContributerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ResearchContributer *ResearchContributerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ResearchContributer.Contract.ResearchContributerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ResearchContributer *ResearchContributerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ResearchContributer.Contract.ResearchContributerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ResearchContributer *ResearchContributerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ResearchContributer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ResearchContributer *ResearchContributerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ResearchContributer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ResearchContributer *ResearchContributerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ResearchContributer.Contract.contract.Transact(opts, method, params...)
}

// AttestationManager is a free data retrieval call binding the contract method 0xb1c16bf1.
//
// Solidity: function attestationManager() view returns(address)
func (_ResearchContributer *ResearchContributerCaller) AttestationManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ResearchContributer.contract.Call(opts, &out, "attestationManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AttestationManager is a free data retrieval call binding the contract method 0xb1c16bf1.
//
// Solidity: function attestationManager() view returns(address)
func (_ResearchContributer *ResearchContributerSession) AttestationManager() (common.Address, error) {
	return _ResearchContributer.Contract.AttestationManager(&_ResearchContributer.CallOpts)
}

// AttestationManager is a free data retrieval call binding the contract method 0xb1c16bf1.
//
// Solidity: function attestationManager() view returns(address)
func (_ResearchContributer *ResearchContributerCallerSession) AttestationManager() (common.Address, error) {
	return _ResearchContributer.Contract.AttestationManager(&_ResearchContributer.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_ResearchContributer *ResearchContributerCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ResearchContributer.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_ResearchContributer *ResearchContributerSession) Registry() (common.Address, error) {
	return _ResearchContributer.Contract.Registry(&_ResearchContributer.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_ResearchContributer *ResearchContributerCallerSession) Registry() (common.Address, error) {
	return _ResearchContributer.Contract.Registry(&_ResearchContributer.CallOpts)
}

// ResearchVerifier is a free data retrieval call binding the contract method 0x6dc468ab.
//
// Solidity: function researchVerifier() view returns(address)
func (_ResearchContributer *ResearchContributerCaller) ResearchVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ResearchContributer.contract.Call(opts, &out, "researchVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ResearchVerifier is a free data retrieval call binding the contract method 0x6dc468ab.
//
// Solidity: function researchVerifier() view returns(address)
func (_ResearchContributer *ResearchContributerSession) ResearchVerifier() (common.Address, error) {
	return _ResearchContributer.Contract.ResearchVerifier(&_ResearchContributer.CallOpts)
}

// ResearchVerifier is a free data retrieval call binding the contract method 0x6dc468ab.
//
// Solidity: function researchVerifier() view returns(address)
func (_ResearchContributer *ResearchContributerCallerSession) ResearchVerifier() (common.Address, error) {
	return _ResearchContributer.Contract.ResearchVerifier(&_ResearchContributer.CallOpts)
}

// UsedNullifiers is a free data retrieval call binding the contract method 0x206137aa.
//
// Solidity: function usedNullifiers(bytes32 ) view returns(bool)
func (_ResearchContributer *ResearchContributerCaller) UsedNullifiers(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _ResearchContributer.contract.Call(opts, &out, "usedNullifiers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedNullifiers is a free data retrieval call binding the contract method 0x206137aa.
//
// Solidity: function usedNullifiers(bytes32 ) view returns(bool)
func (_ResearchContributer *ResearchContributerSession) UsedNullifiers(arg0 [32]byte) (bool, error) {
	return _ResearchContributer.Contract.UsedNullifiers(&_ResearchContributer.CallOpts, arg0)
}

// UsedNullifiers is a free data retrieval call binding the contract method 0x206137aa.
//
// Solidity: function usedNullifiers(bytes32 ) view returns(bool)
func (_ResearchContributer *ResearchContributerCallerSession) UsedNullifiers(arg0 [32]byte) (bool, error) {
	return _ResearchContributer.Contract.UsedNullifiers(&_ResearchContributer.CallOpts, arg0)
}

// Contribute is a paid mutator transaction binding the contract method 0x14595ea5.
//
// Solidity: function contribute(uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 commitment, bytes32 nullifier, uint256 merkleRoot) returns()
func (_ResearchContributer *ResearchContributerTransactor) Contribute(opts *bind.TransactOpts, pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, commitment [32]byte, nullifier [32]byte, merkleRoot *big.Int) (*types.Transaction, error) {
	return _ResearchContributer.contract.Transact(opts, "contribute", pA, pB, pC, commitment, nullifier, merkleRoot)
}

// Contribute is a paid mutator transaction binding the contract method 0x14595ea5.
//
// Solidity: function contribute(uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 commitment, bytes32 nullifier, uint256 merkleRoot) returns()
func (_ResearchContributer *ResearchContributerSession) Contribute(pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, commitment [32]byte, nullifier [32]byte, merkleRoot *big.Int) (*types.Transaction, error) {
	return _ResearchContributer.Contract.Contribute(&_ResearchContributer.TransactOpts, pA, pB, pC, commitment, nullifier, merkleRoot)
}

// Contribute is a paid mutator transaction binding the contract method 0x14595ea5.
//
// Solidity: function contribute(uint256[2] pA, uint256[2][2] pB, uint256[2] pC, bytes32 commitment, bytes32 nullifier, uint256 merkleRoot) returns()
func (_ResearchContributer *ResearchContributerTransactorSession) Contribute(pA [2]*big.Int, pB [2][2]*big.Int, pC [2]*big.Int, commitment [32]byte, nullifier [32]byte, merkleRoot *big.Int) (*types.Transaction, error) {
	return _ResearchContributer.Contract.Contribute(&_ResearchContributer.TransactOpts, pA, pB, pC, commitment, nullifier, merkleRoot)
}

// ResearchContributerContributionSubmittedIterator is returned from FilterContributionSubmitted and is used to iterate over the raw logs and unpacked data for ContributionSubmitted events raised by the ResearchContributer contract.
type ResearchContributerContributionSubmittedIterator struct {
	Event *ResearchContributerContributionSubmitted // Event containing the contract specifics and raw log

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
func (it *ResearchContributerContributionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ResearchContributerContributionSubmitted)
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
		it.Event = new(ResearchContributerContributionSubmitted)
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
func (it *ResearchContributerContributionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ResearchContributerContributionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ResearchContributerContributionSubmitted represents a ContributionSubmitted event raised by the ResearchContributer contract.
type ResearchContributerContributionSubmitted struct {
	Commitment [32]byte
	Nullifier  [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterContributionSubmitted is a free log retrieval operation binding the contract event 0xb7be48326ba4cb9759b8fd532bdc45a9b80fa501c82344a1e5368e41138ab203.
//
// Solidity: event ContributionSubmitted(bytes32 indexed commitment, bytes32 indexed nullifier)
func (_ResearchContributer *ResearchContributerFilterer) FilterContributionSubmitted(opts *bind.FilterOpts, commitment [][32]byte, nullifier [][32]byte) (*ResearchContributerContributionSubmittedIterator, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var nullifierRule []interface{}
	for _, nullifierItem := range nullifier {
		nullifierRule = append(nullifierRule, nullifierItem)
	}

	logs, sub, err := _ResearchContributer.contract.FilterLogs(opts, "ContributionSubmitted", commitmentRule, nullifierRule)
	if err != nil {
		return nil, err
	}
	return &ResearchContributerContributionSubmittedIterator{contract: _ResearchContributer.contract, event: "ContributionSubmitted", logs: logs, sub: sub}, nil
}

// WatchContributionSubmitted is a free log subscription operation binding the contract event 0xb7be48326ba4cb9759b8fd532bdc45a9b80fa501c82344a1e5368e41138ab203.
//
// Solidity: event ContributionSubmitted(bytes32 indexed commitment, bytes32 indexed nullifier)
func (_ResearchContributer *ResearchContributerFilterer) WatchContributionSubmitted(opts *bind.WatchOpts, sink chan<- *ResearchContributerContributionSubmitted, commitment [][32]byte, nullifier [][32]byte) (event.Subscription, error) {

	var commitmentRule []interface{}
	for _, commitmentItem := range commitment {
		commitmentRule = append(commitmentRule, commitmentItem)
	}
	var nullifierRule []interface{}
	for _, nullifierItem := range nullifier {
		nullifierRule = append(nullifierRule, nullifierItem)
	}

	logs, sub, err := _ResearchContributer.contract.WatchLogs(opts, "ContributionSubmitted", commitmentRule, nullifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ResearchContributerContributionSubmitted)
				if err := _ResearchContributer.contract.UnpackLog(event, "ContributionSubmitted", log); err != nil {
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

// ParseContributionSubmitted is a log parse operation binding the contract event 0xb7be48326ba4cb9759b8fd532bdc45a9b80fa501c82344a1e5368e41138ab203.
//
// Solidity: event ContributionSubmitted(bytes32 indexed commitment, bytes32 indexed nullifier)
func (_ResearchContributer *ResearchContributerFilterer) ParseContributionSubmitted(log types.Log) (*ResearchContributerContributionSubmitted, error) {
	event := new(ResearchContributerContributionSubmitted)
	if err := _ResearchContributer.contract.UnpackLog(event, "ContributionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
