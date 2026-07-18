package service

import (
	"context"
	"database/sql"
	"math"
	"time"

	"github.com/barbodimani81/Dragon-Market/internal/auction/domain"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	wdomain "github.com/barbodimani81/Dragon-Market/internal/wallet/domain"
	"github.com/barbodimani81/Dragon-Market/pkg/database"
	"github.com/google/uuid"
)

type AuctionService struct {
	repo gen.Querier
	tx   database.Transactor
}

func NewAuctionService(repo gen.Querier, tx database.Transactor) *AuctionService {
	return &AuctionService{repo: repo, tx: tx}
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

// PlaceAuctionBid processes incoming money commitment within a running transaction block.
func (s *AuctionService) PlaceAuctionBid(ctx context.Context, auctionID uuid.UUID, bidderID uuid.UUID, amount int64) (gen.PlaceBidRow, error) {
	var result gen.PlaceBidRow

	err := s.tx.WithinTransaction(ctx, func(q database.TransactionQuerier) error {
		bidderWallet, err := q.GetWalletForUpdate(ctx, bidderID)
		if err != nil {
			return wdomain.ErrWalletNotFound
		}

		auction, err := q.GetAuctionForBid(ctx, auctionID)
		if err != nil {
			return domain.ErrAuctionNotFound
		}

		if time.Now().After(auction.ExpiresAt) {
			return domain.ErrAuctionExpired
		}

		if auction.SellerID == bidderID {
			return domain.ErrSellerCannotBid
		}

		minRequiredBid := auction.StartPrice
		if auction.CurrentHighestBid.Valid && auction.CurrentHighestBid.Int64 > 0 {
			increment := int64(math.Ceil(float64(auction.CurrentHighestBid.Int64) * 0.05))
			if increment < 1 {
				increment = 1
			}
			minRequiredBid = auction.CurrentHighestBid.Int64 + increment
		}

		if amount < minRequiredBid {
			return domain.ErrBidTooLow
		}

		if bidderWallet.AvailableBalance < amount {
			return wdomain.ErrInsufficientFunds
		}

		if _, err := q.ReserveFunds(ctx, gen.ReserveFundsParams{
			UserID: bidderID,
			Amount: amount,
		}); err != nil {
			return err
		}

		if auction.CurrentHighestBidder.Valid && auction.CurrentHighestBid.Int64 > 0 {
			if _, err := q.ReleaseFunds(ctx, gen.ReleaseFundsParams{
				UserID: auction.CurrentHighestBidder.UUID,
				Amount: auction.CurrentHighestBid.Int64,
			}); err != nil {
				return err
			}
		}

		updated, err := q.PlaceBid(ctx, gen.PlaceBidParams{
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
		if err != nil {
			return err
		}

		result = updated
		return nil
	})

	return result, err
}

func (s *AuctionService) ListActiveAuctions(ctx context.Context) ([]gen.ListActiveAuctionsRow, error) {
	return s.repo.ListActiveAuctions(ctx)
}
