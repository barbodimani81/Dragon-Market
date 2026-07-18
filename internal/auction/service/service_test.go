package service

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/barbodimani81/Dragon-Market/internal/auction/domain"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/barbodimani81/Dragon-Market/pkg/database"
	"github.com/google/uuid"
)

type fakeAuctionTx struct {
	calls    []string
	wallet   gen.Wallet
	auction  gen.GetAuctionForBidRow
	result   gen.PlaceBidRow
}

func (f *fakeAuctionTx) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (gen.Wallet, error) {
	f.calls = append(f.calls, "wallet_for_update:"+userID.String())
	return f.wallet, nil
}

func (f *fakeAuctionTx) GetListingForUpdate(ctx context.Context, id uuid.UUID) (gen.Listing, error) {
	return gen.Listing{}, nil
}

func (f *fakeAuctionTx) GetAuctionForBid(ctx context.Context, id uuid.UUID) (gen.GetAuctionForBidRow, error) {
	f.calls = append(f.calls, "auction_for_bid:"+id.String())
	return f.auction, nil
}

func (f *fakeAuctionTx) ReserveFunds(ctx context.Context, arg gen.ReserveFundsParams) (gen.Wallet, error) {
	f.calls = append(f.calls, "reserve_funds")
	return f.wallet, nil
}

func (f *fakeAuctionTx) ReleaseFunds(ctx context.Context, arg gen.ReleaseFundsParams) (gen.Wallet, error) {
	f.calls = append(f.calls, "release_funds")
	return f.wallet, nil
}

func (f *fakeAuctionTx) DecreaseWalletBalance(ctx context.Context, arg gen.DecreaseWalletBalanceParams) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeAuctionTx) UpdateWalletBalance(ctx context.Context, arg gen.UpdateWalletBalanceParams) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeAuctionTx) PlaceBid(ctx context.Context, arg gen.PlaceBidParams) (gen.PlaceBidRow, error) {
	f.calls = append(f.calls, "place_bid")
	return f.result, nil
}

func (f *fakeAuctionTx) UpdateListingStatus(ctx context.Context, arg gen.UpdateListingStatusParams) (gen.Listing, error) {
	return gen.Listing{}, nil
}

func (f *fakeAuctionTx) TransferItemOwnership(ctx context.Context, arg gen.TransferItemOwnershipParams) (gen.Item, error) {
	return gen.Item{}, nil
}

func (f *fakeAuctionTx) GetItem(ctx context.Context, id uuid.UUID) (gen.Item, error) {
	return gen.Item{}, nil
}

func (f *fakeAuctionTx) GetListing(ctx context.Context, id uuid.UUID) (gen.Listing, error) {
	return gen.Listing{}, nil
}

func (f *fakeAuctionTx) GetWallet(ctx context.Context, userID uuid.UUID) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeAuctionTx) CreateAuction(ctx context.Context, arg gen.CreateAuctionParams) (gen.CreateAuctionRow, error) {
	return gen.CreateAuctionRow{}, nil
}

func (f *fakeAuctionTx) CreateItem(ctx context.Context, arg gen.CreateItemParams) (gen.Item, error) {
	return gen.Item{}, nil
}

func (f *fakeAuctionTx) CreateListing(ctx context.Context, arg gen.CreateListingParams) (gen.Listing, error) {
	return gen.Listing{}, nil
}

func (f *fakeAuctionTx) CreateWallet(ctx context.Context, arg gen.CreateWalletParams) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeAuctionTx) ListActiveAuctions(ctx context.Context) ([]gen.ListActiveAuctionsRow, error) {
	return nil, nil
}

func (f *fakeAuctionTx) ListActiveListings(ctx context.Context) ([]gen.Listing, error) {
	return nil, nil
}

func (f *fakeAuctionTx) ListItemsByOwner(ctx context.Context, ownerID uuid.UUID) ([]gen.Item, error) {
	return nil, nil
}

func (f *fakeAuctionTx) ReleaseFundsCalls() int {
	count := 0
	for _, call := range f.calls {
		if call == "release_funds" {
			count++
		}
	}
	return count
}

type fakeTransactor struct {
	querier database.TransactionQuerier
	calls   int
}

