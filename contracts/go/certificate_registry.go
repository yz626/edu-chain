// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package certificate

import (
	"math/big"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/v3/abi"
	"github.com/FISCO-BCOS/go-sdk/v3/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/v3/types"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
)

// CertificateRegistryABI is the input ABI used to generate the binding from.
const CertificateRegistryABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ownerName\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"CertAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"CertAlreadyRevoked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"CertNotFound\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"CertNotRevoked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ArrayLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ContractNotPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ContractPaused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyParam\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"IssuerNotAuthorized\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"NotCertIssuer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Unauthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"IssuerAuthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"IssuerRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"CertificateIssued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"CertificateRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"by\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"CertificateRestored\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"BatchIssued\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"addIssuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"issuers\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"names\",\"type\":\"string[]\"}],\"name\":\"addIssuerBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"revokeIssuer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"}],\"name\":\"getIssuerInfo\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"authorizedAt\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"}],\"name\":\"issueCertificate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"certIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"certHashes\",\"type\":\"bytes32[]\"}],\"name\":\"issueCertificateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"revokeCertificate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"restoreCertificate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"certExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"}],\"name\":\"getCertificate\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"issuer\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"issuedAt\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"revoked\",\"type\":\"bool\"},{\"internalType\":\"uint64\",\"name\":\"revokedAt\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"revokeReason\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"certId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"}],\"name\":\"verifyCertificate\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"valid\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"revoked\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"certIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"certHashes\",\"type\":\"bytes32[]\"}],\"name\":\"verifyCertificateBatch\",\"outputs\":[{\"internalType\":\"bool[]\",\"name\":\"valids\",\"type\":\"bool[]\"},{\"internalType\":\"bool[]\",\"name\":\"revokeds\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStats\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalIssued\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalRevoked\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// CertificateRegistry is an auto generated Go binding around a Solidity contract.
type CertificateRegistry struct {
	CertificateRegistryCaller     // Read-only binding to the contract
	CertificateRegistryTransactor // Write-only binding to the contract
	CertificateRegistryFilterer   // Log filterer for contract events
}

