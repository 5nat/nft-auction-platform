// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// NFTAuctionMarketMetaData contains all meta data concerning the NFTAuctionMarket contract.
var NFTAuctionMarketMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"ETH_TOKEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"USD_PRECISION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"auctions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"seller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nft\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minBidUsd\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestBid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"highestBidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"highestBidUsd\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"ended\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"cancelled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bidERC20\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bidEth\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"cancelAuction\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createAuction\",\"inputs\":[{\"name\":\"nft\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minBidUsd\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"duration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endAuction\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBidUsdValue\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextAuctionId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceFeeds\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPriceFeed\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AuctionCancelled\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuctionCreated\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nft\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"minBidUsd\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AuctionEnded\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"winner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"bidToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amountUsd\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BidPlaced\",\"inputs\":[{\"name\":\"auctionId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"bidToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amountUsd\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceFeedSet\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"priceFeed\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// NFTAuctionMarketABI is the input ABI used to generate the binding from.
// Deprecated: Use NFTAuctionMarketMetaData.ABI instead.
var NFTAuctionMarketABI = NFTAuctionMarketMetaData.ABI

// NFTAuctionMarket is an auto generated Go binding around an Ethereum contract.
type NFTAuctionMarket struct {
	NFTAuctionMarketCaller     // Read-only binding to the contract
	NFTAuctionMarketTransactor // Write-only binding to the contract
	NFTAuctionMarketFilterer   // Log filterer for contract events
}

// NFTAuctionMarketCaller is an auto generated read-only Go binding around an Ethereum contract.
type NFTAuctionMarketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFTAuctionMarketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NFTAuctionMarketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFTAuctionMarketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NFTAuctionMarketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFTAuctionMarketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NFTAuctionMarketSession struct {
	Contract     *NFTAuctionMarket // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NFTAuctionMarketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NFTAuctionMarketCallerSession struct {
	Contract *NFTAuctionMarketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// NFTAuctionMarketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NFTAuctionMarketTransactorSession struct {
	Contract     *NFTAuctionMarketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// NFTAuctionMarketRaw is an auto generated low-level Go binding around an Ethereum contract.
type NFTAuctionMarketRaw struct {
	Contract *NFTAuctionMarket // Generic contract binding to access the raw methods on
}

// NFTAuctionMarketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NFTAuctionMarketCallerRaw struct {
	Contract *NFTAuctionMarketCaller // Generic read-only contract binding to access the raw methods on
}

// NFTAuctionMarketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NFTAuctionMarketTransactorRaw struct {
	Contract *NFTAuctionMarketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNFTAuctionMarket creates a new instance of NFTAuctionMarket, bound to a specific deployed contract.
func NewNFTAuctionMarket(address common.Address, backend bind.ContractBackend) (*NFTAuctionMarket, error) {
	contract, err := bindNFTAuctionMarket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarket{NFTAuctionMarketCaller: NFTAuctionMarketCaller{contract: contract}, NFTAuctionMarketTransactor: NFTAuctionMarketTransactor{contract: contract}, NFTAuctionMarketFilterer: NFTAuctionMarketFilterer{contract: contract}}, nil
}

// NewNFTAuctionMarketCaller creates a new read-only instance of NFTAuctionMarket, bound to a specific deployed contract.
func NewNFTAuctionMarketCaller(address common.Address, caller bind.ContractCaller) (*NFTAuctionMarketCaller, error) {
	contract, err := bindNFTAuctionMarket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketCaller{contract: contract}, nil
}

// NewNFTAuctionMarketTransactor creates a new write-only instance of NFTAuctionMarket, bound to a specific deployed contract.
func NewNFTAuctionMarketTransactor(address common.Address, transactor bind.ContractTransactor) (*NFTAuctionMarketTransactor, error) {
	contract, err := bindNFTAuctionMarket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketTransactor{contract: contract}, nil
}

// NewNFTAuctionMarketFilterer creates a new log filterer instance of NFTAuctionMarket, bound to a specific deployed contract.
func NewNFTAuctionMarketFilterer(address common.Address, filterer bind.ContractFilterer) (*NFTAuctionMarketFilterer, error) {
	contract, err := bindNFTAuctionMarket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketFilterer{contract: contract}, nil
}

