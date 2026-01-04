package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	GoogleID     sql.NullString `json:"-" db:"google_id"`
	Email        string         `json:"email" db:"email"`
	PasswordHash sql.NullString `json:"-" db:"password_hash"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

type UserWithProfile struct {
	User
	Profile *Profile `json:"profile,omitempty"`
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
