package services

import (
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/pkg/storage"
)

type AdminService struct {
	adminRepo *repository.AdminRepository
	s3Client  *storage.S3Client
}

func NewAdminService(adminRepo *repository.AdminRepository, s3Client *storage.S3Client) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		s3Client:  s3Client,
	}
}

// ListUsers returns paginated list of users
func (s *AdminService) ListUsers(limit, offset int, search string) ([]repository.AdminUser, int, error) {
	return s.adminRepo.ListUsers(limit, offset, search)
}

// GetUserWithImages returns a user with all their images
func (s *AdminService) GetUserWithImages(userID uuid.UUID) (*repository.AdminUserWithImages, error) {
	return s.adminRepo.GetUserWithImages(userID)
}

// ListMessages returns paginated list of messages
func (s *AdminService) ListMessages(limit, offset int, conversationID *uuid.UUID) ([]repository.AdminMessage, int, error) {
	return s.adminRepo.ListMessages(limit, offset, conversationID)
}

// ListVerificationRequests returns verification requests
func (s *AdminService) ListVerificationRequests(status string) ([]repository.AdminVerificationRequest, error) {
	return s.adminRepo.ListVerificationRequests(status)
}

// ApproveVerification approves a verification request
func (s *AdminService) ApproveVerification(requestID uuid.UUID) error {
	return s.adminRepo.ApproveVerification(requestID)
}

// RejectVerification rejects a verification request
func (s *AdminService) RejectVerification(requestID uuid.UUID, reason string) error {
	return s.adminRepo.RejectVerification(requestID, reason)
}

// DeleteProfileImage deletes a profile image from S3 and database
func (s *AdminService) DeleteProfileImage(imageID uuid.UUID) error {
	// Get the image to get the S3 key
	img, err := s.adminRepo.GetProfileImage(imageID)
	if err != nil {
		return err
	}
	if img == nil {
		return nil
	}

	// Delete from S3 if client is available
	if s.s3Client != nil && img.S3Key != "" {
		s.s3Client.DeleteObject(img.S3Key)
	}

	// Delete from database
	return s.adminRepo.DeleteProfileImage(imageID)
}

// GetProfileImage returns a profile image by ID
func (s *AdminService) GetProfileImage(imageID uuid.UUID) (*models.ProfileImage, error) {
	return s.adminRepo.GetProfileImage(imageID)
}

// DeleteUser deletes a user and all related data
func (s *AdminService) DeleteUser(userID uuid.UUID) error {
	return s.adminRepo.DeleteUser(userID)
}

// UpdateUserWealthStatus updates a user's wealth status
func (s *AdminService) UpdateUserWealthStatus(userID uuid.UUID, status string) error {
	return s.adminRepo.UpdateUserWealthStatus(userID, status)
}

// GetStats returns admin dashboard stats
func (s *AdminService) GetStats() (map[string]interface{}, error) {
	return s.adminRepo.GetStats()
}

// ListAllImages returns all images from profile, messages, and verification
func (s *AdminService) ListAllImages(limit, offset int) ([]repository.AdminImage, int, error) {
	return s.adminRepo.ListAllImages(limit, offset)
}

