package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type PresenceRepository struct {
	db *sql.DB
}

func NewPresenceRepository(db *sql.DB) *PresenceRepository {
	return &PresenceRepository{db: db}
}

func (r *PresenceRepository) Upsert(userID uuid.UUID, isOnline bool) error {
	now := time.Now().UTC()
	_, err := r.db.Exec(`
		INSERT INTO user_presence (user_id, is_online, last_seen, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			is_online = $2,
			last_seen = CASE WHEN $2 = false THEN $3 ELSE user_presence.last_seen END,
			updated_at = $4
	`, userID, isOnline, now, now)
	return err
}

func (r *PresenceRepository) SetOnline(userID uuid.UUID) error {
	return r.Upsert(userID, true)
}

func (r *PresenceRepository) SetOffline(userID uuid.UUID) error {
	return r.Upsert(userID, false)
}

func (r *PresenceRepository) GetPresence(userID uuid.UUID) (*models.UserPresence, error) {
	presence := &models.UserPresence{}
	err := r.db.QueryRow(`
		SELECT user_id, is_online, last_seen, updated_at
		FROM user_presence WHERE user_id = $1
	`, userID).Scan(&presence.UserID, &presence.IsOnline, &presence.LastSeen, &presence.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return presence, nil
}

func (r *PresenceRepository) GetOnlineUsers(userIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(userIDs) == 0 {
		return make(map[uuid.UUID]bool), nil
	}

	query := "SELECT user_id, is_online FROM user_presence WHERE user_id = ANY($1)"

	rows, err := r.db.Query(query, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]bool)
	for rows.Next() {
		var userID uuid.UUID
		var isOnline bool
		rows.Scan(&userID, &isOnline)
		result[userID] = isOnline
	}

	return result, nil
}