// bindNFTAuctionMarket binds a generic wrapper to an already deployed contract.
func bindNFTAuctionMarket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NFTAuctionMarketMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFTAuctionMarket *NFTAuctionMarketRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFTAuctionMarket.Contract.NFTAuctionMarketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NFTAuctionMarket *NFTAuctionMarketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.NFTAuctionMarketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NFTAuctionMarket *NFTAuctionMarketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.NFTAuctionMarketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFTAuctionMarket *NFTAuctionMarketCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFTAuctionMarket.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NFTAuctionMarket *NFTAuctionMarketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NFTAuctionMarket *NFTAuctionMarketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.contract.Transact(opts, method, params...)
}

// ETHTOKEN is a free data retrieval call binding the contract method 0x58bc8337.
//
// Solidity: function ETH_TOKEN() view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) ETHTOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "ETH_TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ETHTOKEN is a free data retrieval call binding the contract method 0x58bc8337.
//
// Solidity: function ETH_TOKEN() view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketSession) ETHTOKEN() (common.Address, error) {
	return _NFTAuctionMarket.Contract.ETHTOKEN(&_NFTAuctionMarket.CallOpts)
}

// ETHTOKEN is a free data retrieval call binding the contract method 0x58bc8337.
//
// Solidity: function ETH_TOKEN() view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) ETHTOKEN() (common.Address, error) {
	return _NFTAuctionMarket.Contract.ETHTOKEN(&_NFTAuctionMarket.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_NFTAuctionMarket *NFTAuctionMarketSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _NFTAuctionMarket.Contract.UPGRADEINTERFACEVERSION(&_NFTAuctionMarket.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _NFTAuctionMarket.Contract.UPGRADEINTERFACEVERSION(&_NFTAuctionMarket.CallOpts)
}

// USDPRECISION is a free data retrieval call binding the contract method 0x4bfdfa6f.
//
// Solidity: function USD_PRECISION() view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) USDPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "USD_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// USDPRECISION is a free data retrieval call binding the contract method 0x4bfdfa6f.
//
// Solidity: function USD_PRECISION() view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketSession) USDPRECISION() (*big.Int, error) {
	return _NFTAuctionMarket.Contract.USDPRECISION(&_NFTAuctionMarket.CallOpts)
}

// USDPRECISION is a free data retrieval call binding the contract method 0x4bfdfa6f.
//
// Solidity: function USD_PRECISION() view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) USDPRECISION() (*big.Int, error) {
	return _NFTAuctionMarket.Contract.USDPRECISION(&_NFTAuctionMarket.CallOpts)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address seller, address nft, uint256 tokenId, uint256 minBidUsd, uint256 highestBid, address highestBidder, uint256 highestBidUsd, address bidToken, uint256 endTime, bool ended, bool cancelled)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) Auctions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Seller        common.Address
	Nft           common.Address
	TokenId       *big.Int
	MinBidUsd     *big.Int
	HighestBid    *big.Int
	HighestBidder common.Address
	HighestBidUsd *big.Int
	BidToken      common.Address
	EndTime       *big.Int
	Ended         bool
	Cancelled     bool
}, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "auctions", arg0)

	outstruct := new(struct {
		Seller        common.Address
		Nft           common.Address
		TokenId       *big.Int
		MinBidUsd     *big.Int
		HighestBid    *big.Int
		HighestBidder common.Address
		HighestBidUsd *big.Int
		BidToken      common.Address
		EndTime       *big.Int
		Ended         bool
		Cancelled     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Seller = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Nft = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MinBidUsd = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.HighestBid = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.HighestBidder = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.HighestBidUsd = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.BidToken = *abi.ConvertType(out[7], new(common.Address)).(*common.Address)
	outstruct.EndTime = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.Ended = *abi.ConvertType(out[9], new(bool)).(*bool)
	outstruct.Cancelled = *abi.ConvertType(out[10], new(bool)).(*bool)

	return *outstruct, err

}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address seller, address nft, uint256 tokenId, uint256 minBidUsd, uint256 highestBid, address highestBidder, uint256 highestBidUsd, address bidToken, uint256 endTime, bool ended, bool cancelled)
func (_NFTAuctionMarket *NFTAuctionMarketSession) Auctions(arg0 *big.Int) (struct {
	Seller        common.Address
	Nft           common.Address
	TokenId       *big.Int
	MinBidUsd     *big.Int
	HighestBid    *big.Int
	HighestBidder common.Address
	HighestBidUsd *big.Int
	BidToken      common.Address
	EndTime       *big.Int
	Ended         bool
	Cancelled     bool
}, error) {
	return _NFTAuctionMarket.Contract.Auctions(&_NFTAuctionMarket.CallOpts, arg0)
}

