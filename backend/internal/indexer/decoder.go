package indexer

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	auctionCreatedTopic   = crypto.Keccak256Hash([]byte("AuctionCreated(uint256,address,address,uint256,uint256,uint256)"))
	bidPlacedTopic        = crypto.Keccak256Hash([]byte("BidPlaced(uint256,address,address,uint256,uint256)"))
	auctionEndedTopic     = crypto.Keccak256Hash([]byte("AuctionEnded(uint256,address,address,uint256,uint256)"))
	auctionCancelledTopic = crypto.Keccak256Hash([]byte("AuctionCancelled(uint256)"))
)

func normalizeAddress(addr common.Address) string {
	return strings.ToLower(addr.Hex())
}

func normalizeHash(hash common.Hash) string {
	return strings.ToLower(hash.Hex())
}
