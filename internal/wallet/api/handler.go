package api

import (
	"context"

	"github.com/barbodimani81/Dragon-Market/api/wire"
	"github.com/barbodimani81/Dragon-Market/internal/wallet/service"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{svc: svc}
}

func (h *WalletHandler) GetWallet(ctx context.Context, request wire.GetWalletRequestObject) (wire.GetWalletResponseObject, error) {
	// TODO: Replace with authenticated user ID from context middleware
	mockUserID := uuid.Nil

	wallet, err := h.svc.GetWallet(ctx, mockUserID)
	if err != nil {
		return wire.GetWallet500JSONResponse{Code: "INTERNAL_ERROR", Message: err.Error()}, nil
	}

	return wire.GetWallet200JSONResponse{
		WalletId:         openapi_types.UUID(wallet.ID),
		UserId:           openapi_types.UUID(wallet.UserID),
		AvailableBalance: int(wallet.AvailableBalance),
		ReservedBalance:  int(wallet.ReservedBalance),
		TotalBalance:     int(wallet.AvailableBalance + wallet.ReservedBalance),
		Currency:         wallet.Currency,
	}, nil
}

func (h *WalletHandler) DepositFunds(ctx context.Context, request wire.DepositFundsRequestObject) (wire.DepositFundsResponseObject, error) {
	mockUserID := uuid.Nil
	amount := int64(request.Body.Amount)

	wallet, err := h.svc.Deposit(ctx, mockUserID, amount)
	if err != nil {
		return wire.DepositFunds400JSONResponse{Code: "DEPOSIT_FAILED", Message: err.Error()}, nil
	}

	return wire.DepositFunds200JSONResponse{
		WalletId:         openapi_types.UUID(wallet.ID),
		UserId:           openapi_types.UUID(wallet.UserID),
		AvailableBalance: int(wallet.AvailableBalance),
		ReservedBalance:  int(wallet.ReservedBalance),
		TotalBalance:     int(wallet.AvailableBalance + wallet.ReservedBalance),
		Currency:         wallet.Currency,
	}, nil
}