// Auctions is a free data retrieval call binding the contract method 0x571a26a0.
//
// Solidity: function auctions(uint256 ) view returns(address seller, address nft, uint256 tokenId, uint256 minBidUsd, uint256 highestBid, address highestBidder, uint256 highestBidUsd, address bidToken, uint256 endTime, bool ended, bool cancelled)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) Auctions(arg0 *big.Int) (struct {
	Seller        common.Address
	Nft           common.Address
	TokenId       *big.Int
	MinBidUsd     *big.Int
	HighestBid    *big.Int
	HighestBidder common.Address
	HighestBidUsd *big.Int
	BidToken      common.Address
	EndTime       *big.Int
	Ended         bool
	Cancelled     bool
}, error) {
	return _NFTAuctionMarket.Contract.Auctions(&_NFTAuctionMarket.CallOpts, arg0)
}

// GetBidUsdValue is a free data retrieval call binding the contract method 0xae5e1775.
//
// Solidity: function getBidUsdValue(address token, uint256 amount) view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) GetBidUsdValue(opts *bind.CallOpts, token common.Address, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "getBidUsdValue", token, amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBidUsdValue is a free data retrieval call binding the contract method 0xae5e1775.
//
// Solidity: function getBidUsdValue(address token, uint256 amount) view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketSession) GetBidUsdValue(token common.Address, amount *big.Int) (*big.Int, error) {
	return _NFTAuctionMarket.Contract.GetBidUsdValue(&_NFTAuctionMarket.CallOpts, token, amount)
}

// GetBidUsdValue is a free data retrieval call binding the contract method 0xae5e1775.
//
// Solidity: function getBidUsdValue(address token, uint256 amount) view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) GetBidUsdValue(token common.Address, amount *big.Int) (*big.Int, error) {
	return _NFTAuctionMarket.Contract.GetBidUsdValue(&_NFTAuctionMarket.CallOpts, token, amount)
}

// NextAuctionId is a free data retrieval call binding the contract method 0xfc528482.
//
// Solidity: function nextAuctionId() view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) NextAuctionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "nextAuctionId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextAuctionId is a free data retrieval call binding the contract method 0xfc528482.
//
// Solidity: function nextAuctionId() view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketSession) NextAuctionId() (*big.Int, error) {
	return _NFTAuctionMarket.Contract.NextAuctionId(&_NFTAuctionMarket.CallOpts)
}

// NextAuctionId is a free data retrieval call binding the contract method 0xfc528482.
//
// Solidity: function nextAuctionId() view returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) NextAuctionId() (*big.Int, error) {
	return _NFTAuctionMarket.Contract.NextAuctionId(&_NFTAuctionMarket.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketSession) Owner() (common.Address, error) {
	return _NFTAuctionMarket.Contract.Owner(&_NFTAuctionMarket.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) Owner() (common.Address, error) {
	return _NFTAuctionMarket.Contract.Owner(&_NFTAuctionMarket.CallOpts)
}

// PriceFeeds is a free data retrieval call binding the contract method 0x9dcb511a.
//
// Solidity: function priceFeeds(address ) view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) PriceFeeds(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "priceFeeds", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PriceFeeds is a free data retrieval call binding the contract method 0x9dcb511a.
//
// Solidity: function priceFeeds(address ) view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketSession) PriceFeeds(arg0 common.Address) (common.Address, error) {
	return _NFTAuctionMarket.Contract.PriceFeeds(&_NFTAuctionMarket.CallOpts, arg0)
}

// PriceFeeds is a free data retrieval call binding the contract method 0x9dcb511a.
//
// Solidity: function priceFeeds(address ) view returns(address)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) PriceFeeds(arg0 common.Address) (common.Address, error) {
	return _NFTAuctionMarket.Contract.PriceFeeds(&_NFTAuctionMarket.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_NFTAuctionMarket *NFTAuctionMarketCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _NFTAuctionMarket.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_NFTAuctionMarket *NFTAuctionMarketSession) ProxiableUUID() ([32]byte, error) {
	return _NFTAuctionMarket.Contract.ProxiableUUID(&_NFTAuctionMarket.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_NFTAuctionMarket *NFTAuctionMarketCallerSession) ProxiableUUID() ([32]byte, error) {
	return _NFTAuctionMarket.Contract.ProxiableUUID(&_NFTAuctionMarket.CallOpts)
}

// BidERC20 is a paid mutator transaction binding the contract method 0xbe68ea68.
//
// Solidity: function bidERC20(uint256 auctionId, address token, uint256 amount) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) BidERC20(opts *bind.TransactOpts, auctionId *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "bidERC20", auctionId, token, amount)
}

