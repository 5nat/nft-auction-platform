package tx

import (
	"context"
	"math/big"
	"strings"
	"time"

	auctionmodule "github.com/5nat/nft-auction-platform/backend/internal/modules/auction"
	"github.com/ethereum/go-ethereum/common"
)

const ETHBidToken = "0x0000000000000000000000000000000000000000"

type ServiceConfig struct {
	DefaultChainID         int64
	DefaultContractAddress string
}

type Service struct {
	builder     CalldataBuilder
	config      ServiceConfig
	auctionRepo auctionmodule.Repository
	policy      *auctionmodule.Policy
}

func NewService(builder CalldataBuilder, cfg ServiceConfig, auctionRepo auctionmodule.Repository) *Service {
	cfg.DefaultContractAddress = normalizeAddressString(cfg.DefaultContractAddress)

	return &Service{
		builder:     builder,
		config:      cfg,
		auctionRepo: auctionRepo,
		policy:      auctionmodule.NewPolicy(),
	}
}

func (s *Service) BuildApproveNFTTx(ctx context.Context, req BuildApproveNFTTxRequest) (TransactionRequestDTO, error) {
	normalizedReq, err := s.normalizeApproveNFTRequest(req)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	data, err := s.builder.BuildApproveNFTCalldata(ctx, normalizedReq)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	return TransactionRequestDTO{
		ChainID: normalizedReq.ChainID,
		To:      normalizedReq.NFTContract,
		Data:    data,
		Value:   "0",
	}, nil
}

func (s *Service) BuildCreateAuctionTx(ctx context.Context, req BuildCreateAuctionTxRequest) (TransactionRequestDTO, error) {
	normalizedReq, err := s.normalizeCreateAuctionRequest(req)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	data, err := s.builder.BuildCreateAuctionCalldata(ctx, normalizedReq)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	return TransactionRequestDTO{
		ChainID: normalizedReq.ChainID,
		To:      normalizedReq.ContractAddress,
		Data:    data,
		Value:   "0",
	}, nil
}

func (s *Service) BuildPlaceBidTx(ctx context.Context, req BuildPlaceBidTxRequest) (TransactionRequestDTO, error) {
	normalizedReq, err := s.normalizePlaceBidRequest(req)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	if err := s.ensureCanPlaceBid(ctx, normalizedReq.AuctionID); err != nil {
		return TransactionRequestDTO{}, err
	}

	data, err := s.builder.BuildPlaceBidCalldata(ctx, normalizedReq)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	value := "0"
	if isETHBidToken(normalizedReq.BidToken) {
		value = normalizedReq.Amount
	}

	return TransactionRequestDTO{
		ChainID: normalizedReq.ChainID,
		To:      normalizedReq.ContractAddress,
		Data:    data,
		Value:   value,
	}, nil
}

func (s *Service) BuildCancelAuctionTx(ctx context.Context, req BuildCancelAuctionTxRequest) (TransactionRequestDTO, error) {
	normalizedReq, err := s.normalizeCancelAuctionRequest(req)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	if err := s.ensureCanCancel(ctx, normalizedReq.AuctionID); err != nil {
		return TransactionRequestDTO{}, err
	}

	data, err := s.builder.BuildCancelAuctionCalldata(ctx, req)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	return TransactionRequestDTO{
		ChainID: normalizedReq.ChainID,
		To:      normalizedReq.ContractAddress,
		Data:    data,
		Value:   "0",
	}, nil
}

func (s *Service) BuildEndAuctionTx(ctx context.Context, req BuildEndAuctionTxRequest) (TransactionRequestDTO, error) {
	normalizedReq, err := s.normalizeEndAuctionRequest(req)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	if err := s.ensureCanEnd(ctx, normalizedReq.AuctionID); err != nil {
		return TransactionRequestDTO{}, err
	}

	data, err := s.builder.BuildEndAuctionCalldata(ctx, normalizedReq)
	if err != nil {
		return TransactionRequestDTO{}, err
	}

	return TransactionRequestDTO{
		ChainID: normalizedReq.ChainID,
		To:      normalizedReq.ContractAddress,
		Data:    data,
		Value:   "0",
	}, nil
}

func (s *Service) ensureCanPlaceBid(ctx context.Context, auctionID uint64) error {
	if s.auctionRepo == nil {
		return nil
	}

	a, err := s.auctionRepo.GetAuction(ctx, auctionmodule.GetAuctionQuery{
		AuctionID: auctionID,
	})
	if err != nil {
		return err
	}

	return s.policy.EnsureCanPlaceBid(a, time.Now())
}

