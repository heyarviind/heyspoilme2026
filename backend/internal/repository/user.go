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

func (r *UserRepository) Create(googleID, email string) (*models.User, error) {
	user := &models.User{
		ID:        uuid.New(),
		GoogleID:  sql.NullString{String: googleID, Valid: googleID != ""},
		Email:     email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO users (id, google_id, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.GoogleID, user.Email, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateWithPassword(email, passwordHash string) (*models.User, error) {
	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: sql.NullString{String: passwordHash, Valid: true},
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT id, google_id, email, password_hash, created_at, updated_at
		FROM users WHERE google_id = $1
	`, googleID).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

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
		SELECT id, google_id, email, password_hash, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

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
		SELECT id, google_id, email, password_hash, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.GoogleID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

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
		return user, false, nil
	}

	// Check if user exists with this email (might have signed up with password first)
	user, err = r.FindByEmail(email)
	if err != nil {
		return nil, false, err
	}

	if user != nil {
		// Link Google account to existing user
		_, err = r.db.Exec(`UPDATE users SET google_id = $1, updated_at = $2 WHERE id = $3`,
			googleID, time.Now().UTC(), user.ID)
		if err != nil {
			return nil, false, err
		}
		user.GoogleID = sql.NullString{String: googleID, Valid: true}
		return user, false, nil
	}

	user, err = r.Create(googleID, email)
	if err != nil {
		return nil, false, err
	}

	return user, true, nil
}
