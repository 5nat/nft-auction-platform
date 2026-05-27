package auction

import "time"

type Policy struct{}

func NewPolicy() *Policy {
	return &Policy{}
}

func (p *Policy) EnsureCanPlaceBid(a *Auction, now time.Time) error {
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

func (p *Policy) EnsureCanCancel(a *Auction, now time.Time) error {
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

func (p *Policy) EnsureCanEnd(a *Auction, now time.Time) error {
	if a == nil {
		return ErrAuctionNotFound
	}

	if a.Status != AuctionStatusActive {
		return ErrAuctionNotActive
	}

	if a.EndTime > 0 && uint64(now.Unix()) < a.EndTime {
		return ErrAuctionNotExpired
	}

	return nil
}
