-- name: GetWallet :one
SELECT id, user_id, available_balance, reserved_balance, currency, created_at, updated_at 
FROM wallets 
WHERE user_id = $1 LIMIT 1;

-- name: CreateWallet :one
INSERT INTO wallets (id, user_id, available_balance, reserved_balance, currency)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: UpdateWalletBalance :one
UPDATE wallets 
SET available_balance = available_balance + $2, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: ReserveFunds :one
-- Moves money from available_balance to reserved_balance (e.g., when placing a bid)
UPDATE wallets
SET available_balance = available_balance - $2,
    reserved_balance = reserved_balance + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND available_balance >= $2
RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;

-- name: ReleaseFunds :one
-- Moves money back from reserved_balance to available_balance (e.g., when outbid)
UPDATE wallets
SET available_balance = available_balance + $2,
    reserved_balance = reserved_balance - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND reserved_balance >= $2
RETURNING id, user_id, available_balance, reserved_balance, currency, created_at, updated_at;
