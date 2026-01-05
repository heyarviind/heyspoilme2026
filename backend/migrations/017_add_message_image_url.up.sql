-- Add image_url column to messages table for image messages
ALTER TABLE messages ADD COLUMN IF NOT EXISTS image_url TEXT;

