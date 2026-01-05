-- Add is_fake column to profiles table to mark fake/demo profiles
ALTER TABLE profiles ADD COLUMN is_fake BOOLEAN NOT NULL DEFAULT false;

-- Index for filtering fake profiles
CREATE INDEX idx_profiles_is_fake ON profiles(is_fake);

