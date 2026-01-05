-- Add profile_score column to profiles table for ranking algorithm
ALTER TABLE profiles ADD COLUMN profile_score FLOAT NOT NULL DEFAULT 0;

-- Create index for efficient sorting by score
CREATE INDEX idx_profiles_score ON profiles(profile_score DESC);



