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

// LogAnchorMetaData contains all meta data concerning the LogAnchor contract.
var LogAnchorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"chainHash\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"batchId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorAt\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainHash\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchors\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"chainHash\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Anchored\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"chainHash\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AnchorExists\",\"inputs\":[{\"name\":\"batchId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OnlyOperator\",\"inputs\":[]}]",
}

// LogAnchorABI is the input ABI used to generate the binding from.
// Deprecated: Use LogAnchorMetaData.ABI instead.
var LogAnchorABI = LogAnchorMetaData.ABI

// LogAnchor is an auto generated Go binding around an Ethereum contract.
type LogAnchor struct {
	LogAnchorCaller     // Read-only binding to the contract
	LogAnchorTransactor // Write-only binding to the contract
	LogAnchorFilterer   // Log filterer for contract events
}

// LogAnchorCaller is an auto generated read-only Go binding around an Ethereum contract.
type LogAnchorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LogAnchorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LogAnchorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LogAnchorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LogAnchorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LogAnchorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LogAnchorSession struct {
	Contract     *LogAnchor        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LogAnchorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LogAnchorCallerSession struct {
	Contract *LogAnchorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// LogAnchorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LogAnchorTransactorSession struct {
	Contract     *LogAnchorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// LogAnchorRaw is an auto generated low-level Go binding around an Ethereum contract.
type LogAnchorRaw struct {
	Contract *LogAnchor // Generic contract binding to access the raw methods on
}

// LogAnchorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LogAnchorCallerRaw struct {
	Contract *LogAnchorCaller // Generic read-only contract binding to access the raw methods on
}

// LogAnchorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LogAnchorTransactorRaw struct {
	Contract *LogAnchorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLogAnchor creates a new instance of LogAnchor, bound to a specific deployed contract.
func NewLogAnchor(address common.Address, backend bind.ContractBackend) (*LogAnchor, error) {
	contract, err := bindLogAnchor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LogAnchor{LogAnchorCaller: LogAnchorCaller{contract: contract}, LogAnchorTransactor: LogAnchorTransactor{contract: contract}, LogAnchorFilterer: LogAnchorFilterer{contract: contract}}, nil
}

// NewLogAnchorCaller creates a new read-only instance of LogAnchor, bound to a specific deployed contract.
func NewLogAnchorCaller(address common.Address, caller bind.ContractCaller) (*LogAnchorCaller, error) {
	contract, err := bindLogAnchor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LogAnchorCaller{contract: contract}, nil
}

// NewLogAnchorTransactor creates a new write-only instance of LogAnchor, bound to a specific deployed contract.
func NewLogAnchorTransactor(address common.Address, transactor bind.ContractTransactor) (*LogAnchorTransactor, error) {
	contract, err := bindLogAnchor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LogAnchorTransactor{contract: contract}, nil
}

// NewLogAnchorFilterer creates a new log filterer instance of LogAnchor, bound to a specific deployed contract.
func NewLogAnchorFilterer(address common.Address, filterer bind.ContractFilterer) (*LogAnchorFilterer, error) {
	contract, err := bindLogAnchor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LogAnchorFilterer{contract: contract}, nil
}

// bindLogAnchor binds a generic wrapper to an already deployed contract.
func bindLogAnchor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LogAnchorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LogAnchor *LogAnchorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogAnchor.Contract.LogAnchorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LogAnchor *LogAnchorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogAnchor.Contract.LogAnchorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LogAnchor *LogAnchorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogAnchor.Contract.LogAnchorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LogAnchor *LogAnchorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LogAnchor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LogAnchor *LogAnchorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LogAnchor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LogAnchor *LogAnchorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LogAnchor.Contract.contract.Transact(opts, method, params...)
}

