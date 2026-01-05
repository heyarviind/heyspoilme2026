-- Migration 014: Add wealth_status for premium male members
-- Apply this migration through PocketBase admin or your preferred migration tool

-- UP Migration
-- ============

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

-- Create index for efficient queries on wealth_status (for female browse ranking)
CREATE INDEX IF NOT EXISTS idx_users_wealth_status ON users(wealth_status);

-- Create index for finding expired subscriptions
CREATE INDEX IF NOT EXISTS idx_users_wealth_expires 
ON users(wealth_status_expires_at) 
WHERE wealth_status != 'none' AND wealth_status_expires_at IS NOT NULL;


-- DOWN Migration (if needed)
-- ============
-- DROP INDEX IF EXISTS idx_users_wealth_expires;
-- DROP INDEX IF EXISTS idx_users_wealth_status;
-- ALTER TABLE users DROP COLUMN IF EXISTS wealth_status_expires_at;
-- ALTER TABLE users DROP COLUMN IF EXISTS wealth_status;
-- DROP TYPE IF EXISTS wealth_status_type;


-- Notes:
-- ======
-- wealth_status values and their meanings:
--   none   = Standard member (default, cannot send messages, cannot view message requests)
--   low    = Trusted member (can send messages, can view message requests)
--   medium = Premium member (same as low + priority placement)
--   high   = Elite member (same as medium + highest priority)
--
-- UI labels (don't expose internal names):
--   none   -> "Standard"
--   low    -> "Trusted"
--   medium -> "Premium"
--   high   -> "Elite"
--
-- Rules implemented:
-- 1. Likes: Require person_verified (is_verified in profiles table) for both genders
-- 2. Messaging:
--    - Females can initiate conversations if person_verified
--    - Males cannot initiate conversations (only reply)
--    - Males need wealth_status != 'none' AND person_verified to send messages
-- 3. Inbox visibility:
--    - Females see all messages
--    - Males with wealth_status = 'none' see locked message requests (blurred previews)
--    - Males with wealth_status != 'none' see all messages
-- 4. Female browse ranking:
--    - Primary sort: wealth_status (high > medium > low > none)
--    - Secondary: person_verified (true > false)
--    - Then: online status, activity, profile score



