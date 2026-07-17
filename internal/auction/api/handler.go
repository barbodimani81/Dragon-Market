package api

import (
	"context"
	"time"

	"github.com/barbodimani81/Dragon-Market/api/wire"
	"github.com/barbodimani81/Dragon-Market/internal/auction/service"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AuctionHandler struct {
	db  gen.DBTX
	svc *service.AuctionService
}

func NewAuctionHandler(db gen.DBTX, svc *service.AuctionService) *AuctionHandler {
	return &AuctionHandler{db: db, svc: svc}
}

func (h *AuctionHandler) ListAuctions(ctx context.Context, request wire.ListAuctionsRequestObject) (wire.ListAuctionsResponseObject, error) {
	auctions, err := h.svc.ListActiveAuctions(ctx)
	if err != nil {
		return wire.ListAuctions500JSONResponse{Code: "INTERNAL_ERROR", Message: err.Error()}, nil
	}

	querier := gen.New(h.db)
	var response wire.ListAuctions200JSONResponse
	for _, a := range auctions {
		item, err := querier.GetItem(ctx, a.ItemID)
		if err != nil {
			return wire.ListAuctions500JSONResponse{Code: "INTERNAL_ERROR", Message: err.Error()}, nil
		}

		var currentHighestBid *int
		if a.CurrentHighestBid.Valid {
			v := int(a.CurrentHighestBid.Int64)
			currentHighestBid = &v
		}

		var currentHighestBidder *openapi_types.UUID
		if a.CurrentHighestBidder.Valid {
			v := openapi_types.UUID(a.CurrentHighestBidder.UUID)
			currentHighestBidder = &v
		}

		response = append(response, struct {
			AuctionId            openapi_types.UUID                       `json:"auction_id"`
			CreatedAt            time.Time                                `json:"created_at"`
			CurrentHighestBid    *int                                     `json:"current_highest_bid"`
			CurrentHighestBidder *openapi_types.UUID                      `json:"current_highest_bidder"`
			ExpiresAt            time.Time                                `json:"expires_at"`
			Item                 struct {
				CreatedAt time.Time                                      `json:"created_at"`
				ItemId    openapi_types.UUID                             `json:"item_id"`
				Name      string                                         `json:"name"`
				OwnerId   openapi_types.UUID                             `json:"owner_id"`
				Rarity    wire.ListAuctions200JSONResponseBodyItemRarity `json:"rarity"`
			} `json:"item"`
			SellerId   openapi_types.UUID                    `json:"seller_id"`
			StartPrice int                                   `json:"start_price"`
			Status     wire.ListAuctions200JSONResponseBodyStatus `json:"status"`
		}{
			AuctionId:            openapi_types.UUID(a.ID),
			CreatedAt:            a.CreatedAt,
			CurrentHighestBid:    currentHighestBid,
			CurrentHighestBidder: currentHighestBidder,
			ExpiresAt:            a.ExpiresAt,
			Item: struct {
				CreatedAt time.Time                                      `json:"created_at"`
				ItemId    openapi_types.UUID                             `json:"item_id"`
				Name      string                                         `json:"name"`
				OwnerId   openapi_types.UUID                             `json:"owner_id"`
				Rarity    wire.ListAuctions200JSONResponseBodyItemRarity `json:"rarity"`
			}{
				CreatedAt: item.CreatedAt,
				ItemId:    openapi_types.UUID(item.ID),
				Name:      item.Name,
				OwnerId:   openapi_types.UUID(item.OwnerID),
				Rarity:    wire.ListAuctions200JSONResponseBodyItemRarity(item.Rarity),
			},
			SellerId:   openapi_types.UUID(a.SellerID),
			StartPrice: int(a.StartPrice),
			Status:     wire.ListAuctions200JSONResponseBodyStatus(a.Status),
		})
	}

	return response, nil
}

func (h *AuctionHandler) CreateAuction(ctx context.Context, request wire.CreateAuctionRequestObject) (wire.CreateAuctionResponseObject, error) {
	mockUserID := uuid.Nil
	itemID := uuid.UUID(request.Body.ItemId)

	querier := gen.New(h.db)
	item, err := querier.GetItem(ctx, itemID)
	if err != nil {
		return wire.CreateAuction400JSONResponse{Code: "INVALID_ITEM", Message: err.Error()}, nil
	}

	auction, err := h.svc.CreateAuction(ctx, itemID, mockUserID, int64(request.Body.StartPrice), request.Body.DurationSeconds)
	if err != nil {
		return wire.CreateAuction400JSONResponse{Code: "CREATION_FAILED", Message: err.Error()}, nil
	}

	var currentHighestBid *int
	if auction.CurrentHighestBid.Valid {
		v := int(auction.CurrentHighestBid.Int64)
		currentHighestBid = &v
	}

	var currentHighestBidder *openapi_types.UUID
	if auction.CurrentHighestBidder.Valid {
		v := openapi_types.UUID(auction.CurrentHighestBidder.UUID)
		currentHighestBidder = &v
	}

	return wire.CreateAuction201JSONResponse{
		AuctionId:            openapi_types.UUID(auction.ID),
		CreatedAt:            auction.CreatedAt,
		CurrentHighestBid:    currentHighestBid,
		CurrentHighestBidder: currentHighestBidder,
		ExpiresAt:            auction.ExpiresAt,
		Item: struct {
			CreatedAt time.Time                                         `json:"created_at"`
			ItemId    openapi_types.UUID                                `json:"item_id"`
			Name      string                                            `json:"name"`
			OwnerId   openapi_types.UUID                                `json:"owner_id"`
			Rarity    wire.CreateAuction201JSONResponseBodyItemRarity   `json:"rarity"`
		}{
			CreatedAt: item.CreatedAt,
			ItemId:    openapi_types.UUID(item.ID),
			Name:      item.Name,
			OwnerId:   openapi_types.UUID(item.OwnerID),
			Rarity:    wire.CreateAuction201JSONResponseBodyItemRarity(item.Rarity),
		},
		SellerId:   openapi_types.UUID(auction.SellerID),
		StartPrice: int(auction.StartPrice),
		Status:     wire.CreateAuction201JSONResponseBodyStatus(auction.Status),
	}, nil
}

func (h *AuctionHandler) PlaceBid(ctx context.Context, request wire.PlaceBidRequestObject) (wire.PlaceBidResponseObject, error) {
	mockUserID := uuid.Nil
	listingID := uuid.UUID(request.Id)

	querier := gen.New(h.db)
	bidRow, err := h.svc.PlaceAuctionBid(ctx, querier, listingID, mockUserID, int64(request.Body.Amount))
	if err != nil {
		return wire.PlaceBid400JSONResponse{Code: "BID_FAILED", Message: err.Error()}, nil
	}

	item, err := querier.GetItem(ctx, bidRow.ItemID)
	if err != nil {
		return wire.PlaceBid400JSONResponse{Code: "BID_FAILED", Message: err.Error()}, nil
	}

	var currentHighestBid *int
	if bidRow.CurrentHighestBid.Valid {
		v := int(bidRow.CurrentHighestBid.Int64)
		currentHighestBid = &v
	}

	var currentHighestBidder *openapi_types.UUID
	if bidRow.CurrentHighestBidder.Valid {
		v := openapi_types.UUID(bidRow.CurrentHighestBidder.UUID)
		currentHighestBidder = &v
	}

	return wire.PlaceBid201JSONResponse{
		AuctionId:            openapi_types.UUID(bidRow.ID),
		CreatedAt:            bidRow.CreatedAt,
		CurrentHighestBid:    currentHighestBid,
		CurrentHighestBidder: currentHighestBidder,
		ExpiresAt:            bidRow.ExpiresAt,
		Item: struct {
			CreatedAt time.Time                                    `json:"created_at"`
			ItemId    openapi_types.UUID                           `json:"item_id"`
			Name      string                                       `json:"name"`
			OwnerId   openapi_types.UUID                           `json:"owner_id"`
			Rarity    wire.PlaceBid201JSONResponseBodyItemRarity   `json:"rarity"`
		}{
			CreatedAt: item.CreatedAt,
			ItemId:    openapi_types.UUID(item.ID),
			Name:      item.Name,
			OwnerId:   openapi_types.UUID(item.OwnerID),
			Rarity:    wire.PlaceBid201JSONResponseBodyItemRarity(item.Rarity),
		},
		SellerId:   openapi_types.UUID(bidRow.SellerID),
		StartPrice: int(bidRow.StartPrice),
		Status:     wire.PlaceBid201JSONResponseBodyStatus(bidRow.Status),
	}, nil
}
