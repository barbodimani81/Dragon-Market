-- name: GetWallet :one
SELECT id, user_id, available_balance, reserved_balance, currency, created_at, updated_at
FROM wallets
WHERE user_id = $1 LIMIT 1;

-- name: GetWalletForUpdate :one
SELECT id, user_id, available_balance, reserved_balance, currency, created_at, updated_at
FROM wallets
WHERE user_id = $1 LIMIT 1 FOR UPDATE;

-- name: CreateWallet :one
INSERT INTO wallets (id, user_id, available_balance, reserved_balance, currency)
VALUES ($1, $2, $3, $4, $5)
    RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: UpdateWalletBalance :one
UPDATE wallets
SET available_balance = available_balance + sqlc.arg(amount), updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.arg(user_id)
    RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: DecreaseWalletBalance :one
UPDATE wallets
SET available_balance = available_balance - sqlc.arg(amount), updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.arg(user_id) AND available_balance >= sqlc.arg(amount)
    RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: ReserveFunds :one
UPDATE wallets
SET available_balance = available_balance - sqlc.arg(amount),
    reserved_balance = reserved_balance + sqlc.arg(amount),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.arg(user_id) AND available_balance >= sqlc.arg(amount)
    RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: ReleaseFunds :one
UPDATE wallets
SET available_balance = available_balance + sqlc.arg(amount),
    reserved_balance = reserved_balance - sqlc.arg(amount),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.arg(user_id) AND reserved_balance >= sqlc.arg(amount)
    RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;
