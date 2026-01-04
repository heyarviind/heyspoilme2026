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
		err = r.db.QueryRow(`
			SELECT id, conversation_id, sender_id, content, read_at, created_at
			FROM messages WHERE conversation_id = $1
			ORDER BY created_at DESC LIMIT 1
		`, conv.ID).Scan(&lastMsg.ID, &lastMsg.ConversationID, &lastMsg.SenderID, &lastMsg.Content, &readAt, &lastMsg.CreatedAt)
		if err == nil {
			if readAt.Valid {
				lastMsg.ReadAt = &readAt.Time
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

func (r *MessageRepository) CreateMessage(conversationID, senderID uuid.UUID, content string) (*models.Message, error) {
	msg := &models.Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt)
	if err != nil {
		return nil, err
	}

	r.db.Exec("UPDATE conversations SET updated_at = $1 WHERE id = $2", time.Now().UTC(), conversationID)

	return msg, nil
}

func (r *MessageRepository) GetMessages(conversationID uuid.UUID, limit, offset int) ([]models.Message, error) {
	rows, err := r.db.Query(`
		SELECT id, conversation_id, sender_id, content, read_at, created_at
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
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content, &readAt, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		if readAt.Valid {
			msg.ReadAt = &readAt.Time
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
