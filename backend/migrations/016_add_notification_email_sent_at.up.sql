-- Add notification_email_sent_at column to messages table
ALTER TABLE messages 
ADD COLUMN IF NOT EXISTS notification_email_sent_at TIMESTAMP WITH TIME ZONE;

-- Index for finding messages that haven't had notifications sent
CREATE INDEX IF NOT EXISTS idx_messages_notification_pending 
ON messages(created_at) 
WHERE notification_email_sent_at IS NULL AND read_at IS NULL;

