-- 1. Inventory Items Table
CREATE TABLE IF NOT EXISTS items (
                                     id UUID PRIMARY KEY,
                                     owner_id UUID NOT NULL,
                                     name VARCHAR(255) NOT NULL,
    rarity VARCHAR(50) NOT NULL, -- COMMON, RARE, LEGENDARY
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
                             );

CREATE INDEX IF NOT EXISTS idx_items_owner_id ON items(owner_id);

-- 2. Fixed-Price Marketplace Listings (Common & Rare)
CREATE TABLE IF NOT EXISTS listings (
                                        id UUID PRIMARY KEY,
                                        item_id UUID NOT NULL UNIQUE,
                                        seller_id UUID NOT NULL,
                                        price BIGINT NOT NULL,
                                        status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, SOLD, CANCELLED
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT fk_listings_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    CONSTRAINT check_positive_price CHECK (price > 0)
    );

CREATE INDEX IF NOT EXISTS idx_listings_status ON listings(status) WHERE status = 'ACTIVE';

-- 3. Auctions Table (Legendary)
CREATE TABLE IF NOT EXISTS auctions (
                                        id UUID PRIMARY KEY,
                                        item_id UUID NOT NULL UNIQUE,
                                        seller_id UUID NOT NULL,
                                        start_price BIGINT NOT NULL,
                                        current_highest_bid BIGINT,
                                        current_highest_bidder UUID,
                                        status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, COMPLETED, CANCELLED
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT fk_auctions_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    CONSTRAINT check_positive_start_price CHECK (start_price > 0),
    CONSTRAINT check_bid_limits CHECK (current_highest_bid IS NULL OR current_highest_bid >= start_price)
    );

CREATE INDEX IF NOT EXISTS idx_auctions_status_expires ON auctions(status, expires_at) WHERE status = 'ACTIVE';