// Anchors is a free data retrieval call binding the contract method 0x368b733e.
//
// Solidity: function anchors(uint256 ) view returns(uint256 chainHash, uint256 blockNumber, uint256 blockTimestamp)
func (_LogAnchor *LogAnchorCaller) Anchors(opts *bind.CallOpts, arg0 *big.Int) (struct {
	ChainHash      *big.Int
	BlockNumber    *big.Int
	BlockTimestamp *big.Int
}, error) {
	var out []interface{}
	err := _LogAnchor.contract.Call(opts, &out, "anchors", arg0)

	outstruct := new(struct {
		ChainHash      *big.Int
		BlockNumber    *big.Int
		BlockTimestamp *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ChainHash = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BlockNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestamp = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Anchors is a free data retrieval call binding the contract method 0x368b733e.
//
// Solidity: function anchors(uint256 ) view returns(uint256 chainHash, uint256 blockNumber, uint256 blockTimestamp)
func (_LogAnchor *LogAnchorSession) Anchors(arg0 *big.Int) (struct {
	ChainHash      *big.Int
	BlockNumber    *big.Int
	BlockTimestamp *big.Int
}, error) {
	return _LogAnchor.Contract.Anchors(&_LogAnchor.CallOpts, arg0)
}

// Anchors is a free data retrieval call binding the contract method 0x368b733e.
//
// Solidity: function anchors(uint256 ) view returns(uint256 chainHash, uint256 blockNumber, uint256 blockTimestamp)
func (_LogAnchor *LogAnchorCallerSession) Anchors(arg0 *big.Int) (struct {
	ChainHash      *big.Int
	BlockNumber    *big.Int
	BlockTimestamp *big.Int
}, error) {
	return _LogAnchor.Contract.Anchors(&_LogAnchor.CallOpts, arg0)
}

// BatchCount is a free data retrieval call binding the contract method 0x06f13056.
//
// Solidity: function batchCount() view returns(uint256)
func (_LogAnchor *LogAnchorCaller) BatchCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LogAnchor.contract.Call(opts, &out, "batchCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BatchCount is a free data retrieval call binding the contract method 0x06f13056.
//
// Solidity: function batchCount() view returns(uint256)
func (_LogAnchor *LogAnchorSession) BatchCount() (*big.Int, error) {
	return _LogAnchor.Contract.BatchCount(&_LogAnchor.CallOpts)
}

// BatchCount is a free data retrieval call binding the contract method 0x06f13056.
//
// Solidity: function batchCount() view returns(uint256)
func (_LogAnchor *LogAnchorCallerSession) BatchCount() (*big.Int, error) {
	return _LogAnchor.Contract.BatchCount(&_LogAnchor.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_LogAnchor *LogAnchorCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LogAnchor.contract.Call(opts, &out, "operator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_LogAnchor *LogAnchorSession) Operator() (common.Address, error) {
	return _LogAnchor.Contract.Operator(&_LogAnchor.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_LogAnchor *LogAnchorCallerSession) Operator() (common.Address, error) {
	return _LogAnchor.Contract.Operator(&_LogAnchor.CallOpts)
}

// Anchor is a paid mutator transaction binding the contract method 0xb5c9d046.
//
// Solidity: function anchor(uint256 chainHash) returns(uint256 batchId)
func (_LogAnchor *LogAnchorTransactor) Anchor(opts *bind.TransactOpts, chainHash *big.Int) (*types.Transaction, error) {
	return _LogAnchor.contract.Transact(opts, "anchor", chainHash)
}

// Anchor is a paid mutator transaction binding the contract method 0xb5c9d046.
//
// Solidity: function anchor(uint256 chainHash) returns(uint256 batchId)
func (_LogAnchor *LogAnchorSession) Anchor(chainHash *big.Int) (*types.Transaction, error) {
	return _LogAnchor.Contract.Anchor(&_LogAnchor.TransactOpts, chainHash)
}

// Anchor is a paid mutator transaction binding the contract method 0xb5c9d046.
//
// Solidity: function anchor(uint256 chainHash) returns(uint256 batchId)
func (_LogAnchor *LogAnchorTransactorSession) Anchor(chainHash *big.Int) (*types.Transaction, error) {
	return _LogAnchor.Contract.Anchor(&_LogAnchor.TransactOpts, chainHash)
}

// AnchorAt is a paid mutator transaction binding the contract method 0x41d4a957.
//
// Solidity: function anchorAt(uint256 batchId, uint256 chainHash) returns(uint256)
func (_LogAnchor *LogAnchorTransactor) AnchorAt(opts *bind.TransactOpts, batchId *big.Int, chainHash *big.Int) (*types.Transaction, error) {
	return _LogAnchor.contract.Transact(opts, "anchorAt", batchId, chainHash)
}

// AnchorAt is a paid mutator transaction binding the contract method 0x41d4a957.
//
// Solidity: function anchorAt(uint256 batchId, uint256 chainHash) returns(uint256)
func (_LogAnchor *LogAnchorSession) AnchorAt(batchId *big.Int, chainHash *big.Int) (*types.Transaction, error) {
	return _LogAnchor.Contract.AnchorAt(&_LogAnchor.TransactOpts, batchId, chainHash)
}

// AnchorAt is a paid mutator transaction binding the contract method 0x41d4a957.
//
// Solidity: function anchorAt(uint256 batchId, uint256 chainHash) returns(uint256)
func (_LogAnchor *LogAnchorTransactorSession) AnchorAt(batchId *big.Int, chainHash *big.Int) (*types.Transaction, error) {
	return _LogAnchor.Contract.AnchorAt(&_LogAnchor.TransactOpts, batchId, chainHash)
}

// LogAnchorAnchoredIterator is returned from FilterAnchored and is used to iterate over the raw logs and unpacked data for Anchored events raised by the LogAnchor contract.
type LogAnchorAnchoredIterator struct {
	Event *LogAnchorAnchored // Event containing the contract specifics and raw log

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
func (it *LogAnchorAnchoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LogAnchorAnchored)
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
		it.Event = new(LogAnchorAnchored)
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
func (it *LogAnchorAnchoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LogAnchorAnchoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LogAnchorAnchored represents a Anchored event raised by the LogAnchor contract.
type LogAnchorAnchored struct {
	BatchId     *big.Int
	ChainHash   *big.Int
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAnchored is a free log retrieval operation binding the contract event 0x39fa7bbaf6245532294d9c0862bec3d244a4e488b4f55529a9c60754606500d7.
//
// Solidity: event Anchored(uint256 indexed batchId, uint256 chainHash, uint256 blockNumber)
func (_LogAnchor *LogAnchorFilterer) FilterAnchored(opts *bind.FilterOpts, batchId []*big.Int) (*LogAnchorAnchoredIterator, error) {

	var batchIdRule []interface{}
	for _, batchIdItem := range batchId {
		batchIdRule = append(batchIdRule, batchIdItem)
	}

	logs, sub, err := _LogAnchor.contract.FilterLogs(opts, "Anchored", batchIdRule)
	if err != nil {
		return nil, err
	}
	return &LogAnchorAnchoredIterator{contract: _LogAnchor.contract, event: "Anchored", logs: logs, sub: sub}, nil
}

// WatchAnchored is a free log subscription operation binding the contract event 0x39fa7bbaf6245532294d9c0862bec3d244a4e488b4f55529a9c60754606500d7.
//
// Solidity: event Anchored(uint256 indexed batchId, uint256 chainHash, uint256 blockNumber)
func (_LogAnchor *LogAnchorFilterer) WatchAnchored(opts *bind.WatchOpts, sink chan<- *LogAnchorAnchored, batchId []*big.Int) (event.Subscription, error) {

	var batchIdRule []interface{}
	for _, batchIdItem := range batchId {
		batchIdRule = append(batchIdRule, batchIdItem)
	}

	logs, sub, err := _LogAnchor.contract.WatchLogs(opts, "Anchored", batchIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LogAnchorAnchored)
				if err := _LogAnchor.contract.UnpackLog(event, "Anchored", log); err != nil {
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

// ParseAnchored is a log parse operation binding the contract event 0x39fa7bbaf6245532294d9c0862bec3d244a4e488b4f55529a9c60754606500d7.
//
// Solidity: event Anchored(uint256 indexed batchId, uint256 chainHash, uint256 blockNumber)
func (_LogAnchor *LogAnchorFilterer) ParseAnchored(log types.Log) (*LogAnchorAnchored, error) {
	event := new(LogAnchorAnchored)
	if err := _LogAnchor.contract.UnpackLog(event, "Anchored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
