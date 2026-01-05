-- Remove password_hash column
DROP INDEX IF EXISTS idx_users_email_password;
ALTER TABLE users DROP COLUMN IF EXISTS password_hash;

-- Note: Cannot restore NOT NULL on google_id if there are null values
-- You may need to handle this manually



