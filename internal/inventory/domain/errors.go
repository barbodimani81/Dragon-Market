package domain

import "errors"

var (
	ErrItemNotFound       = errors.New("item not found")
	ErrUnauthorizedAction = errors.New("user does not own this item")
	ErrInvalidRarity      = errors.New("invalid item rarity specified")
	ErrEmptyItemName      = errors.New("item name cannot be empty")
)
