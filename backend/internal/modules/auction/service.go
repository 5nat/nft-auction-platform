package auction

import (
	"context"
	"math"
	"strings"

	"github.com/5nat/nft-auction-platform/backend/internal/infra/persistence/mysql/model"
	"github.com/ethereum/go-ethereum/common"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100

	SortCreatedDesc   = "created_desc"
	SortEndTimeAsc    = "end_time_asc"
	SortLastEventDesc = "last_event_desc"
)

type ServiceConfig struct {
	DefaultChainID         int64
	DefaultContractAddress string
}

type Service struct {
	repo   Repository
	config ServiceConfig
}

func NewService(repo Repository, cfg ServiceConfig) *Service {
	cfg.DefaultContractAddress = normalizeAddressString(cfg.DefaultContractAddress)

	return &Service{
		repo:   repo,
		config: cfg,
	}
}

func (s *Service) ListAuctions(ctx context.Context, query ListAuctionsQuery) (PageResult[AuctionDTO], error) {
	normalizedQuery, err := s.normalizeListAuctionsQuery(query)
	if err != nil {
		return PageResult[AuctionDTO]{}, err
	}

	auction, total, err := s.repo.ListAuctions(ctx, normalizedQuery)
	if err != nil {
		return PageResult[AuctionDTO]{}, err
	}

	items := make([]AuctionDTO, 0, len(auction))
	for _, item := range auction {
		items = append(items, toAuctionDTO(item))
	}

	return PageResult[AuctionDTO]{
		Items: items,
		Meta:  buildPageMeta(normalizedQuery.Page, normalizedQuery.PageSize, total),
	}, nil
}

func (s *Service) GetAuction(ctx context.Context, query GetAuctionQuery) (AuctionDTO, error) {
	normalizedQuery, err := s.normalizeGetAuctionQuery(query)
	if err != nil {
		return AuctionDTO{}, err
	}

	auction, err := s.repo.GetAuction(ctx, normalizedQuery)
	if err != nil {
		return AuctionDTO{}, err
	}

	return toAuctionDTO(*auction), nil
}

func (s *Service) ListBids(ctx context.Context, query ListBidsQuery) (PageResult[BidDTO], error) {
	normalizedQuery, err := s.normalizeListBidsQuery(query)
	if err != nil {
		return PageResult[BidDTO]{}, err
	}

	bids, total, err := s.repo.ListBids(ctx, normalizedQuery)
	if err != nil {
		return PageResult[BidDTO]{}, err
	}

	items := make([]BidDTO, 0, len(bids))
	for _, item := range bids {
		items = append(items, toBidDTO(item))
	}

	return PageResult[BidDTO]{
		Items: items,
		Meta:  buildPageMeta(normalizedQuery.Page, normalizedQuery.PageSize, total),
	}, nil
}

func (s *Service) normalizeListAuctionsQuery(query ListAuctionsQuery) (ListAuctionsQuery, error) {
	query.ChainID = s.resolveChainID(query.ChainID)

	contractAddress, err := s.resolveContractAddress(query.ContractAddress)
	if err != nil {
		return query, err
	}
	query.ContractAddress = contractAddress

	status, err := normalizeStatus(query.Status)
	if err != nil {
		return query, err
	}
	query.Status = status

	seller, err := normalizeOptionalAddress(query.Seller, "seller")
	if err != nil {
		return query, err
	}
	query.Seller = seller

	nft, err := normalizeOptionalAddress(query.NFT, "nft")
	if err != nil {
		return query, err
	}
	query.NFT = nft

	query.Page, query.PageSize = normalizePage(query.Page, query.PageSize)
	query.Sort = normalizeSort(query.Sort)

	return query, nil
}

func (s *Service) normalizeGetAuctionQuery(query GetAuctionQuery) (GetAuctionQuery, error) {
	query.ChainID = s.resolveChainID(query.ChainID)

	contractAddress, err := s.resolveContractAddress(query.ContractAddress)
	if err != nil {
		return query, err
	}
	query.ContractAddress = contractAddress

	return query, nil
}

