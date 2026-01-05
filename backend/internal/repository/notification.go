package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(userID uuid.UUID, notifType models.NotificationType, data *models.NotificationData) (*models.Notification, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	notif := &models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Data:      jsonData,
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	}

	_, err = r.db.Exec(`
		INSERT INTO notifications (id, user_id, type, data, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, notif.ID, notif.UserID, notif.Type, notif.Data, notif.IsRead, notif.CreatedAt)

	if err != nil {
		return nil, err
	}

	return notif, nil
}

func (r *NotificationRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, user_id, type, data, is_read, created_at
		FROM notifications WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notif models.Notification
		err := rows.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Data, &notif.IsRead, &notif.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(notifID, userID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2
	`, notifID, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE notifications SET is_read = true WHERE user_id = $1
	`, userID)
	return err
}

func (r *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false
	`, userID).Scan(&count)
	return count, err
}

func (r *NotificationRepository) DeleteAllForUser(userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM notifications WHERE user_id = $1`, userID)
	return err
}
