-- Remove is_fake column from profiles table
DROP INDEX IF EXISTS idx_profiles_is_fake;
ALTER TABLE profiles DROP COLUMN IF EXISTS is_fake;

