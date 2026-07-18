package database

import (
	"context"
	"database/sql"

	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/google/uuid"
)

type TransactionQuerier interface {
	GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (gen.Wallet, error)
	GetListingForUpdate(ctx context.Context, id uuid.UUID) (gen.Listing, error)
	GetAuctionForBid(ctx context.Context, id uuid.UUID) (gen.GetAuctionForBidRow, error)
	ReserveFunds(ctx context.Context, arg gen.ReserveFundsParams) (gen.Wallet, error)
	ReleaseFunds(ctx context.Context, arg gen.ReleaseFundsParams) (gen.Wallet, error)
	DecreaseWalletBalance(ctx context.Context, arg gen.DecreaseWalletBalanceParams) (gen.Wallet, error)
	UpdateWalletBalance(ctx context.Context, arg gen.UpdateWalletBalanceParams) (gen.Wallet, error)
	PlaceBid(ctx context.Context, arg gen.PlaceBidParams) (gen.PlaceBidRow, error)
	UpdateListingStatus(ctx context.Context, arg gen.UpdateListingStatusParams) (gen.Listing, error)
	TransferItemOwnership(ctx context.Context, arg gen.TransferItemOwnershipParams) (gen.Item, error)
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(q TransactionQuerier) error) error
}

type SQLTransactor struct {
	db *sql.DB
}

func NewSQLTransactor(db *sql.DB) *SQLTransactor {
	return &SQLTransactor{db: db}
}

func (t *SQLTransactor) WithinTransaction(ctx context.Context, fn func(q TransactionQuerier) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := gen.New(tx)
	if err := fn(queries); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
