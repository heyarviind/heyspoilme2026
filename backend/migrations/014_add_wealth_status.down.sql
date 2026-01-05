DROP INDEX IF EXISTS idx_users_wealth_expires;
DROP INDEX IF EXISTS idx_users_wealth_status;
ALTER TABLE users DROP COLUMN IF EXISTS wealth_status_expires_at;
ALTER TABLE users DROP COLUMN IF EXISTS wealth_status;
DROP TYPE IF EXISTS wealth_status_type;

