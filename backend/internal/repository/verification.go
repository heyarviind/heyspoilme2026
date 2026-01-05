package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(userID uuid.UUID, req *models.CreateVerificationRequest, code string) (*models.VerificationRequest, error) {
	verification := &models.VerificationRequest{
		ID:               uuid.New(),
		UserID:           userID,
		DocumentType:     req.DocumentType,
		DocumentURL:      req.DocumentURL,
		VideoURL:         req.VideoURL,
		VerificationCode: code,
		Status:           models.VerificationStatusPending,
		CreatedAt:        time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO verification_requests (id, user_id, document_type, document_url, video_url, verification_code, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, verification.ID, verification.UserID, verification.DocumentType, verification.DocumentURL,
		verification.VideoURL, verification.VerificationCode, verification.Status, verification.CreatedAt)

	if err != nil {
		return nil, err
	}

	return verification, nil
}

func (r *VerificationRepository) GetPendingByUserID(userID uuid.UUID) (*models.VerificationRequest, error) {
	var v models.VerificationRequest
	var rejectionReason sql.NullString
	var reviewedAt sql.NullTime
	var reviewedBy sql.NullString

	err := r.db.QueryRow(`
		SELECT id, user_id, document_type, document_url, video_url, verification_code, status, 
		       rejection_reason, created_at, reviewed_at, reviewed_by
		FROM verification_requests 
		WHERE user_id = $1 AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(&v.ID, &v.UserID, &v.DocumentType, &v.DocumentURL, &v.VideoURL,
		&v.VerificationCode, &v.Status, &rejectionReason, &v.CreatedAt, &reviewedAt, &reviewedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if rejectionReason.Valid {
		v.RejectionReason = &rejectionReason.String
	}
	if reviewedAt.Valid {
		v.ReviewedAt = &reviewedAt.Time
	}
	if reviewedBy.Valid {
		id, _ := uuid.Parse(reviewedBy.String)
		v.ReviewedBy = &id
	}

	return &v, nil
}

func (r *VerificationRepository) GetLatestByUserID(userID uuid.UUID) (*models.VerificationRequest, error) {
	var v models.VerificationRequest
	var rejectionReason sql.NullString
	var reviewedAt sql.NullTime
	var reviewedBy sql.NullString

	err := r.db.QueryRow(`
		SELECT id, user_id, document_type, document_url, video_url, verification_code, status, 
		       rejection_reason, created_at, reviewed_at, reviewed_by
		FROM verification_requests 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID).Scan(&v.ID, &v.UserID, &v.DocumentType, &v.DocumentURL, &v.VideoURL,
		&v.VerificationCode, &v.Status, &rejectionReason, &v.CreatedAt, &reviewedAt, &reviewedBy)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if rejectionReason.Valid {
		v.RejectionReason = &rejectionReason.String
	}
	if reviewedAt.Valid {
		v.ReviewedAt = &reviewedAt.Time
	}
	if reviewedBy.Valid {
		id, _ := uuid.Parse(reviewedBy.String)
		v.ReviewedBy = &id
	}

	return &v, nil
}

func (r *VerificationRepository) HasPendingRequest(userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM verification_requests WHERE user_id = $1 AND status = 'pending')
	`, userID).Scan(&exists)
	return exists, err
}

func (r *VerificationRepository) DeleteByUserID(userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM verification_requests WHERE user_id = $1`, userID)
	return err
}



