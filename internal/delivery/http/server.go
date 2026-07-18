package http

import (
	"context"
	"github.com/barbodimani81/Dragon-Market/api/wire"
	auctionAPI "github.com/barbodimani81/Dragon-Market/internal/auction/api"
	"github.com/barbodimani81/Dragon-Market/internal/auth/api"
	inventoryAPI "github.com/barbodimani81/Dragon-Market/internal/inventory/api"
	marketplaceAPI "github.com/barbodimani81/Dragon-Market/internal/marketplace/api"
	walletAPI "github.com/barbodimani81/Dragon-Market/internal/wallet/api"
)

var _ wire.StrictServerInterface = (*ModularStrictServer)(nil)

type ModularStrictServer struct {
	*walletAPI.WalletHandler
	*inventoryAPI.InventoryHandler
	*marketplaceAPI.MarketplaceHandler
	*auctionAPI.AuctionHandler
	authHandler *api.AuthHandler
}

func NewModularStrictServer(
	w *walletAPI.WalletHandler,
	i *inventoryAPI.InventoryHandler,
	m *marketplaceAPI.MarketplaceHandler,
	a *auctionAPI.AuctionHandler,
) *ModularStrictServer {
	return &ModularStrictServer{
		WalletHandler:      w,
		InventoryHandler:   i,
		MarketplaceHandler: m,
		AuctionHandler:     a,
	}
}

// RegisterUser handles POST /auth/signup
func (s *ModularStrictServer) RegisterUser(ctx context.Context, request wire.RegisterUserRequestObject) (wire.RegisterUserResponseObject, error) {
	return nil, nil
}

// LoginUser handles POST /auth/login
func (s *ModularStrictServer) LoginUser(ctx context.Context, request wire.LoginUserRequestObject) (wire.LoginUserResponseObject, error) {
	return nil, nil
}
