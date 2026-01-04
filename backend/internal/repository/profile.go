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
		ID:         uuid.New(),
		UserID:     userID,
		Gender:     req.Gender,
		Age:        req.Age,
		Bio:        req.Bio,
		City:       req.City,
		State:      req.State,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		IsComplete: true,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	if req.Gender == models.GenderMale && req.SalaryRange != "" {
		profile.SalaryRange = models.NullString{NullString: sql.NullString{String: req.SalaryRange, Valid: true}}
	}

	_, err := r.db.Exec(`
		INSERT INTO profiles (id, user_id, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, profile.ID, profile.UserID, profile.Gender, profile.Age, profile.Bio, profile.SalaryRange,
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
		SELECT id, user_id, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, is_verified, created_at, updated_at
		FROM profiles WHERE user_id = $1
	`, userID).Scan(&profile.ID, &profile.UserID, &profile.Gender, &profile.Age, &profile.Bio,
		&profile.SalaryRange, &profile.City, &profile.State, &profile.Latitude, &profile.Longitude,
		&profile.IsComplete, &profile.IsVerified, &profile.CreatedAt, &profile.UpdatedAt)

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
		SELECT id, user_id, gender, age, bio, salary_range, city, state, latitude, longitude, is_complete, is_verified, created_at, updated_at
		FROM profiles WHERE id = $1
	`, id).Scan(&profile.ID, &profile.UserID, &profile.Gender, &profile.Age, &profile.Bio,
		&profile.SalaryRange, &profile.City, &profile.State, &profile.Latitude, &profile.Longitude,
		&profile.IsComplete, &profile.IsVerified, &profile.CreatedAt, &profile.UpdatedAt)

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

func (r *ProfileRepository) ListProfiles(requestingUserID uuid.UUID, userLat, userLng float64, query *models.ListProfilesQuery) ([]models.ProfileWithImages, int, error) {
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
	var orderBy string
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
		orderBy = "distance ASC"
	} else {
		distanceCalc = "NULL"
		orderBy = "p.created_at DESC"
	}

	if query.MaxDistance > 0 && userLat != 0 && userLng != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("%s <= $%d", distanceCalc, argIndex))
		args = append(args, query.MaxDistance)
		argIndex++
		whereClause = strings.Join(whereClauses, " AND ")
	}

	offset := (query.Page - 1) * query.Limit
	args = append(args, query.Limit, offset)

	mainQuery := fmt.Sprintf(`
		SELECT p.id, p.user_id, p.gender, p.age, p.bio, p.salary_range, p.city, p.state, 
			   p.latitude, p.longitude, p.is_complete, p.is_verified, p.created_at, p.updated_at,
			   %s as distance,
			   COALESCE(up.is_online, false) as is_online,
			   up.last_seen,
			   EXISTS(SELECT 1 FROM likes WHERE liker_id = $1 AND liked_id = p.user_id) as is_liked
		FROM profiles p
		LEFT JOIN user_presence up ON p.user_id = up.user_id
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, distanceCalc, whereClause, orderBy, argIndex, argIndex+1)

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
		err := rows.Scan(&p.ID, &p.UserID, &p.Gender, &p.Age, &p.Bio, &p.SalaryRange,
			&p.City, &p.State, &p.Latitude, &p.Longitude, &p.IsComplete, &p.IsVerified, &p.CreatedAt, &p.UpdatedAt,
			&distance, &p.IsOnline, &lastSeen, &p.IsLiked)
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

	return result, nil
}
