-- name: GetAuctionForBid :one
SELECT id, item_id, seller_id, start_price, current_highest_bid, current_highest_bidder, status, expires_at, created_at 
FROM auctions 
WHERE id = $1 AND status = 'ACTIVE' FOR UPDATE;

-- name: CreateAuction :one
INSERT INTO auctions (id, item_id, seller_id, start_price, status, expires_at)
VALUES (sqlc.arg(id), sqlc.arg(item_id), sqlc.arg(seller_id), sqlc.arg(start_price), 'ACTIVE', sqlc.arg(expires_at))
RETURNING id, item_id, seller_id, start_price, current_highest_bid, current_highest_bidder, status, expires_at, created_at;

-- name: PlaceBid :one
UPDATE auctions 
SET current_highest_bid = sqlc.arg(amount), current_highest_bidder = sqlc.arg(bidder_id), updated_at = CURRENT_TIMESTAMP 
WHERE id = sqlc.arg(id) 
RETURNING id, item_id, seller_id, start_price, current_highest_bid, current_highest_bidder, status, expires_at, created_at;

-- name: ListActiveAuctions :many
SELECT id, item_id, seller_id, start_price, current_highest_bid, current_highest_bidder, status, expires_at, created_at
FROM auctions
WHERE status = 'ACTIVE'
ORDER BY expires_at ASC;
