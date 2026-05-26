package indexer

import "github.com/ethereum/go-ethereum/core/types"

const (
	EventNameAuctionCreated   = "AuctionCreated"
	EventNameBidPlaced        = "BidPlaced"
	EventNameAuctionEnded     = "AuctionEnded"
	EventNameAuctionCancelled = "AuctionCancelled"
)

type EventMeta struct {
	EventName   string
	TxHash      string
	BlockNumber uint64
	BlockHash   string
	LogIndex    uint64
}

func newEventMeta(eventName string, lg types.Log) EventMeta {
	return EventMeta{
		EventName:   eventName,
		TxHash:      normalizeHash(lg.TxHash),
		BlockNumber: lg.BlockNumber,
		BlockHash:   normalizeHash(lg.BlockHash),
		LogIndex:    uint64(lg.Index),
	}
}
