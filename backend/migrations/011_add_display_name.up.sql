-- Add display_name column to profiles table
ALTER TABLE profiles ADD COLUMN display_name VARCHAR(50) NOT NULL DEFAULT '';

