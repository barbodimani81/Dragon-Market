package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/barbodimani81/Dragon-Market/internal/auction/domain"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/google/uuid"
)

type AuctionService struct {
	repo gen.Querier
}

func NewAuctionService(repo gen.Querier) *AuctionService {
	return &AuctionService{repo: repo}
}

// CreateAuction instantiates an auction window for legendary drops
func (s *AuctionService) CreateAuction(ctx context.Context, itemID uuid.UUID, sellerID uuid.UUID, startPrice int64, durationSeconds int) (gen.CreateAuctionRow, error) {
	if durationSeconds < 60 || durationSeconds > 604800 {
		return gen.CreateAuctionRow{}, domain.ErrInvalidDuration
	}

	item, err := s.repo.GetItem(ctx, itemID)
	if err != nil {
		return gen.CreateAuctionRow{}, err
	}

	if item.Rarity != "LEGENDARY" {
		return gen.CreateAuctionRow{}, domain.ErrInvalidRarity
	}

	expiresAt := time.Now().Add(time.Duration(durationSeconds) * time.Second)

	return s.repo.CreateAuction(ctx, gen.CreateAuctionParams{
		ID:         uuid.New(),
		ItemID:     itemID,
		SellerID:   sellerID,
		StartPrice: startPrice,
		ExpiresAt:  expiresAt,
	})
}

// PlaceAuctionBid processes incoming money commitment within a running transaction block
func (s *AuctionService) PlaceAuctionBid(ctx context.Context, tx gen.Querier, auctionID uuid.UUID, bidderID uuid.UUID, amount int64) (gen.PlaceBidRow, error) {
	// 1. Fetch running auction details utilizing safe row locks
	auction, err := tx.GetAuctionForBid(ctx, auctionID)
	if err != nil {
		return gen.PlaceBidRow{}, domain.ErrAuctionNotFound
	}

	// 2. Guard: Time validation check
	if time.Now().After(auction.ExpiresAt) {
		return gen.PlaceBidRow{}, domain.ErrAuctionExpired
	}

	// 3. Guard: Seller restriction rule
	if auction.SellerID == bidderID {
		return gen.PlaceBidRow{}, domain.ErrSellerCannotBid
	}

	// 4. Guard: High-water bid valuation check
	minRequiredBid := auction.StartPrice
	if auction.CurrentHighestBid.Valid && auction.CurrentHighestBid.Int64 > 0 {
		minRequiredBid = auction.CurrentHighestBid.Int64 + 1
	}

	if amount < minRequiredBid {
		return gen.PlaceBidRow{}, domain.ErrBidTooLow
	}

	// 5. Apply the updated state row value with correct nullable types wrapped
	return tx.PlaceBid(ctx, gen.PlaceBidParams{
		ID: auctionID,
		Amount: sql.NullInt64{
			Int64: amount,
			Valid: true,
		},
		BidderID: uuid.NullUUID{
			UUID:  bidderID,
			Valid: true,
		},
	})
}

func (s *AuctionService) ListActiveAuctions(ctx context.Context) ([]gen.ListActiveAuctionsRow, error) {
	return s.repo.ListActiveAuctions(ctx)
}
