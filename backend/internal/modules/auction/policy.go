package auction

import (
	"strings"
	"time"
)

type Actor struct {
	WalletAddress string
	ChainID       int64
}

func NewActor(walletAddress string, chainID int64) Actor {
	return Actor{
		WalletAddress: strings.TrimSpace(walletAddress),
		ChainID:       chainID,
	}
}

func (a Actor) IsZero() bool {
	return a.WalletAddress == "" || a.ChainID <= 0
}

type Policy struct{}

func NewPolicy() *Policy {
	return &Policy{}
}

func (p *Policy) EnsureCanPlaceBid(a *Auction, actor Actor, now time.Time) error {
	if err := p.ensureActiveAuction(a, now); err != nil {
		return err
	}

	if err := ensureActorMatchesAuction(a, actor); err != nil {
		return err
	}

	if sameAddress(actor.WalletAddress, a.Seller) {
		return ErrSellerCannotBid
	}

	return nil
}

func (p *Policy) EnsureCanCancel(a *Auction, actor Actor, now time.Time) error {
	if err := p.ensureActiveAuction(a, now); err != nil {
		return err
	}

	if err := ensureActorMatchesAuction(a, actor); err != nil {
		return err
	}

	if !sameAddress(actor.WalletAddress, a.Seller) {
		return ErrForbidden
	}

	return nil
}

func (p *Policy) EnsureCanEnd(a *Auction, actor Actor, now time.Time) error {
	if a == nil {
		return ErrAuctionNotFound
	}

	if a.Status != AuctionStatusActive {
		return ErrAuctionNotActive
	}

	if err := ensureActorMatchesAuction(a, actor); err != nil {
		return err
	}

	if a.EndTime > 0 && uint64(now.Unix()) < a.EndTime {
		return ErrAuctionNotExpired
	}

	// 当前设计：end auction 允许任何登录用户触发。
	// 原因是很多链上拍卖合约的 end/settle 是 permissionless 的：
	// 只要拍卖已结束，任何人都可以帮忙结算。
	// 如果你的合约限制只有 seller 可以 end，这里再改成 seller 校验。
	return nil
}

func (p *Policy) ensureActiveAuction(a *Auction, now time.Time) error {
	if a == nil {
		return ErrAuctionNotFound
	}

	if a.Status != AuctionStatusActive {
		return ErrAuctionNotActive
	}

	if a.EndTime > 0 && uint64(now.Unix()) >= a.EndTime {
		return ErrAuctionExpired
	}

	return nil
}

func ensureActorMatchesAuction(a *Auction, actor Actor) error {
	if actor.IsZero() {
		return ErrInvalidActor
	}

	if a.ChainID > 0 && actor.ChainID != a.ChainID {
		return ErrChainMismatch
	}

	return nil
}

func sameAddress(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}
