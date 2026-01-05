package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/services"
)

type VerificationHandler struct {
	verificationService *services.VerificationService
}

func NewVerificationHandler(verificationService *services.VerificationService) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

// GenerateCode generates a verification code for the user to speak in their video
func (h *VerificationHandler) GenerateCode(c *gin.Context) {
	code, err := h.verificationService.GenerateVerificationCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate code"})
		return
	}

	c.JSON(http.StatusOK, models.VerificationCodeResponse{
		Code: code,
	})
}

// SubmitVerification handles the verification request submission
func (h *VerificationHandler) SubmitVerification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the code from the request (it should be stored in session/frontend)
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification code is required"})
		return
	}

	verification, err := h.verificationService.SubmitVerification(userID.(uuid.UUID), &req, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Verification request submitted successfully. We will review it within 24-48 hours.",
		"verification": verification,
	})
}

// GetStatus returns the current verification status
func (h *VerificationHandler) GetStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status, err := h.verificationService.GetVerificationStatus(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}


