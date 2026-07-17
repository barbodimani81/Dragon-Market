package api

import (
	"context"
	"time"

	"github.com/barbodimani81/Dragon-Market/api/wire"
	"github.com/barbodimani81/Dragon-Market/internal/inventory/service"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type InventoryHandler struct {
	svc *service.InventoryService
}

func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (h *InventoryHandler) ListInventory(ctx context.Context, request wire.ListInventoryRequestObject) (wire.ListInventoryResponseObject, error) {
	mockUserID := uuid.Nil

	items, err := h.svc.ListUserInventory(ctx, mockUserID)
	if err != nil {
		return wire.ListInventory500JSONResponse{Code: "INTERNAL_ERROR", Message: err.Error()}, nil
	}

	var response wire.ListInventory200JSONResponse
	for _, item := range items {
		response = append(response, struct {
			CreatedAt time.Time                              `json:"created_at"`
			ItemId    openapi_types.UUID                     `json:"item_id"`
			Name      string                                 `json:"name"`
			OwnerId   openapi_types.UUID                     `json:"owner_id"`
			Rarity    wire.ListInventory200JSONResponseBodyRarity `json:"rarity"`
		}{
			CreatedAt: item.CreatedAt,
			ItemId:    openapi_types.UUID(item.ID),
			Name:      item.Name,
			OwnerId:   openapi_types.UUID(item.OwnerID),
			Rarity:    wire.ListInventory200JSONResponseBodyRarity(item.Rarity),
		})
	}

	return response, nil
}
