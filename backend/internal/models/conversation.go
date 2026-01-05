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

// LockedConversation represents a conversation that the user cannot fully view
// due to wealth_status restrictions. Used for "message requests" teaser.
type LockedConversation struct {
	ID            uuid.UUID `json:"id"`
	BlurredPreview string   `json:"blurred_preview"`
	SenderImage   string    `json:"sender_image,omitempty"`
	SenderAge     int       `json:"sender_age,omitempty"`
	SenderCity    string    `json:"sender_city,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// InboxResponse contains both visible and locked conversations
type InboxResponse struct {
	Conversations      []ConversationWithDetails `json:"conversations"`
	LockedCount        int                       `json:"locked_count"`
	LockedPreviews     []LockedConversation      `json:"locked_previews,omitempty"`
	CanViewAllMessages bool                      `json:"can_view_all_messages"`
}