// CertificateRegistryCaller is an auto generated read-only Go binding around a Solidity contract.
type CertificateRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificateRegistryTransactor is an auto generated write-only Go binding around a Solidity contract.
type CertificateRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificateRegistryFilterer is an auto generated log filtering Go binding around a Solidity contract events.
type CertificateRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CertificateRegistrySession is an auto generated Go binding around a Solidity contract,
// with pre-set call and transact options.
type CertificateRegistrySession struct {
	Contract     *CertificateRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CertificateRegistryCallerSession is an auto generated read-only Go binding around a Solidity contract,
// with pre-set call options.
type CertificateRegistryCallerSession struct {
	Contract *CertificateRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CertificateRegistryTransactorSession is an auto generated write-only Go binding around a Solidity contract,
// with pre-set transact options.
type CertificateRegistryTransactorSession struct {
	Contract     *CertificateRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CertificateRegistryRaw is an auto generated low-level Go binding around a Solidity contract.
type CertificateRegistryRaw struct {
	Contract *CertificateRegistry // Generic contract binding to access the raw methods on
}

// CertificateRegistryCallerRaw is an auto generated low-level read-only Go binding around a Solidity contract.
type CertificateRegistryCallerRaw struct {
	Contract *CertificateRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// CertificateRegistryTransactorRaw is an auto generated low-level write-only Go binding around a Solidity contract.
type CertificateRegistryTransactorRaw struct {
	Contract *CertificateRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCertificateRegistry creates a new instance of CertificateRegistry, bound to a specific deployed contract.
func NewCertificateRegistry(address common.Address, backend bind.ContractBackend) (*CertificateRegistry, error) {
	contract, err := bindCertificateRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CertificateRegistry{CertificateRegistryCaller: CertificateRegistryCaller{contract: contract}, CertificateRegistryTransactor: CertificateRegistryTransactor{contract: contract}, CertificateRegistryFilterer: CertificateRegistryFilterer{contract: contract}}, nil
}

// NewCertificateRegistryCaller creates a new read-only instance of CertificateRegistry, bound to a specific deployed contract.
func NewCertificateRegistryCaller(address common.Address, caller bind.ContractCaller) (*CertificateRegistryCaller, error) {
	contract, err := bindCertificateRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CertificateRegistryCaller{contract: contract}, nil
}

// NewCertificateRegistryTransactor creates a new write-only instance of CertificateRegistry, bound to a specific deployed contract.
func NewCertificateRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*CertificateRegistryTransactor, error) {
	contract, err := bindCertificateRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CertificateRegistryTransactor{contract: contract}, nil
}

// NewCertificateRegistryFilterer creates a new log filterer instance of CertificateRegistry, bound to a specific deployed contract.
func NewCertificateRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*CertificateRegistryFilterer, error) {
	contract, err := bindCertificateRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CertificateRegistryFilterer{contract: contract}, nil
}

// bindCertificateRegistry binds a generic wrapper to an already deployed contract.
func bindCertificateRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CertificateRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CertificateRegistry *CertificateRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CertificateRegistry.Contract.CertificateRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CertificateRegistry *CertificateRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.CertificateRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CertificateRegistry *CertificateRegistryRaw) TransactWithResult(opts *bind.TransactOpts, result interface{}, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.CertificateRegistryTransactor.contract.TransactWithResult(opts, result, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CertificateRegistry *CertificateRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CertificateRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CertificateRegistry *CertificateRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CertificateRegistry *CertificateRegistryTransactorRaw) TransactWithResult(opts *bind.TransactOpts, result interface{}, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.contract.TransactWithResult(opts, result, method, params...)
}

// CertExists is a free data retrieval call binding the contract method 0x7e3cb07d.
//
// Solidity: function certExists(bytes32 certId) constant returns(bool)
func (_CertificateRegistry *CertificateRegistryCaller) CertExists(opts *bind.CallOpts, certId [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CertificateRegistry.contract.Call(opts, out, "certExists", certId)
	return *ret0, err
}

// CertExists is a free data retrieval call binding the contract method 0x7e3cb07d.
//
// Solidity: function certExists(bytes32 certId) constant returns(bool)
func (_CertificateRegistry *CertificateRegistrySession) CertExists(certId [32]byte) (bool, error) {
	return _CertificateRegistry.Contract.CertExists(&_CertificateRegistry.CallOpts, certId)
}

// CertExists is a free data retrieval call binding the contract method 0x7e3cb07d.
//
// Solidity: function certExists(bytes32 certId) constant returns(bool)
func (_CertificateRegistry *CertificateRegistryCallerSession) CertExists(certId [32]byte) (bool, error) {
	return _CertificateRegistry.Contract.CertExists(&_CertificateRegistry.CallOpts, certId)
}

// GetCertificate is a free data retrieval call binding the contract method 0xf333fe08.
//
// Solidity: function getCertificate(bytes32 certId) constant returns(bytes32 certHash, address issuer, uint64 issuedAt, bool revoked, uint64 revokedAt, string revokeReason)
func (_CertificateRegistry *CertificateRegistryCaller) GetCertificate(opts *bind.CallOpts, certId [32]byte) (struct {
	CertHash     [32]byte
	Issuer       common.Address
	IssuedAt     uint64
	Revoked      bool
	RevokedAt    uint64
	RevokeReason string
}, error) {
	ret := new(struct {
		CertHash     [32]byte
		Issuer       common.Address
		IssuedAt     uint64
		Revoked      bool
		RevokedAt    uint64
		RevokeReason string
	})
	out := ret
	err := _CertificateRegistry.contract.Call(opts, out, "getCertificate", certId)
	return *ret, err
}

// GetCertificate is a free data retrieval call binding the contract method 0xf333fe08.
//
// Solidity: function getCertificate(bytes32 certId) constant returns(bytes32 certHash, address issuer, uint64 issuedAt, bool revoked, uint64 revokedAt, string revokeReason)
func (_CertificateRegistry *CertificateRegistrySession) GetCertificate(certId [32]byte) (struct {
	CertHash     [32]byte
	Issuer       common.Address
	IssuedAt     uint64
	Revoked      bool
	RevokedAt    uint64
	RevokeReason string
}, error) {
	return _CertificateRegistry.Contract.GetCertificate(&_CertificateRegistry.CallOpts, certId)
}

// GetCertificate is a free data retrieval call binding the contract method 0xf333fe08.
//
// Solidity: function getCertificate(bytes32 certId) constant returns(bytes32 certHash, address issuer, uint64 issuedAt, bool revoked, uint64 revokedAt, string revokeReason)
func (_CertificateRegistry *CertificateRegistryCallerSession) GetCertificate(certId [32]byte) (struct {
	CertHash     [32]byte
	Issuer       common.Address
	IssuedAt     uint64
	Revoked      bool
	RevokedAt    uint64
	RevokeReason string
}, error) {
	return _CertificateRegistry.Contract.GetCertificate(&_CertificateRegistry.CallOpts, certId)
}

// GetIssuerInfo is a free data retrieval call binding the contract method 0x5e9aab70.
//
// Solidity: function getIssuerInfo(address issuer) constant returns(bool authorized, string name, uint64 authorizedAt)
func (_CertificateRegistry *CertificateRegistryCaller) GetIssuerInfo(opts *bind.CallOpts, issuer common.Address) (struct {
	Authorized   bool
	Name         string
	AuthorizedAt uint64
}, error) {
	ret := new(struct {
		Authorized   bool
		Name         string
		AuthorizedAt uint64
	})
	out := ret
	err := _CertificateRegistry.contract.Call(opts, out, "getIssuerInfo", issuer)
	return *ret, err
}

// GetIssuerInfo is a free data retrieval call binding the contract method 0x5e9aab70.
//
// Solidity: function getIssuerInfo(address issuer) constant returns(bool authorized, string name, uint64 authorizedAt)
func (_CertificateRegistry *CertificateRegistrySession) GetIssuerInfo(issuer common.Address) (struct {
	Authorized   bool
	Name         string
	AuthorizedAt uint64
}, error) {
	return _CertificateRegistry.Contract.GetIssuerInfo(&_CertificateRegistry.CallOpts, issuer)
}

// GetIssuerInfo is a free data retrieval call binding the contract method 0x5e9aab70.
//
// Solidity: function getIssuerInfo(address issuer) constant returns(bool authorized, string name, uint64 authorizedAt)
func (_CertificateRegistry *CertificateRegistryCallerSession) GetIssuerInfo(issuer common.Address) (struct {
	Authorized   bool
	Name         string
	AuthorizedAt uint64
}, error) {
	return _CertificateRegistry.Contract.GetIssuerInfo(&_CertificateRegistry.CallOpts, issuer)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() constant returns(uint256 totalIssued, uint256 totalRevoked)
func (_CertificateRegistry *CertificateRegistryCaller) GetStats(opts *bind.CallOpts) (struct {
	TotalIssued  *big.Int
	TotalRevoked *big.Int
}, error) {
	ret := new(struct {
		TotalIssued  *big.Int
		TotalRevoked *big.Int
	})
	out := ret
	err := _CertificateRegistry.contract.Call(opts, out, "getStats")
	return *ret, err
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() constant returns(uint256 totalIssued, uint256 totalRevoked)
func (_CertificateRegistry *CertificateRegistrySession) GetStats() (struct {
	TotalIssued  *big.Int
	TotalRevoked *big.Int
}, error) {
	return _CertificateRegistry.Contract.GetStats(&_CertificateRegistry.CallOpts)
}

// GetStats is a free data retrieval call binding the contract method 0xc59d4847.
//
// Solidity: function getStats() constant returns(uint256 totalIssued, uint256 totalRevoked)
func (_CertificateRegistry *CertificateRegistryCallerSession) GetStats() (struct {
	TotalIssued  *big.Int
	TotalRevoked *big.Int
}, error) {
	return _CertificateRegistry.Contract.GetStats(&_CertificateRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_CertificateRegistry *CertificateRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CertificateRegistry.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_CertificateRegistry *CertificateRegistrySession) Owner() (common.Address, error) {
	return _CertificateRegistry.Contract.Owner(&_CertificateRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_CertificateRegistry *CertificateRegistryCallerSession) Owner() (common.Address, error) {
	return _CertificateRegistry.Contract.Owner(&_CertificateRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_CertificateRegistry *CertificateRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CertificateRegistry.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_CertificateRegistry *CertificateRegistrySession) Paused() (bool, error) {
	return _CertificateRegistry.Contract.Paused(&_CertificateRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_CertificateRegistry *CertificateRegistryCallerSession) Paused() (bool, error) {
	return _CertificateRegistry.Contract.Paused(&_CertificateRegistry.CallOpts)
}

// VerifyCertificate is a free data retrieval call binding the contract method 0x1f75435d.
//
// Solidity: function verifyCertificate(bytes32 certId, bytes32 certHash) constant returns(bool valid, bool revoked)
func (_CertificateRegistry *CertificateRegistryCaller) VerifyCertificate(opts *bind.CallOpts, certId [32]byte, certHash [32]byte) (struct {
	Valid   bool
	Revoked bool
}, error) {
	ret := new(struct {
		Valid   bool
		Revoked bool
	})
	out := ret
	err := _CertificateRegistry.contract.Call(opts, out, "verifyCertificate", certId, certHash)
	return *ret, err
}

// VerifyCertificate is a free data retrieval call binding the contract method 0x1f75435d.
//
// Solidity: function verifyCertificate(bytes32 certId, bytes32 certHash) constant returns(bool valid, bool revoked)
func (_CertificateRegistry *CertificateRegistrySession) VerifyCertificate(certId [32]byte, certHash [32]byte) (struct {
	Valid   bool
	Revoked bool
}, error) {
	return _CertificateRegistry.Contract.VerifyCertificate(&_CertificateRegistry.CallOpts, certId, certHash)
}

// VerifyCertificate is a free data retrieval call binding the contract method 0x1f75435d.
//
// Solidity: function verifyCertificate(bytes32 certId, bytes32 certHash) constant returns(bool valid, bool revoked)
func (_CertificateRegistry *CertificateRegistryCallerSession) VerifyCertificate(certId [32]byte, certHash [32]byte) (struct {
	Valid   bool
	Revoked bool
}, error) {
	return _CertificateRegistry.Contract.VerifyCertificate(&_CertificateRegistry.CallOpts, certId, certHash)
}

// VerifyCertificateBatch is a free data retrieval call binding the contract method 0xba3ba01e.
//
// Solidity: function verifyCertificateBatch(bytes32[] certIds, bytes32[] certHashes) constant returns(bool[] valids, bool[] revokeds)
func (_CertificateRegistry *CertificateRegistryCaller) VerifyCertificateBatch(opts *bind.CallOpts, certIds [][32]byte, certHashes [][32]byte) (struct {
	Valids   []bool
	Revokeds []bool
}, error) {
	ret := new(struct {
		Valids   []bool
		Revokeds []bool
	})
	out := ret
	err := _CertificateRegistry.contract.Call(opts, out, "verifyCertificateBatch", certIds, certHashes)
	return *ret, err
}

// VerifyCertificateBatch is a free data retrieval call binding the contract method 0xba3ba01e.
//
// Solidity: function verifyCertificateBatch(bytes32[] certIds, bytes32[] certHashes) constant returns(bool[] valids, bool[] revokeds)
func (_CertificateRegistry *CertificateRegistrySession) VerifyCertificateBatch(certIds [][32]byte, certHashes [][32]byte) (struct {
	Valids   []bool
	Revokeds []bool
}, error) {
	return _CertificateRegistry.Contract.VerifyCertificateBatch(&_CertificateRegistry.CallOpts, certIds, certHashes)
}

// VerifyCertificateBatch is a free data retrieval call binding the contract method 0xba3ba01e.
//
// Solidity: function verifyCertificateBatch(bytes32[] certIds, bytes32[] certHashes) constant returns(bool[] valids, bool[] revokeds)
func (_CertificateRegistry *CertificateRegistryCallerSession) VerifyCertificateBatch(certIds [][32]byte, certHashes [][32]byte) (struct {
	Valids   []bool
	Revokeds []bool
}, error) {
	return _CertificateRegistry.Contract.VerifyCertificateBatch(&_CertificateRegistry.CallOpts, certIds, certHashes)
}

// AddIssuer is a paid mutator transaction binding the contract method 0x8fc0859a.
//
// Solidity: function addIssuer(address issuer, string name) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) AddIssuer(opts *bind.TransactOpts, issuer common.Address, name string) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "addIssuer", issuer, name)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncAddIssuer(handler func(*types.Receipt, error), opts *bind.TransactOpts, issuer common.Address, name string) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "addIssuer", issuer, name)
}

// AddIssuer is a paid mutator transaction binding the contract method 0x8fc0859a.
//
// Solidity: function addIssuer(address issuer, string name) returns()
func (_CertificateRegistry *CertificateRegistrySession) AddIssuer(issuer common.Address, name string) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.AddIssuer(&_CertificateRegistry.TransactOpts, issuer, name)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncAddIssuer(handler func(*types.Receipt, error), issuer common.Address, name string) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncAddIssuer(handler, &_CertificateRegistry.TransactOpts, issuer, name)
}

// AddIssuer is a paid mutator transaction binding the contract method 0x8fc0859a.
//
// Solidity: function addIssuer(address issuer, string name) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) AddIssuer(issuer common.Address, name string) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.AddIssuer(&_CertificateRegistry.TransactOpts, issuer, name)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncAddIssuer(handler func(*types.Receipt, error), issuer common.Address, name string) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncAddIssuer(handler, &_CertificateRegistry.TransactOpts, issuer, name)
}

// AddIssuerBatch is a paid mutator transaction binding the contract method 0xb5edd9c0.
//
// Solidity: function addIssuerBatch(address[] issuers, string[] names) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) AddIssuerBatch(opts *bind.TransactOpts, issuers []common.Address, names []string) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "addIssuerBatch", issuers, names)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncAddIssuerBatch(handler func(*types.Receipt, error), opts *bind.TransactOpts, issuers []common.Address, names []string) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "addIssuerBatch", issuers, names)
}

// AddIssuerBatch is a paid mutator transaction binding the contract method 0xb5edd9c0.
//
// Solidity: function addIssuerBatch(address[] issuers, string[] names) returns()
func (_CertificateRegistry *CertificateRegistrySession) AddIssuerBatch(issuers []common.Address, names []string) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.AddIssuerBatch(&_CertificateRegistry.TransactOpts, issuers, names)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncAddIssuerBatch(handler func(*types.Receipt, error), issuers []common.Address, names []string) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncAddIssuerBatch(handler, &_CertificateRegistry.TransactOpts, issuers, names)
}

// AddIssuerBatch is a paid mutator transaction binding the contract method 0xb5edd9c0.
//
// Solidity: function addIssuerBatch(address[] issuers, string[] names) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) AddIssuerBatch(issuers []common.Address, names []string) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.AddIssuerBatch(&_CertificateRegistry.TransactOpts, issuers, names)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncAddIssuerBatch(handler func(*types.Receipt, error), issuers []common.Address, names []string) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncAddIssuerBatch(handler, &_CertificateRegistry.TransactOpts, issuers, names)
}

// IssueCertificate is a paid mutator transaction binding the contract method 0x57bcf883.
//
// Solidity: function issueCertificate(bytes32 certId, bytes32 certHash) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) IssueCertificate(opts *bind.TransactOpts, certId [32]byte, certHash [32]byte) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "issueCertificate", certId, certHash)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncIssueCertificate(handler func(*types.Receipt, error), opts *bind.TransactOpts, certId [32]byte, certHash [32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "issueCertificate", certId, certHash)
}

// IssueCertificate is a paid mutator transaction binding the contract method 0x57bcf883.
//
// Solidity: function issueCertificate(bytes32 certId, bytes32 certHash) returns()
func (_CertificateRegistry *CertificateRegistrySession) IssueCertificate(certId [32]byte, certHash [32]byte) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.IssueCertificate(&_CertificateRegistry.TransactOpts, certId, certHash)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncIssueCertificate(handler func(*types.Receipt, error), certId [32]byte, certHash [32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncIssueCertificate(handler, &_CertificateRegistry.TransactOpts, certId, certHash)
}

// IssueCertificate is a paid mutator transaction binding the contract method 0x57bcf883.
//
// Solidity: function issueCertificate(bytes32 certId, bytes32 certHash) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) IssueCertificate(certId [32]byte, certHash [32]byte) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.IssueCertificate(&_CertificateRegistry.TransactOpts, certId, certHash)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncIssueCertificate(handler func(*types.Receipt, error), certId [32]byte, certHash [32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncIssueCertificate(handler, &_CertificateRegistry.TransactOpts, certId, certHash)
}

// IssueCertificateBatch is a paid mutator transaction binding the contract method 0x2d7ac4a5.
//
// Solidity: function issueCertificateBatch(bytes32[] certIds, bytes32[] certHashes) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) IssueCertificateBatch(opts *bind.TransactOpts, certIds [][32]byte, certHashes [][32]byte) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "issueCertificateBatch", certIds, certHashes)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncIssueCertificateBatch(handler func(*types.Receipt, error), opts *bind.TransactOpts, certIds [][32]byte, certHashes [][32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "issueCertificateBatch", certIds, certHashes)
}

// IssueCertificateBatch is a paid mutator transaction binding the contract method 0x2d7ac4a5.
//
// Solidity: function issueCertificateBatch(bytes32[] certIds, bytes32[] certHashes) returns()
func (_CertificateRegistry *CertificateRegistrySession) IssueCertificateBatch(certIds [][32]byte, certHashes [][32]byte) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.IssueCertificateBatch(&_CertificateRegistry.TransactOpts, certIds, certHashes)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncIssueCertificateBatch(handler func(*types.Receipt, error), certIds [][32]byte, certHashes [][32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncIssueCertificateBatch(handler, &_CertificateRegistry.TransactOpts, certIds, certHashes)
}

// IssueCertificateBatch is a paid mutator transaction binding the contract method 0x2d7ac4a5.
//
// Solidity: function issueCertificateBatch(bytes32[] certIds, bytes32[] certHashes) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) IssueCertificateBatch(certIds [][32]byte, certHashes [][32]byte) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.IssueCertificateBatch(&_CertificateRegistry.TransactOpts, certIds, certHashes)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncIssueCertificateBatch(handler func(*types.Receipt, error), certIds [][32]byte, certHashes [][32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncIssueCertificateBatch(handler, &_CertificateRegistry.TransactOpts, certIds, certHashes)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CertificateRegistry *CertificateRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "pause")
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncPause(handler func(*types.Receipt, error), opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CertificateRegistry *CertificateRegistrySession) Pause() (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.Pause(&_CertificateRegistry.TransactOpts)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncPause(handler func(*types.Receipt, error)) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncPause(handler, &_CertificateRegistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) Pause() (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.Pause(&_CertificateRegistry.TransactOpts)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncPause(handler func(*types.Receipt, error)) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncPause(handler, &_CertificateRegistry.TransactOpts)
}

// RestoreCertificate is a paid mutator transaction binding the contract method 0xf3c5e074.
//
// Solidity: function restoreCertificate(bytes32 certId) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) RestoreCertificate(opts *bind.TransactOpts, certId [32]byte) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "restoreCertificate", certId)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncRestoreCertificate(handler func(*types.Receipt, error), opts *bind.TransactOpts, certId [32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "restoreCertificate", certId)
}

// RestoreCertificate is a paid mutator transaction binding the contract method 0xf3c5e074.
//
// Solidity: function restoreCertificate(bytes32 certId) returns()
func (_CertificateRegistry *CertificateRegistrySession) RestoreCertificate(certId [32]byte) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.RestoreCertificate(&_CertificateRegistry.TransactOpts, certId)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncRestoreCertificate(handler func(*types.Receipt, error), certId [32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncRestoreCertificate(handler, &_CertificateRegistry.TransactOpts, certId)
}

// RestoreCertificate is a paid mutator transaction binding the contract method 0xf3c5e074.
//
// Solidity: function restoreCertificate(bytes32 certId) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) RestoreCertificate(certId [32]byte) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.RestoreCertificate(&_CertificateRegistry.TransactOpts, certId)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncRestoreCertificate(handler func(*types.Receipt, error), certId [32]byte) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncRestoreCertificate(handler, &_CertificateRegistry.TransactOpts, certId)
}

// RevokeCertificate is a paid mutator transaction binding the contract method 0x5bf4adfb.
//
// Solidity: function revokeCertificate(bytes32 certId, string reason) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) RevokeCertificate(opts *bind.TransactOpts, certId [32]byte, reason string) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "revokeCertificate", certId, reason)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncRevokeCertificate(handler func(*types.Receipt, error), opts *bind.TransactOpts, certId [32]byte, reason string) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "revokeCertificate", certId, reason)
}

// RevokeCertificate is a paid mutator transaction binding the contract method 0x5bf4adfb.
//
// Solidity: function revokeCertificate(bytes32 certId, string reason) returns()
func (_CertificateRegistry *CertificateRegistrySession) RevokeCertificate(certId [32]byte, reason string) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.RevokeCertificate(&_CertificateRegistry.TransactOpts, certId, reason)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncRevokeCertificate(handler func(*types.Receipt, error), certId [32]byte, reason string) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncRevokeCertificate(handler, &_CertificateRegistry.TransactOpts, certId, reason)
}

// RevokeCertificate is a paid mutator transaction binding the contract method 0x5bf4adfb.
//
// Solidity: function revokeCertificate(bytes32 certId, string reason) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) RevokeCertificate(certId [32]byte, reason string) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.RevokeCertificate(&_CertificateRegistry.TransactOpts, certId, reason)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncRevokeCertificate(handler func(*types.Receipt, error), certId [32]byte, reason string) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncRevokeCertificate(handler, &_CertificateRegistry.TransactOpts, certId, reason)
}

