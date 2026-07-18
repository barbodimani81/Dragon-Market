package service

import (
	"context"

	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/barbodimani81/Dragon-Market/internal/wallet/domain"
	"github.com/google/uuid"
)

type WalletService struct {
	repo gen.Querier
}

func NewWalletService(repo gen.Querier) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetWallet(ctx context.Context, userID uuid.UUID) (gen.Wallet, error) {
	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return gen.Wallet{}, domain.ErrWalletNotFound
	}
	return wallet, nil
}

func (s *WalletService) Deposit(ctx context.Context, userID uuid.UUID, amount int64) (gen.Wallet, error) {
	if amount <= 0 {
		return gen.Wallet{}, domain.ErrNegativeAmount
	}

	// sqlc names this field based on the query input argument order or names
	wallet, err := s.repo.UpdateWalletBalance(ctx, gen.UpdateWalletBalanceParams{
		UserID: userID,
		Amount: amount, // Updates available_balance = available_balance + $2
	})
	if err != nil {
		return gen.Wallet{}, domain.ErrWalletNotFound
	}
	return wallet, nil
}

func (s *WalletService) ReserveFunds(ctx context.Context, userID uuid.UUID, amount int64) (gen.Wallet, error) {
	if amount <= 0 {
		return gen.Wallet{}, domain.ErrNegativeAmount
	}

	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return gen.Wallet{}, domain.ErrWalletNotFound
	}

	if wallet.AvailableBalance < amount {
		return gen.Wallet{}, domain.ErrInsufficientFunds
	}

	return s.repo.ReserveFunds(ctx, gen.ReserveFundsParams{
		UserID: userID,
		Amount: amount,
	})
}

func (s *WalletService) ReleaseFunds(ctx context.Context, userID uuid.UUID, amount int64) (gen.Wallet, error) {
	if amount <= 0 {
		return gen.Wallet{}, domain.ErrNegativeAmount
	}

	wallet, err := s.repo.GetWallet(ctx, userID)
	if err != nil {
		return gen.Wallet{}, domain.ErrWalletNotFound
	}

	if wallet.ReservedBalance < amount {
		return gen.Wallet{}, domain.ErrReservationNotFound
	}

	return s.repo.ReleaseFunds(ctx, gen.ReleaseFundsParams{
		UserID: userID,
		Amount: amount,
	})
}
