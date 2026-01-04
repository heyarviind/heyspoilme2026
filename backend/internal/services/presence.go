package services

import (
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/internal/websocket"
)

type PresenceService struct {
	presenceRepo *repository.PresenceRepository
	hub          *websocket.Hub
}

func NewPresenceService(presenceRepo *repository.PresenceRepository, hub *websocket.Hub) *PresenceService {
	return &PresenceService{
		presenceRepo: presenceRepo,
		hub:          hub,
	}
}

func (s *PresenceService) SetOnline(userID uuid.UUID) error {
	err := s.presenceRepo.SetOnline(userID)
	if err != nil {
		return err
	}

	s.hub.BroadcastToAll(&models.WSMessage{
		Type: models.WSTypePresence,
		Payload: models.WSPresencePayload{
			UserID:   userID,
			IsOnline: true,
		},
	})

	return nil
}

func (s *PresenceService) SetOffline(userID uuid.UUID) error {
	err := s.presenceRepo.SetOffline(userID)
	if err != nil {
		return err
	}

	s.hub.BroadcastToAll(&models.WSMessage{
		Type: models.WSTypePresence,
		Payload: models.WSPresencePayload{
			UserID:   userID,
			IsOnline: false,
		},
	})

	return nil
}

func (s *PresenceService) GetPresence(userID uuid.UUID) (*models.UserPresence, error) {
	return s.presenceRepo.GetPresence(userID)
}

func (s *PresenceService) GetOnlineUsers(userIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	return s.presenceRepo.GetOnlineUsers(userIDs)
}