// RevokeIssuer is a paid mutator transaction binding the contract method 0x00629679.
//
// Solidity: function revokeIssuer(address issuer) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) RevokeIssuer(opts *bind.TransactOpts, issuer common.Address) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "revokeIssuer", issuer)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncRevokeIssuer(handler func(*types.Receipt, error), opts *bind.TransactOpts, issuer common.Address) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "revokeIssuer", issuer)
}

// RevokeIssuer is a paid mutator transaction binding the contract method 0x00629679.
//
// Solidity: function revokeIssuer(address issuer) returns()
func (_CertificateRegistry *CertificateRegistrySession) RevokeIssuer(issuer common.Address) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.RevokeIssuer(&_CertificateRegistry.TransactOpts, issuer)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncRevokeIssuer(handler func(*types.Receipt, error), issuer common.Address) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncRevokeIssuer(handler, &_CertificateRegistry.TransactOpts, issuer)
}

// RevokeIssuer is a paid mutator transaction binding the contract method 0x00629679.
//
// Solidity: function revokeIssuer(address issuer) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) RevokeIssuer(issuer common.Address) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.RevokeIssuer(&_CertificateRegistry.TransactOpts, issuer)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncRevokeIssuer(handler func(*types.Receipt, error), issuer common.Address) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncRevokeIssuer(handler, &_CertificateRegistry.TransactOpts, issuer)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CertificateRegistry *CertificateRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "transferOwnership", newOwner)
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncTransferOwnership(handler func(*types.Receipt, error), opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CertificateRegistry *CertificateRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.TransferOwnership(&_CertificateRegistry.TransactOpts, newOwner)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncTransferOwnership(handler func(*types.Receipt, error), newOwner common.Address) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncTransferOwnership(handler, &_CertificateRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.TransferOwnership(&_CertificateRegistry.TransactOpts, newOwner)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncTransferOwnership(handler func(*types.Receipt, error), newOwner common.Address) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncTransferOwnership(handler, &_CertificateRegistry.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CertificateRegistry *CertificateRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	var ()
	out := &[]interface{}{}
	transaction, receipt, err := _CertificateRegistry.contract.TransactWithResult(opts, out, "unpause")
	return transaction, receipt, err
}

func (_CertificateRegistry *CertificateRegistryTransactor) AsyncUnpause(handler func(*types.Receipt, error), opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CertificateRegistry.contract.AsyncTransact(opts, handler, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CertificateRegistry *CertificateRegistrySession) Unpause() (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.Unpause(&_CertificateRegistry.TransactOpts)
}

func (_CertificateRegistry *CertificateRegistrySession) AsyncUnpause(handler func(*types.Receipt, error)) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncUnpause(handler, &_CertificateRegistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_CertificateRegistry *CertificateRegistryTransactorSession) Unpause() (*types.Transaction, *types.Receipt, error) {
	return _CertificateRegistry.Contract.Unpause(&_CertificateRegistry.TransactOpts)
}

func (_CertificateRegistry *CertificateRegistryTransactorSession) AsyncUnpause(handler func(*types.Receipt, error)) (*types.Transaction, error) {
	return _CertificateRegistry.Contract.AsyncUnpause(handler, &_CertificateRegistry.TransactOpts)
}

// CertificateRegistryBatchIssued represents a BatchIssued event raised by the CertificateRegistry contract.
type CertificateRegistryBatchIssued struct {
	Issuer    common.Address
	Count     *big.Int
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// WatchBatchIssued is a free log subscription operation binding the contract event 0x86ac3e7644108a3f0d73bc6e19cede6a415b457c502ff1dd24145762605761d1.
//
// Solidity: event BatchIssued(address indexed issuer, uint256 count, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchBatchIssued(fromBlock *int64, handler func(int, []types.Log), issuer common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "BatchIssued", issuer)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllBatchIssued(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "BatchIssued")
}

// ParseBatchIssued is a log parse operation binding the contract event 0x86ac3e7644108a3f0d73bc6e19cede6a415b457c502ff1dd24145762605761d1.
//
// Solidity: event BatchIssued(address indexed issuer, uint256 count, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseBatchIssued(log types.Log) (*CertificateRegistryBatchIssued, error) {
	event := new(CertificateRegistryBatchIssued)
	if err := _CertificateRegistry.contract.UnpackLog(event, "BatchIssued", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchBatchIssued is a free log subscription operation binding the contract event 0x86ac3e7644108a3f0d73bc6e19cede6a415b457c502ff1dd24145762605761d1.
//
// Solidity: event BatchIssued(address indexed issuer, uint256 count, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) WatchBatchIssued(fromBlock *int64, handler func(int, []types.Log), issuer common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchBatchIssued(fromBlock, handler, issuer)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllBatchIssued(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllBatchIssued(fromBlock, handler)
}

// ParseBatchIssued is a log parse operation binding the contract event 0x86ac3e7644108a3f0d73bc6e19cede6a415b457c502ff1dd24145762605761d1.
//
// Solidity: event BatchIssued(address indexed issuer, uint256 count, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) ParseBatchIssued(log types.Log) (*CertificateRegistryBatchIssued, error) {
	return _CertificateRegistry.Contract.ParseBatchIssued(log)
}

// CertificateRegistryCertificateIssued represents a CertificateIssued event raised by the CertificateRegistry contract.
type CertificateRegistryCertificateIssued struct {
	CertId    [32]byte
	CertHash  [32]byte
	Issuer    common.Address
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// WatchCertificateIssued is a free log subscription operation binding the contract event 0x0d7d2d60696ebd106065ef5888a39d0920cc07474cf7321481d5f3a9d5c89d16.
//
// Solidity: event CertificateIssued(bytes32 indexed certId, bytes32 indexed certHash, address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchCertificateIssued(fromBlock *int64, handler func(int, []types.Log), certId [32]byte, certHash [32]byte, issuer common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "CertificateIssued", certId, certHash, issuer)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllCertificateIssued(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "CertificateIssued")
}

// ParseCertificateIssued is a log parse operation binding the contract event 0x0d7d2d60696ebd106065ef5888a39d0920cc07474cf7321481d5f3a9d5c89d16.
//
// Solidity: event CertificateIssued(bytes32 indexed certId, bytes32 indexed certHash, address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseCertificateIssued(log types.Log) (*CertificateRegistryCertificateIssued, error) {
	event := new(CertificateRegistryCertificateIssued)
	if err := _CertificateRegistry.contract.UnpackLog(event, "CertificateIssued", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchCertificateIssued is a free log subscription operation binding the contract event 0x0d7d2d60696ebd106065ef5888a39d0920cc07474cf7321481d5f3a9d5c89d16.
//
// Solidity: event CertificateIssued(bytes32 indexed certId, bytes32 indexed certHash, address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) WatchCertificateIssued(fromBlock *int64, handler func(int, []types.Log), certId [32]byte, certHash [32]byte, issuer common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchCertificateIssued(fromBlock, handler, certId, certHash, issuer)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllCertificateIssued(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllCertificateIssued(fromBlock, handler)
}

// ParseCertificateIssued is a log parse operation binding the contract event 0x0d7d2d60696ebd106065ef5888a39d0920cc07474cf7321481d5f3a9d5c89d16.
//
// Solidity: event CertificateIssued(bytes32 indexed certId, bytes32 indexed certHash, address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) ParseCertificateIssued(log types.Log) (*CertificateRegistryCertificateIssued, error) {
	return _CertificateRegistry.Contract.ParseCertificateIssued(log)
}

// CertificateRegistryCertificateRestored represents a CertificateRestored event raised by the CertificateRegistry contract.
type CertificateRegistryCertificateRestored struct {
	CertId    [32]byte
	By        common.Address
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// WatchCertificateRestored is a free log subscription operation binding the contract event 0xb2fa76a791dd2d45f9d10c03a5711d8e2ff2bc219ec1fed97170d9aaf3d2003e.
//
// Solidity: event CertificateRestored(bytes32 indexed certId, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchCertificateRestored(fromBlock *int64, handler func(int, []types.Log), certId [32]byte, by common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "CertificateRestored", certId, by)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllCertificateRestored(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "CertificateRestored")
}

// ParseCertificateRestored is a log parse operation binding the contract event 0xb2fa76a791dd2d45f9d10c03a5711d8e2ff2bc219ec1fed97170d9aaf3d2003e.
//
// Solidity: event CertificateRestored(bytes32 indexed certId, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseCertificateRestored(log types.Log) (*CertificateRegistryCertificateRestored, error) {
	event := new(CertificateRegistryCertificateRestored)
	if err := _CertificateRegistry.contract.UnpackLog(event, "CertificateRestored", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchCertificateRestored is a free log subscription operation binding the contract event 0xb2fa76a791dd2d45f9d10c03a5711d8e2ff2bc219ec1fed97170d9aaf3d2003e.
//
// Solidity: event CertificateRestored(bytes32 indexed certId, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) WatchCertificateRestored(fromBlock *int64, handler func(int, []types.Log), certId [32]byte, by common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchCertificateRestored(fromBlock, handler, certId, by)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllCertificateRestored(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllCertificateRestored(fromBlock, handler)
}

// ParseCertificateRestored is a log parse operation binding the contract event 0xb2fa76a791dd2d45f9d10c03a5711d8e2ff2bc219ec1fed97170d9aaf3d2003e.
//
// Solidity: event CertificateRestored(bytes32 indexed certId, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) ParseCertificateRestored(log types.Log) (*CertificateRegistryCertificateRestored, error) {
	return _CertificateRegistry.Contract.ParseCertificateRestored(log)
}

