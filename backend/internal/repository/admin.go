package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// AdminUser represents a user with their profile and online status for admin view
type AdminUser struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	WealthStatus  string     `json:"wealth_status"`
	CreatedAt     time.Time  `json:"created_at"`
	DisplayName   *string    `json:"display_name,omitempty"`
	Gender        *string    `json:"gender,omitempty"`
	Age           *int       `json:"age,omitempty"`
	City          *string    `json:"city,omitempty"`
	State         *string    `json:"state,omitempty"`
	IsVerified    bool       `json:"is_verified"`
	IsOnline      bool       `json:"is_online"`
	LastSeen      *time.Time `json:"last_seen,omitempty"`
	ImageCount    int        `json:"image_count"`
}

// AdminUserWithImages includes images for a user
type AdminUserWithImages struct {
	AdminUser
	Images []models.ProfileImage `json:"images"`
}

// AdminMessage represents a message for admin view
type AdminMessage struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id"`
	SenderName     *string   `json:"sender_name,omitempty"`
	SenderEmail    string    `json:"sender_email"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// AdminVerificationRequest represents a verification request for admin view
type AdminVerificationRequest struct {
	models.VerificationRequest
	UserEmail   string  `json:"user_email"`
	DisplayName *string `json:"display_name,omitempty"`
}

// ListUsers returns all users with their profile info
func (r *AdminRepository) ListUsers(limit, offset int, search string) ([]AdminUser, int, error) {
	// Build search clause
	whereClause := "1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause = "(u.email ILIKE $1 OR p.display_name ILIKE $1)"
		args = append(args, "%"+search+"%")
		argIndex = 2
	}

	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM users u LEFT JOIN profiles p ON u.id = p.user_id WHERE ` + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add pagination args
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT u.id, u.email, u.email_verified, COALESCE(u.wealth_status, 'none'), u.created_at,
			   p.display_name, p.gender, p.age, p.city, p.state, COALESCE(p.is_verified, false),
			   COALESCE(up.is_online, false), up.last_seen,
			   (SELECT COUNT(*) FROM profile_images WHERE user_id = u.id)
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		LEFT JOIN user_presence up ON u.id = up.user_id
		WHERE %s
		ORDER BY u.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []AdminUser
	for rows.Next() {
		var u AdminUser
		var displayName, gender, city, state sql.NullString
		var age sql.NullInt32
		var lastSeen sql.NullTime

		err := rows.Scan(&u.ID, &u.Email, &u.EmailVerified, &u.WealthStatus, &u.CreatedAt,
			&displayName, &gender, &age, &city, &state, &u.IsVerified,
			&u.IsOnline, &lastSeen, &u.ImageCount)
		if err != nil {
			return nil, 0, err
		}

		if displayName.Valid {
			u.DisplayName = &displayName.String
		}
		if gender.Valid {
			u.Gender = &gender.String
		}
		if age.Valid {
			ageInt := int(age.Int32)
			u.Age = &ageInt
		}
		if city.Valid {
			u.City = &city.String
		}
		if state.Valid {
			u.State = &state.String
		}
		if lastSeen.Valid {
			u.LastSeen = &lastSeen.Time
		}

		users = append(users, u)
	}

	return users, total, nil
}

// GetUserWithImages returns a single user with all their images
func (r *AdminRepository) GetUserWithImages(userID uuid.UUID) (*AdminUserWithImages, error) {
	var u AdminUser
	var displayName, gender, city, state sql.NullString
	var age sql.NullInt32
	var lastSeen sql.NullTime

	err := r.db.QueryRow(`
		SELECT u.id, u.email, u.email_verified, COALESCE(u.wealth_status, 'none'), u.created_at,
			   p.display_name, p.gender, p.age, p.city, p.state, COALESCE(p.is_verified, false),
			   COALESCE(up.is_online, false), up.last_seen,
			   (SELECT COUNT(*) FROM profile_images WHERE user_id = u.id)
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		LEFT JOIN user_presence up ON u.id = up.user_id
		WHERE u.id = $1
	`, userID).Scan(&u.ID, &u.Email, &u.EmailVerified, &u.WealthStatus, &u.CreatedAt,
		&displayName, &gender, &age, &city, &state, &u.IsVerified,
		&u.IsOnline, &lastSeen, &u.ImageCount)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if displayName.Valid {
		u.DisplayName = &displayName.String
	}
	if gender.Valid {
		u.Gender = &gender.String
	}
	if age.Valid {
		ageInt := int(age.Int32)
		u.Age = &ageInt
	}
	if city.Valid {
		u.City = &city.String
	}
	if state.Valid {
		u.State = &state.String
	}
	if lastSeen.Valid {
		u.LastSeen = &lastSeen.Time
	}

	// Get images
	rows, err := r.db.Query(`
		SELECT id, user_id, s3_key, url, is_primary, sort_order, created_at
		FROM profile_images WHERE user_id = $1
		ORDER BY sort_order ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.ProfileImage
	for rows.Next() {
		var img models.ProfileImage
		err := rows.Scan(&img.ID, &img.UserID, &img.S3Key, &img.URL, &img.IsPrimary, &img.SortOrder, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return &AdminUserWithImages{
		AdminUser: u,
		Images:    images,
	}, nil
}

// ListMessages returns messages with sender info
func (r *AdminRepository) ListMessages(limit, offset int, conversationID *uuid.UUID) ([]AdminMessage, int, error) {
	args := []interface{}{}
	whereClause := "1=1"
	argIndex := 1

	if conversationID != nil {
		whereClause = "m.conversation_id = $1"
		args = append(args, *conversationID)
		argIndex = 2
	}

	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM messages m WHERE ` + whereClause
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT m.id, m.conversation_id, m.sender_id, p.display_name, u.email, m.content, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		LEFT JOIN profiles p ON m.sender_id = p.user_id
		WHERE %s
		ORDER BY m.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []AdminMessage
	for rows.Next() {
		var m AdminMessage
		var senderName sql.NullString

		err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &senderName, &m.SenderEmail, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, 0, err
		}

		if senderName.Valid {
			m.SenderName = &senderName.String
		}

		messages = append(messages, m)
	}

	return messages, total, nil
}

