package services

import (
	"database/sql"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
)

type RankingService struct {
	db       *sql.DB
	interval time.Duration
	stopChan chan struct{}
}

func NewRankingService(db *sql.DB) *RankingService {
	return &RankingService{
		db:       db,
		interval: 15 * time.Minute,
		stopChan: make(chan struct{}),
	}
}

// ProfileScoreData contains all the data needed to calculate a profile's static score
type ProfileScoreData struct {
	UserID        uuid.UUID
	EmailVerified bool
	IsVerified    bool
	BioLength     int
	PhotoCount    int
	HasSalary     bool
	LikesReceived int
	ResponseRate  float64
	CreatedAt     time.Time
}

// CalculateStaticScore computes the static score for a single profile
func (s *RankingService) CalculateStaticScore(data *ProfileScoreData) float64 {
	var score float64

	// Email verified: +10
	if data.EmailVerified {
		score += 10
	}

	// Person verified (is_verified): +25
	if data.IsVerified {
		score += 25
	}

	// Profile completeness: 0-20
	// - Photo count: (photo_count/5)*10, capped at 10
	photoScore := math.Min(10, float64(data.PhotoCount)/5*10)
	// - Bio length: (bio_length/300)*5, capped at 5
	bioScore := math.Min(5, float64(data.BioLength)/300*5)
	// - Has salary: +5
	salaryScore := 0.0
	if data.HasSalary {
		salaryScore = 5
	}
	score += photoScore + bioScore + salaryScore

	// Popularity: min(15, likes_received * 0.5)
	popularityScore := math.Min(15, float64(data.LikesReceived)*0.5)
	score += popularityScore

	// Response rate: response_rate_percent * 0.15
	responseScore := data.ResponseRate * 0.15
	score += responseScore

	// New user boost: decay over 7 days
	// 10 * max(0, 1 - days_since_creation/7)
	daysSinceCreation := time.Since(data.CreatedAt).Hours() / 24
	newUserBoost := 10 * math.Max(0, 1-daysSinceCreation/7)
	score += newUserBoost

	return score
}

// CalculateScoreForUser calculates and returns the static score for a specific user
func (s *RankingService) CalculateScoreForUser(userID uuid.UUID) (float64, error) {
	data, err := s.getProfileScoreData(userID)
	if err != nil {
		return 0, err
	}
	if data == nil {
		return 0, nil
	}
	return s.CalculateStaticScore(data), nil
}

// UpdateScoreForUser calculates and updates the profile score for a specific user
func (s *RankingService) UpdateScoreForUser(userID uuid.UUID) error {
	score, err := s.CalculateScoreForUser(userID)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`UPDATE profiles SET profile_score = $1 WHERE user_id = $2`, score, userID)
	return err
}

// UpdateAllScores recalculates and updates scores for all profiles
func (s *RankingService) UpdateAllScores() error {
	log.Printf("[RankingJob] Starting profile score update")
	startTime := time.Now()

	// Get all user IDs with profiles
	rows, err := s.db.Query(`SELECT user_id FROM profiles WHERE is_complete = true`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			continue
		}
		userIDs = append(userIDs, userID)
	}

	updated := 0
	for _, userID := range userIDs {
		if err := s.UpdateScoreForUser(userID); err != nil {
			log.Printf("[RankingJob] Error updating score for user %s: %v", userID, err)
			continue
		}
		updated++
	}

	log.Printf("[RankingJob] Updated %d profile scores in %v", updated, time.Since(startTime))
	return nil
}

// getProfileScoreData fetches all data needed to calculate a profile's score
func (s *RankingService) getProfileScoreData(userID uuid.UUID) (*ProfileScoreData, error) {
	data := &ProfileScoreData{UserID: userID}

	// Get user and profile data
	err := s.db.QueryRow(`
		SELECT u.email_verified, p.is_verified, LENGTH(p.bio), 
		       CASE WHEN p.salary_range IS NOT NULL AND p.salary_range != '' THEN true ELSE false END,
		       p.created_at
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE u.id = $1 AND p.is_complete = true
	`, userID).Scan(&data.EmailVerified, &data.IsVerified, &data.BioLength, &data.HasSalary, &data.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get photo count
	s.db.QueryRow(`SELECT COUNT(*) FROM profile_images WHERE user_id = $1`, userID).Scan(&data.PhotoCount)

	// Get likes received count
	s.db.QueryRow(`SELECT COUNT(*) FROM likes WHERE liked_id = $1`, userID).Scan(&data.LikesReceived)

	// Calculate response rate
	// Response rate = (messages sent in reply / conversations where user received a message) * 100
	var conversationsReceived, conversationsReplied int

	// Count conversations where this user received at least one message (is not the only sender)
	s.db.QueryRow(`
		SELECT COUNT(DISTINCT cp.conversation_id)
		FROM conversation_participants cp
		JOIN messages m ON cp.conversation_id = m.conversation_id AND m.sender_id != cp.user_id
		WHERE cp.user_id = $1
	`, userID).Scan(&conversationsReceived)

	// Count conversations where this user sent at least one reply
	s.db.QueryRow(`
		SELECT COUNT(DISTINCT m.conversation_id)
		FROM messages m
		JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id AND cp.user_id = $1
		WHERE m.sender_id = $1
		AND EXISTS (
			SELECT 1 FROM messages m2 
			WHERE m2.conversation_id = m.conversation_id 
			AND m2.sender_id != $1 
			AND m2.created_at < m.created_at
		)
	`, userID).Scan(&conversationsReplied)

	if conversationsReceived > 0 {
		data.ResponseRate = float64(conversationsReplied) / float64(conversationsReceived) * 100
	} else {
		// Default to 50% if no conversations to measure
		data.ResponseRate = 50
	}

	return data, nil
}

// Start begins the background job that updates profile scores
func (s *RankingService) Start() {
	log.Printf("[RankingJob] Starting profile ranking job (updating every %v)", s.interval)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately on start
	s.UpdateAllScores()

	for {
		select {
		case <-ticker.C:
			s.UpdateAllScores()
		case <-s.stopChan:
			log.Printf("[RankingJob] Stopping ranking job")
			return
		}
	}
}

// Stop stops the background job
func (s *RankingService) Stop() {
	close(s.stopChan)
}