// CertificateRegistryCertificateRevoked represents a CertificateRevoked event raised by the CertificateRegistry contract.
type CertificateRegistryCertificateRevoked struct {
	CertId    [32]byte
	Reason    string
	By        common.Address
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// WatchCertificateRevoked is a free log subscription operation binding the contract event 0x35855721d51b96b861fa066aad5e6cf076f0e1fef225502edc83eb65873ca6da.
//
// Solidity: event CertificateRevoked(bytes32 indexed certId, string reason, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchCertificateRevoked(fromBlock *int64, handler func(int, []types.Log), certId [32]byte, by common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "CertificateRevoked", certId, by)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllCertificateRevoked(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "CertificateRevoked")
}

// ParseCertificateRevoked is a log parse operation binding the contract event 0x35855721d51b96b861fa066aad5e6cf076f0e1fef225502edc83eb65873ca6da.
//
// Solidity: event CertificateRevoked(bytes32 indexed certId, string reason, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseCertificateRevoked(log types.Log) (*CertificateRegistryCertificateRevoked, error) {
	event := new(CertificateRegistryCertificateRevoked)
	if err := _CertificateRegistry.contract.UnpackLog(event, "CertificateRevoked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchCertificateRevoked is a free log subscription operation binding the contract event 0x35855721d51b96b861fa066aad5e6cf076f0e1fef225502edc83eb65873ca6da.
//
// Solidity: event CertificateRevoked(bytes32 indexed certId, string reason, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) WatchCertificateRevoked(fromBlock *int64, handler func(int, []types.Log), certId [32]byte, by common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchCertificateRevoked(fromBlock, handler, certId, by)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllCertificateRevoked(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllCertificateRevoked(fromBlock, handler)
}

