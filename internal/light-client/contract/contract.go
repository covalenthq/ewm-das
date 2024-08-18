// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"threshold\",\"type\":\"uint64\"}],\"name\":\"BlockHeightSubmissionThresholdChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockSpecimenHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"resultHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"storageURL\",\"type\":\"string\"}],\"name\":\"BlockResultProductionProofSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validatorBitMap\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"blockSpecimenHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"resulthash\",\"type\":\"bytes32\"}],\"name\":\"BlockResultQuorum\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"newBlockResultRewardAllocation\",\"type\":\"uint128\"}],\"name\":\"BlockResultRewardChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockOnTargetChain\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockOnCurrentChain\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"secondsPerBlockTargetChain\",\"type\":\"uint256\"}],\"name\":\"ChainSyncDataChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxSubmissions\",\"type\":\"uint256\"}],\"name\":\"MaxSubmissionsPerBlockHeightChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"newStakeRequirement\",\"type\":\"uint128\"}],\"name\":\"MinimumRequiredStakeChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"nthBlock\",\"type\":\"uint64\"}],\"name\":\"NthBlockChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"activeOperatorCount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"OperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"QuorumNotReached\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newSessionDuration\",\"type\":\"uint64\"}],\"name\":\"ResultSessionDurationChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"minSubmissions\",\"type\":\"uint64\"}],\"name\":\"ResultSessionMinSubmissionChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newQuorumThreshold\",\"type\":\"uint256\"}],\"name\":\"ResultSessionQuorumChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"secondsPerBlockCurrentChain\",\"type\":\"uint64\"}],\"name\":\"SecondsPerBlockCurrentChainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"}],\"name\":\"SessionStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newStakingManager\",\"type\":\"address\"}],\"name\":\"StakingManagerChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"ValidatorDisabled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"ValidatorEnabled\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"AUDITOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BLOCK_RESULT_PRODUCER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOVERNANCE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"auditor\",\"type\":\"address\"}],\"name\":\"addAuditor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"addBRPOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"governor\",\"type\":\"address\"}],\"name\":\"addGovernor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"disableValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"enableValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"}],\"name\":\"finalizeResultSession\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_brps\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"__governors\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"__auditors\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBRPRoleData\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"requiredStake\",\"type\":\"uint128\"},{\"internalType\":\"address[]\",\"name\":\"activeMembers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"}],\"name\":\"getChainData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockOnTargetChain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockOnCurrentChain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"secondsPerBlockTargetChain\",\"type\":\"uint256\"},{\"internalType\":\"uint128\",\"name\":\"allowedThreshold\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"maxSubmissionsPerBlockHeight\",\"type\":\"uint128\"},{\"internalType\":\"uint64\",\"name\":\"nthBlock\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"getEnabledOperatorCount\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMetadata\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"stakingManager\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"blockResultRewardAllocation\",\"type\":\"uint128\"},{\"internalType\":\"uint64\",\"name\":\"blockResultSessionDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minSubmissionsRequired\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"blockResultQuorum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"secondsPerBlockCurrentChain\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"getOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"resulthash\",\"type\":\"bytes32\"}],\"name\":\"getURLS\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialGovernor\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingManager\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isSessionOpen\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"validatorId\",\"type\":\"uint128\"}],\"name\":\"isValidatorEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operatorRoles\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"auditor\",\"type\":\"address\"}],\"name\":\"removeAuditor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"removeBRPOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"governor\",\"type\":\"address\"}],\"name\":\"removeGovernor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"threshold\",\"type\":\"uint64\"}],\"name\":\"setBlockHeightSubmissionsThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"newBlockResultReward\",\"type\":\"uint128\"}],\"name\":\"setBlockResultReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newSessionDuration\",\"type\":\"uint64\"}],\"name\":\"setBlockResultSessionDuration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"blockOnTargetChain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockOnCurrentChain\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"secondsPerBlockTargetChain\",\"type\":\"uint256\"}],\"name\":\"setChainSyncData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"maxSubmissions\",\"type\":\"uint64\"}],\"name\":\"setMaxSubmissionsPerBlockHeight\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"minSubmissions\",\"type\":\"uint64\"}],\"name\":\"setMinSubmissionsRequired\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"n\",\"type\":\"uint64\"}],\"name\":\"setNthBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"quorum\",\"type\":\"uint256\"}],\"name\":\"setQuorumThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"secondsPerBlockCurrentChain\",\"type\":\"uint64\"}],\"name\":\"setSecondsPerBlockCurrentChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"stakingManagerAddress\",\"type\":\"address\"}],\"name\":\"setStakingManagerAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainId\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"blockHeight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"blockSpecimenHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"resultHash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"storageURL\",\"type\":\"string\"}],\"name\":\"submitBlockResultProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorIDs\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// AUDITORROLE is a free data retrieval call binding the contract method 0x6e1d616e.
//
// Solidity: function AUDITOR_ROLE() view returns(bytes32)
func (_Contract *ContractCaller) AUDITORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "AUDITOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AUDITORROLE is a free data retrieval call binding the contract method 0x6e1d616e.
//
// Solidity: function AUDITOR_ROLE() view returns(bytes32)
func (_Contract *ContractSession) AUDITORROLE() ([32]byte, error) {
	return _Contract.Contract.AUDITORROLE(&_Contract.CallOpts)
}

// AUDITORROLE is a free data retrieval call binding the contract method 0x6e1d616e.
//
// Solidity: function AUDITOR_ROLE() view returns(bytes32)
func (_Contract *ContractCallerSession) AUDITORROLE() ([32]byte, error) {
	return _Contract.Contract.AUDITORROLE(&_Contract.CallOpts)
}

// BLOCKRESULTPRODUCERROLE is a free data retrieval call binding the contract method 0x41c1278d.
//
// Solidity: function BLOCK_RESULT_PRODUCER_ROLE() view returns(bytes32)
func (_Contract *ContractCaller) BLOCKRESULTPRODUCERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "BLOCK_RESULT_PRODUCER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BLOCKRESULTPRODUCERROLE is a free data retrieval call binding the contract method 0x41c1278d.
//
// Solidity: function BLOCK_RESULT_PRODUCER_ROLE() view returns(bytes32)
func (_Contract *ContractSession) BLOCKRESULTPRODUCERROLE() ([32]byte, error) {
	return _Contract.Contract.BLOCKRESULTPRODUCERROLE(&_Contract.CallOpts)
}

// BLOCKRESULTPRODUCERROLE is a free data retrieval call binding the contract method 0x41c1278d.
//
// Solidity: function BLOCK_RESULT_PRODUCER_ROLE() view returns(bytes32)
func (_Contract *ContractCallerSession) BLOCKRESULTPRODUCERROLE() ([32]byte, error) {
	return _Contract.Contract.BLOCKRESULTPRODUCERROLE(&_Contract.CallOpts)
}

// GOVERNANCEROLE is a free data retrieval call binding the contract method 0xf36c8f5c.
//
// Solidity: function GOVERNANCE_ROLE() view returns(bytes32)
func (_Contract *ContractCaller) GOVERNANCEROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "GOVERNANCE_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GOVERNANCEROLE is a free data retrieval call binding the contract method 0xf36c8f5c.
//
// Solidity: function GOVERNANCE_ROLE() view returns(bytes32)
func (_Contract *ContractSession) GOVERNANCEROLE() ([32]byte, error) {
	return _Contract.Contract.GOVERNANCEROLE(&_Contract.CallOpts)
}

// GOVERNANCEROLE is a free data retrieval call binding the contract method 0xf36c8f5c.
//
// Solidity: function GOVERNANCE_ROLE() view returns(bytes32)
func (_Contract *ContractCallerSession) GOVERNANCEROLE() ([32]byte, error) {
	return _Contract.Contract.GOVERNANCEROLE(&_Contract.CallOpts)
}

