package services

import (
	"errors"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
	userRepo    *repository.UserRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository, userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

func (s *ProfileService) CreateProfile(userID uuid.UUID, req *models.CreateProfileRequest) (*models.Profile, error) {
	if req.Gender == models.GenderMale && req.SalaryRange == "" {
		return nil, errors.New("salary range is required for male users")
	}

	existing, _ := s.profileRepo.FindByUserID(userID)
	if existing != nil {
		return nil, errors.New("profile already exists")
	}

	return s.profileRepo.Create(userID, req)
}

func (s *ProfileService) GetProfile(userID uuid.UUID) (*models.Profile, error) {
	return s.profileRepo.FindByUserID(userID)
}

func (s *ProfileService) GetProfileWithDetails(profileUserID, requestingUserID uuid.UUID) (*models.ProfileWithImages, error) {
	return s.profileRepo.GetProfileWithDetails(profileUserID, requestingUserID)
}

func (s *ProfileService) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error) {
	return s.profileRepo.Update(userID, req)
}

func (s *ProfileService) ListProfiles(requestingUserID uuid.UUID, query *models.ListProfilesQuery) ([]models.ProfileWithImages, int, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 50 {
		query.Limit = 20
	}

	// Get user's profile for location, use defaults if not found
	var userLat, userLng float64 = 0, 0
	profile, err := s.profileRepo.FindByUserID(requestingUserID)
	if err == nil && profile != nil {
		userLat = profile.Latitude
		userLng = profile.Longitude
	}

	return s.profileRepo.ListProfiles(requestingUserID, userLat, userLng, query)
}

func (s *ProfileService) AddProfileImage(userID uuid.UUID, s3Key, url string, isPrimary bool) (*models.ProfileImage, error) {
	return s.profileRepo.AddImage(userID, s3Key, url, isPrimary)
}

func (s *ProfileService) DeleteProfileImage(imageID, userID uuid.UUID) error {
	return s.profileRepo.DeleteImage(imageID, userID)
}

func (s *ProfileService) GetProfileImages(userID uuid.UUID) ([]models.ProfileImage, error) {
	return s.profileRepo.GetImages(userID)
}
