package repository

import (
	"database/sql"
	"fmt"
	"strings"
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
	Bio           *string    `json:"bio,omitempty"`
	City          *string    `json:"city,omitempty"`
	State         *string    `json:"state,omitempty"`
	Latitude      *float64   `json:"latitude,omitempty"`
	Longitude     *float64   `json:"longitude,omitempty"`
	IsVerified    bool       `json:"is_verified"`
	IsFake        bool       `json:"is_fake"`
	IsOnline      bool       `json:"is_online"`
	LastSeen      *time.Time `json:"last_seen,omitempty"`
	ImageCount    int        `json:"image_count"`
	PrimaryImage  *string    `json:"primary_image,omitempty"`
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
			   p.display_name, p.gender, p.age, p.bio, p.city, p.state, p.latitude, p.longitude,
			   COALESCE(p.is_verified, false), COALESCE(p.is_fake, false),
			   COALESCE(up.is_online, false), up.last_seen,
			   (SELECT COUNT(*) FROM profile_images WHERE user_id = u.id),
			   (SELECT url FROM profile_images WHERE user_id = u.id ORDER BY is_primary DESC, sort_order ASC LIMIT 1)
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
		var displayName, gender, bio, city, state, primaryImage sql.NullString
		var age sql.NullInt32
		var latitude, longitude sql.NullFloat64
		var lastSeen sql.NullTime

		err := rows.Scan(&u.ID, &u.Email, &u.EmailVerified, &u.WealthStatus, &u.CreatedAt,
			&displayName, &gender, &age, &bio, &city, &state, &latitude, &longitude,
			&u.IsVerified, &u.IsFake, &u.IsOnline, &lastSeen, &u.ImageCount, &primaryImage)
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
		if bio.Valid {
			u.Bio = &bio.String
		}
		if city.Valid {
			u.City = &city.String
		}
		if state.Valid {
			u.State = &state.String
		}
		if latitude.Valid {
			u.Latitude = &latitude.Float64
		}
		if longitude.Valid {
			u.Longitude = &longitude.Float64
		}
		if lastSeen.Valid {
			u.LastSeen = &lastSeen.Time
		}
		if primaryImage.Valid {
			u.PrimaryImage = &primaryImage.String
		}

		users = append(users, u)
	}

	return users, total, nil
}

