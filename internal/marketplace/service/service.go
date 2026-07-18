package service

import (
	"context"
	"sort"

	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/barbodimani81/Dragon-Market/internal/marketplace/domain"
	wdomain "github.com/barbodimani81/Dragon-Market/internal/wallet/domain"
	"github.com/barbodimani81/Dragon-Market/pkg/database"
	"github.com/google/uuid"
)

type MarketplaceService struct {
	repo gen.Querier
	tx   database.Transactor
}

func NewMarketplaceService(repo gen.Querier, tx database.Transactor) *MarketplaceService {
	return &MarketplaceService{repo: repo, tx: tx}
}

// GetListing fetches a specific active marketplace listing
func (s *MarketplaceService) GetListing(ctx context.Context, listingID uuid.UUID) (gen.Listing, error) {
	listing, err := s.repo.GetListing(ctx, listingID)
	if err != nil {
		return gen.Listing{}, domain.ErrListingNotFound
	}
	return listing, nil
}

// CreateListing processes and saves a brand new fixed-price item listing
func (s *MarketplaceService) CreateListing(ctx context.Context, itemID uuid.UUID, sellerID uuid.UUID, price int64) (gen.Listing, error) {
	if price <= 0 {
		return gen.Listing{}, domain.ErrNegativePrice
	}

	// 1. Fetch item from the repository to check its rarity constraints
	item, err := s.repo.GetItem(ctx, itemID)
	if err != nil {
		return gen.Listing{}, err
	}

	// 2. Business Rule: Guard against Legendary items on the fixed marketplace
	if item.Rarity == "LEGENDARY" {
		return gen.Listing{}, domain.ErrInvalidListingRarity
	}

	// 3. Save listing to database
	return s.repo.CreateListing(ctx, gen.CreateListingParams{
		ID:       uuid.New(),
		ItemID:   itemID,
		SellerID: sellerID,
		Price:    price,
	})
}

func (s *MarketplaceService) ListActiveListings(ctx context.Context) ([]gen.Listing, error) {
	return s.repo.ListActiveListings(ctx)
}

// BuyListing handles the atomic exchange of money for an item at a fixed price.
func (s *MarketplaceService) BuyListing(ctx context.Context, listingID uuid.UUID, buyerID uuid.UUID) (gen.Listing, error) {
	var updatedListing gen.Listing

	err := s.tx.WithinTransaction(ctx, func(q database.TransactionQuerier) error {
		listing, err := q.GetListingForUpdate(ctx, listingID)
		if err != nil {
			return domain.ErrListingNotFound
		}

		if listing.Status != "ACTIVE" {
			return domain.ErrListingNotActive
		}

		if listing.SellerID == buyerID {
			return domain.ErrSellerIsBuyer
		}

		walletIDs := []uuid.UUID{buyerID, listing.SellerID}
		sort.Slice(walletIDs, func(i, j int) bool {
			return walletIDs[i].String() < walletIDs[j].String()
		})

		lockedWallets := make(map[uuid.UUID]gen.Wallet, 2)
		for _, walletID := range walletIDs {
			wallet, err := q.GetWalletForUpdate(ctx, walletID)
			if err != nil {
				return wdomain.ErrWalletNotFound
			}
			lockedWallets[walletID] = wallet
		}

		buyerWallet := lockedWallets[buyerID]
		if buyerWallet.AvailableBalance < listing.Price {
			return wdomain.ErrInsufficientFunds
		}

		if _, err := q.DecreaseWalletBalance(ctx, gen.DecreaseWalletBalanceParams{
			UserID: buyerID,
			Amount: listing.Price,
		}); err != nil {
			return wdomain.ErrInsufficientFunds
		}

		if _, err := q.UpdateWalletBalance(ctx, gen.UpdateWalletBalanceParams{
			UserID: listing.SellerID,
			Amount: listing.Price,
		}); err != nil {
			return err
		}

		if _, err := q.TransferItemOwnership(ctx, gen.TransferItemOwnershipParams{
			ID:         listing.ItemID,
			NewOwnerID: buyerID,
		}); err != nil {
			return err
		}

		updated, err := q.UpdateListingStatus(ctx, gen.UpdateListingStatusParams{
			ID:     listingID,
			Status: "SOLD",
		})
		if err != nil {
			return err
		}

		updatedListing = updated
		return nil
	})

	return updatedListing, err
}