func (f *fakeTransactor) WithinTransaction(ctx context.Context, fn func(q database.TransactionQuerier) error) error {
	f.calls++
	return fn(f.querier)
}

func TestPlaceAuctionBid_TransactionFlow(t *testing.T) {
	bidderID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	sellerID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	previousBidderID := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	auctionID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	itemID := uuid.MustParse("00000000-0000-0000-0000-000000000011")

	txQuerier := &fakeAuctionTx{
		wallet: gen.Wallet{
			ID:               uuid.New(),
			UserID:           bidderID,
			AvailableBalance: 500,
			ReservedBalance:  0,
		},
		auction: gen.GetAuctionForBidRow{
			ID:                   auctionID,
			ItemID:               itemID,
			SellerID:             sellerID,
			StartPrice:           100,
			CurrentHighestBid:    sql.NullInt64{Int64: 120, Valid: true},
			CurrentHighestBidder: uuid.NullUUID{UUID: previousBidderID, Valid: true},
			Status:               "ACTIVE",
			ExpiresAt:            time.Now().Add(time.Hour),
		},
		result: gen.PlaceBidRow{
			ID:                   auctionID,
			ItemID:               itemID,
			SellerID:             sellerID,
			StartPrice:           100,
			CurrentHighestBid:    sql.NullInt64{Int64: 130, Valid: true},
			CurrentHighestBidder: uuid.NullUUID{UUID: bidderID, Valid: true},
			Status:               "ACTIVE",
			ExpiresAt:            time.Now().Add(time.Hour),
			CreatedAt:            time.Now(),
		},
	}

	svc := NewAuctionService(nil, &fakeTransactor{querier: txQuerier})

	got, err := svc.PlaceAuctionBid(context.Background(), auctionID, bidderID, 130)
	if err != nil {
		t.Fatalf("PlaceAuctionBid() error = %v", err)
	}

	if !reflect.DeepEqual(got, txQuerier.result) {
		t.Fatalf("PlaceAuctionBid() got = %#v, want %#v", got, txQuerier.result)
	}

	wantCalls := []string{
		"wallet_for_update:" + bidderID.String(),
		"auction_for_bid:" + auctionID.String(),
		"reserve_funds",
		"release_funds",
		"place_bid",
	}
	if !reflect.DeepEqual(txQuerier.calls, wantCalls) {
		t.Fatalf("transaction calls = %#v, want %#v", txQuerier.calls, wantCalls)
	}
}

func TestPlaceAuctionBid_TooLowStopsBeforeMutation(t *testing.T) {
	bidderID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	sellerID := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	auctionID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	itemID := uuid.MustParse("00000000-0000-0000-0000-000000000011")

	txQuerier := &fakeAuctionTx{
		wallet: gen.Wallet{
			ID:               uuid.New(),
			UserID:           bidderID,
			AvailableBalance: 500,
			ReservedBalance:  0,
		},
		auction: gen.GetAuctionForBidRow{
			ID:                   auctionID,
			ItemID:               itemID,
			SellerID:             sellerID,
			StartPrice:           100,
			CurrentHighestBid:    sql.NullInt64{Int64: 120, Valid: true},
			CurrentHighestBidder: uuid.NullUUID{Valid: false},
			Status:               "ACTIVE",
			ExpiresAt:            time.Now().Add(time.Hour),
		},
	}

	svc := NewAuctionService(nil, &fakeTransactor{querier: txQuerier})

	_, err := svc.PlaceAuctionBid(context.Background(), auctionID, bidderID, 125)
	if !errors.Is(err, domain.ErrBidTooLow) {
		t.Fatalf("PlaceAuctionBid() error = %v, want %v", err, domain.ErrBidTooLow)
	}

	wantCalls := []string{
		"wallet_for_update:" + bidderID.String(),
		"auction_for_bid:" + auctionID.String(),
	}
	if !reflect.DeepEqual(txQuerier.calls, wantCalls) {
		t.Fatalf("transaction calls = %#v, want %#v", txQuerier.calls, wantCalls)
	}
}

var _ database.TransactionQuerier = (*fakeAuctionTx)(nil)