// BidERC20 is a paid mutator transaction binding the contract method 0xbe68ea68.
//
// Solidity: function bidERC20(uint256 auctionId, address token, uint256 amount) returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) BidERC20(auctionId *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.BidERC20(&_NFTAuctionMarket.TransactOpts, auctionId, token, amount)
}

// BidERC20 is a paid mutator transaction binding the contract method 0xbe68ea68.
//
// Solidity: function bidERC20(uint256 auctionId, address token, uint256 amount) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) BidERC20(auctionId *big.Int, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.BidERC20(&_NFTAuctionMarket.TransactOpts, auctionId, token, amount)
}

// BidEth is a paid mutator transaction binding the contract method 0x55c623c6.
//
// Solidity: function bidEth(uint256 auctionId) payable returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) BidEth(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "bidEth", auctionId)
}

// BidEth is a paid mutator transaction binding the contract method 0x55c623c6.
//
// Solidity: function bidEth(uint256 auctionId) payable returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) BidEth(auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.BidEth(&_NFTAuctionMarket.TransactOpts, auctionId)
}

// BidEth is a paid mutator transaction binding the contract method 0x55c623c6.
//
// Solidity: function bidEth(uint256 auctionId) payable returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) BidEth(auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.BidEth(&_NFTAuctionMarket.TransactOpts, auctionId)
}

// CancelAuction is a paid mutator transaction binding the contract method 0x96b5a755.
//
// Solidity: function cancelAuction(uint256 auctionId) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) CancelAuction(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "cancelAuction", auctionId)
}

// CancelAuction is a paid mutator transaction binding the contract method 0x96b5a755.
//
// Solidity: function cancelAuction(uint256 auctionId) returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) CancelAuction(auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.CancelAuction(&_NFTAuctionMarket.TransactOpts, auctionId)
}

// CancelAuction is a paid mutator transaction binding the contract method 0x96b5a755.
//
// Solidity: function cancelAuction(uint256 auctionId) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) CancelAuction(auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.CancelAuction(&_NFTAuctionMarket.TransactOpts, auctionId)
}

// CreateAuction is a paid mutator transaction binding the contract method 0x61beb1d7.
//
// Solidity: function createAuction(address nft, uint256 tokenId, uint256 minBidUsd, uint256 duration) returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) CreateAuction(opts *bind.TransactOpts, nft common.Address, tokenId *big.Int, minBidUsd *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "createAuction", nft, tokenId, minBidUsd, duration)
}

// CreateAuction is a paid mutator transaction binding the contract method 0x61beb1d7.
//
// Solidity: function createAuction(address nft, uint256 tokenId, uint256 minBidUsd, uint256 duration) returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketSession) CreateAuction(nft common.Address, tokenId *big.Int, minBidUsd *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.CreateAuction(&_NFTAuctionMarket.TransactOpts, nft, tokenId, minBidUsd, duration)
}

// CreateAuction is a paid mutator transaction binding the contract method 0x61beb1d7.
//
// Solidity: function createAuction(address nft, uint256 tokenId, uint256 minBidUsd, uint256 duration) returns(uint256)
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) CreateAuction(nft common.Address, tokenId *big.Int, minBidUsd *big.Int, duration *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.CreateAuction(&_NFTAuctionMarket.TransactOpts, nft, tokenId, minBidUsd, duration)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 auctionId) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) EndAuction(opts *bind.TransactOpts, auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "endAuction", auctionId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 auctionId) returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) EndAuction(auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.EndAuction(&_NFTAuctionMarket.TransactOpts, auctionId)
}

