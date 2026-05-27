package auction

import "errors"

var (
	ErrAuctionNotFound   = errors.New("auction not found")
	ErrAuctionNotActive  = errors.New("auction not active")
	ErrAuctionExpired    = errors.New("auction expired")
	ErrAuctionNotExpired = errors.New("auction not expired")
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
