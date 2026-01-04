package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPresence struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	IsOnline  bool      `json:"is_online" db:"is_online"`
	LastSeen  time.Time `json:"last_seen" db:"last_seen"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
