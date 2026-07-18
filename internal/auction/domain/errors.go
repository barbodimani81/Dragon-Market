package domain

import "errors"

var (
	ErrAuctionNotFound   = errors.New("auction not found or inactive")
	ErrAuctionExpired    = errors.New("this auction has already expired")
	ErrBidTooLow         = errors.New("bid amount must be higher than the current highest bid")
	ErrSellerCannotBid   = errors.New("you cannot place a bid on your own auction")
	ErrInvalidDuration   = errors.New("auction duration must be between 1 minute and 7 days")
	ErrInvalidRarity     = errors.New("only LEGENDARY items can be placed up for auction")
)
