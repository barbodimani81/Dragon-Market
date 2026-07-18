package api

import (
	"context"
	"errors"

	"github.com/barbodimani81/Dragon-Market/api/wire"
	"github.com/barbodimani81/Dragon-Market/internal/auth/domain"
	"github.com/barbodimani81/Dragon-Market/internal/auth/service"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterUser handles user signup parsing and processing
func (h *AuthHandler) RegisterUser(ctx context.Context, request wire.RegisterUserRequestObject) (wire.RegisterUserResponseObject, error) {
	_, err := h.authService.SignUp(ctx, request.Body.Username, request.Body.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			// Returns the generated 409 wrapper struct
			return wire.RegisterUser409Response{}, nil
		default:
			// Returns the generated 400 wrapper struct
			return wire.RegisterUser400Response{}, nil
		}
	}

	return wire.RegisterUser201JSONResponse{
		Message: "User created successfully",
	}, nil
}

// LoginUser handles credentials matching and creates a session token mock payload
func (h *AuthHandler) LoginUser(ctx context.Context, request wire.LoginUserRequestObject) (wire.LoginUserResponseObject, error) {
	user, err := h.authService.LogIn(ctx, request.Body.Username, request.Body.Password)
	if err != nil {
		// Returns the generated 401 wrapper struct
		return wire.LoginUser401Response{}, nil
	}

	mockToken := "mock-jwt-session-token-for-" + user.Username

	return wire.LoginUser200JSONResponse{
		Message: "Authentication successful",
		Token:   mockToken,
		UserId:  openapi_types.UUID(user.ID),
	}, nil
}
