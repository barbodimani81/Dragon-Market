package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	wdomain "github.com/barbodimani81/Dragon-Market/internal/wallet/domain"
	"github.com/barbodimani81/Dragon-Market/pkg/database"
	"github.com/google/uuid"
)

type fakeMarketplaceTx struct {
	calls         []string
	listing       gen.Listing
	buyerWallet    gen.Wallet
	sellerWallet   gen.Wallet
	updatedListing gen.Listing
}

type fakeTransactor struct {
	querier database.TransactionQuerier
	calls   int
}

func (f *fakeTransactor) WithinTransaction(ctx context.Context, fn func(q database.TransactionQuerier) error) error {
	f.calls++
	return fn(f.querier)
}

func (f *fakeMarketplaceTx) GetWalletForUpdate(ctx context.Context, userID uuid.UUID) (gen.Wallet, error) {
	f.calls = append(f.calls, "wallet_for_update:"+userID.String())
	switch userID {
	case f.buyerWallet.UserID:
		return f.buyerWallet, nil
	case f.sellerWallet.UserID:
		return f.sellerWallet, nil
	default:
		return gen.Wallet{}, nil
	}
}

func (f *fakeMarketplaceTx) GetListingForUpdate(ctx context.Context, id uuid.UUID) (gen.Listing, error) {
	f.calls = append(f.calls, "listing_for_update:"+id.String())
	return f.listing, nil
}

func (f *fakeMarketplaceTx) GetAuctionForBid(ctx context.Context, id uuid.UUID) (gen.GetAuctionForBidRow, error) {
	return gen.GetAuctionForBidRow{}, nil
}

func (f *fakeMarketplaceTx) ReserveFunds(ctx context.Context, arg gen.ReserveFundsParams) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeMarketplaceTx) ReleaseFunds(ctx context.Context, arg gen.ReleaseFundsParams) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeMarketplaceTx) DecreaseWalletBalance(ctx context.Context, arg gen.DecreaseWalletBalanceParams) (gen.Wallet, error) {
	f.calls = append(f.calls, "decrease_wallet_balance")
	return f.buyerWallet, nil
}

func (f *fakeMarketplaceTx) UpdateWalletBalance(ctx context.Context, arg gen.UpdateWalletBalanceParams) (gen.Wallet, error) {
	f.calls = append(f.calls, "update_wallet_balance")
	return f.sellerWallet, nil
}

func (f *fakeMarketplaceTx) PlaceBid(ctx context.Context, arg gen.PlaceBidParams) (gen.PlaceBidRow, error) {
	return gen.PlaceBidRow{}, nil
}

func (f *fakeMarketplaceTx) UpdateListingStatus(ctx context.Context, arg gen.UpdateListingStatusParams) (gen.Listing, error) {
	f.calls = append(f.calls, "update_listing_status")
	return f.updatedListing, nil
}

func (f *fakeMarketplaceTx) TransferItemOwnership(ctx context.Context, arg gen.TransferItemOwnershipParams) (gen.Item, error) {
	f.calls = append(f.calls, "transfer_item_ownership")
	return gen.Item{ID: arg.ID, OwnerID: arg.NewOwnerID}, nil
}

func (f *fakeMarketplaceTx) GetItem(ctx context.Context, id uuid.UUID) (gen.Item, error) {
	return gen.Item{}, nil
}

func (f *fakeMarketplaceTx) GetListing(ctx context.Context, id uuid.UUID) (gen.Listing, error) {
	return gen.Listing{}, nil
}

func (f *fakeMarketplaceTx) GetWallet(ctx context.Context, userID uuid.UUID) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeMarketplaceTx) CreateAuction(ctx context.Context, arg gen.CreateAuctionParams) (gen.CreateAuctionRow, error) {
	return gen.CreateAuctionRow{}, nil
}

func (f *fakeMarketplaceTx) CreateItem(ctx context.Context, arg gen.CreateItemParams) (gen.Item, error) {
	return gen.Item{}, nil
}

func (f *fakeMarketplaceTx) CreateListing(ctx context.Context, arg gen.CreateListingParams) (gen.Listing, error) {
	return gen.Listing{}, nil
}

func (f *fakeMarketplaceTx) CreateWallet(ctx context.Context, arg gen.CreateWalletParams) (gen.Wallet, error) {
	return gen.Wallet{}, nil
}

