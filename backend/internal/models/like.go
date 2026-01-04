package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID `json:"id" db:"id"`
	LikerID   uuid.UUID `json:"liker_id" db:"liker_id"`
	LikedID   uuid.UUID `json:"liked_id" db:"liked_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type LikeWithProfile struct {
	Like
	Profile *ProfileWithImages `json:"profile,omitempty"`
}