// EndAuction is a paid mutator transaction binding the contract method 0xb9a2de3a.
//
// Solidity: function endAuction(uint256 auctionId) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) EndAuction(auctionId *big.Int) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.EndAuction(&_NFTAuctionMarket.TransactOpts, auctionId)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address initialOwner) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "initialize", initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address initialOwner) returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) Initialize(initialOwner common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.Initialize(&_NFTAuctionMarket.TransactOpts, initialOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address initialOwner) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) Initialize(initialOwner common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.Initialize(&_NFTAuctionMarket.TransactOpts, initialOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) RenounceOwnership() (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.RenounceOwnership(&_NFTAuctionMarket.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.RenounceOwnership(&_NFTAuctionMarket.TransactOpts)
}

// SetPriceFeed is a paid mutator transaction binding the contract method 0x76e11286.
//
// Solidity: function setPriceFeed(address token, address priceFeed) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) SetPriceFeed(opts *bind.TransactOpts, token common.Address, priceFeed common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "setPriceFeed", token, priceFeed)
}

// SetPriceFeed is a paid mutator transaction binding the contract method 0x76e11286.
//
// Solidity: function setPriceFeed(address token, address priceFeed) returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) SetPriceFeed(token common.Address, priceFeed common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.SetPriceFeed(&_NFTAuctionMarket.TransactOpts, token, priceFeed)
}

// SetPriceFeed is a paid mutator transaction binding the contract method 0x76e11286.
//
// Solidity: function setPriceFeed(address token, address priceFeed) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) SetPriceFeed(token common.Address, priceFeed common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.SetPriceFeed(&_NFTAuctionMarket.TransactOpts, token, priceFeed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.TransferOwnership(&_NFTAuctionMarket.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.TransferOwnership(&_NFTAuctionMarket.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NFTAuctionMarket.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_NFTAuctionMarket *NFTAuctionMarketSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.UpgradeToAndCall(&_NFTAuctionMarket.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_NFTAuctionMarket *NFTAuctionMarketTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _NFTAuctionMarket.Contract.UpgradeToAndCall(&_NFTAuctionMarket.TransactOpts, newImplementation, data)
}