// ListVerificationRequests returns pending verification requests
func (r *AdminRepository) ListVerificationRequests(status string) ([]AdminVerificationRequest, error) {
	statusFilter := "v.status = $1"
	args := []interface{}{status}
	if status == "" {
		statusFilter = "1=1"
		args = []interface{}{}
	}

	query := `
		SELECT v.id, v.user_id, v.document_type, v.document_url, v.video_url, v.verification_code, 
			   v.status, v.rejection_reason, v.created_at, v.reviewed_at, v.reviewed_by,
			   u.email, p.display_name
		FROM verification_requests v
		JOIN users u ON v.user_id = u.id
		LEFT JOIN profiles p ON v.user_id = p.user_id
		WHERE ` + statusFilter + `
		ORDER BY v.created_at DESC
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []AdminVerificationRequest
	for rows.Next() {
		var req AdminVerificationRequest
		var rejectionReason, displayName sql.NullString
		var reviewedAt sql.NullTime
		var reviewedBy sql.NullString

		err := rows.Scan(&req.ID, &req.UserID, &req.DocumentType, &req.DocumentURL, &req.VideoURL,
			&req.VerificationCode, &req.Status, &rejectionReason, &req.CreatedAt, &reviewedAt, &reviewedBy,
			&req.UserEmail, &displayName)
		if err != nil {
			return nil, err
		}

		if rejectionReason.Valid {
			req.RejectionReason = &rejectionReason.String
		}
		if reviewedAt.Valid {
			req.ReviewedAt = &reviewedAt.Time
		}
		if reviewedBy.Valid {
			id, _ := uuid.Parse(reviewedBy.String)
			req.ReviewedBy = &id
		}
		if displayName.Valid {
			req.DisplayName = &displayName.String
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// ApproveVerification approves a verification request
func (r *AdminRepository) ApproveVerification(requestID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the user_id from the request
	var userID uuid.UUID
	err = tx.QueryRow(`SELECT user_id FROM verification_requests WHERE id = $1`, requestID).Scan(&userID)
	if err != nil {
		return err
	}

	// Update the verification request status
	_, err = tx.Exec(`
		UPDATE verification_requests 
		SET status = 'approved', reviewed_at = $1 
		WHERE id = $2
	`, time.Now().UTC(), requestID)
	if err != nil {
		return err
	}

	// Update the profile is_verified status
	_, err = tx.Exec(`
		UPDATE profiles SET is_verified = true, updated_at = $1 WHERE user_id = $2
	`, time.Now().UTC(), userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RejectVerification rejects a verification request
func (r *AdminRepository) RejectVerification(requestID uuid.UUID, reason string) error {
	_, err := r.db.Exec(`
		UPDATE verification_requests 
		SET status = 'rejected', rejection_reason = $1, reviewed_at = $2 
		WHERE id = $3
	`, reason, time.Now().UTC(), requestID)
	return err
}

// DeleteProfileImage deletes a profile image by its ID
func (r *AdminRepository) DeleteProfileImage(imageID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM profile_images WHERE id = $1`, imageID)
	return err
}