func (s *Service) ensureCanCancel(ctx context.Context, auctionID uint64) error {
	if s.auctionRepo == nil {
		return nil
	}

	a, err := s.auctionRepo.GetAuction(ctx, auctionmodule.GetAuctionQuery{
		AuctionID: auctionID,
	})
	if err != nil {
		return err
	}

	return s.policy.EnsureCanCancel(a, time.Now())
}

func (s *Service) ensureCanEnd(ctx context.Context, auctionID uint64) error {
	if s.auctionRepo == nil {
		return nil
	}

	a, err := s.auctionRepo.GetAuction(ctx, auctionmodule.GetAuctionQuery{
		AuctionID: auctionID,
	})
	if err != nil {
		return err
	}

	return s.policy.EnsureCanEnd(a, time.Now())
}

func (s *Service) normalizeApproveNFTRequest(req BuildApproveNFTTxRequest) (BuildApproveNFTTxRequest, error) {
	req.ChainID = s.resolveChainID(req.ChainID)

	nftContract, err := normalizeRequiredAddress(req.NFTContract, ErrInvalidNFTContract)
	if err != nil {
		return req, err
	}
	req.NFTContract = nftContract

	if !isValidUintString(req.TokenID) {
		return req, ErrInvalidTokenID
	}

	operator := strings.TrimSpace(req.Operator)
	if operator == "" {
		operator = s.config.DefaultContractAddress
	}

	operator, err = normalizeRequiredAddress(operator, ErrInvalidOperator)
	if err != nil {
		return req, err
	}
	req.Operator = operator

	return req, nil
}

func (s *Service) normalizeCreateAuctionRequest(req BuildCreateAuctionTxRequest) (BuildCreateAuctionTxRequest, error) {
	req.ChainID = s.resolveChainID(req.ChainID)

	contractAddress, err := s.resolveContractAddress(req.ContractAddress)
	if err != nil {
		return req, err
	}
	req.ContractAddress = contractAddress

	nftContract, err := normalizeRequiredAddress(req.NFTContract, ErrInvalidNFTContract)
	if err != nil {
		return req, err
	}
	req.NFTContract = nftContract

	if !isValidUintString(req.TokenID) {
		return req, ErrInvalidTokenID
	}

	if !isValidUintString(req.MinBidUSD) {
		return req, ErrInvalidAmount
	}

	if req.Duration == 0 {
		return req, ErrInvalidDuration
	}

	return req, nil
}

func (s *Service) normalizePlaceBidRequest(req BuildPlaceBidTxRequest) (BuildPlaceBidTxRequest, error) {
	req.ChainID = s.resolveChainID(req.ChainID)

	contractAddress, err := s.resolveContractAddress(req.ContractAddress)
	if err != nil {
		return req, err
	}
	req.ContractAddress = contractAddress

	if req.AuctionID == 0 {
		return req, ErrInvalidAuctionID
	}

	req.BidToken = normalizeBidToken(req.BidToken)

	if !isETHBidToken(req.BidToken) {
		bidToken, err := normalizeRequiredAddress(req.BidToken, ErrInvalidBidToken)
		if err != nil {
			return req, err
		}
		req.BidToken = bidToken
	}

	if !isValidUintString(req.Amount) {
		return req, ErrInvalidAmount
	}

	return req, nil
}

func (s *Service) normalizeCancelAuctionRequest(req BuildCancelAuctionTxRequest) (BuildCancelAuctionTxRequest, error) {
	req.ChainID = s.resolveChainID(req.ChainID)

	contractAddress, err := s.resolveContractAddress(req.ContractAddress)
	if err != nil {
		return req, err
	}
	req.ContractAddress = contractAddress

	if req.AuctionID == 0 {
		return req, ErrInvalidAuctionID
	}

	return req, nil
}

func (s *Service) normalizeEndAuctionRequest(req BuildEndAuctionTxRequest) (BuildEndAuctionTxRequest, error) {
	req.ChainID = s.resolveChainID(req.ChainID)

	contractAddress, err := s.resolveContractAddress(req.ContractAddress)
	if err != nil {
		return req, err
	}
	req.ContractAddress = contractAddress

	if req.AuctionID == 0 {
		return req, ErrInvalidAuctionID
	}

	return req, nil
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
		if s.config.DefaultContractAddress == "" {
			return "", ErrInvalidContractAddress
		}

		return s.config.DefaultContractAddress, nil
	}

	return normalizeRequiredAddress(contractAddress, ErrInvalidContractAddress)
}

func normalizeRequiredAddress(address string, err error) (string, error) {
	address = strings.TrimSpace(address)
	if !common.IsHexAddress(address) {
		return "", err
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

func isValidUintString(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}

	n, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return false
	}

	return n.Sign() >= 0
}

func normalizeBidToken(token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return ETHBidToken
	}

	return normalizeAddressString(token)
}

func isETHBidToken(token string) bool {
	return common.HexToAddress(token) == common.Address{}
}
