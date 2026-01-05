DROP INDEX IF EXISTS idx_users_verification_token;
DROP INDEX IF EXISTS idx_users_email_verified;
ALTER TABLE users DROP COLUMN IF EXISTS verification_token_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS verification_token;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;

