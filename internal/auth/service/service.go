package service

import (
	"context"

	"github.com/barbodimani81/Dragon-Market/internal/auth/domain"
	gen "github.com/barbodimani81/Dragon-Market/internal/database/gen"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo gen.Querier
}

func NewAuthService(repo gen.Querier) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

// SignUp takes raw inputs, hashes the secret safely, and commits the profile to the DB
func (s *AuthService) SignUp(ctx context.Context, username, password string) (gen.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return gen.User{}, domain.ErrInternalServer
	}

	user, err := s.repo.CreateUser(ctx, gen.CreateUserParams{
		Username:     username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		// Detect typical unique constraint violation from postgres drivers
		return gen.User{}, domain.ErrUserAlreadyExists
	}

	return user, nil
}

// LogIn checks credentials against saved hashes and verifies the session match
func (s *AuthService) LogIn(ctx context.Context, username, password string) (gen.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return gen.User{}, domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return gen.User{}, domain.ErrInvalidCredentials
	}

	return user, nil
}
