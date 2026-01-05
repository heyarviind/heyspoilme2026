-- Add verified status for premium male profiles
ALTER TABLE profiles ADD COLUMN is_verified BOOLEAN DEFAULT FALSE;

-- Index for faster verified profile queries
CREATE INDEX idx_profiles_verified ON profiles(is_verified) WHERE is_verified = true;



