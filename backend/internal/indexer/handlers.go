package indexer

import (
	"context"
	"fmt"

	"github.com/5nat/nft-auction-platform/backend/internal/model"
	"github.com/ethereum/go-ethereum/core/types"
)

func (idx *Indexer) processAuctionCreated(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseAuctionCreated(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse AuctionCreated: %w", err)
	}

	inserted := false

	err = idx.repo.WithTx(ctx, func(repo EventRepository) error {
		var insertProcessedLogErr error

		inserted, insertProcessedLogErr = repo.InsertProcessedLog(ctx, lg, "AuctionCreated")
		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}

		if !inserted {
			return nil
		}

		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		if !event.EndTime.IsUint64() {
			return fmt.Errorf("end time overflow: %s", event.EndTime.String())
		}

		auction := model.Auction{
			ChainID:            idx.chainID,
			ContractAddress:    normalizeAddress(idx.marketAddress),
			AuctionID:          event.AuctionId.Uint64(),
			Seller:             normalizeAddress(event.Seller),
			NFTContract:        normalizeAddress(event.Nft),
			TokenID:            event.TokenId.String(),
			MinBidUSD:          event.MinBidUsd.String(),
			HighestBidder:      "",
			HighestBidToken:    "",
			HighestBidAmount:   "0",
			HighestBidUSD:      "0",
			EndTime:            event.EndTime.Uint64(),
			Status:             model.AuctionStatusActive,
			CreatedTxHash:      normalizeHash(lg.TxHash),
			CreatedBlockNumber: lg.BlockNumber,
		}

		return repo.CreateAuction(ctx, auction)
	})

	if err != nil {
		return false, err
	}

	if !inserted {
		idx.logger.Debug(
			"AuctionCreated already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
		return false, nil
	}

	idx.logger.Info(
		"AuctionCreated indexed",
		"auction_id", event.AuctionId.String(),
		"seller", event.Seller.Hex(),
		"nft", event.Nft.Hex(),
		"token_id", event.TokenId.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) processBidPlaced(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseBidPlaced(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse BidPlaced: %w", err)
	}

	inserted := false

	err = idx.repo.WithTx(ctx, func(repo EventRepository) error {
		var insertProcessedLogErr error
		inserted, insertProcessedLogErr = repo.InsertProcessedLog(ctx, lg, "BidPlaced")

		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}

		if !inserted {
			return nil
		}

		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		auctionID := event.AuctionId.Uint64()

		bid := model.Bid{
			ChainID:         idx.chainID,
			ContractAddress: normalizeAddress(idx.marketAddress),
			AuctionID:       auctionID,
			Bidder:          normalizeAddress(event.Bidder),
			BidToken:        normalizeAddress(event.BidToken),
			Amount:          event.Amount.String(),
			AmountUSD:       event.AmountUsd.String(),
			TxHash:          normalizeHash(lg.TxHash),
			LogIndex:        uint64(lg.Index),
			BlockNumber:     lg.BlockNumber,
			BlockHash:       normalizeHash(lg.BlockHash),
		}

		if createBidErr := repo.CreateBid(ctx, bid); createBidErr != nil {
			return createBidErr
		}

		return repo.UpdateAuctionHighestBid(ctx, UpdateAuctionHighestBidInput{
			AuctionID: auctionID,
			Bidder:    normalizeAddress(event.Bidder),
			BidToken:  normalizeAddress(event.BidToken),
			Amount:    event.Amount.String(),
			AmountUSD: event.AmountUsd.String(),
		})
	})

	if err != nil {
		return false, err
	}

	if !inserted {
		idx.logger.Debug(
			"BIdPlaced already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
	}

	idx.logger.Info(
		"BidPlaced indexed",
		"auction_id", event.AuctionId.String(),
		"bidder", event.Bidder.Hex(),
		"bid_token", event.BidToken.Hex(),
		"amount", event.Amount.String(),
		"amount_usd", event.AmountUsd.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) processAuctionEnded(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseAuctionEnded(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse AuctionEnded: %w", err)
	}

	inserted := false

	err = idx.repo.WithTx(ctx, func(repo EventRepository) error {
		var insertProcessedLogErr error

		inserted, insertProcessedLogErr = repo.InsertProcessedLog(ctx, lg, "AuctionEnded")
		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}

		if !inserted {
			return nil
		}

		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		auctionID := event.AuctionId.Uint64()

		return repo.MarkAuctionEnded(ctx, MarkAuctionEndedInput{
			AuctionID: auctionID,
			Winner:    normalizeAddress(event.Winner),
			BidToken:  normalizeAddress(event.BidToken),
			Amount:    event.Amount.String(),
			AmountUSD: event.AmountUsd.String(),
		})
	})

	if err != nil {
		return false, err
	}

	if !inserted {
		idx.logger.Debug(
			"AuctionEnded already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
		return false, nil
	}

	idx.logger.Info(
		"AuctionEnded indexed",
		"auction_id", event.AuctionId.String(),
		"winner", event.Winner.Hex(),
		"bid_token", event.BidToken.Hex(),
		"amount", event.Amount.String(),
		"amount_usd", event.AmountUsd.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}

func (idx *Indexer) processAuctionCancelled(ctx context.Context, lg types.Log) (bool, error) {
	event, err := idx.market.ParseAuctionCancelled(lg)
	if err != nil {
		return false, fmt.Errorf("failed to parse AuctionCancelled: %w", err)
	}

	inserted := false

	err = idx.repo.WithTx(ctx, func(repo EventRepository) error {
		var insertProcessedLogErr error

		inserted, insertProcessedLogErr = repo.InsertProcessedLog(ctx, lg, "AuctionCancelled")
		if insertProcessedLogErr != nil {
			return insertProcessedLogErr
		}

		if !inserted {
			return nil
		}

		if !event.AuctionId.IsUint64() {
			return fmt.Errorf("auction id overflow: %s", event.AuctionId.String())
		}

		auctionID := event.AuctionId.Uint64()

		return repo.MarkAuctionCancelled(ctx, MarkAuctionCancelledInput{
			AuctionID: auctionID,
		})
	})

	if err != nil {
		return false, err
	}
	if !inserted {
		idx.logger.Debug(
			"AuctionCancelled already processed",
			"auction_id", event.AuctionId.String(),
			"tx_hash", lg.TxHash.Hex(),
			"log_index", lg.Index,
		)
	}

	idx.logger.Info(
		"AuctionCancelled indexed",
		"auction_id", event.AuctionId.String(),
		"block_number", lg.BlockNumber,
		"tx_hash", lg.TxHash.Hex(),
		"log_index", lg.Index,
	)

	return true, nil
}
