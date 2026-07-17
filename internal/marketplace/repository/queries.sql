-- name: GetListing :one
SELECT id, item_id, seller_id, price, status, created_at, updated_at 
FROM listings 
WHERE id = $1 LIMIT 1;

-- name: ListActiveListings :many
SELECT id, item_id, seller_id, price, status, created_at, updated_at
FROM listings
WHERE status = 'ACTIVE'
ORDER BY created_at DESC;

-- name: CreateListing :one
INSERT INTO listings (id, item_id, seller_id, price, status) 
VALUES (sqlc.arg(id), sqlc.arg(item_id), sqlc.arg(seller_id), sqlc.arg(price), 'ACTIVE') 
RETURNING id, item_id, seller_id, price, status, created_at, updated_at;

-- name: UpdateListingStatus :one
UPDATE listings
SET status = sqlc.arg(status), updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING id, item_id, seller_id, price, status, created_at, updated_at;
