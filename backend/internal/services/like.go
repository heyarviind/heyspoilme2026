package services

import (
	"errors"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/internal/websocket"
)

var (
	ErrNotVerified = errors.New("person verification required to like profiles")
)

type LikeService struct {
	likeRepo         *repository.LikeRepository
	notificationRepo *repository.NotificationRepository
	profileRepo      *repository.ProfileRepository
	hub              *websocket.Hub
}

func NewLikeService(likeRepo *repository.LikeRepository, notificationRepo *repository.NotificationRepository, profileRepo *repository.ProfileRepository, hub *websocket.Hub) *LikeService {
	return &LikeService{
		likeRepo:         likeRepo,
		notificationRepo: notificationRepo,
		profileRepo:      profileRepo,
		hub:              hub,
	}
}

func (s *LikeService) LikeProfile(likerID, likedID uuid.UUID, likerName, likerImage string) (*models.Like, error) {
	// Check if liker is person_verified (is_verified in profiles table)
	likerProfile, err := s.profileRepo.FindByUserID(likerID)
	if err != nil || likerProfile == nil {
		return nil, errors.New("profile not found")
	}
	if !likerProfile.IsVerified {
		return nil, ErrNotVerified
	}

	exists, _ := s.likeRepo.Exists(likerID, likedID)
	if exists {
		return nil, nil
	}

	like, err := s.likeRepo.Create(likerID, likedID)
	if err != nil {
		return nil, err
	}

	notifData := &models.NotificationData{
		FromUserID:    likerID,
		FromUserName:  likerName,
		FromUserImage: likerImage,
	}
	notification, err := s.notificationRepo.Create(likedID, models.NotificationTypeLike, notifData)
	if err == nil && notification != nil {
		s.hub.BroadcastToUser(likedID, &models.WSMessage{
			Type:    models.WSTypeNotification,
			Payload: notification,
		})
	}

	return like, nil
}

func (s *LikeService) UnlikeProfile(likerID, likedID uuid.UUID) error {
	return s.likeRepo.Delete(likerID, likedID)
}

func (s *LikeService) GetReceivedLikes(userID uuid.UUID, limit, offset int) ([]models.Like, int, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}
	return s.likeRepo.GetReceivedLikes(userID, limit, offset)
}

func (s *LikeService) GetGivenLikes(userID uuid.UUID, limit, offset int) ([]models.Like, int, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}
	return s.likeRepo.GetGivenLikes(userID, limit, offset)
}

func (s *LikeService) IsLiked(likerID, likedID uuid.UUID) (bool, error) {
	return s.likeRepo.Exists(likerID, likedID)
}
