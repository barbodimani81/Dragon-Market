-- 1. Create the base users structure
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

-- 2. Link wallets to users by appending the foreign key constraint
-- Using an ALTER sequence to maintain structural continuity with migration 0001
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS user_id UUID REFERENCES users(id) ON DELETE CASCADE;
