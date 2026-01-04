package services

import (
	"errors"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/internal/websocket"
)

type ChatService struct {
	messageRepo *repository.MessageRepository
	profileRepo *repository.ProfileRepository
	hub         *websocket.Hub
}

func NewChatService(messageRepo *repository.MessageRepository, profileRepo *repository.ProfileRepository, hub *websocket.Hub) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		profileRepo: profileRepo,
		hub:         hub,
	}
}

func (s *ChatService) CreateConversation(senderID uuid.UUID, req *models.CreateConversationRequest) (*models.ConversationWithDetails, error) {
	senderProfile, err := s.profileRepo.FindByUserID(senderID)
	if err != nil || senderProfile == nil {
		return nil, errors.New("sender profile not found")
	}

	recipientProfile, err := s.profileRepo.FindByUserID(req.RecipientID)
	if err != nil || recipientProfile == nil {
		return nil, errors.New("recipient profile not found")
	}

	existingConv, _ := s.messageRepo.FindConversationBetweenUsers(senderID, req.RecipientID)
	if existingConv != nil {
		return nil, errors.New("conversation already exists")
	}

	if senderProfile.Gender == models.GenderMale {
		return nil, errors.New("males cannot initiate conversations")
	}

	conv, err := s.messageRepo.CreateConversation(senderID, []uuid.UUID{senderID, req.RecipientID})
	if err != nil {
		return nil, err
	}

	msg, err := s.messageRepo.CreateMessage(conv.ID, senderID, req.Message)
	if err != nil {
		return nil, err
	}

	s.hub.BroadcastToUser(req.RecipientID, &models.WSMessage{
		Type:    models.WSTypeMessage,
		Payload: msg,
	})

	result := &models.ConversationWithDetails{
		Conversation: *conv,
		Participants: []uuid.UUID{senderID, req.RecipientID},
		LastMessage:  msg,
	}

	return result, nil
}

func (s *ChatService) GetConversations(userID uuid.UUID) ([]models.ConversationWithDetails, error) {
	conversations, err := s.messageRepo.GetUserConversations(userID)
	if err != nil {
		return nil, err
	}

	for i := range conversations {
		for _, participantID := range conversations[i].Participants {
			if participantID != userID {
				profile, _ := s.profileRepo.GetProfileWithDetails(participantID, userID)
				conversations[i].OtherUser = profile
				break
			}
		}
	}

	return conversations, nil
}

func (s *ChatService) SendMessage(conversationID, senderID uuid.UUID, content string) (*models.Message, error) {
	inConv, err := s.messageRepo.IsUserInConversation(conversationID, senderID)
	if err != nil || !inConv {
		return nil, errors.New("not authorized to send message in this conversation")
	}

	msg, err := s.messageRepo.CreateMessage(conversationID, senderID, content)
	if err != nil {
		return nil, err
	}

	participants, _ := s.messageRepo.GetConversationParticipants(conversationID)
	for _, participantID := range participants {
		if participantID != senderID {
			s.hub.BroadcastToUser(participantID, &models.WSMessage{
				Type:    models.WSTypeMessage,
				Payload: msg,
			})
		}
	}

	return msg, nil
}

func (s *ChatService) GetMessages(conversationID, userID uuid.UUID, limit, offset int) ([]models.Message, error) {
	inConv, err := s.messageRepo.IsUserInConversation(conversationID, userID)
	if err != nil || !inConv {
		return nil, errors.New("not authorized to view this conversation")
	}

	s.messageRepo.MarkMessagesAsRead(conversationID, userID)

	if limit < 1 || limit > 100 {
		limit = 50
	}

	return s.messageRepo.GetMessages(conversationID, limit, offset)
}

func (s *ChatService) BroadcastTyping(conversationID, userID uuid.UUID, isTyping bool) error {
	participants, err := s.messageRepo.GetConversationParticipants(conversationID)
	if err != nil {
		return err
	}

	msgType := models.WSTypeTyping
	if !isTyping {
		msgType = models.WSTypeStopTyping
	}

	for _, participantID := range participants {
		if participantID != userID {
			s.hub.BroadcastToUser(participantID, &models.WSMessage{
				Type: msgType,
				Payload: models.WSTypingPayload{
					ConversationID: conversationID,
					UserID:         userID,
				},
			})
		}
	}

	return nil
}