// GetAllOperators is a free data retrieval call binding the contract method 0xd911c632.
//
// Solidity: function getAllOperators() view returns(address[] _brps, address[] __governors, address[] __auditors)
func (_Contract *ContractCaller) GetAllOperators(opts *bind.CallOpts) (struct {
	Brps      []common.Address
	Governors []common.Address
	Auditors  []common.Address
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getAllOperators")

	outstruct := new(struct {
		Brps      []common.Address
		Governors []common.Address
		Auditors  []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Brps = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Governors = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.Auditors = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetAllOperators is a free data retrieval call binding the contract method 0xd911c632.
//
// Solidity: function getAllOperators() view returns(address[] _brps, address[] __governors, address[] __auditors)
func (_Contract *ContractSession) GetAllOperators() (struct {
	Brps      []common.Address
	Governors []common.Address
	Auditors  []common.Address
}, error) {
	return _Contract.Contract.GetAllOperators(&_Contract.CallOpts)
}

// GetAllOperators is a free data retrieval call binding the contract method 0xd911c632.
//
// Solidity: function getAllOperators() view returns(address[] _brps, address[] __governors, address[] __auditors)
func (_Contract *ContractCallerSession) GetAllOperators() (struct {
	Brps      []common.Address
	Governors []common.Address
	Auditors  []common.Address
}, error) {
	return _Contract.Contract.GetAllOperators(&_Contract.CallOpts)
}

// GetBRPRoleData is a free data retrieval call binding the contract method 0x920036b7.
//
// Solidity: function getBRPRoleData() view returns(uint128 requiredStake, address[] activeMembers)
func (_Contract *ContractCaller) GetBRPRoleData(opts *bind.CallOpts) (struct {
	RequiredStake *big.Int
	ActiveMembers []common.Address
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getBRPRoleData")

	outstruct := new(struct {
		RequiredStake *big.Int
		ActiveMembers []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RequiredStake = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ActiveMembers = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetBRPRoleData is a free data retrieval call binding the contract method 0x920036b7.
//
// Solidity: function getBRPRoleData() view returns(uint128 requiredStake, address[] activeMembers)
func (_Contract *ContractSession) GetBRPRoleData() (struct {
	RequiredStake *big.Int
	ActiveMembers []common.Address
}, error) {
	return _Contract.Contract.GetBRPRoleData(&_Contract.CallOpts)
}

// GetBRPRoleData is a free data retrieval call binding the contract method 0x920036b7.
//
// Solidity: function getBRPRoleData() view returns(uint128 requiredStake, address[] activeMembers)
func (_Contract *ContractCallerSession) GetBRPRoleData() (struct {
	RequiredStake *big.Int
	ActiveMembers []common.Address
}, error) {
	return _Contract.Contract.GetBRPRoleData(&_Contract.CallOpts)
}

// GetChainData is a free data retrieval call binding the contract method 0x54cfa69f.
//
// Solidity: function getChainData(uint64 chainId) view returns(uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain, uint128 allowedThreshold, uint128 maxSubmissionsPerBlockHeight, uint64 nthBlock)
func (_Contract *ContractCaller) GetChainData(opts *bind.CallOpts, chainId uint64) (struct {
	BlockOnTargetChain           *big.Int
	BlockOnCurrentChain          *big.Int
	SecondsPerBlockTargetChain   *big.Int
	AllowedThreshold             *big.Int
	MaxSubmissionsPerBlockHeight *big.Int
	NthBlock                     uint64
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getChainData", chainId)

	outstruct := new(struct {
		BlockOnTargetChain           *big.Int
		BlockOnCurrentChain          *big.Int
		SecondsPerBlockTargetChain   *big.Int
		AllowedThreshold             *big.Int
		MaxSubmissionsPerBlockHeight *big.Int
		NthBlock                     uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BlockOnTargetChain = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BlockOnCurrentChain = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SecondsPerBlockTargetChain = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.AllowedThreshold = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.MaxSubmissionsPerBlockHeight = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.NthBlock = *abi.ConvertType(out[5], new(uint64)).(*uint64)

	return *outstruct, err

}

// GetChainData is a free data retrieval call binding the contract method 0x54cfa69f.
//
// Solidity: function getChainData(uint64 chainId) view returns(uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain, uint128 allowedThreshold, uint128 maxSubmissionsPerBlockHeight, uint64 nthBlock)
func (_Contract *ContractSession) GetChainData(chainId uint64) (struct {
	BlockOnTargetChain           *big.Int
	BlockOnCurrentChain          *big.Int
	SecondsPerBlockTargetChain   *big.Int
	AllowedThreshold             *big.Int
	MaxSubmissionsPerBlockHeight *big.Int
	NthBlock                     uint64
}, error) {
	return _Contract.Contract.GetChainData(&_Contract.CallOpts, chainId)
}

// GetChainData is a free data retrieval call binding the contract method 0x54cfa69f.
//
// Solidity: function getChainData(uint64 chainId) view returns(uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain, uint128 allowedThreshold, uint128 maxSubmissionsPerBlockHeight, uint64 nthBlock)
func (_Contract *ContractCallerSession) GetChainData(chainId uint64) (struct {
	BlockOnTargetChain           *big.Int
	BlockOnCurrentChain          *big.Int
	SecondsPerBlockTargetChain   *big.Int
	AllowedThreshold             *big.Int
	MaxSubmissionsPerBlockHeight *big.Int
	NthBlock                     uint64
}, error) {
	return _Contract.Contract.GetChainData(&_Contract.CallOpts, chainId)
}

// GetEnabledOperatorCount is a free data retrieval call binding the contract method 0x43b845b5.
//
// Solidity: function getEnabledOperatorCount(uint128 validatorId) view returns(uint128)
func (_Contract *ContractCaller) GetEnabledOperatorCount(opts *bind.CallOpts, validatorId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getEnabledOperatorCount", validatorId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEnabledOperatorCount is a free data retrieval call binding the contract method 0x43b845b5.
//
// Solidity: function getEnabledOperatorCount(uint128 validatorId) view returns(uint128)
func (_Contract *ContractSession) GetEnabledOperatorCount(validatorId *big.Int) (*big.Int, error) {
	return _Contract.Contract.GetEnabledOperatorCount(&_Contract.CallOpts, validatorId)
}

// GetEnabledOperatorCount is a free data retrieval call binding the contract method 0x43b845b5.
//
// Solidity: function getEnabledOperatorCount(uint128 validatorId) view returns(uint128)
func (_Contract *ContractCallerSession) GetEnabledOperatorCount(validatorId *big.Int) (*big.Int, error) {
	return _Contract.Contract.GetEnabledOperatorCount(&_Contract.CallOpts, validatorId)
}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns(address stakingManager, uint128 blockResultRewardAllocation, uint64 blockResultSessionDuration, uint64 minSubmissionsRequired, uint256 blockResultQuorum, uint256 secondsPerBlockCurrentChain)
func (_Contract *ContractCaller) GetMetadata(opts *bind.CallOpts) (struct {
	StakingManager              common.Address
	BlockResultRewardAllocation *big.Int
	BlockResultSessionDuration  uint64
	MinSubmissionsRequired      uint64
	BlockResultQuorum           *big.Int
	SecondsPerBlockCurrentChain *big.Int
}, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getMetadata")

	outstruct := new(struct {
		StakingManager              common.Address
		BlockResultRewardAllocation *big.Int
		BlockResultSessionDuration  uint64
		MinSubmissionsRequired      uint64
		BlockResultQuorum           *big.Int
		SecondsPerBlockCurrentChain *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StakingManager = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.BlockResultRewardAllocation = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockResultSessionDuration = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.MinSubmissionsRequired = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.BlockResultQuorum = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.SecondsPerBlockCurrentChain = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns(address stakingManager, uint128 blockResultRewardAllocation, uint64 blockResultSessionDuration, uint64 minSubmissionsRequired, uint256 blockResultQuorum, uint256 secondsPerBlockCurrentChain)
func (_Contract *ContractSession) GetMetadata() (struct {
	StakingManager              common.Address
	BlockResultRewardAllocation *big.Int
	BlockResultSessionDuration  uint64
	MinSubmissionsRequired      uint64
	BlockResultQuorum           *big.Int
	SecondsPerBlockCurrentChain *big.Int
}, error) {
	return _Contract.Contract.GetMetadata(&_Contract.CallOpts)
}

// GetMetadata is a free data retrieval call binding the contract method 0x7a5b4f59.
//
// Solidity: function getMetadata() view returns(address stakingManager, uint128 blockResultRewardAllocation, uint64 blockResultSessionDuration, uint64 minSubmissionsRequired, uint256 blockResultQuorum, uint256 secondsPerBlockCurrentChain)
func (_Contract *ContractCallerSession) GetMetadata() (struct {
	StakingManager              common.Address
	BlockResultRewardAllocation *big.Int
	BlockResultSessionDuration  uint64
	MinSubmissionsRequired      uint64
	BlockResultQuorum           *big.Int
	SecondsPerBlockCurrentChain *big.Int
}, error) {
	return _Contract.Contract.GetMetadata(&_Contract.CallOpts)
}

// GetOperators is a free data retrieval call binding the contract method 0xd3a8b2a8.
//
// Solidity: function getOperators(uint128 validatorId) view returns(address[])
func (_Contract *ContractCaller) GetOperators(opts *bind.CallOpts, validatorId *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getOperators", validatorId)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperators is a free data retrieval call binding the contract method 0xd3a8b2a8.
//
// Solidity: function getOperators(uint128 validatorId) view returns(address[])
func (_Contract *ContractSession) GetOperators(validatorId *big.Int) ([]common.Address, error) {
	return _Contract.Contract.GetOperators(&_Contract.CallOpts, validatorId)
}

// GetOperators is a free data retrieval call binding the contract method 0xd3a8b2a8.
//
// Solidity: function getOperators(uint128 validatorId) view returns(address[])
func (_Contract *ContractCallerSession) GetOperators(validatorId *big.Int) ([]common.Address, error) {
	return _Contract.Contract.GetOperators(&_Contract.CallOpts, validatorId)
}

// GetURLS is a free data retrieval call binding the contract method 0xd5839da9.
//
// Solidity: function getURLS(bytes32 resulthash) view returns(string[])
func (_Contract *ContractCaller) GetURLS(opts *bind.CallOpts, resulthash [32]byte) ([]string, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getURLS", resulthash)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetURLS is a free data retrieval call binding the contract method 0xd5839da9.
//
// Solidity: function getURLS(bytes32 resulthash) view returns(string[])
func (_Contract *ContractSession) GetURLS(resulthash [32]byte) ([]string, error) {
	return _Contract.Contract.GetURLS(&_Contract.CallOpts, resulthash)
}

// GetURLS is a free data retrieval call binding the contract method 0xd5839da9.
//
// Solidity: function getURLS(bytes32 resulthash) view returns(string[])
func (_Contract *ContractCallerSession) GetURLS(resulthash [32]byte) ([]string, error) {
	return _Contract.Contract.GetURLS(&_Contract.CallOpts, resulthash)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address operator) view returns(bool)
func (_Contract *ContractCaller) IsEnabled(opts *bind.CallOpts, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "isEnabled", operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address operator) view returns(bool)
func (_Contract *ContractSession) IsEnabled(operator common.Address) (bool, error) {
	return _Contract.Contract.IsEnabled(&_Contract.CallOpts, operator)
}

// IsEnabled is a free data retrieval call binding the contract method 0x9015d371.
//
// Solidity: function isEnabled(address operator) view returns(bool)
func (_Contract *ContractCallerSession) IsEnabled(operator common.Address) (bool, error) {
	return _Contract.Contract.IsEnabled(&_Contract.CallOpts, operator)
}

// IsSessionOpen is a free data retrieval call binding the contract method 0xba75fd2f.
//
// Solidity: function isSessionOpen(uint64 chainId, uint64 blockHeight, address operator) view returns(bool)
func (_Contract *ContractCaller) IsSessionOpen(opts *bind.CallOpts, chainId uint64, blockHeight uint64, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "isSessionOpen", chainId, blockHeight, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSessionOpen is a free data retrieval call binding the contract method 0xba75fd2f.
//
// Solidity: function isSessionOpen(uint64 chainId, uint64 blockHeight, address operator) view returns(bool)
func (_Contract *ContractSession) IsSessionOpen(chainId uint64, blockHeight uint64, operator common.Address) (bool, error) {
	return _Contract.Contract.IsSessionOpen(&_Contract.CallOpts, chainId, blockHeight, operator)
}

// IsSessionOpen is a free data retrieval call binding the contract method 0xba75fd2f.
//
// Solidity: function isSessionOpen(uint64 chainId, uint64 blockHeight, address operator) view returns(bool)
func (_Contract *ContractCallerSession) IsSessionOpen(chainId uint64, blockHeight uint64, operator common.Address) (bool, error) {
	return _Contract.Contract.IsSessionOpen(&_Contract.CallOpts, chainId, blockHeight, operator)
}

// IsValidatorEnabled is a free data retrieval call binding the contract method 0x429a481b.
//
// Solidity: function isValidatorEnabled(uint128 validatorId) view returns(bool)
func (_Contract *ContractCaller) IsValidatorEnabled(opts *bind.CallOpts, validatorId *big.Int) (bool, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "isValidatorEnabled", validatorId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorEnabled is a free data retrieval call binding the contract method 0x429a481b.
//
// Solidity: function isValidatorEnabled(uint128 validatorId) view returns(bool)
func (_Contract *ContractSession) IsValidatorEnabled(validatorId *big.Int) (bool, error) {
	return _Contract.Contract.IsValidatorEnabled(&_Contract.CallOpts, validatorId)
}

// IsValidatorEnabled is a free data retrieval call binding the contract method 0x429a481b.
//
// Solidity: function isValidatorEnabled(uint128 validatorId) view returns(bool)
func (_Contract *ContractCallerSession) IsValidatorEnabled(validatorId *big.Int) (bool, error) {
	return _Contract.Contract.IsValidatorEnabled(&_Contract.CallOpts, validatorId)
}

// OperatorRoles is a free data retrieval call binding the contract method 0x6ab9d8e8.
//
// Solidity: function operatorRoles(address ) view returns(bytes32)
func (_Contract *ContractCaller) OperatorRoles(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "operatorRoles", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OperatorRoles is a free data retrieval call binding the contract method 0x6ab9d8e8.
//
// Solidity: function operatorRoles(address ) view returns(bytes32)
func (_Contract *ContractSession) OperatorRoles(arg0 common.Address) ([32]byte, error) {
	return _Contract.Contract.OperatorRoles(&_Contract.CallOpts, arg0)
}

// OperatorRoles is a free data retrieval call binding the contract method 0x6ab9d8e8.
//
// Solidity: function operatorRoles(address ) view returns(bytes32)
func (_Contract *ContractCallerSession) OperatorRoles(arg0 common.Address) ([32]byte, error) {
	return _Contract.Contract.OperatorRoles(&_Contract.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// ValidatorIDs is a free data retrieval call binding the contract method 0x0d92f4ed.
//
// Solidity: function validatorIDs(address ) view returns(uint128)
func (_Contract *ContractCaller) ValidatorIDs(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "validatorIDs", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorIDs is a free data retrieval call binding the contract method 0x0d92f4ed.
//
// Solidity: function validatorIDs(address ) view returns(uint128)
func (_Contract *ContractSession) ValidatorIDs(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.ValidatorIDs(&_Contract.CallOpts, arg0)
}

// ValidatorIDs is a free data retrieval call binding the contract method 0x0d92f4ed.
//
// Solidity: function validatorIDs(address ) view returns(uint128)
func (_Contract *ContractCallerSession) ValidatorIDs(arg0 common.Address) (*big.Int, error) {
	return _Contract.Contract.ValidatorIDs(&_Contract.CallOpts, arg0)
}

// AddAuditor is a paid mutator transaction binding the contract method 0xe429cef1.
//
// Solidity: function addAuditor(address auditor) returns()
func (_Contract *ContractTransactor) AddAuditor(opts *bind.TransactOpts, auditor common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addAuditor", auditor)
}

// AddAuditor is a paid mutator transaction binding the contract method 0xe429cef1.
//
// Solidity: function addAuditor(address auditor) returns()
func (_Contract *ContractSession) AddAuditor(auditor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddAuditor(&_Contract.TransactOpts, auditor)
}

// AddAuditor is a paid mutator transaction binding the contract method 0xe429cef1.
//
// Solidity: function addAuditor(address auditor) returns()
func (_Contract *ContractTransactorSession) AddAuditor(auditor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddAuditor(&_Contract.TransactOpts, auditor)
}

// AddBRPOperator is a paid mutator transaction binding the contract method 0x222388c2.
//
// Solidity: function addBRPOperator(address operator, uint128 validatorId) returns()
func (_Contract *ContractTransactor) AddBRPOperator(opts *bind.TransactOpts, operator common.Address, validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addBRPOperator", operator, validatorId)
}

// AddBRPOperator is a paid mutator transaction binding the contract method 0x222388c2.
//
// Solidity: function addBRPOperator(address operator, uint128 validatorId) returns()
func (_Contract *ContractSession) AddBRPOperator(operator common.Address, validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.AddBRPOperator(&_Contract.TransactOpts, operator, validatorId)
}

// AddBRPOperator is a paid mutator transaction binding the contract method 0x222388c2.
//
// Solidity: function addBRPOperator(address operator, uint128 validatorId) returns()
func (_Contract *ContractTransactorSession) AddBRPOperator(operator common.Address, validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.AddBRPOperator(&_Contract.TransactOpts, operator, validatorId)
}

// AddGovernor is a paid mutator transaction binding the contract method 0x3c4a25d0.
//
// Solidity: function addGovernor(address governor) returns()
func (_Contract *ContractTransactor) AddGovernor(opts *bind.TransactOpts, governor common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addGovernor", governor)
}

// AddGovernor is a paid mutator transaction binding the contract method 0x3c4a25d0.
//
// Solidity: function addGovernor(address governor) returns()
func (_Contract *ContractSession) AddGovernor(governor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddGovernor(&_Contract.TransactOpts, governor)
}

// AddGovernor is a paid mutator transaction binding the contract method 0x3c4a25d0.
//
// Solidity: function addGovernor(address governor) returns()
func (_Contract *ContractTransactorSession) AddGovernor(governor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.AddGovernor(&_Contract.TransactOpts, governor)
}

// DisableValidator is a paid mutator transaction binding the contract method 0x42e45079.
//
// Solidity: function disableValidator(uint128 validatorId) returns()
func (_Contract *ContractTransactor) DisableValidator(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "disableValidator", validatorId)
}

// DisableValidator is a paid mutator transaction binding the contract method 0x42e45079.
//
// Solidity: function disableValidator(uint128 validatorId) returns()
func (_Contract *ContractSession) DisableValidator(validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DisableValidator(&_Contract.TransactOpts, validatorId)
}

// DisableValidator is a paid mutator transaction binding the contract method 0x42e45079.
//
// Solidity: function disableValidator(uint128 validatorId) returns()
func (_Contract *ContractTransactorSession) DisableValidator(validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.DisableValidator(&_Contract.TransactOpts, validatorId)
}

// EnableValidator is a paid mutator transaction binding the contract method 0xedff4b61.
//
// Solidity: function enableValidator(uint128 validatorId) returns()
func (_Contract *ContractTransactor) EnableValidator(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "enableValidator", validatorId)
}

// EnableValidator is a paid mutator transaction binding the contract method 0xedff4b61.
//
// Solidity: function enableValidator(uint128 validatorId) returns()
func (_Contract *ContractSession) EnableValidator(validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.EnableValidator(&_Contract.TransactOpts, validatorId)
}

// EnableValidator is a paid mutator transaction binding the contract method 0xedff4b61.
//
// Solidity: function enableValidator(uint128 validatorId) returns()
func (_Contract *ContractTransactorSession) EnableValidator(validatorId *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.EnableValidator(&_Contract.TransactOpts, validatorId)
}

// FinalizeResultSession is a paid mutator transaction binding the contract method 0xd09796bd.
//
// Solidity: function finalizeResultSession(uint64 chainId, uint64 blockHeight) returns()
func (_Contract *ContractTransactor) FinalizeResultSession(opts *bind.TransactOpts, chainId uint64, blockHeight uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "finalizeResultSession", chainId, blockHeight)
}

// FinalizeResultSession is a paid mutator transaction binding the contract method 0xd09796bd.
//
// Solidity: function finalizeResultSession(uint64 chainId, uint64 blockHeight) returns()
func (_Contract *ContractSession) FinalizeResultSession(chainId uint64, blockHeight uint64) (*types.Transaction, error) {
	return _Contract.Contract.FinalizeResultSession(&_Contract.TransactOpts, chainId, blockHeight)
}

// FinalizeResultSession is a paid mutator transaction binding the contract method 0xd09796bd.
//
// Solidity: function finalizeResultSession(uint64 chainId, uint64 blockHeight) returns()
func (_Contract *ContractTransactorSession) FinalizeResultSession(chainId uint64, blockHeight uint64) (*types.Transaction, error) {
	return _Contract.Contract.FinalizeResultSession(&_Contract.TransactOpts, chainId, blockHeight)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address initialGovernor, address stakingManager) returns()
func (_Contract *ContractTransactor) Initialize(opts *bind.TransactOpts, initialGovernor common.Address, stakingManager common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "initialize", initialGovernor, stakingManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address initialGovernor, address stakingManager) returns()
func (_Contract *ContractSession) Initialize(initialGovernor common.Address, stakingManager common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, initialGovernor, stakingManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address initialGovernor, address stakingManager) returns()
func (_Contract *ContractTransactorSession) Initialize(initialGovernor common.Address, stakingManager common.Address) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, initialGovernor, stakingManager)
}

// RemoveAuditor is a paid mutator transaction binding the contract method 0xe6116cfd.
//
// Solidity: function removeAuditor(address auditor) returns()
func (_Contract *ContractTransactor) RemoveAuditor(opts *bind.TransactOpts, auditor common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "removeAuditor", auditor)
}

// RemoveAuditor is a paid mutator transaction binding the contract method 0xe6116cfd.
//
// Solidity: function removeAuditor(address auditor) returns()
func (_Contract *ContractSession) RemoveAuditor(auditor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveAuditor(&_Contract.TransactOpts, auditor)
}

// RemoveAuditor is a paid mutator transaction binding the contract method 0xe6116cfd.
//
// Solidity: function removeAuditor(address auditor) returns()
func (_Contract *ContractTransactorSession) RemoveAuditor(auditor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveAuditor(&_Contract.TransactOpts, auditor)
}

// RemoveBRPOperator is a paid mutator transaction binding the contract method 0x6b7511cb.
//
// Solidity: function removeBRPOperator(address operator) returns()
func (_Contract *ContractTransactor) RemoveBRPOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "removeBRPOperator", operator)
}

// RemoveBRPOperator is a paid mutator transaction binding the contract method 0x6b7511cb.
//
// Solidity: function removeBRPOperator(address operator) returns()
func (_Contract *ContractSession) RemoveBRPOperator(operator common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveBRPOperator(&_Contract.TransactOpts, operator)
}

// RemoveBRPOperator is a paid mutator transaction binding the contract method 0x6b7511cb.
//
// Solidity: function removeBRPOperator(address operator) returns()
func (_Contract *ContractTransactorSession) RemoveBRPOperator(operator common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveBRPOperator(&_Contract.TransactOpts, operator)
}

// RemoveGovernor is a paid mutator transaction binding the contract method 0xeecdac88.
//
// Solidity: function removeGovernor(address governor) returns()
func (_Contract *ContractTransactor) RemoveGovernor(opts *bind.TransactOpts, governor common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "removeGovernor", governor)
}

// RemoveGovernor is a paid mutator transaction binding the contract method 0xeecdac88.
//
// Solidity: function removeGovernor(address governor) returns()
func (_Contract *ContractSession) RemoveGovernor(governor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveGovernor(&_Contract.TransactOpts, governor)
}

// RemoveGovernor is a paid mutator transaction binding the contract method 0xeecdac88.
//
// Solidity: function removeGovernor(address governor) returns()
func (_Contract *ContractTransactorSession) RemoveGovernor(governor common.Address) (*types.Transaction, error) {
	return _Contract.Contract.RemoveGovernor(&_Contract.TransactOpts, governor)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// SetBlockHeightSubmissionsThreshold is a paid mutator transaction binding the contract method 0x2c58ed42.
//
// Solidity: function setBlockHeightSubmissionsThreshold(uint64 chainId, uint64 threshold) returns()
func (_Contract *ContractTransactor) SetBlockHeightSubmissionsThreshold(opts *bind.TransactOpts, chainId uint64, threshold uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setBlockHeightSubmissionsThreshold", chainId, threshold)
}

// SetBlockHeightSubmissionsThreshold is a paid mutator transaction binding the contract method 0x2c58ed42.
//
// Solidity: function setBlockHeightSubmissionsThreshold(uint64 chainId, uint64 threshold) returns()
func (_Contract *ContractSession) SetBlockHeightSubmissionsThreshold(chainId uint64, threshold uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetBlockHeightSubmissionsThreshold(&_Contract.TransactOpts, chainId, threshold)
}

// SetBlockHeightSubmissionsThreshold is a paid mutator transaction binding the contract method 0x2c58ed42.
//
// Solidity: function setBlockHeightSubmissionsThreshold(uint64 chainId, uint64 threshold) returns()
func (_Contract *ContractTransactorSession) SetBlockHeightSubmissionsThreshold(chainId uint64, threshold uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetBlockHeightSubmissionsThreshold(&_Contract.TransactOpts, chainId, threshold)
}

// SetBlockResultReward is a paid mutator transaction binding the contract method 0x5a5e84f2.
//
// Solidity: function setBlockResultReward(uint128 newBlockResultReward) returns()
func (_Contract *ContractTransactor) SetBlockResultReward(opts *bind.TransactOpts, newBlockResultReward *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setBlockResultReward", newBlockResultReward)
}

// SetBlockResultReward is a paid mutator transaction binding the contract method 0x5a5e84f2.
//
// Solidity: function setBlockResultReward(uint128 newBlockResultReward) returns()
func (_Contract *ContractSession) SetBlockResultReward(newBlockResultReward *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetBlockResultReward(&_Contract.TransactOpts, newBlockResultReward)
}

// SetBlockResultReward is a paid mutator transaction binding the contract method 0x5a5e84f2.
//
// Solidity: function setBlockResultReward(uint128 newBlockResultReward) returns()
func (_Contract *ContractTransactorSession) SetBlockResultReward(newBlockResultReward *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetBlockResultReward(&_Contract.TransactOpts, newBlockResultReward)
}

// SetBlockResultSessionDuration is a paid mutator transaction binding the contract method 0x55791fea.
//
// Solidity: function setBlockResultSessionDuration(uint64 newSessionDuration) returns()
func (_Contract *ContractTransactor) SetBlockResultSessionDuration(opts *bind.TransactOpts, newSessionDuration uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setBlockResultSessionDuration", newSessionDuration)
}

// SetBlockResultSessionDuration is a paid mutator transaction binding the contract method 0x55791fea.
//
// Solidity: function setBlockResultSessionDuration(uint64 newSessionDuration) returns()
func (_Contract *ContractSession) SetBlockResultSessionDuration(newSessionDuration uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetBlockResultSessionDuration(&_Contract.TransactOpts, newSessionDuration)
}

// SetBlockResultSessionDuration is a paid mutator transaction binding the contract method 0x55791fea.
//
// Solidity: function setBlockResultSessionDuration(uint64 newSessionDuration) returns()
func (_Contract *ContractTransactorSession) SetBlockResultSessionDuration(newSessionDuration uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetBlockResultSessionDuration(&_Contract.TransactOpts, newSessionDuration)
}

// SetChainSyncData is a paid mutator transaction binding the contract method 0x99146284.
//
// Solidity: function setChainSyncData(uint64 chainId, uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain) returns()
func (_Contract *ContractTransactor) SetChainSyncData(opts *bind.TransactOpts, chainId uint64, blockOnTargetChain *big.Int, blockOnCurrentChain *big.Int, secondsPerBlockTargetChain *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setChainSyncData", chainId, blockOnTargetChain, blockOnCurrentChain, secondsPerBlockTargetChain)
}

// SetChainSyncData is a paid mutator transaction binding the contract method 0x99146284.
//
// Solidity: function setChainSyncData(uint64 chainId, uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain) returns()
func (_Contract *ContractSession) SetChainSyncData(chainId uint64, blockOnTargetChain *big.Int, blockOnCurrentChain *big.Int, secondsPerBlockTargetChain *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetChainSyncData(&_Contract.TransactOpts, chainId, blockOnTargetChain, blockOnCurrentChain, secondsPerBlockTargetChain)
}

// SetChainSyncData is a paid mutator transaction binding the contract method 0x99146284.
//
// Solidity: function setChainSyncData(uint64 chainId, uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain) returns()
func (_Contract *ContractTransactorSession) SetChainSyncData(chainId uint64, blockOnTargetChain *big.Int, blockOnCurrentChain *big.Int, secondsPerBlockTargetChain *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetChainSyncData(&_Contract.TransactOpts, chainId, blockOnTargetChain, blockOnCurrentChain, secondsPerBlockTargetChain)
}

// SetMaxSubmissionsPerBlockHeight is a paid mutator transaction binding the contract method 0x67585e44.
//
// Solidity: function setMaxSubmissionsPerBlockHeight(uint64 chainId, uint64 maxSubmissions) returns()
func (_Contract *ContractTransactor) SetMaxSubmissionsPerBlockHeight(opts *bind.TransactOpts, chainId uint64, maxSubmissions uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setMaxSubmissionsPerBlockHeight", chainId, maxSubmissions)
}

// SetMaxSubmissionsPerBlockHeight is a paid mutator transaction binding the contract method 0x67585e44.
//
// Solidity: function setMaxSubmissionsPerBlockHeight(uint64 chainId, uint64 maxSubmissions) returns()
func (_Contract *ContractSession) SetMaxSubmissionsPerBlockHeight(chainId uint64, maxSubmissions uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetMaxSubmissionsPerBlockHeight(&_Contract.TransactOpts, chainId, maxSubmissions)
}

// SetMaxSubmissionsPerBlockHeight is a paid mutator transaction binding the contract method 0x67585e44.
//
// Solidity: function setMaxSubmissionsPerBlockHeight(uint64 chainId, uint64 maxSubmissions) returns()
func (_Contract *ContractTransactorSession) SetMaxSubmissionsPerBlockHeight(chainId uint64, maxSubmissions uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetMaxSubmissionsPerBlockHeight(&_Contract.TransactOpts, chainId, maxSubmissions)
}

// SetMinSubmissionsRequired is a paid mutator transaction binding the contract method 0x93742b56.
//
// Solidity: function setMinSubmissionsRequired(uint64 minSubmissions) returns()
func (_Contract *ContractTransactor) SetMinSubmissionsRequired(opts *bind.TransactOpts, minSubmissions uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setMinSubmissionsRequired", minSubmissions)
}

// SetMinSubmissionsRequired is a paid mutator transaction binding the contract method 0x93742b56.
//
// Solidity: function setMinSubmissionsRequired(uint64 minSubmissions) returns()
func (_Contract *ContractSession) SetMinSubmissionsRequired(minSubmissions uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetMinSubmissionsRequired(&_Contract.TransactOpts, minSubmissions)
}

// SetMinSubmissionsRequired is a paid mutator transaction binding the contract method 0x93742b56.
//
// Solidity: function setMinSubmissionsRequired(uint64 minSubmissions) returns()
func (_Contract *ContractTransactorSession) SetMinSubmissionsRequired(minSubmissions uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetMinSubmissionsRequired(&_Contract.TransactOpts, minSubmissions)
}

// SetNthBlock is a paid mutator transaction binding the contract method 0xe3201409.
//
// Solidity: function setNthBlock(uint64 chainId, uint64 n) returns()
func (_Contract *ContractTransactor) SetNthBlock(opts *bind.TransactOpts, chainId uint64, n uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setNthBlock", chainId, n)
}

// SetNthBlock is a paid mutator transaction binding the contract method 0xe3201409.
//
// Solidity: function setNthBlock(uint64 chainId, uint64 n) returns()
func (_Contract *ContractSession) SetNthBlock(chainId uint64, n uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetNthBlock(&_Contract.TransactOpts, chainId, n)
}

// SetNthBlock is a paid mutator transaction binding the contract method 0xe3201409.
//
// Solidity: function setNthBlock(uint64 chainId, uint64 n) returns()
func (_Contract *ContractTransactorSession) SetNthBlock(chainId uint64, n uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetNthBlock(&_Contract.TransactOpts, chainId, n)
}

// SetQuorumThreshold is a paid mutator transaction binding the contract method 0x4524c7e1.
//
// Solidity: function setQuorumThreshold(uint256 quorum) returns()
func (_Contract *ContractTransactor) SetQuorumThreshold(opts *bind.TransactOpts, quorum *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setQuorumThreshold", quorum)
}

// SetQuorumThreshold is a paid mutator transaction binding the contract method 0x4524c7e1.
//
// Solidity: function setQuorumThreshold(uint256 quorum) returns()
func (_Contract *ContractSession) SetQuorumThreshold(quorum *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetQuorumThreshold(&_Contract.TransactOpts, quorum)
}

// SetQuorumThreshold is a paid mutator transaction binding the contract method 0x4524c7e1.
//
// Solidity: function setQuorumThreshold(uint256 quorum) returns()
func (_Contract *ContractTransactorSession) SetQuorumThreshold(quorum *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetQuorumThreshold(&_Contract.TransactOpts, quorum)
}

// SetSecondsPerBlockCurrentChain is a paid mutator transaction binding the contract method 0x2ba719ad.
//
// Solidity: function setSecondsPerBlockCurrentChain(uint64 secondsPerBlockCurrentChain) returns()
func (_Contract *ContractTransactor) SetSecondsPerBlockCurrentChain(opts *bind.TransactOpts, secondsPerBlockCurrentChain uint64) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setSecondsPerBlockCurrentChain", secondsPerBlockCurrentChain)
}

// SetSecondsPerBlockCurrentChain is a paid mutator transaction binding the contract method 0x2ba719ad.
//
// Solidity: function setSecondsPerBlockCurrentChain(uint64 secondsPerBlockCurrentChain) returns()
func (_Contract *ContractSession) SetSecondsPerBlockCurrentChain(secondsPerBlockCurrentChain uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetSecondsPerBlockCurrentChain(&_Contract.TransactOpts, secondsPerBlockCurrentChain)
}

// SetSecondsPerBlockCurrentChain is a paid mutator transaction binding the contract method 0x2ba719ad.
//
// Solidity: function setSecondsPerBlockCurrentChain(uint64 secondsPerBlockCurrentChain) returns()
func (_Contract *ContractTransactorSession) SetSecondsPerBlockCurrentChain(secondsPerBlockCurrentChain uint64) (*types.Transaction, error) {
	return _Contract.Contract.SetSecondsPerBlockCurrentChain(&_Contract.TransactOpts, secondsPerBlockCurrentChain)
}

// SetStakingManagerAddress is a paid mutator transaction binding the contract method 0x37e15bce.
//
// Solidity: function setStakingManagerAddress(address stakingManagerAddress) returns()
func (_Contract *ContractTransactor) SetStakingManagerAddress(opts *bind.TransactOpts, stakingManagerAddress common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setStakingManagerAddress", stakingManagerAddress)
}

// SetStakingManagerAddress is a paid mutator transaction binding the contract method 0x37e15bce.
//
// Solidity: function setStakingManagerAddress(address stakingManagerAddress) returns()
func (_Contract *ContractSession) SetStakingManagerAddress(stakingManagerAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.SetStakingManagerAddress(&_Contract.TransactOpts, stakingManagerAddress)
}

// SetStakingManagerAddress is a paid mutator transaction binding the contract method 0x37e15bce.
//
// Solidity: function setStakingManagerAddress(address stakingManagerAddress) returns()
func (_Contract *ContractTransactorSession) SetStakingManagerAddress(stakingManagerAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.SetStakingManagerAddress(&_Contract.TransactOpts, stakingManagerAddress)
}

// SubmitBlockResultProof is a paid mutator transaction binding the contract method 0x6543413a.
//
// Solidity: function submitBlockResultProof(uint64 chainId, uint64 blockHeight, bytes32 blockSpecimenHash, bytes32 resultHash, string storageURL) returns()
func (_Contract *ContractTransactor) SubmitBlockResultProof(opts *bind.TransactOpts, chainId uint64, blockHeight uint64, blockSpecimenHash [32]byte, resultHash [32]byte, storageURL string) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "submitBlockResultProof", chainId, blockHeight, blockSpecimenHash, resultHash, storageURL)
}

// SubmitBlockResultProof is a paid mutator transaction binding the contract method 0x6543413a.
//
// Solidity: function submitBlockResultProof(uint64 chainId, uint64 blockHeight, bytes32 blockSpecimenHash, bytes32 resultHash, string storageURL) returns()
func (_Contract *ContractSession) SubmitBlockResultProof(chainId uint64, blockHeight uint64, blockSpecimenHash [32]byte, resultHash [32]byte, storageURL string) (*types.Transaction, error) {
	return _Contract.Contract.SubmitBlockResultProof(&_Contract.TransactOpts, chainId, blockHeight, blockSpecimenHash, resultHash, storageURL)
}

// SubmitBlockResultProof is a paid mutator transaction binding the contract method 0x6543413a.
//
// Solidity: function submitBlockResultProof(uint64 chainId, uint64 blockHeight, bytes32 blockSpecimenHash, bytes32 resultHash, string storageURL) returns()
func (_Contract *ContractTransactorSession) SubmitBlockResultProof(chainId uint64, blockHeight uint64, blockSpecimenHash [32]byte, resultHash [32]byte, storageURL string) (*types.Transaction, error) {
	return _Contract.Contract.SubmitBlockResultProof(&_Contract.TransactOpts, chainId, blockHeight, blockSpecimenHash, resultHash, storageURL)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// ContractBlockHeightSubmissionThresholdChangedIterator is returned from FilterBlockHeightSubmissionThresholdChanged and is used to iterate over the raw logs and unpacked data for BlockHeightSubmissionThresholdChanged events raised by the Contract contract.
type ContractBlockHeightSubmissionThresholdChangedIterator struct {
	Event *ContractBlockHeightSubmissionThresholdChanged // Event containing the contract specifics and raw log

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
func (it *ContractBlockHeightSubmissionThresholdChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBlockHeightSubmissionThresholdChanged)
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
		it.Event = new(ContractBlockHeightSubmissionThresholdChanged)
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
func (it *ContractBlockHeightSubmissionThresholdChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBlockHeightSubmissionThresholdChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBlockHeightSubmissionThresholdChanged represents a BlockHeightSubmissionThresholdChanged event raised by the Contract contract.
type ContractBlockHeightSubmissionThresholdChanged struct {
	ChainId   uint64
	Threshold uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBlockHeightSubmissionThresholdChanged is a free log retrieval operation binding the contract event 0x4e4c0afc3a2b327c2f061f8ff5190a491f1042ba8f292a887bab97840947b7a9.
//
// Solidity: event BlockHeightSubmissionThresholdChanged(uint64 indexed chainId, uint64 threshold)
func (_Contract *ContractFilterer) FilterBlockHeightSubmissionThresholdChanged(opts *bind.FilterOpts, chainId []uint64) (*ContractBlockHeightSubmissionThresholdChangedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BlockHeightSubmissionThresholdChanged", chainIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractBlockHeightSubmissionThresholdChangedIterator{contract: _Contract.contract, event: "BlockHeightSubmissionThresholdChanged", logs: logs, sub: sub}, nil
}

// WatchBlockHeightSubmissionThresholdChanged is a free log subscription operation binding the contract event 0x4e4c0afc3a2b327c2f061f8ff5190a491f1042ba8f292a887bab97840947b7a9.
//
// Solidity: event BlockHeightSubmissionThresholdChanged(uint64 indexed chainId, uint64 threshold)
func (_Contract *ContractFilterer) WatchBlockHeightSubmissionThresholdChanged(opts *bind.WatchOpts, sink chan<- *ContractBlockHeightSubmissionThresholdChanged, chainId []uint64) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BlockHeightSubmissionThresholdChanged", chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBlockHeightSubmissionThresholdChanged)
				if err := _Contract.contract.UnpackLog(event, "BlockHeightSubmissionThresholdChanged", log); err != nil {
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

// ParseBlockHeightSubmissionThresholdChanged is a log parse operation binding the contract event 0x4e4c0afc3a2b327c2f061f8ff5190a491f1042ba8f292a887bab97840947b7a9.
//
// Solidity: event BlockHeightSubmissionThresholdChanged(uint64 indexed chainId, uint64 threshold)
func (_Contract *ContractFilterer) ParseBlockHeightSubmissionThresholdChanged(log types.Log) (*ContractBlockHeightSubmissionThresholdChanged, error) {
	event := new(ContractBlockHeightSubmissionThresholdChanged)
	if err := _Contract.contract.UnpackLog(event, "BlockHeightSubmissionThresholdChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBlockResultProductionProofSubmittedIterator is returned from FilterBlockResultProductionProofSubmitted and is used to iterate over the raw logs and unpacked data for BlockResultProductionProofSubmitted events raised by the Contract contract.
type ContractBlockResultProductionProofSubmittedIterator struct {
	Event *ContractBlockResultProductionProofSubmitted // Event containing the contract specifics and raw log

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
func (it *ContractBlockResultProductionProofSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBlockResultProductionProofSubmitted)
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
		it.Event = new(ContractBlockResultProductionProofSubmitted)
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
func (it *ContractBlockResultProductionProofSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBlockResultProductionProofSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBlockResultProductionProofSubmitted represents a BlockResultProductionProofSubmitted event raised by the Contract contract.
type ContractBlockResultProductionProofSubmitted struct {
	ChainId           uint64
	BlockHeight       uint64
	BlockSpecimenHash [32]byte
	ResultHash        [32]byte
	StorageURL        string
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterBlockResultProductionProofSubmitted is a free log retrieval operation binding the contract event 0x508f479d80dbbd88c8372648a5a1cf88212d9a3e8aa0f7f516c32c6ae970ebbc.
//
// Solidity: event BlockResultProductionProofSubmitted(uint64 chainId, uint64 blockHeight, bytes32 blockSpecimenHash, bytes32 resultHash, string storageURL)
func (_Contract *ContractFilterer) FilterBlockResultProductionProofSubmitted(opts *bind.FilterOpts) (*ContractBlockResultProductionProofSubmittedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BlockResultProductionProofSubmitted")
	if err != nil {
		return nil, err
	}
	return &ContractBlockResultProductionProofSubmittedIterator{contract: _Contract.contract, event: "BlockResultProductionProofSubmitted", logs: logs, sub: sub}, nil
}

// WatchBlockResultProductionProofSubmitted is a free log subscription operation binding the contract event 0x508f479d80dbbd88c8372648a5a1cf88212d9a3e8aa0f7f516c32c6ae970ebbc.
//
// Solidity: event BlockResultProductionProofSubmitted(uint64 chainId, uint64 blockHeight, bytes32 blockSpecimenHash, bytes32 resultHash, string storageURL)
func (_Contract *ContractFilterer) WatchBlockResultProductionProofSubmitted(opts *bind.WatchOpts, sink chan<- *ContractBlockResultProductionProofSubmitted) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BlockResultProductionProofSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBlockResultProductionProofSubmitted)
				if err := _Contract.contract.UnpackLog(event, "BlockResultProductionProofSubmitted", log); err != nil {
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

// ParseBlockResultProductionProofSubmitted is a log parse operation binding the contract event 0x508f479d80dbbd88c8372648a5a1cf88212d9a3e8aa0f7f516c32c6ae970ebbc.
//
// Solidity: event BlockResultProductionProofSubmitted(uint64 chainId, uint64 blockHeight, bytes32 blockSpecimenHash, bytes32 resultHash, string storageURL)
func (_Contract *ContractFilterer) ParseBlockResultProductionProofSubmitted(log types.Log) (*ContractBlockResultProductionProofSubmitted, error) {
	event := new(ContractBlockResultProductionProofSubmitted)
	if err := _Contract.contract.UnpackLog(event, "BlockResultProductionProofSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBlockResultQuorumIterator is returned from FilterBlockResultQuorum and is used to iterate over the raw logs and unpacked data for BlockResultQuorum events raised by the Contract contract.
type ContractBlockResultQuorumIterator struct {
	Event *ContractBlockResultQuorum // Event containing the contract specifics and raw log

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
func (it *ContractBlockResultQuorumIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBlockResultQuorum)
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
		it.Event = new(ContractBlockResultQuorum)
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
func (it *ContractBlockResultQuorumIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBlockResultQuorumIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBlockResultQuorum represents a BlockResultQuorum event raised by the Contract contract.
type ContractBlockResultQuorum struct {
	ChainId           uint64
	BlockHeight       uint64
	ValidatorBitMap   *big.Int
	BlockSpecimenHash [32]byte
	Resulthash        [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterBlockResultQuorum is a free log retrieval operation binding the contract event 0xab2d4a77e5823abb8d239b8a7c25c6f53f7239af1ccdb18dfd27c7da0db24a22.
//
// Solidity: event BlockResultQuorum(uint64 indexed chainId, uint64 indexed blockHeight, uint256 validatorBitMap, bytes32 indexed blockSpecimenHash, bytes32 resulthash)
func (_Contract *ContractFilterer) FilterBlockResultQuorum(opts *bind.FilterOpts, chainId []uint64, blockHeight []uint64, blockSpecimenHash [][32]byte) (*ContractBlockResultQuorumIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var blockHeightRule []interface{}
	for _, blockHeightItem := range blockHeight {
		blockHeightRule = append(blockHeightRule, blockHeightItem)
	}

	var blockSpecimenHashRule []interface{}
	for _, blockSpecimenHashItem := range blockSpecimenHash {
		blockSpecimenHashRule = append(blockSpecimenHashRule, blockSpecimenHashItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BlockResultQuorum", chainIdRule, blockHeightRule, blockSpecimenHashRule)
	if err != nil {
		return nil, err
	}
	return &ContractBlockResultQuorumIterator{contract: _Contract.contract, event: "BlockResultQuorum", logs: logs, sub: sub}, nil
}

// WatchBlockResultQuorum is a free log subscription operation binding the contract event 0xab2d4a77e5823abb8d239b8a7c25c6f53f7239af1ccdb18dfd27c7da0db24a22.
//
// Solidity: event BlockResultQuorum(uint64 indexed chainId, uint64 indexed blockHeight, uint256 validatorBitMap, bytes32 indexed blockSpecimenHash, bytes32 resulthash)
func (_Contract *ContractFilterer) WatchBlockResultQuorum(opts *bind.WatchOpts, sink chan<- *ContractBlockResultQuorum, chainId []uint64, blockHeight []uint64, blockSpecimenHash [][32]byte) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var blockHeightRule []interface{}
	for _, blockHeightItem := range blockHeight {
		blockHeightRule = append(blockHeightRule, blockHeightItem)
	}

	var blockSpecimenHashRule []interface{}
	for _, blockSpecimenHashItem := range blockSpecimenHash {
		blockSpecimenHashRule = append(blockSpecimenHashRule, blockSpecimenHashItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BlockResultQuorum", chainIdRule, blockHeightRule, blockSpecimenHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBlockResultQuorum)
				if err := _Contract.contract.UnpackLog(event, "BlockResultQuorum", log); err != nil {
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

// ParseBlockResultQuorum is a log parse operation binding the contract event 0xab2d4a77e5823abb8d239b8a7c25c6f53f7239af1ccdb18dfd27c7da0db24a22.
//
// Solidity: event BlockResultQuorum(uint64 indexed chainId, uint64 indexed blockHeight, uint256 validatorBitMap, bytes32 indexed blockSpecimenHash, bytes32 resulthash)
func (_Contract *ContractFilterer) ParseBlockResultQuorum(log types.Log) (*ContractBlockResultQuorum, error) {
	event := new(ContractBlockResultQuorum)
	if err := _Contract.contract.UnpackLog(event, "BlockResultQuorum", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractBlockResultRewardChangedIterator is returned from FilterBlockResultRewardChanged and is used to iterate over the raw logs and unpacked data for BlockResultRewardChanged events raised by the Contract contract.
type ContractBlockResultRewardChangedIterator struct {
	Event *ContractBlockResultRewardChanged // Event containing the contract specifics and raw log

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
func (it *ContractBlockResultRewardChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBlockResultRewardChanged)
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
		it.Event = new(ContractBlockResultRewardChanged)
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
func (it *ContractBlockResultRewardChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBlockResultRewardChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBlockResultRewardChanged represents a BlockResultRewardChanged event raised by the Contract contract.
type ContractBlockResultRewardChanged struct {
	NewBlockResultRewardAllocation *big.Int
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterBlockResultRewardChanged is a free log retrieval operation binding the contract event 0xa425d7d9b858a4250625446e4504257053202ae93ba0d1dea68dce7ec87a05c5.
//
// Solidity: event BlockResultRewardChanged(uint128 newBlockResultRewardAllocation)
func (_Contract *ContractFilterer) FilterBlockResultRewardChanged(opts *bind.FilterOpts) (*ContractBlockResultRewardChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BlockResultRewardChanged")
	if err != nil {
		return nil, err
	}
	return &ContractBlockResultRewardChangedIterator{contract: _Contract.contract, event: "BlockResultRewardChanged", logs: logs, sub: sub}, nil
}

// WatchBlockResultRewardChanged is a free log subscription operation binding the contract event 0xa425d7d9b858a4250625446e4504257053202ae93ba0d1dea68dce7ec87a05c5.
//
// Solidity: event BlockResultRewardChanged(uint128 newBlockResultRewardAllocation)
func (_Contract *ContractFilterer) WatchBlockResultRewardChanged(opts *bind.WatchOpts, sink chan<- *ContractBlockResultRewardChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BlockResultRewardChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBlockResultRewardChanged)
				if err := _Contract.contract.UnpackLog(event, "BlockResultRewardChanged", log); err != nil {
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

// ParseBlockResultRewardChanged is a log parse operation binding the contract event 0xa425d7d9b858a4250625446e4504257053202ae93ba0d1dea68dce7ec87a05c5.
//
// Solidity: event BlockResultRewardChanged(uint128 newBlockResultRewardAllocation)
func (_Contract *ContractFilterer) ParseBlockResultRewardChanged(log types.Log) (*ContractBlockResultRewardChanged, error) {
	event := new(ContractBlockResultRewardChanged)
	if err := _Contract.contract.UnpackLog(event, "BlockResultRewardChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractChainSyncDataChangedIterator is returned from FilterChainSyncDataChanged and is used to iterate over the raw logs and unpacked data for ChainSyncDataChanged events raised by the Contract contract.
type ContractChainSyncDataChangedIterator struct {
	Event *ContractChainSyncDataChanged // Event containing the contract specifics and raw log

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
func (it *ContractChainSyncDataChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractChainSyncDataChanged)
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
		it.Event = new(ContractChainSyncDataChanged)
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
func (it *ContractChainSyncDataChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractChainSyncDataChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractChainSyncDataChanged represents a ChainSyncDataChanged event raised by the Contract contract.
type ContractChainSyncDataChanged struct {
	ChainId                    uint64
	BlockOnTargetChain         *big.Int
	BlockOnCurrentChain        *big.Int
	SecondsPerBlockTargetChain *big.Int
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterChainSyncDataChanged is a free log retrieval operation binding the contract event 0xfd97af399d19e6be9256c99c8e52b1809cdbc4dc96816739612b6fd4e6d940b0.
//
// Solidity: event ChainSyncDataChanged(uint64 indexed chainId, uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain)
func (_Contract *ContractFilterer) FilterChainSyncDataChanged(opts *bind.FilterOpts, chainId []uint64) (*ContractChainSyncDataChangedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ChainSyncDataChanged", chainIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractChainSyncDataChangedIterator{contract: _Contract.contract, event: "ChainSyncDataChanged", logs: logs, sub: sub}, nil
}

// WatchChainSyncDataChanged is a free log subscription operation binding the contract event 0xfd97af399d19e6be9256c99c8e52b1809cdbc4dc96816739612b6fd4e6d940b0.
//
// Solidity: event ChainSyncDataChanged(uint64 indexed chainId, uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain)
func (_Contract *ContractFilterer) WatchChainSyncDataChanged(opts *bind.WatchOpts, sink chan<- *ContractChainSyncDataChanged, chainId []uint64) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ChainSyncDataChanged", chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractChainSyncDataChanged)
				if err := _Contract.contract.UnpackLog(event, "ChainSyncDataChanged", log); err != nil {
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

// ParseChainSyncDataChanged is a log parse operation binding the contract event 0xfd97af399d19e6be9256c99c8e52b1809cdbc4dc96816739612b6fd4e6d940b0.
//
// Solidity: event ChainSyncDataChanged(uint64 indexed chainId, uint256 blockOnTargetChain, uint256 blockOnCurrentChain, uint256 secondsPerBlockTargetChain)
func (_Contract *ContractFilterer) ParseChainSyncDataChanged(log types.Log) (*ContractChainSyncDataChanged, error) {
	event := new(ContractChainSyncDataChanged)
	if err := _Contract.contract.UnpackLog(event, "ChainSyncDataChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Contract contract.
type ContractInitializedIterator struct {
	Event *ContractInitialized // Event containing the contract specifics and raw log

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
func (it *ContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractInitialized)
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
		it.Event = new(ContractInitialized)
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
func (it *ContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractInitialized represents a Initialized event raised by the Contract contract.
type ContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Contract *ContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContractInitializedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContractInitializedIterator{contract: _Contract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Contract *ContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContractInitialized) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractInitialized)
				if err := _Contract.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Contract *ContractFilterer) ParseInitialized(log types.Log) (*ContractInitialized, error) {
	event := new(ContractInitialized)
	if err := _Contract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractMaxSubmissionsPerBlockHeightChangedIterator is returned from FilterMaxSubmissionsPerBlockHeightChanged and is used to iterate over the raw logs and unpacked data for MaxSubmissionsPerBlockHeightChanged events raised by the Contract contract.
type ContractMaxSubmissionsPerBlockHeightChangedIterator struct {
	Event *ContractMaxSubmissionsPerBlockHeightChanged // Event containing the contract specifics and raw log

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
func (it *ContractMaxSubmissionsPerBlockHeightChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractMaxSubmissionsPerBlockHeightChanged)
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
		it.Event = new(ContractMaxSubmissionsPerBlockHeightChanged)
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
func (it *ContractMaxSubmissionsPerBlockHeightChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractMaxSubmissionsPerBlockHeightChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractMaxSubmissionsPerBlockHeightChanged represents a MaxSubmissionsPerBlockHeightChanged event raised by the Contract contract.
type ContractMaxSubmissionsPerBlockHeightChanged struct {
	MaxSubmissions *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMaxSubmissionsPerBlockHeightChanged is a free log retrieval operation binding the contract event 0x1bca1fb481202bb14258ce1030d54e9e7bafc8b696d96b9eb733826e58a3a030.
//
// Solidity: event MaxSubmissionsPerBlockHeightChanged(uint256 maxSubmissions)
func (_Contract *ContractFilterer) FilterMaxSubmissionsPerBlockHeightChanged(opts *bind.FilterOpts) (*ContractMaxSubmissionsPerBlockHeightChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "MaxSubmissionsPerBlockHeightChanged")
	if err != nil {
		return nil, err
	}
	return &ContractMaxSubmissionsPerBlockHeightChangedIterator{contract: _Contract.contract, event: "MaxSubmissionsPerBlockHeightChanged", logs: logs, sub: sub}, nil
}

// WatchMaxSubmissionsPerBlockHeightChanged is a free log subscription operation binding the contract event 0x1bca1fb481202bb14258ce1030d54e9e7bafc8b696d96b9eb733826e58a3a030.
//
// Solidity: event MaxSubmissionsPerBlockHeightChanged(uint256 maxSubmissions)
func (_Contract *ContractFilterer) WatchMaxSubmissionsPerBlockHeightChanged(opts *bind.WatchOpts, sink chan<- *ContractMaxSubmissionsPerBlockHeightChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "MaxSubmissionsPerBlockHeightChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractMaxSubmissionsPerBlockHeightChanged)
				if err := _Contract.contract.UnpackLog(event, "MaxSubmissionsPerBlockHeightChanged", log); err != nil {
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

// ParseMaxSubmissionsPerBlockHeightChanged is a log parse operation binding the contract event 0x1bca1fb481202bb14258ce1030d54e9e7bafc8b696d96b9eb733826e58a3a030.
//
// Solidity: event MaxSubmissionsPerBlockHeightChanged(uint256 maxSubmissions)
func (_Contract *ContractFilterer) ParseMaxSubmissionsPerBlockHeightChanged(log types.Log) (*ContractMaxSubmissionsPerBlockHeightChanged, error) {
	event := new(ContractMaxSubmissionsPerBlockHeightChanged)
	if err := _Contract.contract.UnpackLog(event, "MaxSubmissionsPerBlockHeightChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractMinimumRequiredStakeChangedIterator is returned from FilterMinimumRequiredStakeChanged and is used to iterate over the raw logs and unpacked data for MinimumRequiredStakeChanged events raised by the Contract contract.
type ContractMinimumRequiredStakeChangedIterator struct {
	Event *ContractMinimumRequiredStakeChanged // Event containing the contract specifics and raw log

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
func (it *ContractMinimumRequiredStakeChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractMinimumRequiredStakeChanged)
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
		it.Event = new(ContractMinimumRequiredStakeChanged)
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
func (it *ContractMinimumRequiredStakeChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractMinimumRequiredStakeChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractMinimumRequiredStakeChanged represents a MinimumRequiredStakeChanged event raised by the Contract contract.
type ContractMinimumRequiredStakeChanged struct {
	NewStakeRequirement *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterMinimumRequiredStakeChanged is a free log retrieval operation binding the contract event 0xb6c040bb0324b47cbf9a620cce03b311e24597626a57322173d5d5465f739d27.
//
// Solidity: event MinimumRequiredStakeChanged(uint128 newStakeRequirement)
func (_Contract *ContractFilterer) FilterMinimumRequiredStakeChanged(opts *bind.FilterOpts) (*ContractMinimumRequiredStakeChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "MinimumRequiredStakeChanged")
	if err != nil {
		return nil, err
	}
	return &ContractMinimumRequiredStakeChangedIterator{contract: _Contract.contract, event: "MinimumRequiredStakeChanged", logs: logs, sub: sub}, nil
}

// WatchMinimumRequiredStakeChanged is a free log subscription operation binding the contract event 0xb6c040bb0324b47cbf9a620cce03b311e24597626a57322173d5d5465f739d27.
//
// Solidity: event MinimumRequiredStakeChanged(uint128 newStakeRequirement)
func (_Contract *ContractFilterer) WatchMinimumRequiredStakeChanged(opts *bind.WatchOpts, sink chan<- *ContractMinimumRequiredStakeChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "MinimumRequiredStakeChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractMinimumRequiredStakeChanged)
				if err := _Contract.contract.UnpackLog(event, "MinimumRequiredStakeChanged", log); err != nil {
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

// ParseMinimumRequiredStakeChanged is a log parse operation binding the contract event 0xb6c040bb0324b47cbf9a620cce03b311e24597626a57322173d5d5465f739d27.
//
// Solidity: event MinimumRequiredStakeChanged(uint128 newStakeRequirement)
func (_Contract *ContractFilterer) ParseMinimumRequiredStakeChanged(log types.Log) (*ContractMinimumRequiredStakeChanged, error) {
	event := new(ContractMinimumRequiredStakeChanged)
	if err := _Contract.contract.UnpackLog(event, "MinimumRequiredStakeChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractNthBlockChangedIterator is returned from FilterNthBlockChanged and is used to iterate over the raw logs and unpacked data for NthBlockChanged events raised by the Contract contract.
type ContractNthBlockChangedIterator struct {
	Event *ContractNthBlockChanged // Event containing the contract specifics and raw log

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
func (it *ContractNthBlockChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractNthBlockChanged)
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
		it.Event = new(ContractNthBlockChanged)
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
func (it *ContractNthBlockChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractNthBlockChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractNthBlockChanged represents a NthBlockChanged event raised by the Contract contract.
type ContractNthBlockChanged struct {
	ChainId  uint64
	NthBlock uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNthBlockChanged is a free log retrieval operation binding the contract event 0xbbfa9310306e8a8485d109f8be6b0a808473ce55d2e94b8ca3447c9ddb2854b4.
//
// Solidity: event NthBlockChanged(uint64 indexed chainId, uint64 indexed nthBlock)
func (_Contract *ContractFilterer) FilterNthBlockChanged(opts *bind.FilterOpts, chainId []uint64, nthBlock []uint64) (*ContractNthBlockChangedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var nthBlockRule []interface{}
	for _, nthBlockItem := range nthBlock {
		nthBlockRule = append(nthBlockRule, nthBlockItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "NthBlockChanged", chainIdRule, nthBlockRule)
	if err != nil {
		return nil, err
	}
	return &ContractNthBlockChangedIterator{contract: _Contract.contract, event: "NthBlockChanged", logs: logs, sub: sub}, nil
}

// WatchNthBlockChanged is a free log subscription operation binding the contract event 0xbbfa9310306e8a8485d109f8be6b0a808473ce55d2e94b8ca3447c9ddb2854b4.
//
// Solidity: event NthBlockChanged(uint64 indexed chainId, uint64 indexed nthBlock)
func (_Contract *ContractFilterer) WatchNthBlockChanged(opts *bind.WatchOpts, sink chan<- *ContractNthBlockChanged, chainId []uint64, nthBlock []uint64) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var nthBlockRule []interface{}
	for _, nthBlockItem := range nthBlock {
		nthBlockRule = append(nthBlockRule, nthBlockItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "NthBlockChanged", chainIdRule, nthBlockRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractNthBlockChanged)
				if err := _Contract.contract.UnpackLog(event, "NthBlockChanged", log); err != nil {
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

// ParseNthBlockChanged is a log parse operation binding the contract event 0xbbfa9310306e8a8485d109f8be6b0a808473ce55d2e94b8ca3447c9ddb2854b4.
//
// Solidity: event NthBlockChanged(uint64 indexed chainId, uint64 indexed nthBlock)
func (_Contract *ContractFilterer) ParseNthBlockChanged(log types.Log) (*ContractNthBlockChanged, error) {
	event := new(ContractNthBlockChanged)
	if err := _Contract.contract.UnpackLog(event, "NthBlockChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the Contract contract.
type ContractOperatorAddedIterator struct {
	Event *ContractOperatorAdded // Event containing the contract specifics and raw log

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
func (it *ContractOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOperatorAdded)
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
		it.Event = new(ContractOperatorAdded)
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
func (it *ContractOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOperatorAdded represents a OperatorAdded event raised by the Contract contract.
type ContractOperatorAdded struct {
	Operator    common.Address
	ValidatorId *big.Int
	Role        [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x797ca55fc7be0f65c71f10996f7a16f801094f8ae3811874afc5a39730772a42.
//
// Solidity: event OperatorAdded(address operator, uint128 validatorId, bytes32 role)
func (_Contract *ContractFilterer) FilterOperatorAdded(opts *bind.FilterOpts) (*ContractOperatorAddedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return &ContractOperatorAddedIterator{contract: _Contract.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x797ca55fc7be0f65c71f10996f7a16f801094f8ae3811874afc5a39730772a42.
//
// Solidity: event OperatorAdded(address operator, uint128 validatorId, bytes32 role)
func (_Contract *ContractFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *ContractOperatorAdded) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OperatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOperatorAdded)
				if err := _Contract.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ParseOperatorAdded is a log parse operation binding the contract event 0x797ca55fc7be0f65c71f10996f7a16f801094f8ae3811874afc5a39730772a42.
//
// Solidity: event OperatorAdded(address operator, uint128 validatorId, bytes32 role)
func (_Contract *ContractFilterer) ParseOperatorAdded(log types.Log) (*ContractOperatorAdded, error) {
	event := new(ContractOperatorAdded)
	if err := _Contract.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractOperatorRemovedIterator is returned from FilterOperatorRemoved and is used to iterate over the raw logs and unpacked data for OperatorRemoved events raised by the Contract contract.
type ContractOperatorRemovedIterator struct {
	Event *ContractOperatorRemoved // Event containing the contract specifics and raw log

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
func (it *ContractOperatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOperatorRemoved)
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
		it.Event = new(ContractOperatorRemoved)
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
func (it *ContractOperatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOperatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOperatorRemoved represents a OperatorRemoved event raised by the Contract contract.
type ContractOperatorRemoved struct {
	Operator            common.Address
	ValidatorId         *big.Int
	ActiveOperatorCount *big.Int
	Role                [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterOperatorRemoved is a free log retrieval operation binding the contract event 0xca6d116f31cc5d708ae73029e1d63d1be48afdfd99819a3405bfbaf229748d90.
//
// Solidity: event OperatorRemoved(address operator, uint128 validatorId, uint128 activeOperatorCount, bytes32 role)
func (_Contract *ContractFilterer) FilterOperatorRemoved(opts *bind.FilterOpts) (*ContractOperatorRemovedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OperatorRemoved")
	if err != nil {
		return nil, err
	}
	return &ContractOperatorRemovedIterator{contract: _Contract.contract, event: "OperatorRemoved", logs: logs, sub: sub}, nil
}

// WatchOperatorRemoved is a free log subscription operation binding the contract event 0xca6d116f31cc5d708ae73029e1d63d1be48afdfd99819a3405bfbaf229748d90.
//
// Solidity: event OperatorRemoved(address operator, uint128 validatorId, uint128 activeOperatorCount, bytes32 role)
func (_Contract *ContractFilterer) WatchOperatorRemoved(opts *bind.WatchOpts, sink chan<- *ContractOperatorRemoved) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OperatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOperatorRemoved)
				if err := _Contract.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
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

// ParseOperatorRemoved is a log parse operation binding the contract event 0xca6d116f31cc5d708ae73029e1d63d1be48afdfd99819a3405bfbaf229748d90.
//
// Solidity: event OperatorRemoved(address operator, uint128 validatorId, uint128 activeOperatorCount, bytes32 role)
func (_Contract *ContractFilterer) ParseOperatorRemoved(log types.Log) (*ContractOperatorRemoved, error) {
	event := new(ContractOperatorRemoved)
	if err := _Contract.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Contract contract.
type ContractOwnershipTransferredIterator struct {
	Event *ContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnershipTransferred)
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
		it.Event = new(ContractOwnershipTransferred)
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
func (it *ContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnershipTransferred represents a OwnershipTransferred event raised by the Contract contract.
type ContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnershipTransferredIterator{contract: _Contract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnershipTransferred)
				if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) ParseOwnershipTransferred(log types.Log) (*ContractOwnershipTransferred, error) {
	event := new(ContractOwnershipTransferred)
	if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractQuorumNotReachedIterator is returned from FilterQuorumNotReached and is used to iterate over the raw logs and unpacked data for QuorumNotReached events raised by the Contract contract.
type ContractQuorumNotReachedIterator struct {
	Event *ContractQuorumNotReached // Event containing the contract specifics and raw log

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
func (it *ContractQuorumNotReachedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractQuorumNotReached)
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
		it.Event = new(ContractQuorumNotReached)
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
func (it *ContractQuorumNotReachedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractQuorumNotReachedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractQuorumNotReached represents a QuorumNotReached event raised by the Contract contract.
type ContractQuorumNotReached struct {
	ChainId     uint64
	BlockHeight uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterQuorumNotReached is a free log retrieval operation binding the contract event 0x398fd8f638a7242217f011fd0720a06747f7a85b7d28d7276684b841baea4021.
//
// Solidity: event QuorumNotReached(uint64 indexed chainId, uint64 blockHeight)
func (_Contract *ContractFilterer) FilterQuorumNotReached(opts *bind.FilterOpts, chainId []uint64) (*ContractQuorumNotReachedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "QuorumNotReached", chainIdRule)
	if err != nil {
		return nil, err
	}
	return &ContractQuorumNotReachedIterator{contract: _Contract.contract, event: "QuorumNotReached", logs: logs, sub: sub}, nil
}

// WatchQuorumNotReached is a free log subscription operation binding the contract event 0x398fd8f638a7242217f011fd0720a06747f7a85b7d28d7276684b841baea4021.
//
// Solidity: event QuorumNotReached(uint64 indexed chainId, uint64 blockHeight)
func (_Contract *ContractFilterer) WatchQuorumNotReached(opts *bind.WatchOpts, sink chan<- *ContractQuorumNotReached, chainId []uint64) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "QuorumNotReached", chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractQuorumNotReached)
				if err := _Contract.contract.UnpackLog(event, "QuorumNotReached", log); err != nil {
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

// ParseQuorumNotReached is a log parse operation binding the contract event 0x398fd8f638a7242217f011fd0720a06747f7a85b7d28d7276684b841baea4021.
//
// Solidity: event QuorumNotReached(uint64 indexed chainId, uint64 blockHeight)
func (_Contract *ContractFilterer) ParseQuorumNotReached(log types.Log) (*ContractQuorumNotReached, error) {
	event := new(ContractQuorumNotReached)
	if err := _Contract.contract.UnpackLog(event, "QuorumNotReached", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractResultSessionDurationChangedIterator is returned from FilterResultSessionDurationChanged and is used to iterate over the raw logs and unpacked data for ResultSessionDurationChanged events raised by the Contract contract.
type ContractResultSessionDurationChangedIterator struct {
	Event *ContractResultSessionDurationChanged // Event containing the contract specifics and raw log

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
func (it *ContractResultSessionDurationChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractResultSessionDurationChanged)
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
		it.Event = new(ContractResultSessionDurationChanged)
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
func (it *ContractResultSessionDurationChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractResultSessionDurationChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractResultSessionDurationChanged represents a ResultSessionDurationChanged event raised by the Contract contract.
type ContractResultSessionDurationChanged struct {
	NewSessionDuration uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterResultSessionDurationChanged is a free log retrieval operation binding the contract event 0x80972f2be3a50171fc4fb48963f365cbd47dc216ace7f3628a584503a80a9f97.
//
// Solidity: event ResultSessionDurationChanged(uint64 newSessionDuration)
func (_Contract *ContractFilterer) FilterResultSessionDurationChanged(opts *bind.FilterOpts) (*ContractResultSessionDurationChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ResultSessionDurationChanged")
	if err != nil {
		return nil, err
	}
	return &ContractResultSessionDurationChangedIterator{contract: _Contract.contract, event: "ResultSessionDurationChanged", logs: logs, sub: sub}, nil
}

// WatchResultSessionDurationChanged is a free log subscription operation binding the contract event 0x80972f2be3a50171fc4fb48963f365cbd47dc216ace7f3628a584503a80a9f97.
//
// Solidity: event ResultSessionDurationChanged(uint64 newSessionDuration)
func (_Contract *ContractFilterer) WatchResultSessionDurationChanged(opts *bind.WatchOpts, sink chan<- *ContractResultSessionDurationChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ResultSessionDurationChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractResultSessionDurationChanged)
				if err := _Contract.contract.UnpackLog(event, "ResultSessionDurationChanged", log); err != nil {
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

// ParseResultSessionDurationChanged is a log parse operation binding the contract event 0x80972f2be3a50171fc4fb48963f365cbd47dc216ace7f3628a584503a80a9f97.
//
// Solidity: event ResultSessionDurationChanged(uint64 newSessionDuration)
func (_Contract *ContractFilterer) ParseResultSessionDurationChanged(log types.Log) (*ContractResultSessionDurationChanged, error) {
	event := new(ContractResultSessionDurationChanged)
	if err := _Contract.contract.UnpackLog(event, "ResultSessionDurationChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractResultSessionMinSubmissionChangedIterator is returned from FilterResultSessionMinSubmissionChanged and is used to iterate over the raw logs and unpacked data for ResultSessionMinSubmissionChanged events raised by the Contract contract.
type ContractResultSessionMinSubmissionChangedIterator struct {
	Event *ContractResultSessionMinSubmissionChanged // Event containing the contract specifics and raw log

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
func (it *ContractResultSessionMinSubmissionChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractResultSessionMinSubmissionChanged)
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
		it.Event = new(ContractResultSessionMinSubmissionChanged)
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
func (it *ContractResultSessionMinSubmissionChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractResultSessionMinSubmissionChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractResultSessionMinSubmissionChanged represents a ResultSessionMinSubmissionChanged event raised by the Contract contract.
type ContractResultSessionMinSubmissionChanged struct {
	MinSubmissions uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterResultSessionMinSubmissionChanged is a free log retrieval operation binding the contract event 0xebd426f3d605bdf7e97ba7cd4a971371661140ad0acc8e0081f067d2004a7176.
//
// Solidity: event ResultSessionMinSubmissionChanged(uint64 minSubmissions)
func (_Contract *ContractFilterer) FilterResultSessionMinSubmissionChanged(opts *bind.FilterOpts) (*ContractResultSessionMinSubmissionChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ResultSessionMinSubmissionChanged")
	if err != nil {
		return nil, err
	}
	return &ContractResultSessionMinSubmissionChangedIterator{contract: _Contract.contract, event: "ResultSessionMinSubmissionChanged", logs: logs, sub: sub}, nil
}

// WatchResultSessionMinSubmissionChanged is a free log subscription operation binding the contract event 0xebd426f3d605bdf7e97ba7cd4a971371661140ad0acc8e0081f067d2004a7176.
//
// Solidity: event ResultSessionMinSubmissionChanged(uint64 minSubmissions)
func (_Contract *ContractFilterer) WatchResultSessionMinSubmissionChanged(opts *bind.WatchOpts, sink chan<- *ContractResultSessionMinSubmissionChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ResultSessionMinSubmissionChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractResultSessionMinSubmissionChanged)
				if err := _Contract.contract.UnpackLog(event, "ResultSessionMinSubmissionChanged", log); err != nil {
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

// ParseResultSessionMinSubmissionChanged is a log parse operation binding the contract event 0xebd426f3d605bdf7e97ba7cd4a971371661140ad0acc8e0081f067d2004a7176.
//
// Solidity: event ResultSessionMinSubmissionChanged(uint64 minSubmissions)
func (_Contract *ContractFilterer) ParseResultSessionMinSubmissionChanged(log types.Log) (*ContractResultSessionMinSubmissionChanged, error) {
	event := new(ContractResultSessionMinSubmissionChanged)
	if err := _Contract.contract.UnpackLog(event, "ResultSessionMinSubmissionChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractResultSessionQuorumChangedIterator is returned from FilterResultSessionQuorumChanged and is used to iterate over the raw logs and unpacked data for ResultSessionQuorumChanged events raised by the Contract contract.
type ContractResultSessionQuorumChangedIterator struct {
	Event *ContractResultSessionQuorumChanged // Event containing the contract specifics and raw log

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
func (it *ContractResultSessionQuorumChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractResultSessionQuorumChanged)
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
		it.Event = new(ContractResultSessionQuorumChanged)
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
func (it *ContractResultSessionQuorumChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractResultSessionQuorumChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractResultSessionQuorumChanged represents a ResultSessionQuorumChanged event raised by the Contract contract.
type ContractResultSessionQuorumChanged struct {
	NewQuorumThreshold *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterResultSessionQuorumChanged is a free log retrieval operation binding the contract event 0xe8c66e2621de650a92131e007d8bbc4cbf3bb8d4df7471d1f93eb20d70039a7c.
//
// Solidity: event ResultSessionQuorumChanged(uint256 newQuorumThreshold)
func (_Contract *ContractFilterer) FilterResultSessionQuorumChanged(opts *bind.FilterOpts) (*ContractResultSessionQuorumChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ResultSessionQuorumChanged")
	if err != nil {
		return nil, err
	}
	return &ContractResultSessionQuorumChangedIterator{contract: _Contract.contract, event: "ResultSessionQuorumChanged", logs: logs, sub: sub}, nil
}

// WatchResultSessionQuorumChanged is a free log subscription operation binding the contract event 0xe8c66e2621de650a92131e007d8bbc4cbf3bb8d4df7471d1f93eb20d70039a7c.
//
// Solidity: event ResultSessionQuorumChanged(uint256 newQuorumThreshold)
func (_Contract *ContractFilterer) WatchResultSessionQuorumChanged(opts *bind.WatchOpts, sink chan<- *ContractResultSessionQuorumChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ResultSessionQuorumChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractResultSessionQuorumChanged)
				if err := _Contract.contract.UnpackLog(event, "ResultSessionQuorumChanged", log); err != nil {
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

// ParseResultSessionQuorumChanged is a log parse operation binding the contract event 0xe8c66e2621de650a92131e007d8bbc4cbf3bb8d4df7471d1f93eb20d70039a7c.
//
// Solidity: event ResultSessionQuorumChanged(uint256 newQuorumThreshold)
func (_Contract *ContractFilterer) ParseResultSessionQuorumChanged(log types.Log) (*ContractResultSessionQuorumChanged, error) {
	event := new(ContractResultSessionQuorumChanged)
	if err := _Contract.contract.UnpackLog(event, "ResultSessionQuorumChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSecondsPerBlockCurrentChainChangedIterator is returned from FilterSecondsPerBlockCurrentChainChanged and is used to iterate over the raw logs and unpacked data for SecondsPerBlockCurrentChainChanged events raised by the Contract contract.
type ContractSecondsPerBlockCurrentChainChangedIterator struct {
	Event *ContractSecondsPerBlockCurrentChainChanged // Event containing the contract specifics and raw log

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
func (it *ContractSecondsPerBlockCurrentChainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSecondsPerBlockCurrentChainChanged)
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
		it.Event = new(ContractSecondsPerBlockCurrentChainChanged)
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
func (it *ContractSecondsPerBlockCurrentChainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSecondsPerBlockCurrentChainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSecondsPerBlockCurrentChainChanged represents a SecondsPerBlockCurrentChainChanged event raised by the Contract contract.
type ContractSecondsPerBlockCurrentChainChanged struct {
	SecondsPerBlockCurrentChain uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSecondsPerBlockCurrentChainChanged is a free log retrieval operation binding the contract event 0x52eb144349cf62d6190a9e1cbb6a601848aa63df834dd2a2e75bb0be3fef86f4.
//
// Solidity: event SecondsPerBlockCurrentChainChanged(uint64 indexed secondsPerBlockCurrentChain)
func (_Contract *ContractFilterer) FilterSecondsPerBlockCurrentChainChanged(opts *bind.FilterOpts, secondsPerBlockCurrentChain []uint64) (*ContractSecondsPerBlockCurrentChainChangedIterator, error) {

	var secondsPerBlockCurrentChainRule []interface{}
	for _, secondsPerBlockCurrentChainItem := range secondsPerBlockCurrentChain {
		secondsPerBlockCurrentChainRule = append(secondsPerBlockCurrentChainRule, secondsPerBlockCurrentChainItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "SecondsPerBlockCurrentChainChanged", secondsPerBlockCurrentChainRule)
	if err != nil {
		return nil, err
	}
	return &ContractSecondsPerBlockCurrentChainChangedIterator{contract: _Contract.contract, event: "SecondsPerBlockCurrentChainChanged", logs: logs, sub: sub}, nil
}

// WatchSecondsPerBlockCurrentChainChanged is a free log subscription operation binding the contract event 0x52eb144349cf62d6190a9e1cbb6a601848aa63df834dd2a2e75bb0be3fef86f4.
//
// Solidity: event SecondsPerBlockCurrentChainChanged(uint64 indexed secondsPerBlockCurrentChain)
func (_Contract *ContractFilterer) WatchSecondsPerBlockCurrentChainChanged(opts *bind.WatchOpts, sink chan<- *ContractSecondsPerBlockCurrentChainChanged, secondsPerBlockCurrentChain []uint64) (event.Subscription, error) {

	var secondsPerBlockCurrentChainRule []interface{}
	for _, secondsPerBlockCurrentChainItem := range secondsPerBlockCurrentChain {
		secondsPerBlockCurrentChainRule = append(secondsPerBlockCurrentChainRule, secondsPerBlockCurrentChainItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "SecondsPerBlockCurrentChainChanged", secondsPerBlockCurrentChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSecondsPerBlockCurrentChainChanged)
				if err := _Contract.contract.UnpackLog(event, "SecondsPerBlockCurrentChainChanged", log); err != nil {
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

// ParseSecondsPerBlockCurrentChainChanged is a log parse operation binding the contract event 0x52eb144349cf62d6190a9e1cbb6a601848aa63df834dd2a2e75bb0be3fef86f4.
//
// Solidity: event SecondsPerBlockCurrentChainChanged(uint64 indexed secondsPerBlockCurrentChain)
func (_Contract *ContractFilterer) ParseSecondsPerBlockCurrentChainChanged(log types.Log) (*ContractSecondsPerBlockCurrentChainChanged, error) {
	event := new(ContractSecondsPerBlockCurrentChainChanged)
	if err := _Contract.contract.UnpackLog(event, "SecondsPerBlockCurrentChainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractSessionStartedIterator is returned from FilterSessionStarted and is used to iterate over the raw logs and unpacked data for SessionStarted events raised by the Contract contract.
type ContractSessionStartedIterator struct {
	Event *ContractSessionStarted // Event containing the contract specifics and raw log

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
func (it *ContractSessionStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractSessionStarted)
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
		it.Event = new(ContractSessionStarted)
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
func (it *ContractSessionStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractSessionStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractSessionStarted represents a SessionStarted event raised by the Contract contract.
type ContractSessionStarted struct {
	ChainId     uint64
	BlockHeight uint64
	Deadline    uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSessionStarted is a free log retrieval operation binding the contract event 0x8b1f889addbfa41db5227bae3b091bd5c8b9a9122f874dfe54ba2f75aabe1f4c.
//
// Solidity: event SessionStarted(uint64 indexed chainId, uint64 indexed blockHeight, uint64 deadline)
func (_Contract *ContractFilterer) FilterSessionStarted(opts *bind.FilterOpts, chainId []uint64, blockHeight []uint64) (*ContractSessionStartedIterator, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var blockHeightRule []interface{}
	for _, blockHeightItem := range blockHeight {
		blockHeightRule = append(blockHeightRule, blockHeightItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "SessionStarted", chainIdRule, blockHeightRule)
	if err != nil {
		return nil, err
	}
	return &ContractSessionStartedIterator{contract: _Contract.contract, event: "SessionStarted", logs: logs, sub: sub}, nil
}

// WatchSessionStarted is a free log subscription operation binding the contract event 0x8b1f889addbfa41db5227bae3b091bd5c8b9a9122f874dfe54ba2f75aabe1f4c.
//
// Solidity: event SessionStarted(uint64 indexed chainId, uint64 indexed blockHeight, uint64 deadline)
func (_Contract *ContractFilterer) WatchSessionStarted(opts *bind.WatchOpts, sink chan<- *ContractSessionStarted, chainId []uint64, blockHeight []uint64) (event.Subscription, error) {

	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}
	var blockHeightRule []interface{}
	for _, blockHeightItem := range blockHeight {
		blockHeightRule = append(blockHeightRule, blockHeightItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "SessionStarted", chainIdRule, blockHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractSessionStarted)
				if err := _Contract.contract.UnpackLog(event, "SessionStarted", log); err != nil {
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

// ParseSessionStarted is a log parse operation binding the contract event 0x8b1f889addbfa41db5227bae3b091bd5c8b9a9122f874dfe54ba2f75aabe1f4c.
//
// Solidity: event SessionStarted(uint64 indexed chainId, uint64 indexed blockHeight, uint64 deadline)
func (_Contract *ContractFilterer) ParseSessionStarted(log types.Log) (*ContractSessionStarted, error) {
	event := new(ContractSessionStarted)
	if err := _Contract.contract.UnpackLog(event, "SessionStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractStakingManagerChangedIterator is returned from FilterStakingManagerChanged and is used to iterate over the raw logs and unpacked data for StakingManagerChanged events raised by the Contract contract.
type ContractStakingManagerChangedIterator struct {
	Event *ContractStakingManagerChanged // Event containing the contract specifics and raw log

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
func (it *ContractStakingManagerChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractStakingManagerChanged)
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
		it.Event = new(ContractStakingManagerChanged)
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
func (it *ContractStakingManagerChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractStakingManagerChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractStakingManagerChanged represents a StakingManagerChanged event raised by the Contract contract.
type ContractStakingManagerChanged struct {
	NewStakingManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStakingManagerChanged is a free log retrieval operation binding the contract event 0xf725afeae606c3f3c4c0ac3963e5c76a046bc4f386be98100c54e55bf5aeab36.
//
// Solidity: event StakingManagerChanged(address newStakingManager)
func (_Contract *ContractFilterer) FilterStakingManagerChanged(opts *bind.FilterOpts) (*ContractStakingManagerChangedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "StakingManagerChanged")
	if err != nil {
		return nil, err
	}
	return &ContractStakingManagerChangedIterator{contract: _Contract.contract, event: "StakingManagerChanged", logs: logs, sub: sub}, nil
}

// WatchStakingManagerChanged is a free log subscription operation binding the contract event 0xf725afeae606c3f3c4c0ac3963e5c76a046bc4f386be98100c54e55bf5aeab36.
//
// Solidity: event StakingManagerChanged(address newStakingManager)
func (_Contract *ContractFilterer) WatchStakingManagerChanged(opts *bind.WatchOpts, sink chan<- *ContractStakingManagerChanged) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "StakingManagerChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractStakingManagerChanged)
				if err := _Contract.contract.UnpackLog(event, "StakingManagerChanged", log); err != nil {
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

// ParseStakingManagerChanged is a log parse operation binding the contract event 0xf725afeae606c3f3c4c0ac3963e5c76a046bc4f386be98100c54e55bf5aeab36.
//
// Solidity: event StakingManagerChanged(address newStakingManager)
func (_Contract *ContractFilterer) ParseStakingManagerChanged(log types.Log) (*ContractStakingManagerChanged, error) {
	event := new(ContractStakingManagerChanged)
	if err := _Contract.contract.UnpackLog(event, "StakingManagerChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractValidatorDisabledIterator is returned from FilterValidatorDisabled and is used to iterate over the raw logs and unpacked data for ValidatorDisabled events raised by the Contract contract.
type ContractValidatorDisabledIterator struct {
	Event *ContractValidatorDisabled // Event containing the contract specifics and raw log

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
func (it *ContractValidatorDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractValidatorDisabled)
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
		it.Event = new(ContractValidatorDisabled)
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
func (it *ContractValidatorDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractValidatorDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractValidatorDisabled represents a ValidatorDisabled event raised by the Contract contract.
type ContractValidatorDisabled struct {
	ValidatorId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterValidatorDisabled is a free log retrieval operation binding the contract event 0xf97fbe9a37a2eae093c44e9bbcd9afde23ba64215e8fd80b9c49b1fbd4d58e54.
//
// Solidity: event ValidatorDisabled(uint128 validatorId)
func (_Contract *ContractFilterer) FilterValidatorDisabled(opts *bind.FilterOpts) (*ContractValidatorDisabledIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ValidatorDisabled")
	if err != nil {
		return nil, err
	}
	return &ContractValidatorDisabledIterator{contract: _Contract.contract, event: "ValidatorDisabled", logs: logs, sub: sub}, nil
}

// WatchValidatorDisabled is a free log subscription operation binding the contract event 0xf97fbe9a37a2eae093c44e9bbcd9afde23ba64215e8fd80b9c49b1fbd4d58e54.
//
// Solidity: event ValidatorDisabled(uint128 validatorId)
func (_Contract *ContractFilterer) WatchValidatorDisabled(opts *bind.WatchOpts, sink chan<- *ContractValidatorDisabled) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ValidatorDisabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractValidatorDisabled)
				if err := _Contract.contract.UnpackLog(event, "ValidatorDisabled", log); err != nil {
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

// ParseValidatorDisabled is a log parse operation binding the contract event 0xf97fbe9a37a2eae093c44e9bbcd9afde23ba64215e8fd80b9c49b1fbd4d58e54.
//
// Solidity: event ValidatorDisabled(uint128 validatorId)
func (_Contract *ContractFilterer) ParseValidatorDisabled(log types.Log) (*ContractValidatorDisabled, error) {
	event := new(ContractValidatorDisabled)
	if err := _Contract.contract.UnpackLog(event, "ValidatorDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractValidatorEnabledIterator is returned from FilterValidatorEnabled and is used to iterate over the raw logs and unpacked data for ValidatorEnabled events raised by the Contract contract.
type ContractValidatorEnabledIterator struct {
	Event *ContractValidatorEnabled // Event containing the contract specifics and raw log

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
func (it *ContractValidatorEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractValidatorEnabled)
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
		it.Event = new(ContractValidatorEnabled)
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
func (it *ContractValidatorEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractValidatorEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractValidatorEnabled represents a ValidatorEnabled event raised by the Contract contract.
type ContractValidatorEnabled struct {
	ValidatorId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterValidatorEnabled is a free log retrieval operation binding the contract event 0x553b029ba5c74688a5da732136d246f722502db24ed6b4aaf1cdc9f2f9ef23ef.
//
// Solidity: event ValidatorEnabled(uint128 validatorId)
func (_Contract *ContractFilterer) FilterValidatorEnabled(opts *bind.FilterOpts) (*ContractValidatorEnabledIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ValidatorEnabled")
	if err != nil {
		return nil, err
	}
	return &ContractValidatorEnabledIterator{contract: _Contract.contract, event: "ValidatorEnabled", logs: logs, sub: sub}, nil
}

// WatchValidatorEnabled is a free log subscription operation binding the contract event 0x553b029ba5c74688a5da732136d246f722502db24ed6b4aaf1cdc9f2f9ef23ef.
//
// Solidity: event ValidatorEnabled(uint128 validatorId)
func (_Contract *ContractFilterer) WatchValidatorEnabled(opts *bind.WatchOpts, sink chan<- *ContractValidatorEnabled) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ValidatorEnabled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractValidatorEnabled)
				if err := _Contract.contract.UnpackLog(event, "ValidatorEnabled", log); err != nil {
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

// ParseValidatorEnabled is a log parse operation binding the contract event 0x553b029ba5c74688a5da732136d246f722502db24ed6b4aaf1cdc9f2f9ef23ef.
//
// Solidity: event ValidatorEnabled(uint128 validatorId)
func (_Contract *ContractFilterer) ParseValidatorEnabled(log types.Log) (*ContractValidatorEnabled, error) {
	event := new(ContractValidatorEnabled)
	if err := _Contract.contract.UnpackLog(event, "ValidatorEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