// ParseCertificateRevoked is a log parse operation binding the contract event 0x35855721d51b96b861fa066aad5e6cf076f0e1fef225502edc83eb65873ca6da.
//
// Solidity: event CertificateRevoked(bytes32 indexed certId, string reason, address indexed by, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) ParseCertificateRevoked(log types.Log) (*CertificateRegistryCertificateRevoked, error) {
	return _CertificateRegistry.Contract.ParseCertificateRevoked(log)
}

// CertificateRegistryIssuerAuthorized represents a IssuerAuthorized event raised by the CertificateRegistry contract.
type CertificateRegistryIssuerAuthorized struct {
	Issuer    common.Address
	Name      string
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// WatchIssuerAuthorized is a free log subscription operation binding the contract event 0xd02cb7461b5b3281f7cdd043a7505986ac4e55952cffa4c7d9b6224b9da73522.
//
// Solidity: event IssuerAuthorized(address indexed issuer, string name, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchIssuerAuthorized(fromBlock *int64, handler func(int, []types.Log), issuer common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "IssuerAuthorized", issuer)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllIssuerAuthorized(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "IssuerAuthorized")
}

// ParseIssuerAuthorized is a log parse operation binding the contract event 0xd02cb7461b5b3281f7cdd043a7505986ac4e55952cffa4c7d9b6224b9da73522.
//
// Solidity: event IssuerAuthorized(address indexed issuer, string name, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseIssuerAuthorized(log types.Log) (*CertificateRegistryIssuerAuthorized, error) {
	event := new(CertificateRegistryIssuerAuthorized)
	if err := _CertificateRegistry.contract.UnpackLog(event, "IssuerAuthorized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchIssuerAuthorized is a free log subscription operation binding the contract event 0xd02cb7461b5b3281f7cdd043a7505986ac4e55952cffa4c7d9b6224b9da73522.
//
// Solidity: event IssuerAuthorized(address indexed issuer, string name, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) WatchIssuerAuthorized(fromBlock *int64, handler func(int, []types.Log), issuer common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchIssuerAuthorized(fromBlock, handler, issuer)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllIssuerAuthorized(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllIssuerAuthorized(fromBlock, handler)
}

