package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(likerID, likedID uuid.UUID) (*models.Like, error) {
	like := &models.Like{
		ID:        uuid.New(),
		LikerID:   likerID,
		LikedID:   likedID,
		CreatedAt: time.Now().UTC(),
	}

	_, err := r.db.Exec(`
		INSERT INTO likes (id, liker_id, liked_id, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (liker_id, liked_id) DO NOTHING
	`, like.ID, like.LikerID, like.LikedID, like.CreatedAt)

	if err != nil {
		return nil, err
	}

	return like, nil
}

func (r *LikeRepository) Delete(likerID, likedID uuid.UUID) error {
	_, err := r.db.Exec(`
		DELETE FROM likes WHERE liker_id = $1 AND liked_id = $2
	`, likerID, likedID)
	return err
}

func (r *LikeRepository) Exists(likerID, likedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM likes WHERE liker_id = $1 AND liked_id = $2)
	`, likerID, likedID).Scan(&exists)
	return exists, err
}

func (r *LikeRepository) GetReceivedLikes(userID uuid.UUID, limit, offset int) ([]models.Like, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM likes WHERE liked_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, liker_id, liked_id, created_at
		FROM likes WHERE liked_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var likes []models.Like
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.LikerID, &like.LikedID, &like.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		likes = append(likes, like)
	}

	return likes, total, nil
}

func (r *LikeRepository) GetGivenLikes(userID uuid.UUID, limit, offset int) ([]models.Like, int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM likes WHERE liker_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT id, liker_id, liked_id, created_at
		FROM likes WHERE liker_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var likes []models.Like
	for rows.Next() {
		var like models.Like
		err := rows.Scan(&like.ID, &like.LikerID, &like.LikedID, &like.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		likes = append(likes, like)
	}

	return likes, total, nil
}
