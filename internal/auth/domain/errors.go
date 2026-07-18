package domain

import "errors"

var (
	// ErrUserAlreadyExists is thrown when a username is already taken
	ErrUserAlreadyExists = errors.New("username is already registered in the system")

	// ErrInvalidCredentials is thrown when password checks fail during login
	ErrInvalidCredentials = errors.New("invalid username or password credentials provided")

	// ErrInternalServer is a fallback for unhandled database glitches
	ErrInternalServer = errors.New("an internal authentication error occurred")
)