// ParseIssuerAuthorized is a log parse operation binding the contract event 0xd02cb7461b5b3281f7cdd043a7505986ac4e55952cffa4c7d9b6224b9da73522.
//
// Solidity: event IssuerAuthorized(address indexed issuer, string name, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) ParseIssuerAuthorized(log types.Log) (*CertificateRegistryIssuerAuthorized, error) {
	return _CertificateRegistry.Contract.ParseIssuerAuthorized(log)
}

// CertificateRegistryIssuerRevoked represents a IssuerRevoked event raised by the CertificateRegistry contract.
type CertificateRegistryIssuerRevoked struct {
	Issuer    common.Address
	Timestamp uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// WatchIssuerRevoked is a free log subscription operation binding the contract event 0xcd83be3e28cace3428b652f82ed8327b1c83008fd6f8e906c77554cab875ba74.
//
// Solidity: event IssuerRevoked(address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchIssuerRevoked(fromBlock *int64, handler func(int, []types.Log), issuer common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "IssuerRevoked", issuer)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllIssuerRevoked(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "IssuerRevoked")
}

// ParseIssuerRevoked is a log parse operation binding the contract event 0xcd83be3e28cace3428b652f82ed8327b1c83008fd6f8e906c77554cab875ba74.
//
// Solidity: event IssuerRevoked(address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseIssuerRevoked(log types.Log) (*CertificateRegistryIssuerRevoked, error) {
	event := new(CertificateRegistryIssuerRevoked)
	if err := _CertificateRegistry.contract.UnpackLog(event, "IssuerRevoked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchIssuerRevoked is a free log subscription operation binding the contract event 0xcd83be3e28cace3428b652f82ed8327b1c83008fd6f8e906c77554cab875ba74.
//
// Solidity: event IssuerRevoked(address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) WatchIssuerRevoked(fromBlock *int64, handler func(int, []types.Log), issuer common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchIssuerRevoked(fromBlock, handler, issuer)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllIssuerRevoked(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllIssuerRevoked(fromBlock, handler)
}

