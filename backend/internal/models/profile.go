package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type SalaryRange string

const (
	SalaryRange5To10   SalaryRange = "5-10 LPA"
	SalaryRange10To20  SalaryRange = "10-20 LPA"
	SalaryRange20To50  SalaryRange = "20-50 LPA"
	SalaryRange50Plus  SalaryRange = "50+ LPA"
)

// NullString is a wrapper around sql.NullString that properly marshals to JSON
type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}

type Profile struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	DisplayName  string     `json:"display_name" db:"display_name"`
	Gender       Gender     `json:"gender" db:"gender"`
	Age          int        `json:"age" db:"age"`
	Bio          string     `json:"bio" db:"bio"`
	SalaryRange  NullString `json:"salary_range,omitempty" db:"salary_range"`
	City         string     `json:"city" db:"city"`
	State        string     `json:"state" db:"state"`
	Latitude     float64    `json:"latitude" db:"latitude"`
	Longitude    float64    `json:"longitude" db:"longitude"`
	IsComplete   bool       `json:"is_complete" db:"is_complete"`
	IsVerified   bool       `json:"is_verified" db:"is_verified"`
	IsFake       bool       `json:"is_fake" db:"is_fake"`
	ProfileScore float64    `json:"profile_score" db:"profile_score"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type ProfileImage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	S3Key     string    `json:"s3_key" db:"s3_key"`
	URL       string    `json:"url" db:"url"`
	IsPrimary bool      `json:"is_primary" db:"is_primary"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ProfileWithImages struct {
	Profile
	Images       []ProfileImage `json:"images"`
	IsOnline     bool           `json:"is_online"`
	LastSeen     *time.Time     `json:"last_seen,omitempty"`
	Distance     *float64       `json:"distance_km,omitempty"`
	IsLiked      bool           `json:"is_liked"`
	WealthStatus string         `json:"wealth_status,omitempty"`
	HasLikedMe   bool           `json:"has_liked_me"`
}

type CreateProfileRequest struct {
	DisplayName string  `json:"display_name" binding:"required,min=2,max=50"`
	Gender      Gender  `json:"gender" binding:"required"`
	Age         int     `json:"age" binding:"required,gte=21,lte=100"`
	Bio         string  `json:"bio" binding:"required,max=500"`
	SalaryRange string  `json:"salary_range,omitempty"`
	City        string  `json:"city" binding:"required"`
	State       string  `json:"state" binding:"required"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type UpdateProfileRequest struct {
	DisplayName *string  `json:"display_name,omitempty"`
	Age         *int     `json:"age,omitempty"`
	Bio         *string  `json:"bio,omitempty"`
	SalaryRange *string  `json:"salary_range,omitempty"`
	City        *string  `json:"city,omitempty"`
	State       *string  `json:"state,omitempty"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}

type ListProfilesQuery struct {
	Page        int     `form:"page,default=1"`
	Limit       int     `form:"limit,default=20"`
	Gender      string  `form:"gender"`
	City        string  `form:"city"`
	State       string  `form:"state"`
	MinAge      int     `form:"min_age"`
	MaxAge      int     `form:"max_age"`
	MaxDistance float64 `form:"max_distance"`
	OnlineOnly  bool    `form:"online_only"`
}