// NFTAuctionMarketAuctionCancelledIterator is returned from FilterAuctionCancelled and is used to iterate over the raw logs and unpacked data for AuctionCancelled events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketAuctionCancelledIterator struct {
	Event *NFTAuctionMarketAuctionCancelled // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketAuctionCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketAuctionCancelled)
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
		it.Event = new(NFTAuctionMarketAuctionCancelled)
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
func (it *NFTAuctionMarketAuctionCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketAuctionCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketAuctionCancelled represents a AuctionCancelled event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketAuctionCancelled struct {
	AuctionId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuctionCancelled is a free log retrieval operation binding the contract event 0x2809c7e17bf978fbc7194c0a694b638c4215e9140cacc6c38ca36010b45697df.
//
// Solidity: event AuctionCancelled(uint256 indexed auctionId)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterAuctionCancelled(opts *bind.FilterOpts, auctionId []*big.Int) (*NFTAuctionMarketAuctionCancelledIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "AuctionCancelled", auctionIdRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketAuctionCancelledIterator{contract: _NFTAuctionMarket.contract, event: "AuctionCancelled", logs: logs, sub: sub}, nil
}

// WatchAuctionCancelled is a free log subscription operation binding the contract event 0x2809c7e17bf978fbc7194c0a694b638c4215e9140cacc6c38ca36010b45697df.
//
// Solidity: event AuctionCancelled(uint256 indexed auctionId)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchAuctionCancelled(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketAuctionCancelled, auctionId []*big.Int) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "AuctionCancelled", auctionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketAuctionCancelled)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "AuctionCancelled", log); err != nil {
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

// ParseAuctionCancelled is a log parse operation binding the contract event 0x2809c7e17bf978fbc7194c0a694b638c4215e9140cacc6c38ca36010b45697df.
//
// Solidity: event AuctionCancelled(uint256 indexed auctionId)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseAuctionCancelled(log types.Log) (*NFTAuctionMarketAuctionCancelled, error) {
	event := new(NFTAuctionMarketAuctionCancelled)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "AuctionCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketAuctionCreatedIterator is returned from FilterAuctionCreated and is used to iterate over the raw logs and unpacked data for AuctionCreated events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketAuctionCreatedIterator struct {
	Event *NFTAuctionMarketAuctionCreated // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketAuctionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketAuctionCreated)
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
		it.Event = new(NFTAuctionMarketAuctionCreated)
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
func (it *NFTAuctionMarketAuctionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketAuctionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketAuctionCreated represents a AuctionCreated event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketAuctionCreated struct {
	AuctionId *big.Int
	Seller    common.Address
	Nft       common.Address
	TokenId   *big.Int
	MinBidUsd *big.Int
	EndTime   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuctionCreated is a free log retrieval operation binding the contract event 0x06b9e486c68303eb64052e0493f906f3d93a1b7149b6b8dcff221aebd16c3513.
//
// Solidity: event AuctionCreated(uint256 indexed auctionId, address indexed seller, address indexed nft, uint256 tokenId, uint256 minBidUsd, uint256 endTime)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterAuctionCreated(opts *bind.FilterOpts, auctionId []*big.Int, seller []common.Address, nft []common.Address) (*NFTAuctionMarketAuctionCreatedIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}
	var nftRule []interface{}
	for _, nftItem := range nft {
		nftRule = append(nftRule, nftItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "AuctionCreated", auctionIdRule, sellerRule, nftRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketAuctionCreatedIterator{contract: _NFTAuctionMarket.contract, event: "AuctionCreated", logs: logs, sub: sub}, nil
}

// WatchAuctionCreated is a free log subscription operation binding the contract event 0x06b9e486c68303eb64052e0493f906f3d93a1b7149b6b8dcff221aebd16c3513.
//
// Solidity: event AuctionCreated(uint256 indexed auctionId, address indexed seller, address indexed nft, uint256 tokenId, uint256 minBidUsd, uint256 endTime)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchAuctionCreated(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketAuctionCreated, auctionId []*big.Int, seller []common.Address, nft []common.Address) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}
	var nftRule []interface{}
	for _, nftItem := range nft {
		nftRule = append(nftRule, nftItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "AuctionCreated", auctionIdRule, sellerRule, nftRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketAuctionCreated)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "AuctionCreated", log); err != nil {
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

// ParseAuctionCreated is a log parse operation binding the contract event 0x06b9e486c68303eb64052e0493f906f3d93a1b7149b6b8dcff221aebd16c3513.
//
// Solidity: event AuctionCreated(uint256 indexed auctionId, address indexed seller, address indexed nft, uint256 tokenId, uint256 minBidUsd, uint256 endTime)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseAuctionCreated(log types.Log) (*NFTAuctionMarketAuctionCreated, error) {
	event := new(NFTAuctionMarketAuctionCreated)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "AuctionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketAuctionEndedIterator is returned from FilterAuctionEnded and is used to iterate over the raw logs and unpacked data for AuctionEnded events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketAuctionEndedIterator struct {
	Event *NFTAuctionMarketAuctionEnded // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketAuctionEndedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketAuctionEnded)
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
		it.Event = new(NFTAuctionMarketAuctionEnded)
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
func (it *NFTAuctionMarketAuctionEndedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketAuctionEndedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketAuctionEnded represents a AuctionEnded event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketAuctionEnded struct {
	AuctionId *big.Int
	Winner    common.Address
	BidToken  common.Address
	Amount    *big.Int
	AmountUsd *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAuctionEnded is a free log retrieval operation binding the contract event 0x596165d0521c3cb4157fad2621686f086daed4663acb3d03441a92b9277f5683.
//
// Solidity: event AuctionEnded(uint256 indexed auctionId, address indexed winner, address indexed bidToken, uint256 amount, uint256 amountUsd)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterAuctionEnded(opts *bind.FilterOpts, auctionId []*big.Int, winner []common.Address, bidToken []common.Address) (*NFTAuctionMarketAuctionEndedIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}
	var bidTokenRule []interface{}
	for _, bidTokenItem := range bidToken {
		bidTokenRule = append(bidTokenRule, bidTokenItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "AuctionEnded", auctionIdRule, winnerRule, bidTokenRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketAuctionEndedIterator{contract: _NFTAuctionMarket.contract, event: "AuctionEnded", logs: logs, sub: sub}, nil
}

// WatchAuctionEnded is a free log subscription operation binding the contract event 0x596165d0521c3cb4157fad2621686f086daed4663acb3d03441a92b9277f5683.
//
// Solidity: event AuctionEnded(uint256 indexed auctionId, address indexed winner, address indexed bidToken, uint256 amount, uint256 amountUsd)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchAuctionEnded(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketAuctionEnded, auctionId []*big.Int, winner []common.Address, bidToken []common.Address) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}
	var bidTokenRule []interface{}
	for _, bidTokenItem := range bidToken {
		bidTokenRule = append(bidTokenRule, bidTokenItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "AuctionEnded", auctionIdRule, winnerRule, bidTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketAuctionEnded)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "AuctionEnded", log); err != nil {
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

// ParseAuctionEnded is a log parse operation binding the contract event 0x596165d0521c3cb4157fad2621686f086daed4663acb3d03441a92b9277f5683.
//
// Solidity: event AuctionEnded(uint256 indexed auctionId, address indexed winner, address indexed bidToken, uint256 amount, uint256 amountUsd)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseAuctionEnded(log types.Log) (*NFTAuctionMarketAuctionEnded, error) {
	event := new(NFTAuctionMarketAuctionEnded)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "AuctionEnded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketBidPlacedIterator is returned from FilterBidPlaced and is used to iterate over the raw logs and unpacked data for BidPlaced events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketBidPlacedIterator struct {
	Event *NFTAuctionMarketBidPlaced // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketBidPlacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketBidPlaced)
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
		it.Event = new(NFTAuctionMarketBidPlaced)
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
func (it *NFTAuctionMarketBidPlacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketBidPlacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketBidPlaced represents a BidPlaced event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketBidPlaced struct {
	AuctionId *big.Int
	Bidder    common.Address
	BidToken  common.Address
	Amount    *big.Int
	AmountUsd *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBidPlaced is a free log retrieval operation binding the contract event 0x2808decb743a25d04efe1bd3dc192acde3be644e2f6ad1dce5d3c46643e1c602.
//
// Solidity: event BidPlaced(uint256 indexed auctionId, address indexed bidder, address indexed bidToken, uint256 amount, uint256 amountUsd)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterBidPlaced(opts *bind.FilterOpts, auctionId []*big.Int, bidder []common.Address, bidToken []common.Address) (*NFTAuctionMarketBidPlacedIterator, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var bidTokenRule []interface{}
	for _, bidTokenItem := range bidToken {
		bidTokenRule = append(bidTokenRule, bidTokenItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "BidPlaced", auctionIdRule, bidderRule, bidTokenRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketBidPlacedIterator{contract: _NFTAuctionMarket.contract, event: "BidPlaced", logs: logs, sub: sub}, nil
}

// WatchBidPlaced is a free log subscription operation binding the contract event 0x2808decb743a25d04efe1bd3dc192acde3be644e2f6ad1dce5d3c46643e1c602.
//
// Solidity: event BidPlaced(uint256 indexed auctionId, address indexed bidder, address indexed bidToken, uint256 amount, uint256 amountUsd)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchBidPlaced(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketBidPlaced, auctionId []*big.Int, bidder []common.Address, bidToken []common.Address) (event.Subscription, error) {

	var auctionIdRule []interface{}
	for _, auctionIdItem := range auctionId {
		auctionIdRule = append(auctionIdRule, auctionIdItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var bidTokenRule []interface{}
	for _, bidTokenItem := range bidToken {
		bidTokenRule = append(bidTokenRule, bidTokenItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "BidPlaced", auctionIdRule, bidderRule, bidTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketBidPlaced)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "BidPlaced", log); err != nil {
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

// ParseBidPlaced is a log parse operation binding the contract event 0x2808decb743a25d04efe1bd3dc192acde3be644e2f6ad1dce5d3c46643e1c602.
//
// Solidity: event BidPlaced(uint256 indexed auctionId, address indexed bidder, address indexed bidToken, uint256 amount, uint256 amountUsd)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseBidPlaced(log types.Log) (*NFTAuctionMarketBidPlaced, error) {
	event := new(NFTAuctionMarketBidPlaced)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "BidPlaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketInitializedIterator struct {
	Event *NFTAuctionMarketInitialized // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketInitialized)
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
		it.Event = new(NFTAuctionMarketInitialized)
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
func (it *NFTAuctionMarketInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketInitialized represents a Initialized event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterInitialized(opts *bind.FilterOpts) (*NFTAuctionMarketInitializedIterator, error) {

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketInitializedIterator{contract: _NFTAuctionMarket.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketInitialized) (event.Subscription, error) {

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketInitialized)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseInitialized(log types.Log) (*NFTAuctionMarketInitialized, error) {
	event := new(NFTAuctionMarketInitialized)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketOwnershipTransferredIterator struct {
	Event *NFTAuctionMarketOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketOwnershipTransferred)
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
		it.Event = new(NFTAuctionMarketOwnershipTransferred)
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
func (it *NFTAuctionMarketOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketOwnershipTransferred represents a OwnershipTransferred event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NFTAuctionMarketOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketOwnershipTransferredIterator{contract: _NFTAuctionMarket.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketOwnershipTransferred)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseOwnershipTransferred(log types.Log) (*NFTAuctionMarketOwnershipTransferred, error) {
	event := new(NFTAuctionMarketOwnershipTransferred)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketPriceFeedSetIterator is returned from FilterPriceFeedSet and is used to iterate over the raw logs and unpacked data for PriceFeedSet events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketPriceFeedSetIterator struct {
	Event *NFTAuctionMarketPriceFeedSet // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketPriceFeedSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketPriceFeedSet)
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
		it.Event = new(NFTAuctionMarketPriceFeedSet)
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
func (it *NFTAuctionMarketPriceFeedSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketPriceFeedSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketPriceFeedSet represents a PriceFeedSet event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketPriceFeedSet struct {
	Token     common.Address
	PriceFeed common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceFeedSet is a free log retrieval operation binding the contract event 0xd2d8394cf7549a5ddbc2ba3dd7b2de8d53c891472d1f2907008ed6a10045fdae.
//
// Solidity: event PriceFeedSet(address indexed token, address indexed priceFeed)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterPriceFeedSet(opts *bind.FilterOpts, token []common.Address, priceFeed []common.Address) (*NFTAuctionMarketPriceFeedSetIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var priceFeedRule []interface{}
	for _, priceFeedItem := range priceFeed {
		priceFeedRule = append(priceFeedRule, priceFeedItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "PriceFeedSet", tokenRule, priceFeedRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketPriceFeedSetIterator{contract: _NFTAuctionMarket.contract, event: "PriceFeedSet", logs: logs, sub: sub}, nil
}

// WatchPriceFeedSet is a free log subscription operation binding the contract event 0xd2d8394cf7549a5ddbc2ba3dd7b2de8d53c891472d1f2907008ed6a10045fdae.
//
// Solidity: event PriceFeedSet(address indexed token, address indexed priceFeed)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchPriceFeedSet(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketPriceFeedSet, token []common.Address, priceFeed []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var priceFeedRule []interface{}
	for _, priceFeedItem := range priceFeed {
		priceFeedRule = append(priceFeedRule, priceFeedItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "PriceFeedSet", tokenRule, priceFeedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketPriceFeedSet)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "PriceFeedSet", log); err != nil {
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

// ParsePriceFeedSet is a log parse operation binding the contract event 0xd2d8394cf7549a5ddbc2ba3dd7b2de8d53c891472d1f2907008ed6a10045fdae.
//
// Solidity: event PriceFeedSet(address indexed token, address indexed priceFeed)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParsePriceFeedSet(log types.Log) (*NFTAuctionMarketPriceFeedSet, error) {
	event := new(NFTAuctionMarketPriceFeedSet)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "PriceFeedSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFTAuctionMarketUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the NFTAuctionMarket contract.
type NFTAuctionMarketUpgradedIterator struct {
	Event *NFTAuctionMarketUpgraded // Event containing the contract specifics and raw log

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
func (it *NFTAuctionMarketUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFTAuctionMarketUpgraded)
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
		it.Event = new(NFTAuctionMarketUpgraded)
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
func (it *NFTAuctionMarketUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFTAuctionMarketUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFTAuctionMarketUpgraded represents a Upgraded event raised by the NFTAuctionMarket contract.
type NFTAuctionMarketUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*NFTAuctionMarketUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &NFTAuctionMarketUpgradedIterator{contract: _NFTAuctionMarket.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *NFTAuctionMarketUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _NFTAuctionMarket.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFTAuctionMarketUpgraded)
				if err := _NFTAuctionMarket.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_NFTAuctionMarket *NFTAuctionMarketFilterer) ParseUpgraded(log types.Log) (*NFTAuctionMarketUpgraded, error) {
	event := new(NFTAuctionMarketUpgraded)
	if err := _NFTAuctionMarket.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
