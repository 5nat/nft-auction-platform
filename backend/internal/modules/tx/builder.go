package tx

import "context"

type CalldataBuilder interface {
	BuildApproveNFTCalldata(ctx context.Context, req BuildApproveNFTTxRequest) (string, error)
	BuildCreateAuctionCalldata(ctx context.Context, req BuildCreateAuctionTxRequest) (string, error)
	BuildPlaceBidCalldata(ctx context.Context, req BuildPlaceBidTxRequest) (string, error)
	BuildCancelAuctionCalldata(ctx context.Context, req BuildCancelAuctionTxRequest) (string, error)
	BuildEndAuctionCalldata(ctx context.Context, req BuildEndAuctionTxRequest) (string, error)
}
