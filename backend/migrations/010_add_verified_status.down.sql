-- Remove verified status
DROP INDEX IF EXISTS idx_profiles_verified;
ALTER TABLE profiles DROP COLUMN IF EXISTS is_verified;