// ParseIssuerRevoked is a log parse operation binding the contract event 0xcd83be3e28cace3428b652f82ed8327b1c83008fd6f8e906c77554cab875ba74.
//
// Solidity: event IssuerRevoked(address indexed issuer, uint64 timestamp)
func (_CertificateRegistry *CertificateRegistrySession) ParseIssuerRevoked(log types.Log) (*CertificateRegistryIssuerRevoked, error) {
	return _CertificateRegistry.Contract.ParseIssuerRevoked(log)
}

// CertificateRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the CertificateRegistry contract.
type CertificateRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchOwnershipTransferred(fromBlock *int64, handler func(int, []types.Log), previousOwner common.Address, newOwner common.Address) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "OwnershipTransferred", previousOwner, newOwner)
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllOwnershipTransferred(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "OwnershipTransferred")
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*CertificateRegistryOwnershipTransferred, error) {
	event := new(CertificateRegistryOwnershipTransferred)
	if err := _CertificateRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CertificateRegistry *CertificateRegistrySession) WatchOwnershipTransferred(fromBlock *int64, handler func(int, []types.Log), previousOwner common.Address, newOwner common.Address) (string, error) {
	return _CertificateRegistry.Contract.WatchOwnershipTransferred(fromBlock, handler, previousOwner, newOwner)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllOwnershipTransferred(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllOwnershipTransferred(fromBlock, handler)
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CertificateRegistry *CertificateRegistrySession) ParseOwnershipTransferred(log types.Log) (*CertificateRegistryOwnershipTransferred, error) {
	return _CertificateRegistry.Contract.ParseOwnershipTransferred(log)
}

// CertificateRegistryPaused represents a Paused event raised by the CertificateRegistry contract.
type CertificateRegistryPaused struct {
	By  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address by)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchPaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "Paused")
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllPaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "Paused")
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address by)
func (_CertificateRegistry *CertificateRegistryFilterer) ParsePaused(log types.Log) (*CertificateRegistryPaused, error) {
	event := new(CertificateRegistryPaused)
	if err := _CertificateRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address by)
