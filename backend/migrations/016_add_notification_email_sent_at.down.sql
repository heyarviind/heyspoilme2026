DROP INDEX IF EXISTS idx_messages_notification_pending;
ALTER TABLE messages DROP COLUMN IF EXISTS notification_email_sent_at;

