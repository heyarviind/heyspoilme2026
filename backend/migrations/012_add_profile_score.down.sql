-- Remove profile_score column from profiles table
DROP INDEX IF EXISTS idx_profiles_score;
ALTER TABLE profiles DROP COLUMN profile_score;


