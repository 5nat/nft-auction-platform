package auction

import "errors"

var (
	ErrAuctionNotFound   = errors.New("auction not found")
	ErrAuctionNotActive  = errors.New("auction not active")
	ErrAuctionExpired    = errors.New("auction expired")
	ErrAuctionNotExpired = errors.New("auction not expired")
	ErrForbidden         = errors.New("auction action forbidden")
	ErrSellerCannotBid   = errors.New("seller cannot bid on own auction")
	ErrInvalidActor      = errors.New("invalid auction actor")
	ErrChainMismatch     = errors.New("actor chain id does not match auction chain id")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}
