package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(userID uuid.UUID, req *models.CreateProfileRequest) (*models.Profile, error) {
	profile := &models.Profile{
		ID:          uuid.New(),
		UserID:      userID,
		DisplayName: req.DisplayName,
		Gender:      req.Gender,
		Age:         req.Age,
		Bio:         req.Bio,
		City:        req.City,
		State:       req.State,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		IsComplete:  true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if req.Gender == models.GenderMale && req.SalaryRange != "" {
		profile.SalaryRange = models.NullString{NullString: sql.NullString{String: req.SalaryRange, Valid: true}}
	}

	_, err := r.db.Exec(`
		INSERT INTO profiles (id, user_id, display_name, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, profile.ID, profile.UserID, profile.DisplayName, profile.Gender, profile.Age, profile.Bio, profile.SalaryRange,
		profile.City, profile.State, profile.Latitude, profile.Longitude, profile.IsComplete,
		profile.CreatedAt, profile.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) FindByUserID(userID uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	err := r.db.QueryRow(`
		SELECT id, user_id, display_name, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, is_verified, profile_score, created_at, updated_at
		FROM profiles WHERE user_id = $1
	`, userID).Scan(&profile.ID, &profile.UserID, &profile.DisplayName, &profile.Gender, &profile.Age, &profile.Bio,
		&profile.SalaryRange, &profile.City, &profile.State, &profile.Latitude, &profile.Longitude,
		&profile.IsComplete, &profile.IsVerified, &profile.ProfileScore, &profile.CreatedAt, &profile.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) FindByID(id uuid.UUID) (*models.Profile, error) {
	profile := &models.Profile{}
	err := r.db.QueryRow(`
		SELECT id, user_id, display_name, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, is_verified, profile_score, created_at, updated_at
		FROM profiles WHERE id = $1
	`, id).Scan(&profile.ID, &profile.UserID, &profile.DisplayName, &profile.Gender, &profile.Age, &profile.Bio,
		&profile.SalaryRange, &profile.City, &profile.State, &profile.Latitude, &profile.Longitude,
		&profile.IsComplete, &profile.IsVerified, &profile.ProfileScore, &profile.CreatedAt, &profile.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *ProfileRepository) Update(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error) {
	setClauses := []string{"updated_at = $1"}
	args := []interface{}{time.Now().UTC()}
	argIndex := 2

	if req.DisplayName != nil {
		setClauses = append(setClauses, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, *req.DisplayName)
		argIndex++
	}
	if req.Age != nil {
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", argIndex))
		args = append(args, *req.Age)
		argIndex++
	}
	if req.Bio != nil {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argIndex))
		args = append(args, *req.Bio)
		argIndex++
	}
	if req.SalaryRange != nil {
		setClauses = append(setClauses, fmt.Sprintf("salary_range = $%d", argIndex))
		args = append(args, *req.SalaryRange)
		argIndex++
	}
	if req.City != nil {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *req.City)
		argIndex++
	}
	if req.State != nil {
		setClauses = append(setClauses, fmt.Sprintf("state = $%d", argIndex))
		args = append(args, *req.State)
		argIndex++
	}
	if req.Latitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("latitude = $%d", argIndex))
		args = append(args, *req.Latitude)
		argIndex++
	}
	if req.Longitude != nil {
		setClauses = append(setClauses, fmt.Sprintf("longitude = $%d", argIndex))
		args = append(args, *req.Longitude)
		argIndex++
	}

	args = append(args, userID)
	query := fmt.Sprintf("UPDATE profiles SET %s WHERE user_id = $%d", strings.Join(setClauses, ", "), argIndex)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.FindByUserID(userID)
}

func (r *ProfileRepository) GetImages(userID uuid.UUID) ([]models.ProfileImage, error) {
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

	return images, nil
}

func (r *ProfileRepository) AddImage(userID uuid.UUID, s3Key, url string, isPrimary bool) (*models.ProfileImage, error) {
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

func (r *ProfileRepository) DeleteImage(imageID, userID uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM profile_images WHERE id = $1 AND user_id = $2", imageID, userID)
	return err
}

func (r *ProfileRepository) GetImageS3Keys(userID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(`SELECT s3_key FROM profile_images WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (r *ProfileRepository) DeleteAllImages(userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM profile_images WHERE user_id = $1`, userID)
	return err
}

func (r *ProfileRepository) Delete(userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM profiles WHERE user_id = $1`, userID)
	return err
}

func (r *ProfileRepository) ListProfiles(requestingUserID uuid.UUID, userLat, userLng float64, query *models.ListProfilesQuery) ([]models.ProfileWithImages, int, error) {
	// Determine if the requesting user is browsing males (i.e., they're female)
	// This affects ranking: females see males ranked by wealth_status first
	var requestingUserGender models.Gender
	r.db.QueryRow(`SELECT gender FROM profiles WHERE user_id = $1`, requestingUserID).Scan(&requestingUserGender)
	isFemaleViewingMales := requestingUserGender == models.GenderFemale && query.Gender == "male"

	whereClauses := []string{"p.user_id != $1", "p.is_complete = true"}
	args := []interface{}{requestingUserID}
	argIndex := 2

	if query.Gender != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.gender = $%d", argIndex))
		args = append(args, query.Gender)
		argIndex++
	}
	if query.City != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.city ILIKE $%d", argIndex))
		args = append(args, "%"+query.City+"%")
		argIndex++
	}
	if query.State != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.state ILIKE $%d", argIndex))
		args = append(args, "%"+query.State+"%")
		argIndex++
	}
	if query.MinAge > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("p.age >= $%d", argIndex))
		args = append(args, query.MinAge)
		argIndex++
	}
	if query.MaxAge > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("p.age <= $%d", argIndex))
		args = append(args, query.MaxAge)
		argIndex++
	}
	if query.OnlineOnly {
		whereClauses = append(whereClauses, "EXISTS(SELECT 1 FROM user_presence up2 WHERE up2.user_id = p.user_id AND up2.is_online = true)")
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Count total before adding distance params
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM profiles p WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Distance calculation - use CASE to handle zero lat/lng
	var distanceCalc string
	if userLat != 0 || userLng != 0 {
		distanceCalc = fmt.Sprintf(`
			CASE WHEN p.latitude = 0 AND p.longitude = 0 THEN 99999
			ELSE (6371 * acos(
				LEAST(1.0, GREATEST(-1.0,
					cos(radians($%d)) * cos(radians(p.latitude)) *
					cos(radians(p.longitude) - radians($%d)) +
					sin(radians($%d)) * sin(radians(p.latitude))
				))
			))
			END
		`, argIndex, argIndex+1, argIndex+2)
		args = append(args, userLat, userLng, userLat)
		argIndex += 3
	} else {
		distanceCalc = "NULL"
	}

	if query.MaxDistance > 0 && userLat != 0 && userLng != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("%s <= $%d", distanceCalc, argIndex))
		args = append(args, query.MaxDistance)
		argIndex++
		whereClause = strings.Join(whereClauses, " AND ")
	}

	offset := (query.Page - 1) * query.Limit
	args = append(args, query.Limit, offset)

	// Dynamic scoring calculation:
	// - Online status: +15 if online, +5 if seen in last hour
	// - Mutual interest: +30 if they already liked the requesting user
	// - Distance penalty: -min(20, distance_km * 0.1) for distance > 50km
	var distancePenaltyCalc string
	if userLat != 0 || userLng != 0 {
		distancePenaltyCalc = fmt.Sprintf(`
			CASE 
				WHEN %s < 50 THEN 0
				ELSE -LEAST(20, %s * 0.1)
			END
		`, distanceCalc, distanceCalc)
	} else {
		distancePenaltyCalc = "0"
	}

	// Base scoring formula
	var finalScoreCalc string
	if isFemaleViewingMales {
		// Female viewing males: wealth_status is primary factor
		// wealth_status: high=300, medium=200, low=100, none=0
		// Then: person_verified (+50), distance bucket, activity, profile_score
		finalScoreCalc = fmt.Sprintf(`
			CASE COALESCE(u.wealth_status, 'none')
				WHEN 'high' THEN 300
				WHEN 'medium' THEN 200
				WHEN 'low' THEN 100
				ELSE 0 END +
			CASE WHEN p.is_verified THEN 50 ELSE 0 END +
			CASE WHEN COALESCE(up.is_online, false) THEN 30
				 WHEN up.last_seen > NOW() - INTERVAL '1 hour' THEN 20
				 WHEN up.last_seen > NOW() - INTERVAL '1 day' THEN 10
				 ELSE 0 END +
			CASE WHEN EXISTS(SELECT 1 FROM likes WHERE liker_id = p.user_id AND liked_id = $1) THEN 40
				 ELSE 0 END +
			(p.profile_score * 0.5) +
			%s
		`, distancePenaltyCalc)
	} else {
		// Default ranking (male viewing females, or any other case)
		finalScoreCalc = fmt.Sprintf(`
			p.profile_score +
			CASE WHEN COALESCE(up.is_online, false) THEN 15
				 WHEN up.last_seen > NOW() - INTERVAL '1 hour' THEN 5
				 ELSE 0 END +
			CASE WHEN EXISTS(SELECT 1 FROM likes WHERE liker_id = p.user_id AND liked_id = $1) THEN 30
				 ELSE 0 END +
			%s
		`, distancePenaltyCalc)
	}

	mainQuery := fmt.Sprintf(`
		SELECT p.id, p.user_id, p.display_name, p.gender, p.age, p.bio, p.salary_range, p.city, p.state, 
			   p.latitude, p.longitude, p.is_complete, p.is_verified, p.profile_score, p.created_at, p.updated_at,
			   %s as distance,
			   COALESCE(up.is_online, false) as is_online,
			   up.last_seen,
			   EXISTS(SELECT 1 FROM likes WHERE liker_id = $1 AND liked_id = p.user_id) as is_liked,
			   EXISTS(SELECT 1 FROM likes WHERE liker_id = p.user_id AND liked_id = $1) as has_liked_me,
			   COALESCE(u.wealth_status, 'none') as wealth_status,
			   (%s) as final_score
		FROM profiles p
		LEFT JOIN user_presence up ON p.user_id = up.user_id
		LEFT JOIN users u ON p.user_id = u.id
		WHERE %s
		ORDER BY final_score DESC
		LIMIT $%d OFFSET $%d
	`, distanceCalc, finalScoreCalc, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(mainQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var profiles []models.ProfileWithImages
	for rows.Next() {
		var p models.ProfileWithImages
		var distance sql.NullFloat64
		var lastSeen sql.NullTime
		var finalScore float64 // scanned but not stored - used only for ordering
		err := rows.Scan(&p.ID, &p.UserID, &p.DisplayName, &p.Gender, &p.Age, &p.Bio, &p.SalaryRange,
			&p.City, &p.State, &p.Latitude, &p.Longitude, &p.IsComplete, &p.IsVerified, &p.ProfileScore, &p.CreatedAt, &p.UpdatedAt,
			&distance, &p.IsOnline, &lastSeen, &p.IsLiked, &p.HasLikedMe, &p.WealthStatus, &finalScore)
		if err != nil {
			return nil, 0, err
		}
		if distance.Valid {
			p.Distance = &distance.Float64
		}
		if lastSeen.Valid {
			p.LastSeen = &lastSeen.Time
		}

		images, _ := r.GetImages(p.UserID)
		p.Images = images

		profiles = append(profiles, p)
	}

	return profiles, total, nil
}

func (r *ProfileRepository) GetProfileWithDetails(profileUserID, requestingUserID uuid.UUID) (*models.ProfileWithImages, error) {
	profile, err := r.FindByUserID(profileUserID)
	if err != nil || profile == nil {
		return nil, err
	}

	result := &models.ProfileWithImages{
		Profile: *profile,
	}

	images, _ := r.GetImages(profileUserID)
	result.Images = images

	var isOnline bool
	var lastSeen sql.NullTime
	r.db.QueryRow(`
		SELECT COALESCE(is_online, false), last_seen 
		FROM user_presence WHERE user_id = $1
	`, profileUserID).Scan(&isOnline, &lastSeen)
	result.IsOnline = isOnline
	if lastSeen.Valid {
		result.LastSeen = &lastSeen.Time
	}

	var isLiked bool
	r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM likes WHERE liker_id = $1 AND liked_id = $2)
	`, requestingUserID, profileUserID).Scan(&isLiked)
	result.IsLiked = isLiked

	// Check if profile owner has liked the requesting user
	var hasLikedMe bool
	r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM likes WHERE liker_id = $1 AND liked_id = $2)
	`, profileUserID, requestingUserID).Scan(&hasLikedMe)
	result.HasLikedMe = hasLikedMe

	// Get wealth status from users table
	var wealthStatus sql.NullString
	r.db.QueryRow(`
		SELECT COALESCE(wealth_status, 'none') FROM users WHERE id = $1
	`, profileUserID).Scan(&wealthStatus)
	if wealthStatus.Valid {
		result.WealthStatus = wealthStatus.String
	}

	return result, nil
}
