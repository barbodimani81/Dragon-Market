package api

import (
	"context"
	"time"

	"github.com/barbodimani81/Dragon-Market/api/wire"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/barbodimani81/Dragon-Market/internal/marketplace/service"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type MarketplaceHandler struct {
	db  gen.DBTX
	svc *service.MarketplaceService
}

func NewMarketplaceHandler(db gen.DBTX, svc *service.MarketplaceService) *MarketplaceHandler {
	return &MarketplaceHandler{db: db, svc: svc}
}

func (h *MarketplaceHandler) ListListings(ctx context.Context, request wire.ListListingsRequestObject) (wire.ListListingsResponseObject, error) {
	listings, err := h.svc.ListActiveListings(ctx)
	if err != nil {
		return wire.ListListings500JSONResponse{Code: "INTERNAL_ERROR", Message: err.Error()}, nil
	}

	querier := gen.New(h.db)
	var response wire.ListListings200JSONResponse
	for _, listing := range listings {
		item, err := querier.GetItem(ctx, listing.ItemID)
		if err != nil {
			return wire.ListListings500JSONResponse{Code: "INTERNAL_ERROR", Message: err.Error()}, nil
		}

		response = append(response, struct {
			CreatedAt time.Time `json:"created_at"`
			Item      struct {
				CreatedAt time.Time                                      `json:"created_at"`
				ItemId    openapi_types.UUID                             `json:"item_id"`
				Name      string                                         `json:"name"`
				OwnerId   openapi_types.UUID                             `json:"owner_id"`
				Rarity    wire.ListListings200JSONResponseBodyItemRarity `json:"rarity"`
			} `json:"item"`
			ListingId openapi_types.UUID                           `json:"listing_id"`
			Price     int                                          `json:"price"`
			SellerId  openapi_types.UUID                           `json:"seller_id"`
			Status    wire.ListListings200JSONResponseBodyStatus   `json:"status"`
		}{
			CreatedAt: listing.CreatedAt,
			Item: struct {
				CreatedAt time.Time                                      `json:"created_at"`
				ItemId    openapi_types.UUID                             `json:"item_id"`
				Name      string                                         `json:"name"`
				OwnerId   openapi_types.UUID                             `json:"owner_id"`
				Rarity    wire.ListListings200JSONResponseBodyItemRarity `json:"rarity"`
			}{
				CreatedAt: item.CreatedAt,
				ItemId:    openapi_types.UUID(item.ID),
				Name:      item.Name,
				OwnerId:   openapi_types.UUID(item.OwnerID),
				Rarity:    wire.ListListings200JSONResponseBodyItemRarity(item.Rarity),
			},
			ListingId: openapi_types.UUID(listing.ID),
			Price:     int(listing.Price),
			SellerId:  openapi_types.UUID(listing.SellerID),
			Status:    wire.ListListings200JSONResponseBodyStatus(listing.Status),
		})
	}

	return response, nil
}

func (h *MarketplaceHandler) CreateListing(ctx context.Context, request wire.CreateListingRequestObject) (wire.CreateListingResponseObject, error) {
	mockUserID := uuid.Nil
	itemID := uuid.UUID(request.Body.ItemId)

	querier := gen.New(h.db)
	item, err := querier.GetItem(ctx, itemID)
	if err != nil {
		return wire.CreateListing400JSONResponse{Code: "INVALID_ITEM", Message: err.Error()}, nil
	}

	listing, err := h.svc.CreateListing(ctx, itemID, mockUserID, int64(request.Body.Price))
	if err != nil {
		return wire.CreateListing400JSONResponse{Code: "CREATION_FAILED", Message: err.Error()}, nil
	}

	return wire.CreateListing201JSONResponse{
		CreatedAt: listing.CreatedAt,
		Item: struct {
			CreatedAt time.Time                                        `json:"created_at"`
			ItemId    openapi_types.UUID                               `json:"item_id"`
			Name      string                                           `json:"name"`
			OwnerId   openapi_types.UUID                               `json:"owner_id"`
			Rarity    wire.CreateListing201JSONResponseBodyItemRarity  `json:"rarity"`
		}{
			CreatedAt: item.CreatedAt,
			ItemId:    openapi_types.UUID(item.ID),
			Name:      item.Name,
			OwnerId:   openapi_types.UUID(item.OwnerID),
			Rarity:    wire.CreateListing201JSONResponseBodyItemRarity(item.Rarity),
		},
		ListingId: openapi_types.UUID(listing.ID),
		Price:     int(listing.Price),
		SellerId:  openapi_types.UUID(listing.SellerID),
		Status:    wire.CreateListing201JSONResponseBodyStatus(listing.Status),
	}, nil
}

func (h *MarketplaceHandler) BuyListing(ctx context.Context, request wire.BuyListingRequestObject) (wire.BuyListingResponseObject, error) {
	mockUserID := uuid.Nil
	listingID := uuid.UUID(request.Id)

	listing, err := h.svc.BuyListing(ctx, listingID, mockUserID)
	if err != nil {
		return wire.BuyListing400JSONResponse{Code: "PURCHASE_FAILED", Message: err.Error()}, nil
	}

	querier := gen.New(h.db)
	item, err := querier.GetItem(ctx, listing.ItemID)
	if err != nil {
		return wire.BuyListing400JSONResponse{Code: "PURCHASE_FAILED", Message: err.Error()}, nil
	}

	return wire.BuyListing200JSONResponse{
		CreatedAt: listing.CreatedAt,
		Item: struct {
			CreatedAt time.Time                                      `json:"created_at"`
			ItemId    openapi_types.UUID                             `json:"item_id"`
			Name      string                                         `json:"name"`
			OwnerId   openapi_types.UUID                             `json:"owner_id"`
			Rarity    wire.BuyListing200JSONResponseBodyItemRarity   `json:"rarity"`
		}{
			CreatedAt: item.CreatedAt,
			ItemId:    openapi_types.UUID(item.ID),
			Name:      item.Name,
			OwnerId:   openapi_types.UUID(item.OwnerID),
			Rarity:    wire.BuyListing200JSONResponseBodyItemRarity(item.Rarity),
		},
		ListingId: openapi_types.UUID(listing.ID),
		Price:     int(listing.Price),
		SellerId:  openapi_types.UUID(listing.SellerID),
		Status:    wire.BuyListing200JSONResponseBodyStatus(listing.Status),
	}, nil
}
