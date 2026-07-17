package service

import (
	"context"

	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/barbodimani81/Dragon-Market/internal/marketplace/domain"
	"github.com/google/uuid"
)

type MarketplaceService struct {
	repo gen.Querier
}

func NewMarketplaceService(repo gen.Querier) *MarketplaceService {
	return &MarketplaceService{repo: repo}
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
		return gen.Listing{}, error(err) // maps cleanly back or can use inventory domain errors
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

// BuyListing handles the atomic exchange of money for an item at a fixed price
func (s *MarketplaceService) BuyListing(ctx context.Context, tx gen.Querier, listingID uuid.UUID, buyerID uuid.UUID) (gen.Listing, error) {
	// 1. Fetch the listing
	listing, err := tx.GetListing(ctx, listingID)
	if err != nil {
		return gen.Listing{}, domain.ErrListingNotFound
	}

	// 2. Business Rule: Ensure listing is still active
	if listing.Status != "ACTIVE" {
		return gen.Listing{}, domain.ErrListingNotActive
	}

	// 3. Business Rule: Prevent buying your own item
	if listing.SellerID == buyerID {
		return gen.Listing{}, domain.ErrSellerIsBuyer
	}

	// 4. Update the listing status to 'SOLD'
	updatedListing, err := tx.UpdateListingStatus(ctx, gen.UpdateListingStatusParams{
		ID:     listingID,
		Status: "SOLD",
	})
	if err != nil {
		return gen.Listing{}, err
	}

	// 5. Transfer item ownership to the buyer
	_, err = tx.TransferItemOwnership(ctx, gen.TransferItemOwnershipParams{
		ID:         listing.ItemID,
		NewOwnerID: buyerID,
	})
	if err != nil {
		return gen.Listing{}, err
	}

	// NOTE: In the orchestrator layer (like your controller/handler), you will pair this 
	// transaction block with calls to s.walletRepo.UpdateWalletBalance to deduct funds from the buyer
	// and add funds to the seller.

	return updatedListing, nil
}