// GetUserWithImages returns a single user with all their images
func (r *AdminRepository) GetUserWithImages(userID uuid.UUID) (*AdminUserWithImages, error) {
	var u AdminUser
	var displayName, gender, bio, city, state sql.NullString
	var age sql.NullInt32
	var latitude, longitude sql.NullFloat64
	var lastSeen sql.NullTime

	err := r.db.QueryRow(`
		SELECT u.id, u.email, u.email_verified, COALESCE(u.wealth_status, 'none'), u.created_at,
			   p.display_name, p.gender, p.age, p.bio, p.city, p.state, p.latitude, p.longitude,
			   COALESCE(p.is_verified, false), COALESCE(p.is_fake, false),
			   COALESCE(up.is_online, false), up.last_seen,
			   (SELECT COUNT(*) FROM profile_images WHERE user_id = u.id)
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		LEFT JOIN user_presence up ON u.id = up.user_id
		WHERE u.id = $1
	`, userID).Scan(&u.ID, &u.Email, &u.EmailVerified, &u.WealthStatus, &u.CreatedAt,
		&displayName, &gender, &age, &bio, &city, &state, &latitude, &longitude,
		&u.IsVerified, &u.IsFake, &u.IsOnline, &lastSeen, &u.ImageCount)

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
	if bio.Valid {
		u.Bio = &bio.String
	}
	if city.Valid {
		u.City = &city.String
	}
	if state.Valid {
		u.State = &state.String
	}
	if latitude.Valid {
		u.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		u.Longitude = &longitude.Float64
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

// AddUserProfileImage adds a profile image for a user (used by admin)
func (r *AdminRepository) AddUserProfileImage(userID uuid.UUID, s3Key, url string, isPrimary bool) (*models.ProfileImage, error) {
	var maxOrder int
	r.db.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM profile_images WHERE user_id = $1", userID).Scan(&maxOrder)

	img := &models.ProfileImage{
		ID:        uuid.New(),
		UserID:    userID,
		S3Key:     s3Key,
		URL:       url,
		IsPrimary: isPrimary,
		SortOrder: maxOrder + 1,
		CreatedAt: time.Now().UTC(),
	}

	if isPrimary {
		r.db.Exec("UPDATE profile_images SET is_primary = false WHERE user_id = $1", userID)
	}

	_, err := r.db.Exec(`
		INSERT INTO profile_images (id, user_id, s3_key, url, is_primary, sort_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, img.ID, img.UserID, img.S3Key, img.URL, img.IsPrimary, img.SortOrder, img.CreatedAt)

	if err != nil {
		return nil, err
	}

	return img, nil
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

// UpdateUserPresence updates a user's online status and last_seen time
func (r *AdminRepository) UpdateUserPresence(userID uuid.UUID, isOnline bool) error {
	now := time.Now().UTC()

	// Try to update existing record
	result, err := r.db.Exec(`
		UPDATE user_presence SET is_online = $1, last_seen = $2 WHERE user_id = $3
	`, isOnline, now, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Insert new record
		_, err = r.db.Exec(`
			INSERT INTO user_presence (user_id, is_online, last_seen) VALUES ($1, $2, $3)
		`, userID, isOnline, now)
		return err
	}

	return nil
}

// UpdateUserProfile updates a user's profile fields
func (r *AdminRepository) UpdateUserProfile(userID uuid.UUID, displayName *string, age *int, bio *string, city *string, state *string, latitude *float64, longitude *float64) error {
	setClauses := []string{"updated_at = $1"}
	args := []interface{}{time.Now().UTC()}
	argIndex := 2

	if displayName != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, *displayName)
		argIndex++
	}
	if age != nil {
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", argIndex))
		args = append(args, *age)
		argIndex++
	}
	if bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argIndex))
		args = append(args, *bio)
		argIndex++
	}
	if city != nil {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *city)
		argIndex++
	}
	if state != nil {
		setClauses = append(setClauses, fmt.Sprintf("state = $%d", argIndex))
		args = append(args, *state)
		argIndex++
	}
	if latitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("latitude = $%d", argIndex))
		args = append(args, *latitude)
		argIndex++
	}
	if longitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("longitude = $%d", argIndex))
		args = append(args, *longitude)
		argIndex++
	}

	if len(setClauses) == 1 {
		// Only updated_at, nothing to update
		return nil
	}

	args = append(args, userID)
	query := fmt.Sprintf("UPDATE profiles SET %s WHERE user_id = $%d", 
		strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	return err
}

// SendMessageAsUser sends a message as a specific user to another user
func (r *AdminRepository) SendMessageAsUser(senderID, recipientID uuid.UUID, content string) error {
	now := time.Now().UTC()

	// Find or create conversation between these users
	var conversationID uuid.UUID
	err := r.db.QueryRow(`
		SELECT c.id FROM conversations c
		JOIN conversation_participants cp1 ON c.id = cp1.conversation_id AND cp1.user_id = $1
		JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id = $2
		LIMIT 1
	`, senderID, recipientID).Scan(&conversationID)

	if err == sql.ErrNoRows {
		// Create new conversation
		conversationID = uuid.New()
		_, err = r.db.Exec(`
			INSERT INTO conversations (id, initiated_by, created_at) VALUES ($1, $2, $3)
		`, conversationID, senderID, now)
		if err != nil {
			return err
		}

		// Add participants
		_, err = r.db.Exec(`
			INSERT INTO conversation_participants (conversation_id, user_id, joined_at)
			VALUES ($1, $2, $3), ($1, $4, $3)
		`, conversationID, senderID, now, recipientID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Insert message
	messageID := uuid.New()
	_, err = r.db.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, messageID, conversationID, senderID, content, now)
	if err != nil {
		return err
	}

	// Update conversation last_message_at
	_, err = r.db.Exec(`
		UPDATE conversations SET last_message_at = $1 WHERE id = $2
	`, now, conversationID)

	return err
}

