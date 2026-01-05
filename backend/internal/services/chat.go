package services

import (
	"errors"

	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/internal/websocket"
)

var (
	ErrMaleCannotInitiate     = errors.New("males cannot initiate conversations")
	ErrVerificationRequired   = errors.New("identity verification required to send messages")
	ErrWealthStatusRequired   = errors.New("subscription required to send messages")
	ErrMessageNotVisible      = errors.New("message not visible to recipient")
)

type ChatService struct {
	messageRepo *repository.MessageRepository
	profileRepo *repository.ProfileRepository
	userRepo    *repository.UserRepository
	hub         *websocket.Hub
}

func NewChatService(messageRepo *repository.MessageRepository, profileRepo *repository.ProfileRepository, userRepo *repository.UserRepository, hub *websocket.Hub) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		profileRepo: profileRepo,
		userRepo:    userRepo,
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

	// Sender must be person_verified to message
	if !senderProfile.IsVerified {
		return nil, ErrVerificationRequired
	}

	// Get sender's user info for wealth status check
	senderUser, err := s.userRepo.FindByID(senderID)
	if err != nil || senderUser == nil {
		return nil, errors.New("sender not found")
	}

	// Males cannot initiate conversations
	if senderProfile.Gender == models.GenderMale {
		return nil, ErrMaleCannotInitiate
	}

	// Female sending to male: check recipient's wealth status
	// Message is stored but we track visibility
	recipientUser, err := s.userRepo.FindByID(req.RecipientID)
	if err != nil || recipientUser == nil {
		return nil, errors.New("recipient not found")
	}

	conv, err := s.messageRepo.CreateConversation(senderID, []uuid.UUID{senderID, req.RecipientID})
	if err != nil {
		return nil, err
	}

	msg, err := s.messageRepo.CreateMessage(conv.ID, senderID, req.Message, nil)
	if err != nil {
		return nil, err
	}

	// Only send real-time notification if recipient can view messages
	// (male with wealth_status != 'none', or female)
	recipientCanView := recipientProfile.Gender == models.GenderFemale || recipientUser.WealthStatus.CanViewMessages()
	if recipientCanView {
		s.hub.BroadcastToUser(req.RecipientID, &models.WSMessage{
			Type:    models.WSTypeMessage,
			Payload: msg,
		})
	}

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

// GetInbox returns the inbox with locked message support for males with wealth_status=none
func (s *ChatService) GetInbox(userID uuid.UUID) (*models.InboxResponse, error) {
	userProfile, err := s.profileRepo.FindByUserID(userID)
	if err != nil || userProfile == nil {
		return nil, errors.New("profile not found")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	conversations, err := s.messageRepo.GetUserConversations(userID)
	if err != nil {
		return nil, err
	}

	// Females can view all messages
	// Males need wealth_status != 'none' to view all messages
	canViewAll := userProfile.Gender == models.GenderFemale || user.WealthStatus.CanViewMessages()

	response := &models.InboxResponse{
		CanViewAllMessages: canViewAll,
	}

	if canViewAll {
		// User can view all conversations
		for i := range conversations {
			for _, participantID := range conversations[i].Participants {
				if participantID != userID {
					profile, _ := s.profileRepo.GetProfileWithDetails(participantID, userID)
					conversations[i].OtherUser = profile
					break
				}
			}
		}
		response.Conversations = conversations
		response.LockedCount = 0
	} else {
		// Male with wealth_status=none - show locked previews
		response.Conversations = []models.ConversationWithDetails{}
		response.LockedCount = len(conversations)

		// Create blurred previews (up to 5)
		maxPreviews := 5
		if len(conversations) < maxPreviews {
			maxPreviews = len(conversations)
		}

		for i := 0; i < maxPreviews; i++ {
			conv := conversations[i]
			var otherUserID uuid.UUID
			for _, participantID := range conv.Participants {
				if participantID != userID {
					otherUserID = participantID
					break
				}
			}

			profile, _ := s.profileRepo.FindByUserID(otherUserID)
			images, _ := s.profileRepo.GetImages(otherUserID)

			lockedConv := models.LockedConversation{
				ID:        conv.ID,
				CreatedAt: conv.CreatedAt,
			}

			if profile != nil {
				lockedConv.SenderAge = profile.Age
				lockedConv.SenderCity = profile.City
			}

			if len(images) > 0 {
				lockedConv.SenderImage = images[0].URL
			}

			// Create blurred preview (first few words + "...")
			if conv.LastMessage != nil && conv.LastMessage.Content != "" {
				preview := conv.LastMessage.Content
				if len(preview) > 30 {
					preview = preview[:30]
				}
				// Blur by showing only first 10 chars + "..."
				if len(preview) > 10 {
					lockedConv.BlurredPreview = preview[:10] + "..."
				} else {
					lockedConv.BlurredPreview = preview + "..."
				}
			} else {
				lockedConv.BlurredPreview = "New message..."
			}

			response.LockedPreviews = append(response.LockedPreviews, lockedConv)
		}
	}

	return response, nil
}

func (s *ChatService) SendMessage(conversationID, senderID uuid.UUID, content string, imageURL *string) (*models.Message, error) {
	inConv, err := s.messageRepo.IsUserInConversation(conversationID, senderID)
	if err != nil || !inConv {
		return nil, errors.New("not authorized to send message in this conversation")
	}

	// Validate: message must have either content or image
	if content == "" && (imageURL == nil || *imageURL == "") {
		return nil, errors.New("message must contain text or an image")
	}

	// Get sender's profile and user for verification
	senderProfile, err := s.profileRepo.FindByUserID(senderID)
	if err != nil || senderProfile == nil {
		return nil, errors.New("sender profile not found")
	}

	senderUser, err := s.userRepo.FindByID(senderID)
	if err != nil || senderUser == nil {
		return nil, errors.New("sender not found")
	}

	// Sender must be person_verified
	if !senderProfile.IsVerified {
		return nil, ErrVerificationRequired
	}

	// Males must have wealth_status != 'none' to send messages
	if senderProfile.Gender == models.GenderMale && !senderUser.WealthStatus.CanMessage() {
		return nil, ErrWealthStatusRequired
	}

	msg, err := s.messageRepo.CreateMessage(conversationID, senderID, content, imageURL)
	if err != nil {
		return nil, err
	}

	participants, _ := s.messageRepo.GetConversationParticipants(conversationID)
	for _, participantID := range participants {
		if participantID != senderID {
			// Get recipient info to check if they can view this message
			recipientProfile, _ := s.profileRepo.FindByUserID(participantID)
			recipientUser, _ := s.userRepo.FindByID(participantID)

			// Only broadcast if recipient can view
			// Females can always view, males need wealth_status
			recipientCanView := recipientProfile != nil && (
				recipientProfile.Gender == models.GenderFemale ||
				(recipientUser != nil && recipientUser.WealthStatus.CanViewMessages()))

			if recipientCanView {
				s.hub.BroadcastToUser(participantID, &models.WSMessage{
					Type:    models.WSTypeMessage,
					Payload: msg,
				})
			}
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

func (s *ChatService) GetUnreadCount(userID uuid.UUID) (int, error) {
	return s.messageRepo.GetUnreadMessageCount(userID)
}