// GetProfileImage returns a profile image by its ID
func (r *AdminRepository) GetProfileImage(imageID uuid.UUID) (*models.ProfileImage, error) {
	var img models.ProfileImage
	err := r.db.QueryRow(`
		SELECT id, user_id, s3_key, url, is_primary, sort_order, created_at
		FROM profile_images WHERE id = $1
	`, imageID).Scan(&img.ID, &img.UserID, &img.S3Key, &img.URL, &img.IsPrimary, &img.SortOrder, &img.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &img, nil
}

// DeleteUser deletes a user and all related data
func (r *AdminRepository) DeleteUser(userID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete in order of dependencies
	// 1. Delete messages in conversations the user is part of
	_, err = tx.Exec(`
		DELETE FROM messages WHERE conversation_id IN (
			SELECT conversation_id FROM conversation_participants WHERE user_id = $1
		)
	`, userID)
	if err != nil {
		return err
	}

	// 2. Delete conversation participants
	_, err = tx.Exec(`
		DELETE FROM conversation_participants WHERE conversation_id IN (
			SELECT conversation_id FROM conversation_participants WHERE user_id = $1
		)
	`, userID)
	if err != nil {
		return err
	}

	// 3. Delete conversations
	_, err = tx.Exec(`
		DELETE FROM conversations WHERE id IN (
			SELECT id FROM conversations WHERE initiated_by = $1
		) OR id NOT IN (SELECT conversation_id FROM conversation_participants)
	`, userID)
	if err != nil {
		return err
	}

	// 4. Delete likes
	_, err = tx.Exec(`DELETE FROM likes WHERE liker_id = $1 OR liked_id = $1`, userID)
	if err != nil {
		return err
	}

	// 5. Delete notifications
	_, err = tx.Exec(`DELETE FROM notifications WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// 6. Delete profile images
	_, err = tx.Exec(`DELETE FROM profile_images WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// 7. Delete verification requests
	_, err = tx.Exec(`DELETE FROM verification_requests WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// 8. Delete user presence
	_, err = tx.Exec(`DELETE FROM user_presence WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// 9. Delete profile
	_, err = tx.Exec(`DELETE FROM profiles WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// 10. Delete user
	_, err = tx.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateUserWealthStatus updates a user's wealth status (trusted member)
func (r *AdminRepository) UpdateUserWealthStatus(userID uuid.UUID, status string) error {
	_, err := r.db.Exec(`
		UPDATE users SET wealth_status = $1, updated_at = $2 WHERE id = $3
	`, status, time.Now().UTC(), userID)
	return err
}

// GetStats returns admin dashboard stats
func (r *AdminRepository) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total users
	var totalUsers int
	r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	stats["total_users"] = totalUsers

	// Online users
	var onlineUsers int
	r.db.QueryRow(`SELECT COUNT(*) FROM user_presence WHERE is_online = true`).Scan(&onlineUsers)
	stats["online_users"] = onlineUsers

	// Verified users (person verified)
	var verifiedUsers int
	r.db.QueryRow(`SELECT COUNT(*) FROM profiles WHERE is_verified = true`).Scan(&verifiedUsers)
	stats["verified_users"] = verifiedUsers

	// Pending verifications
	var pendingVerifications int
	r.db.QueryRow(`SELECT COUNT(*) FROM verification_requests WHERE status = 'pending'`).Scan(&pendingVerifications)
	stats["pending_verifications"] = pendingVerifications

	// Total messages
	var totalMessages int
	r.db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&totalMessages)
	stats["total_messages"] = totalMessages

	// Trusted members (wealth_status = 'low')
	var trustedMembers int
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE wealth_status = 'low'`).Scan(&trustedMembers)
	stats["trusted_members"] = trustedMembers

	// Premium members
	var premiumMembers int
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE wealth_status = 'medium'`).Scan(&premiumMembers)
	stats["premium_members"] = premiumMembers

	// Elite members
	var eliteMembers int
	r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE wealth_status = 'high'`).Scan(&eliteMembers)
	stats["elite_members"] = eliteMembers

	return stats, nil
}

// AdminImage represents an image from any source for admin view
type AdminImage struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Source    string    `json:"source"` // "profile", "message", "verification_doc", "verification_video"
	UserID    uuid.UUID `json:"user_id"`
	UserEmail string    `json:"user_email"`
	CreatedAt time.Time `json:"created_at"`
}

// ListAllImages returns all images from profile_images and verification_requests
func (r *AdminRepository) ListAllImages(limit, offset int) ([]AdminImage, int, error) {
	// Count total from all sources
	var totalProfile, totalVerification int
	r.db.QueryRow(`SELECT COUNT(*) FROM profile_images`).Scan(&totalProfile)
	r.db.QueryRow(`SELECT COUNT(*) * 2 FROM verification_requests`).Scan(&totalVerification) // document + video
	total := totalProfile + totalVerification

	// Union query for all images (excluding messages as image_url column may not exist)
	query := `
		SELECT id, url, source, user_id, user_email, created_at FROM (
			SELECT pi.id::text, pi.url, 'profile' as source, pi.user_id, u.email as user_email, pi.created_at
			FROM profile_images pi
			JOIN users u ON pi.user_id = u.id
			
			UNION ALL
			
			SELECT v.id::text || '_doc', v.document_url as url, 'verification_doc' as source, v.user_id, u.email as user_email, v.created_at
			FROM verification_requests v
			JOIN users u ON v.user_id = u.id
			
			UNION ALL
			
			SELECT v.id::text || '_video', v.video_url as url, 'verification_video' as source, v.user_id, u.email as user_email, v.created_at
			FROM verification_requests v
			JOIN users u ON v.user_id = u.id
		) combined
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var images []AdminImage
	for rows.Next() {
		var img AdminImage
		err := rows.Scan(&img.ID, &img.URL, &img.Source, &img.UserID, &img.UserEmail, &img.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		images = append(images, img)
	}

	return images, total, nil
}

