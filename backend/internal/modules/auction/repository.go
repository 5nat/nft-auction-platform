package auction

import (
	"context"
)

type Repository interface {
	ListAuctions(ctx context.Context, query ListAuctionsQuery) ([]Auction, int64, error)
	GetAuction(ctx context.Context, query GetAuctionQuery) (*Auction, error)
	ListBids(ctx context.Context, query ListBidsQuery) ([]Bid, int64, error)
}
