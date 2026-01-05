package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// WealthStatus represents the payment/subscription tier for male users
type WealthStatus string

const (
	WealthStatusNone   WealthStatus = "none"
	WealthStatusLow    WealthStatus = "low"
	WealthStatusMedium WealthStatus = "medium"
	WealthStatusHigh   WealthStatus = "high"
)

// WealthStatusLabel returns the user-facing label for the status
func (w WealthStatus) Label() string {
	switch w {
	case WealthStatusLow:
		return "Trusted"
	case WealthStatusMedium:
		return "Premium"
	case WealthStatusHigh:
		return "Elite"
	default:
		return "Standard"
	}
}

// CanMessage returns true if this wealth status allows messaging
func (w WealthStatus) CanMessage() bool {
	return w == WealthStatusLow || w == WealthStatusMedium || w == WealthStatusHigh
}

// CanViewMessages returns true if this wealth status allows viewing message requests
func (w WealthStatus) CanViewMessages() bool {
	return w.CanMessage()
}

type User struct {
	ID                         uuid.UUID      `json:"id" db:"id"`
	GoogleID                   sql.NullString `json:"-" db:"google_id"`
	Email                      string         `json:"email" db:"email"`
	PasswordHash               sql.NullString `json:"-" db:"password_hash"`
	EmailVerified              bool           `json:"email_verified" db:"email_verified"`
	VerificationToken          sql.NullString `json:"-" db:"verification_token"`
	VerificationTokenExpiresAt sql.NullTime   `json:"-" db:"verification_token_expires_at"`
	WealthStatus               WealthStatus   `json:"wealth_status" db:"wealth_status"`
	WealthStatusExpiresAt      sql.NullTime   `json:"wealth_status_expires_at,omitempty" db:"wealth_status_expires_at"`
	CreatedAt                  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time      `json:"updated_at" db:"updated_at"`
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