func (_CertificateRegistry *CertificateRegistrySession) WatchPaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchPaused(fromBlock, handler)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllPaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllPaused(fromBlock, handler)
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address by)
func (_CertificateRegistry *CertificateRegistrySession) ParsePaused(log types.Log) (*CertificateRegistryPaused, error) {
	return _CertificateRegistry.Contract.ParsePaused(log)
}

// CertificateRegistryUnpaused represents a Unpaused event raised by the CertificateRegistry contract.
type CertificateRegistryUnpaused struct {
	By  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address by)
func (_CertificateRegistry *CertificateRegistryFilterer) WatchUnpaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "Unpaused")
}

func (_CertificateRegistry *CertificateRegistryFilterer) WatchAllUnpaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.contract.WatchLogs(fromBlock, handler, "Unpaused")
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address by)
func (_CertificateRegistry *CertificateRegistryFilterer) ParseUnpaused(log types.Log) (*CertificateRegistryUnpaused, error) {
	event := new(CertificateRegistryUnpaused)
	if err := _CertificateRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address by)
func (_CertificateRegistry *CertificateRegistrySession) WatchUnpaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchUnpaused(fromBlock, handler)
}

func (_CertificateRegistry *CertificateRegistrySession) WatchAllUnpaused(fromBlock *int64, handler func(int, []types.Log)) (string, error) {
	return _CertificateRegistry.Contract.WatchAllUnpaused(fromBlock, handler)
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address by)
func (_CertificateRegistry *CertificateRegistrySession) ParseUnpaused(log types.Log) (*CertificateRegistryUnpaused, error) {
	return _CertificateRegistry.Contract.ParseUnpaused(log)
}
