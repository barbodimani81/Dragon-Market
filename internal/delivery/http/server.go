package http

import (
	"github.com/barbodimani81/Dragon-Market/api/wire"
	auctionAPI "github.com/barbodimani81/Dragon-Market/internal/auction/api"
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
