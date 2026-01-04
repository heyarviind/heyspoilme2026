package models

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID          uuid.UUID `json:"id" db:"id"`
	InitiatedBy uuid.UUID `json:"initiated_by" db:"initiated_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ConversationParticipant struct {
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
}

type ConversationWithDetails struct {
	Conversation
	Participants []uuid.UUID        `json:"participants"`
	OtherUser    *ProfileWithImages `json:"other_user,omitempty"`
	LastMessage  *Message           `json:"last_message,omitempty"`
	UnreadCount  int                `json:"unread_count"`
}

type CreateConversationRequest struct {
	RecipientID uuid.UUID `json:"recipient_id" binding:"required"`
	Message     string    `json:"message" binding:"required,max=2000"`
}
