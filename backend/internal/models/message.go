package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	ConversationID uuid.UUID  `json:"conversation_id" db:"conversation_id"`
	SenderID       uuid.UUID  `json:"sender_id" db:"sender_id"`
	Content        string     `json:"content" db:"content"`
	ImageURL       *string    `json:"image_url,omitempty" db:"image_url"`
	ReadAt         *time.Time `json:"read_at,omitempty" db:"read_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

type MessageWithSender struct {
	Message
	SenderName  string `json:"sender_name"`
	SenderImage string `json:"sender_image,omitempty"`
}

type SendMessageRequest struct {
	Content  string  `json:"content" binding:"max=2000"`
	ImageURL *string `json:"image_url,omitempty"`
}

type WSMessageType string

const (
	WSTypeMessage      WSMessageType = "message"
	WSTypeTyping       WSMessageType = "typing"
	WSTypeStopTyping   WSMessageType = "stop_typing"
	WSTypeReadReceipt  WSMessageType = "read_receipt"
	WSTypeNotification WSMessageType = "notification"
	WSTypePresence     WSMessageType = "presence"
)

type WSMessage struct {
	Type    WSMessageType `json:"type"`
	Payload interface{}   `json:"payload"`
}

type WSTypingPayload struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
}

type WSPresencePayload struct {
	UserID   uuid.UUID `json:"user_id"`
	IsOnline bool      `json:"is_online"`
}

// MessageNotificationInfo contains info needed to send email notification for unread message
type MessageNotificationInfo struct {
	MessageID      uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	SenderName     string
	Content        string
	CreatedAt      time.Time
	RecipientID    uuid.UUID
	RecipientEmail string
}
