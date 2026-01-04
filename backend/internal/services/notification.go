package services

import (
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

func (s *NotificationService) GetNotifications(userID uuid.UUID, limit, offset int) ([]models.Notification, int, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}
	return s.notificationRepo.GetByUserID(userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(notifID, userID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(notifID, userID)
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int, error) {
	return s.notificationRepo.GetUnreadCount(userID)
}

func (s *NotificationService) CreateNotification(userID uuid.UUID, notifType models.NotificationType, data *models.NotificationData) (*models.Notification, error) {
	return s.notificationRepo.Create(userID, notifType, data)
}
