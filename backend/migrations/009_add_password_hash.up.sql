-- Make google_id nullable and add password_hash
ALTER TABLE users ALTER COLUMN google_id DROP NOT NULL;
ALTER TABLE users ADD COLUMN password_hash VARCHAR(255);

-- Create index for faster email lookups during password auth
CREATE INDEX IF NOT EXISTS idx_users_email_password ON users(email) WHERE password_hash IS NOT NULL;



