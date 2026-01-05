package services

import (
	"log"
	"strings"

	"github.com/google/uuid"

	"heyspoilme/internal/repository"
	"heyspoilme/pkg/storage"
)

type AccountService struct {
	userRepo         *repository.UserRepository
	profileRepo      *repository.ProfileRepository
	messageRepo      *repository.MessageRepository
	likeRepo         *repository.LikeRepository
	notificationRepo *repository.NotificationRepository
	presenceRepo     *repository.PresenceRepository
	s3Client         *storage.S3Client
}

func NewAccountService(
	userRepo *repository.UserRepository,
	profileRepo *repository.ProfileRepository,
	messageRepo *repository.MessageRepository,
	likeRepo *repository.LikeRepository,
	notificationRepo *repository.NotificationRepository,
	presenceRepo *repository.PresenceRepository,
	s3Client *storage.S3Client,
) *AccountService {
	return &AccountService{
		userRepo:         userRepo,
		profileRepo:      profileRepo,
		messageRepo:      messageRepo,
		likeRepo:         likeRepo,
		notificationRepo: notificationRepo,
		presenceRepo:     presenceRepo,
		s3Client:         s3Client,
	}
}

// DeleteAccount permanently deletes a user account and all associated data
func (s *AccountService) DeleteAccount(userID uuid.UUID) error {
	log.Printf("[Account] Starting account deletion for user: %s", userID)

	// 1. Get all image S3 keys before deleting database records
	var s3KeysToDelete []string

	// Profile images
	profileImageKeys, err := s.profileRepo.GetImageS3Keys(userID)
	if err != nil {
		log.Printf("[Account] Warning: failed to get profile image keys: %v", err)
	} else {
		s3KeysToDelete = append(s3KeysToDelete, profileImageKeys...)
	}

	// Message images (need to extract S3 key from URL)
	messageImageURLs, err := s.messageRepo.GetUserMessageImageURLs(userID)
	if err != nil {
		log.Printf("[Account] Warning: failed to get message image URLs: %v", err)
	} else {
		for _, url := range messageImageURLs {
			key := extractS3KeyFromURL(url)
			if key != "" {
				s3KeysToDelete = append(s3KeysToDelete, key)
			}
		}
	}

	log.Printf("[Account] Found %d S3 objects to delete", len(s3KeysToDelete))

	// 2. Delete database records in order (respecting foreign key constraints)

	// Delete notifications
	if err := s.notificationRepo.DeleteAllForUser(userID); err != nil {
		log.Printf("[Account] Warning: failed to delete notifications: %v", err)
	}

	// Delete likes (both given and received)
	if err := s.likeRepo.DeleteAllForUser(userID); err != nil {
		log.Printf("[Account] Warning: failed to delete likes: %v", err)
	}

	// Delete conversations, messages, and participants
	if err := s.messageRepo.DeleteUserConversations(userID); err != nil {
		log.Printf("[Account] Warning: failed to delete conversations: %v", err)
	}

	// Delete presence
	if err := s.presenceRepo.Delete(userID); err != nil {
		log.Printf("[Account] Warning: failed to delete presence: %v", err)
	}

	// Delete profile images from database
	if err := s.profileRepo.DeleteAllImages(userID); err != nil {
		log.Printf("[Account] Warning: failed to delete profile images: %v", err)
	}

	// Delete profile
	if err := s.profileRepo.Delete(userID); err != nil {
		log.Printf("[Account] Warning: failed to delete profile: %v", err)
	}

	// Delete user
	if err := s.userRepo.Delete(userID); err != nil {
		log.Printf("[Account] Error: failed to delete user: %v", err)
		return err
	}

	// 3. Delete S3 objects asynchronously
	if s.s3Client != nil {
		go func() {
			for _, key := range s3KeysToDelete {
				if err := s.s3Client.DeleteObject(key); err != nil {
					log.Printf("[Account] Warning: failed to delete S3 object %s: %v", key, err)
				} else {
					log.Printf("[Account] Deleted S3 object: %s", key)
				}
			}
			log.Printf("[Account] Finished deleting %d S3 objects for user: %s", len(s3KeysToDelete), userID)
		}()
	}

	log.Printf("[Account] Account deletion completed for user: %s", userID)
	return nil
}

// extractS3KeyFromURL extracts the S3 key from a full URL
// Example: https://cdn.heyspoilme.com/chat/xxx/yyy/file.webp -> chat/xxx/yyy/file.webp
func extractS3KeyFromURL(url string) string {
	// Find common prefixes
	prefixes := []string{
		"/profiles/",
		"/chat/",
	}

	for _, prefix := range prefixes {
		if idx := strings.Index(url, prefix); idx != -1 {
			return url[idx+1:] // Skip the leading slash
		}
	}

	return ""
}


