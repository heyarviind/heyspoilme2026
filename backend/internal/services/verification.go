package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
)

type VerificationService struct {
	verificationRepo *repository.VerificationRepository
	profileRepo      *repository.ProfileRepository
}

func NewVerificationService(verificationRepo *repository.VerificationRepository, profileRepo *repository.ProfileRepository) *VerificationService {
	return &VerificationService{
		verificationRepo: verificationRepo,
		profileRepo:      profileRepo,
	}
}

// GenerateVerificationCode creates a random 6-digit code for video verification
func (s *VerificationService) GenerateVerificationCode() (string, error) {
	// Generate a random 6-digit number
	max := big.NewInt(999999)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	// Format with leading zeros to ensure 6 digits
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// SubmitVerification creates a new verification request
func (s *VerificationService) SubmitVerification(userID uuid.UUID, req *models.CreateVerificationRequest, code string) (*models.VerificationRequest, error) {
	// Check if profile exists
	profile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	// Check if already verified
	if profile.IsVerified {
		return nil, errors.New("profile is already verified")
	}

	// Check if there's a pending request
	hasPending, err := s.verificationRepo.HasPendingRequest(userID)
	if err != nil {
		return nil, err
	}
	if hasPending {
		return nil, errors.New("you already have a pending verification request")
	}

	// Validate document type
	switch req.DocumentType {
	case models.DocumentTypeAadhar, models.DocumentTypePassport, models.DocumentTypeDrivingLicense:
		// Valid
	default:
		return nil, errors.New("invalid document type")
	}

	return s.verificationRepo.Create(userID, req, code)
}

// GetVerificationStatus returns the current verification status for a user
func (s *VerificationService) GetVerificationStatus(userID uuid.UUID) (*models.VerificationStatusResponse, error) {
	// First check if profile is already verified
	profile, err := s.profileRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if profile != nil && profile.IsVerified {
		return &models.VerificationStatusResponse{
			Status: models.VerificationStatusApproved,
		}, nil
	}

	// Check for any verification request
	req, err := s.verificationRepo.GetLatestByUserID(userID)
	if err != nil {
		return nil, err
	}

	if req == nil {
		return &models.VerificationStatusResponse{
			Status: "none",
		}, nil
	}

	return &models.VerificationStatusResponse{
		Status:          req.Status,
		RejectionReason: req.RejectionReason,
		CreatedAt:       &req.CreatedAt,
	}, nil
}



