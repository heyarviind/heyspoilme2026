-- Create wealth_status enum type
DO $$ BEGIN
    CREATE TYPE wealth_status_type AS ENUM ('none', 'low', 'medium', 'high');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Add wealth_status column to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS wealth_status wealth_status_type DEFAULT 'none';

-- Add wealth_status_expires_at for subscription expiry tracking
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS wealth_status_expires_at TIMESTAMP WITH TIME ZONE;

-- Create index for efficient queries on wealth_status
CREATE INDEX IF NOT EXISTS idx_users_wealth_status ON users(wealth_status);

-- Create index for finding expired subscriptions
CREATE INDEX IF NOT EXISTS idx_users_wealth_expires 
ON users(wealth_status_expires_at) 
WHERE wealth_status != 'none' AND wealth_status_expires_at IS NOT NULL;