func (f *fakeMarketplaceTx) ListActiveAuctions(ctx context.Context) ([]gen.ListActiveAuctionsRow, error) {
	return nil, nil
}

func (f *fakeMarketplaceTx) ListActiveListings(ctx context.Context) ([]gen.Listing, error) {
	return nil, nil
}

func (f *fakeMarketplaceTx) ListItemsByOwner(ctx context.Context, ownerID uuid.UUID) ([]gen.Item, error) {
	return nil, nil
}

var _ database.TransactionQuerier = (*fakeMarketplaceTx)(nil)

func TestBuyListing_TransactionFlow(t *testing.T) {
	buyerID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	sellerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	listingID := uuid.MustParse("00000000-0000-0000-0000-000000000020")
	itemID := uuid.MustParse("00000000-0000-0000-0000-000000000021")

	txQuerier := &fakeMarketplaceTx{
		listing: gen.Listing{
			ID:        listingID,
			ItemID:    itemID,
			SellerID:  sellerID,
			Price:     150,
			Status:    "ACTIVE",
			CreatedAt: time.Now().Add(-time.Hour),
		},
		buyerWallet: gen.Wallet{
			ID:               uuid.New(),
			UserID:           buyerID,
			AvailableBalance: 500,
		},
		sellerWallet: gen.Wallet{
			ID:               uuid.New(),
			UserID:           sellerID,
			AvailableBalance: 0,
		},
		updatedListing: gen.Listing{
			ID:        listingID,
			ItemID:    itemID,
			SellerID:  sellerID,
			Price:     150,
			Status:    "SOLD",
			CreatedAt: time.Now().Add(-time.Hour),
		},
	}

	svc := NewMarketplaceService(nil, &fakeTransactor{querier: txQuerier})

	got, err := svc.BuyListing(context.Background(), listingID, buyerID)
	if err != nil {
		t.Fatalf("BuyListing() error = %v", err)
	}

	if !reflect.DeepEqual(got, txQuerier.updatedListing) {
		t.Fatalf("BuyListing() got = %#v, want %#v", got, txQuerier.updatedListing)
	}

	wantCalls := []string{
		"listing_for_update:" + listingID.String(),
		"wallet_for_update:" + sellerID.String(),
		"wallet_for_update:" + buyerID.String(),
		"decrease_wallet_balance",
		"update_wallet_balance",
		"transfer_item_ownership",
		"update_listing_status",
	}
	if !reflect.DeepEqual(txQuerier.calls, wantCalls) {
		t.Fatalf("transaction calls = %#v, want %#v", txQuerier.calls, wantCalls)
	}
}

func TestBuyListing_InsufficientFundsStopsBeforeMutation(t *testing.T) {
	buyerID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	sellerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	listingID := uuid.MustParse("00000000-0000-0000-0000-000000000020")
	itemID := uuid.MustParse("00000000-0000-0000-0000-000000000021")

	txQuerier := &fakeMarketplaceTx{
		listing: gen.Listing{
			ID:        listingID,
			ItemID:    itemID,
			SellerID:  sellerID,
			Price:     150,
			Status:    "ACTIVE",
			CreatedAt: time.Now().Add(-time.Hour),
		},
		buyerWallet: gen.Wallet{
			ID:               uuid.New(),
			UserID:           buyerID,
			AvailableBalance: 100,
		},
		sellerWallet: gen.Wallet{
			ID:               uuid.New(),
			UserID:           sellerID,
			AvailableBalance: 0,
		},
	}

	svc := NewMarketplaceService(nil, &fakeTransactor{querier: txQuerier})

	_, err := svc.BuyListing(context.Background(), listingID, buyerID)
	if !errors.Is(err, wdomain.ErrInsufficientFunds) {
		t.Fatalf("BuyListing() error = %v, want %v", err, wdomain.ErrInsufficientFunds)
	}

	wantCalls := []string{
		"listing_for_update:" + listingID.String(),
		"wallet_for_update:" + sellerID.String(),
		"wallet_for_update:" + buyerID.String(),
	}
	if !reflect.DeepEqual(txQuerier.calls, wantCalls) {
		t.Fatalf("transaction calls = %#v, want %#v", txQuerier.calls, wantCalls)
	}
}