func (s *Service) normalizeListBidsQuery(query ListBidsQuery) (ListBidsQuery, error) {
	query.ChainID = s.resolveChainID(query.ChainID)

	contractAddress, err := s.resolveContractAddress(query.ContractAddress)
	if err != nil {
		return query, err
	}

	query.ContractAddress = contractAddress

	query.Page, query.PageSize = normalizePage(query.Page, query.PageSize)

	return query, nil
}

func (s *Service) resolveChainID(chainID int64) int64 {
	if chainID != 0 {
		return chainID
	}

	return s.config.DefaultChainID
}

func (s *Service) resolveContractAddress(contractAddress string) (string, error) {
	contractAddress = strings.TrimSpace(contractAddress)
	if contractAddress == "" {
		return s.config.DefaultContractAddress, nil
	}

	return normalizeRequiredAddress(contractAddress, "contract_address")
}

func normalizeSort(sort string) string {
	sort = strings.ToLower(strings.TrimSpace(sort))

	switch sort {
	case SortCreatedDesc, SortEndTimeAsc, SortLastEventDesc:
		return sort
	default:
		return SortCreatedDesc
	}
}

func normalizePage(page int, PageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}

	if PageSize <= 0 {
		PageSize = DefaultPageSize
	}

	if PageSize > MaxPageSize {
		PageSize = MaxPageSize
	}

	return page, PageSize
}

func normalizeStatus(status string) (string, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status == "" {
		return "", nil
	}

	switch status {
	case model.AuctionStatusActive, model.AuctionStatusEnded, model.AuctionStatusCancelled:
		return status, nil
	default:
		return "", NewValidationError("invalid status")
	}
}

func normalizeOptionalAddress(address string, field string) (string, error) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", nil
	}

	return normalizeRequiredAddress(address, field)
}

func normalizeRequiredAddress(address string, field string) (string, error) {
	address = strings.TrimSpace(address)
	if !common.IsHexAddress(address) {
		return "", NewValidationError("invalid" + field)
	}

	return normalizeAddressString(address), nil
}

func normalizeAddressString(address string) string {
	address = strings.TrimSpace(address)
	if address == "" {
		return ""
	}

	return strings.ToLower(common.HexToAddress(address).Hex())
}

func buildPageMeta(page int, pageSize int, total int64) PageMeta {
	totalPages := 0
	if pageSize > 0 && total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	return PageMeta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

func toAuctionDTO(item model.Auction) AuctionDTO {
	return AuctionDTO{
		ChainID:         item.ChainID,
		ContractAddress: item.ContractAddress,
		AuctionID:       item.AuctionID,

		Seller:      item.Seller,
		NFTContract: item.NFTContract,
		TokenID:     item.TokenID,

		MinBidUSD: item.MinBidUSD,

		HighestBidder:    item.HighestBidder,
		HighestBidToken:  item.HighestBidToken,
		HighestBidAmount: item.HighestBidAmount,
		HighestBidUSD:    item.HighestBidUSD,

		Status:  item.Status,
		EndTime: item.EndTime,

		CreatedTxHash:      item.CreatedTxHash,
		CreatedBlockNumber: item.CreatedBlockNumber,
		CreatedBlockHash:   item.CreatedBlockHash,
		CreatedLogIndex:    item.CreatedLogIndex,

		LastEventName:        item.LastEventName,
		LastEventTxHash:      item.LastEventTxHash,
		LastEventBlockNumber: item.LastEventBlockNumber,
		LastEventBlockHash:   item.LastEventBlockHash,
		LastEventLogIndex:    item.LastEventLogIndex,

		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func toBidDTO(item model.Bid) BidDTO {
	return BidDTO{
		ChainID:         item.ChainID,
		ContractAddress: item.ContractAddress,
		AuctionID:       item.AuctionID,

		Bidder:   item.Bidder,
		BidToken: item.BidToken,

		Amount:    item.Amount,
		AmountUSD: item.AmountUSD,

		TxHash:      item.TxHash,
		LogIndex:    item.LogIndex,
		BlockNumber: item.BlockNumber,
		BlockHash:   item.BlockHash,

		CreatedAt: item.CreatedAt,
	}
}
