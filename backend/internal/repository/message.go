package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) CreateConversation(initiatedBy uuid.UUID, participants []uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{
		ID:          uuid.New(),
		InitiatedBy: initiatedBy,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO conversations (id, initiated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`, conv.ID, conv.InitiatedBy, conv.CreatedAt, conv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, userID := range participants {
		_, err = tx.Exec(`
			INSERT INTO conversation_participants (conversation_id, user_id)
			VALUES ($1, $2)
		`, conv.ID, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *MessageRepository) FindConversationBetweenUsers(userID1, userID2 uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{}
	err := r.db.QueryRow(`
		SELECT c.id, c.initiated_by, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp1 ON c.id = cp1.conversation_id AND cp1.user_id = $1
		JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id = $2
	`, userID1, userID2).Scan(&conv.ID, &conv.InitiatedBy, &conv.CreatedAt, &conv.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *MessageRepository) GetUserConversations(userID uuid.UUID) ([]models.ConversationWithDetails, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.initiated_by, c.created_at, c.updated_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1
		ORDER BY c.updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.ConversationWithDetails
	for rows.Next() {
		var conv models.ConversationWithDetails
		err := rows.Scan(&conv.ID, &conv.InitiatedBy, &conv.CreatedAt, &conv.UpdatedAt)
		if err != nil {
			return nil, err
		}

		partRows, err := r.db.Query(`
			SELECT user_id FROM conversation_participants WHERE conversation_id = $1
		`, conv.ID)
		if err != nil {
			return nil, err
		}
		for partRows.Next() {
			var partID uuid.UUID
			partRows.Scan(&partID)
			conv.Participants = append(conv.Participants, partID)
		}
		partRows.Close()

		var lastMsg models.Message
		var readAt sql.NullTime
		var imageURL sql.NullString
		err = r.db.QueryRow(`
			SELECT id, conversation_id, sender_id, content, image_url, read_at, created_at
			FROM messages WHERE conversation_id = $1
			ORDER BY created_at DESC LIMIT 1
		`, conv.ID).Scan(&lastMsg.ID, &lastMsg.ConversationID, &lastMsg.SenderID, &lastMsg.Content, &imageURL, &readAt, &lastMsg.CreatedAt)
		if err == nil {
			if readAt.Valid {
				lastMsg.ReadAt = &readAt.Time
			}
			if imageURL.Valid {
				lastMsg.ImageURL = &imageURL.String
			}
			conv.LastMessage = &lastMsg
		}

		r.db.QueryRow(`
			SELECT COUNT(*) FROM messages 
			WHERE conversation_id = $1 AND sender_id != $2 AND read_at IS NULL
		`, conv.ID, userID).Scan(&conv.UnreadCount)

		conversations = append(conversations, conv)
	}

	return conversations, nil
}

func (r *MessageRepository) GetConversation(conversationID uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{}
	err := r.db.QueryRow(`
		SELECT id, initiated_by, created_at, updated_at
		FROM conversations WHERE id = $1
	`, conversationID).Scan(&conv.ID, &conv.InitiatedBy, &conv.CreatedAt, &conv.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return conv, nil
}

func (r *MessageRepository) IsUserInConversation(conversationID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2)
	`, conversationID, userID).Scan(&exists)
	return exists, err
}

func (r *MessageRepository) GetConversationParticipants(conversationID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(`
		SELECT user_id FROM conversation_participants WHERE conversation_id = $1
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		rows.Scan(&userID)
		participants = append(participants, userID)
	}

	return participants, nil
}

func (r *MessageRepository) CreateMessage(conversationID, senderID uuid.UUID, content string, imageURL *string) (*models.Message, error) {
	msg := &models.Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		ImageURL:       imageURL,
		CreatedAt:      time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, content, image_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.ImageURL, msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	r.db.Exec("UPDATE conversations SET updated_at = $1 WHERE id = $2", time.Now().UTC(), conversationID)

	return msg, nil
}

func (r *MessageRepository) GetMessages(conversationID uuid.UUID, limit, offset int) ([]models.Message, error) {
	rows, err := r.db.Query(`
		SELECT id, conversation_id, sender_id, content, image_url, read_at, created_at
		FROM messages WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var readAt sql.NullTime
		var imageURL sql.NullString
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &imageURL, &readAt, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		if readAt.Valid {
			msg.ReadAt = &readAt.Time
		}
		if imageURL.Valid {
			msg.ImageURL = &imageURL.String
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (r *MessageRepository) MarkMessagesAsRead(conversationID, userID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE messages SET read_at = $1
		WHERE conversation_id = $2 AND sender_id != $3 AND read_at IS NULL
	`, time.Now().UTC(), conversationID, userID)
	return err
}

func (r *MessageRepository) GetUnreadMessageCount(userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM messages m
		JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id
		WHERE cp.user_id = $1 AND m.sender_id != $1 AND m.read_at IS NULL
	`, userID).Scan(&count)
	return count, err
}

// GetUserMessageImageURLs returns all image URLs from messages sent by the user
func (r *MessageRepository) GetUserMessageImageURLs(userID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT image_url FROM messages 
		WHERE sender_id = $1 AND image_url IS NOT NULL AND image_url != ''
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// GetUnreadMessagesNeedingNotification returns messages that are unread for more than the specified duration
// and haven't had a notification email sent yet
func (r *MessageRepository) GetUnreadMessagesNeedingNotification(olderThan time.Duration) ([]models.MessageNotificationInfo, error) {
	cutoffTime := time.Now().UTC().Add(-olderThan)
	
	rows, err := r.db.Query(`
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.created_at,
			   u.email as recipient_email, u.id as recipient_id,
			   p.display_name as sender_name
		FROM messages m
		JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id AND cp.user_id != m.sender_id
		JOIN users u ON cp.user_id = u.id
		LEFT JOIN profiles p ON m.sender_id = p.user_id
		WHERE m.read_at IS NULL 
		  AND m.created_at < $1
		  AND m.notification_email_sent_at IS NULL
		  AND u.email_verified = true
	`, cutoffTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.MessageNotificationInfo
	for rows.Next() {
		var n models.MessageNotificationInfo
		var senderName sql.NullString
		err := rows.Scan(&n.MessageID, &n.ConversationID, &n.SenderID, &n.Content, &n.CreatedAt,
			&n.RecipientEmail, &n.RecipientID, &senderName)
		if err != nil {
			return nil, err
		}
		if senderName.Valid {
			n.SenderName = senderName.String
		} else {
			n.SenderName = "Someone"
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// MarkNotificationEmailSent marks that a notification email has been sent for a message
func (r *MessageRepository) MarkNotificationEmailSent(messageID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE messages SET notification_email_sent_at = $1 WHERE id = $2
	`, time.Now().UTC(), messageID)
	return err
}

// DeleteUserConversations deletes all conversations where the user is a participant
// This also handles cleanup of related messages and participants
func (r *MessageRepository) DeleteUserConversations(userID uuid.UUID) error {
	// Get all conversations the user is in
	rows, err := r.db.Query(`
		SELECT conversation_id FROM conversation_participants WHERE user_id = $1
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var conversationIDs []uuid.UUID
	for rows.Next() {
		var convID uuid.UUID
		if err := rows.Scan(&convID); err != nil {
			return err
		}
		conversationIDs = append(conversationIDs, convID)
	}

	for _, convID := range conversationIDs {
		// Delete messages in the conversation
		if _, err := r.db.Exec(`DELETE FROM messages WHERE conversation_id = $1`, convID); err != nil {
			return err
		}
		// Delete participants
		if _, err := r.db.Exec(`DELETE FROM conversation_participants WHERE conversation_id = $1`, convID); err != nil {
			return err
		}
		// Delete the conversation
		if _, err := r.db.Exec(`DELETE FROM conversations WHERE id = $1`, convID); err != nil {
			return err
		}
	}

	return nil
}
