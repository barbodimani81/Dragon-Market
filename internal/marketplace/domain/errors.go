package domain

import "errors"

var (
	ErrListingNotFound      = errors.New("marketplace listing not found")
	ErrInvalidListingRarity = errors.New("only COMMON and RARE items can be listed on the fixed-price marketplace")
	ErrSellerIsBuyer        = errors.New("you cannot purchase your own item listing")
	ErrListingNotActive     = errors.New("this listing is no longer active")
	ErrNegativePrice        = errors.New("listing price must be greater than zero")
)
