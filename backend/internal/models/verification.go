package models

import (
	"time"

	"github.com/google/uuid"
)

type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusApproved VerificationStatus = "approved"
	VerificationStatusRejected VerificationStatus = "rejected"
)

type DocumentType string

const (
	DocumentTypeAadhar         DocumentType = "aadhar"
	DocumentTypePassport       DocumentType = "passport"
	DocumentTypeDrivingLicense DocumentType = "driving_license"
)

type VerificationRequest struct {
	ID               uuid.UUID          `json:"id" db:"id"`
	UserID           uuid.UUID          `json:"user_id" db:"user_id"`
	DocumentType     DocumentType       `json:"document_type" db:"document_type"`
	DocumentURL      string             `json:"document_url" db:"document_url"`
	VideoURL         string             `json:"video_url" db:"video_url"`
	VerificationCode string             `json:"verification_code" db:"verification_code"`
	Status           VerificationStatus `json:"status" db:"status"`
	RejectionReason  *string            `json:"rejection_reason,omitempty" db:"rejection_reason"`
	CreatedAt        time.Time          `json:"created_at" db:"created_at"`
	ReviewedAt       *time.Time         `json:"reviewed_at,omitempty" db:"reviewed_at"`
	ReviewedBy       *uuid.UUID         `json:"reviewed_by,omitempty" db:"reviewed_by"`
}

type CreateVerificationRequest struct {
	DocumentType DocumentType `json:"document_type" binding:"required"`
	DocumentURL  string       `json:"document_url" binding:"required"`
	VideoURL     string       `json:"video_url" binding:"required"`
}

type VerificationCodeResponse struct {
	Code string `json:"code"`
}

type VerificationStatusResponse struct {
	Status          VerificationStatus `json:"status"`
	RejectionReason *string            `json:"rejection_reason,omitempty"`
	CreatedAt       *time.Time         `json:"created_at,omitempty"`
}


