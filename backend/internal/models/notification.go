package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeLike        NotificationType = "new_like"
	NotificationTypeMessage     NotificationType = "new_message"
	NotificationTypeProfileView NotificationType = "profile_view"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	UserID    uuid.UUID        `json:"user_id" db:"user_id"`
	Type      NotificationType `json:"type" db:"type"`
	Data      json.RawMessage  `json:"data" db:"data"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

type NotificationData struct {
	FromUserID     uuid.UUID `json:"from_user_id,omitempty"`
	FromUserName   string    `json:"from_user_name,omitempty"`
	FromUserImage  string    `json:"from_user_image,omitempty"`
	ConversationID uuid.UUID `json:"conversation_id,omitempty"`
	MessagePreview string    `json:"message_preview,omitempty"`
}

type NotificationWithDetails struct {
	Notification
	ParsedData NotificationData `json:"parsed_data"`
}



