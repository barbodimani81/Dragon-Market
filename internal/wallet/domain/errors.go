package domain

import "errors"

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientFunds   = errors.New("insufficient available balance")
	ErrNegativeAmount      = errors.New("amount must be greater than zero")
	ErrReservationNotFound = errors.New("no active reservation found to release")
)
