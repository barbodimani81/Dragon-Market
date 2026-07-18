-- 1. Strip the foreign key constraint column back off the wallets engine
ALTER TABLE wallets DROP COLUMN IF EXISTS user_id;

-- 2. Obliterate the users space entirely
DROP TABLE IF EXISTS users;
