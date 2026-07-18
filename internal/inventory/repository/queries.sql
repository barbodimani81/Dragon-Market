-- name: GetItem :one
SELECT id, owner_id, name, rarity, created_at, updated_at 
FROM items 
WHERE id = $1 LIMIT 1;

-- name: CreateItem :one
INSERT INTO items (id, owner_id, name, rarity) 
VALUES (sqlc.arg(id), sqlc.arg(owner_id), sqlc.arg(name), sqlc.arg(rarity)) 
RETURNING id, owner_id, name, rarity, created_at, updated_at;

-- name: TransferItemOwnership :one
UPDATE items 
SET owner_id = sqlc.arg(new_owner_id), updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, owner_id, name, rarity, created_at, updated_at;

-- name: ListItemsByOwner :many
SELECT id, owner_id, name, rarity, created_at, updated_at
FROM items
WHERE owner_id = $1
ORDER BY created_at DESC;
