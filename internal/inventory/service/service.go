package service

import (
	"context"
	"strings"

	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"github.com/barbodimani81/Dragon-Market/internal/inventory/domain"
	"github.com/google/uuid"
)

type InventoryService struct {
	repo gen.Querier
}

func NewInventoryService(repo gen.Querier) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) GetItem(ctx context.Context, itemID uuid.UUID) (gen.Item, error) {
	item, err := s.repo.GetItem(ctx, itemID)
	if err != nil {
		return gen.Item{}, domain.ErrItemNotFound
	}
	return item, nil
}

func (s *InventoryService) MintItem(ctx context.Context, ownerID uuid.UUID, name string, rarity string) (gen.Item, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return gen.Item{}, domain.ErrEmptyItemName
	}

	cleanRarity := strings.ToUpper(strings.TrimSpace(rarity))
	if cleanRarity != "COMMON" && cleanRarity != "RARE" && cleanRarity != "LEGENDARY" {
		return gen.Item{}, domain.ErrInvalidRarity
	}

	return s.repo.CreateItem(ctx, gen.CreateItemParams{
		ID:      uuid.New(),
		OwnerID: ownerID,
		Name:    cleanName,
		Rarity:  cleanRarity,
	})
}

func (s *InventoryService) ListUserInventory(ctx context.Context, ownerID uuid.UUID) ([]gen.Item, error) {
	return s.repo.ListItemsByOwner(ctx, ownerID)
}
