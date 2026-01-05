package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/services"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) CreateConversation(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversation, err := h.chatService.CreateConversation(userID, &req)
	if err != nil {
		if errors.Is(err, services.ErrVerificationRequired) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "person_verification_required",
				"message":         "Complete identity verification to send messages",
				"person_verified": false,
			})
			return
		}
		if errors.Is(err, services.ErrMaleCannotInitiate) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "male_cannot_initiate",
				"message": "Only women can start conversations",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conversation)
}

func (h *ChatHandler) GetConversations(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	conversations, err := h.chatService.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	conversationIDStr := c.Param("id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation ID"})
		return
	}

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.chatService.SendMessage(conversationID, userID, req.Content, req.ImageURL)
	if err != nil {
		if errors.Is(err, services.ErrVerificationRequired) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "person_verification_required",
				"message":         "Complete identity verification to send messages",
				"person_verified": false,
			})
			return
		}
		if errors.Is(err, services.ErrWealthStatusRequired) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "subscription_required",
				"message":       "Upgrade to send messages",
				"wealth_status": "none",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	conversationIDStr := c.Param("id")
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation ID"})
		return
	}

	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	messages, err := h.chatService.GetMessages(conversationID, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *ChatHandler) GetUnreadCount(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	count, err := h.chatService.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetInbox returns the user's inbox with locked message support for males
func (h *ChatHandler) GetInbox(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	inbox, err := h.chatService.GetInbox(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inbox)
}
