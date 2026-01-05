package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(googleID, email string, emailVerified bool) (*models.User, error) {
	user := &models.User{
		ID:            uuid.New(),
		GoogleID:      sql.NullString{String: googleID, Valid: googleID != ""},
		Email:         email,
		EmailVerified: emailVerified,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO users (id, google_id, email, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.GoogleID, user.Email, user.EmailVerified, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateWithPassword(email, passwordHash, verificationToken string, tokenExpiresAt time.Time) (*models.User, error) {
	user := &models.User{
		ID:                         uuid.New(),
		Email:                      email,
		PasswordHash:               sql.NullString{String: passwordHash, Valid: true},
		EmailVerified:              false,
		VerificationToken:          sql.NullString{String: verificationToken, Valid: true},
		VerificationTokenExpiresAt: sql.NullTime{Time: tokenExpiresAt, Valid: true},
		CreatedAt:                  time.Now().UTC(),
		UpdatedAt:                  time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO users (id, email, password_hash, email_verified, verification_token, verification_token_expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, user.ID, user.Email, user.PasswordHash, user.EmailVerified, user.VerificationToken, user.VerificationTokenExpiresAt, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, google_id, email, password_hash, email_verified, verification_token, verification_token_expires_at, 
		       COALESCE(wealth_status, 'none'), wealth_status_expires_at, created_at, updated_at
		FROM users WHERE google_id = $1
	`, googleID).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.EmailVerified, 
		&user.VerificationToken, &user.VerificationTokenExpiresAt, &user.WealthStatus, &user.WealthStatusExpiresAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, google_id, email, password_hash, email_verified, verification_token, verification_token_expires_at,
		       COALESCE(wealth_status, 'none'), wealth_status_expires_at, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.EmailVerified, 
		&user.VerificationToken, &user.VerificationTokenExpiresAt, &user.WealthStatus, &user.WealthStatusExpiresAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, google_id, email, password_hash, email_verified, verification_token, verification_token_expires_at,
		       COALESCE(wealth_status, 'none'), wealth_status_expires_at, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.EmailVerified, 
		&user.VerificationToken, &user.VerificationTokenExpiresAt, &user.WealthStatus, &user.WealthStatusExpiresAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindOrCreate(googleID, email string) (*models.User, bool, error) {
	user, err := r.FindByGoogleID(googleID)
	if err != nil {
		return nil, false, err
	}

	if user != nil {
		// Ensure Google users are always verified
		if !user.EmailVerified {
			_, err = r.db.Exec(`UPDATE users SET email_verified = true, updated_at = $1 WHERE id = $2`,
				time.Now().UTC(), user.ID)
			if err != nil {
				return nil, false, err
			}
			user.EmailVerified = true
		}
		return user, false, nil
	}

	// Check if user exists with this email (might have signed up with password first)
	user, err = r.FindByEmail(email)
	if err != nil {
		return nil, false, err
	}

	if user != nil {
		// Link Google account to existing user and mark as verified
		_, err = r.db.Exec(`UPDATE users SET google_id = $1, email_verified = true, verification_token = NULL, verification_token_expires_at = NULL, updated_at = $2 WHERE id = $3`,
			googleID, time.Now().UTC(), user.ID)
		if err != nil {
			return nil, false, err
		}
		user.GoogleID = sql.NullString{String: googleID, Valid: true}
		user.EmailVerified = true
		return user, false, nil
	}

	// Create new user with Google - automatically verified
	user, err = r.Create(googleID, email, true)
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}

func (r *UserRepository) FindByVerificationToken(token string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, google_id, email, password_hash, email_verified, verification_token, verification_token_expires_at,
		       COALESCE(wealth_status, 'none'), wealth_status_expires_at, created_at, updated_at
		FROM users WHERE verification_token = $1
	`, token).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.EmailVerified, 
		&user.VerificationToken, &user.VerificationTokenExpiresAt, &user.WealthStatus, &user.WealthStatusExpiresAt,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) VerifyEmail(userID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE users SET email_verified = true, verification_token = NULL, verification_token_expires_at = NULL, updated_at = $1 WHERE id = $2
	`, time.Now().UTC(), userID)
	return err
}

func (r *UserRepository) UpdateVerificationToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		UPDATE users SET verification_token = $1, verification_token_expires_at = $2, updated_at = $3 WHERE id = $4
	`, token, expiresAt, time.Now().UTC(), userID)
	return err
}

func (r *UserRepository) Delete(userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	return err
}

// UpdateWealthStatus updates the user's wealth status and optional expiry
func (r *UserRepository) UpdateWealthStatus(userID uuid.UUID, status models.WealthStatus, expiresAt *time.Time) error {
	if expiresAt != nil {
		_, err := r.db.Exec(`
			UPDATE users SET wealth_status = $1, wealth_status_expires_at = $2, updated_at = $3 WHERE id = $4
		`, status, expiresAt, time.Now().UTC(), userID)
		return err
	}
	_, err := r.db.Exec(`
		UPDATE users SET wealth_status = $1, wealth_status_expires_at = NULL, updated_at = $2 WHERE id = $3
	`, status, time.Now().UTC(), userID)
	return err
}

// GetUserWithGender returns user info along with their gender from profile (for messaging rules)
func (r *UserRepository) GetUserWithGender(userID uuid.UUID) (*models.User, models.Gender, bool, error) {
	user, err := r.FindByID(userID)
	if err != nil || user == nil {
		return nil, "", false, err
	}
	
	var gender models.Gender
	var personVerified bool
	err = r.db.QueryRow(`
		SELECT gender, is_verified FROM profiles WHERE user_id = $1
	`, userID).Scan(&gender, &personVerified)
	
	if err == sql.ErrNoRows {
		return user, "", false, nil
	}
	if err != nil {
		return nil, "", false, err
	}
	
	return user, gender, personVerified, nil
}

// ExpireWealthStatuses downgrades users whose wealth status has expired
func (r *UserRepository) ExpireWealthStatuses() (int64, error) {
	result, err := r.db.Exec(`
		UPDATE users 
		SET wealth_status = 'none', wealth_status_expires_at = NULL, updated_at = $1
		WHERE wealth_status != 'none' 
		AND wealth_status_expires_at IS NOT NULL 
		AND wealth_status_expires_at < $1
	`, time.Now().UTC())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
