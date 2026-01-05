package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/services"
	"heyspoilme/pkg/storage"
)

type AdminHandler struct {
	adminService       *services.AdminService
	featureFlagService *services.FeatureFlagService
	s3Client           *storage.S3Client
	adminCode1         string
	adminCode2         string
}

func NewAdminHandler(adminService *services.AdminService, featureFlagService *services.FeatureFlagService, s3Client *storage.S3Client, adminCode1, adminCode2 string) *AdminHandler {
	return &AdminHandler{
		adminService:       adminService,
		featureFlagService: featureFlagService,
		s3Client:           s3Client,
		adminCode1:         adminCode1,
		adminCode2:         adminCode2,
	}
}

// ValidateAdminAccess middleware to validate admin path codes
func (h *AdminHandler) ValidateAdminAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		code1 := c.Param("code1")
		code2 := c.Param("code2")

		if code1 != h.adminCode1 || code2 != h.adminCode2 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetStats returns admin dashboard statistics
func (h *AdminHandler) GetStats(c *gin.Context) {
	stats, err := h.adminService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ListUsers returns paginated list of users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page := 1
	limit := 20
	search := c.Query("search")

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	offset := (page - 1) * limit
	users, total, err := h.adminService.ListUsers(limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetUser returns a single user with their images
func (h *AdminHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.adminService.GetUserWithImages(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListMessages returns paginated list of messages
func (h *AdminHandler) ListMessages(c *gin.Context) {
	page := 1
	limit := 50

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	var conversationID *uuid.UUID
	if convIDStr := c.Query("conversation_id"); convIDStr != "" {
		if id, err := uuid.Parse(convIDStr); err == nil {
			conversationID = &id
		}
	}

	offset := (page - 1) * limit
	messages, total, err := h.adminService.ListMessages(limit, offset, conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// ListVerificationRequests returns verification requests
func (h *AdminHandler) ListVerificationRequests(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		status = "pending"
	}

	requests, err := h.adminService.ListVerificationRequests(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
	})
}

// ApproveVerification approves a verification request
func (h *AdminHandler) ApproveVerification(c *gin.Context) {
	requestIDStr := c.Param("requestId")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	if err := h.adminService.ApproveVerification(requestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification approved"})
}

// RejectVerification rejects a verification request
func (h *AdminHandler) RejectVerification(c *gin.Context) {
	requestIDStr := c.Param("requestId")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason is required"})
		return
	}

	if err := h.adminService.RejectVerification(requestID, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification rejected"})
}

// DeleteProfileImage deletes a profile image
func (h *AdminHandler) DeleteProfileImage(c *gin.Context) {
	imageIDStr := c.Param("imageId")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	if err := h.adminService.DeleteProfileImage(imageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image deleted"})
}

// DeleteUser deletes a user and all related data
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.adminService.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// UpdateUserWealthStatus updates a user's wealth status
func (h *AdminHandler) UpdateUserWealthStatus(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	// Validate status
	validStatuses := map[string]bool{"none": true, "low": true, "medium": true, "high": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status, must be one of: none, low, medium, high"})
		return
	}

	if err := h.adminService.UpdateUserWealthStatus(userID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wealth status updated"})
}

// ListAllImages returns all images from profile, messages, and verification
func (h *AdminHandler) ListAllImages(c *gin.Context) {
	page := 1
	limit := 50

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	offset := (page - 1) * limit
	images, total, err := h.adminService.ListAllImages(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"images": images,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// GetFeatureFlags returns all feature flags
func (h *AdminHandler) GetFeatureFlags(c *gin.Context) {
	flags := h.featureFlagService.GetAllWithDefaults()
	c.JSON(http.StatusOK, gin.H{"flags": flags})
}

// UpdateFeatureFlag updates a feature flag
func (h *AdminHandler) UpdateFeatureFlag(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.featureFlagService.SetFlag(key, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "feature flag updated",
		"key":     key,
		"enabled": req.Enabled,
	})
}

// GetPublicFeatureFlags returns public feature flag status (no auth required)
func (h *AdminHandler) GetPublicFeatureFlags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"restrictions_enabled": h.featureFlagService.RestrictionsEnabled(),
	})
}

// GetUserUploadURL returns a presigned URL for uploading an image to a user's profile
func (h *AdminHandler) GetUserUploadURL(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		ContentType string `json:"content_type" binding:"required"`
		FileExt     string `json:"file_ext" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !validTypes[req.ContentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type"})
		return
	}

	if h.s3Client == nil {
		log.Println("[AdminUpload] ERROR: S3 client is nil - storage not configured")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "image upload not configured"})
		return
	}

	fileKey := fmt.Sprintf("profiles/%s/%s%s", userID.String(), uuid.New().String(), req.FileExt)
	log.Printf("[AdminUpload] Generating presigned URL for user %s, key: %s", userID, fileKey)

	uploadURL, err := h.s3Client.GetPresignedUploadURL(fileKey, req.ContentType)
	if err != nil {
		log.Printf("[AdminUpload] ERROR generating presigned URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate presigned URL"})
		return
	}

	publicURL := h.s3Client.GetPublicURL(fileKey)

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"s3_key":     fileKey,
		"public_url": publicURL,
	})
}

// AddUserImage adds a profile image for a user (after uploading to S3)
func (h *AdminHandler) AddUserImage(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		S3Key     string `json:"s3_key" binding:"required"`
		URL       string `json:"url" binding:"required"`
		IsPrimary bool   `json:"is_primary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := h.adminService.AddUserProfileImage(userID, req.S3Key, req.URL, req.IsPrimary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, image)
}

// UpdateUserPresence updates a user's online status
func (h *AdminHandler) UpdateUserPresence(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		IsOnline bool `json:"is_online"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.UpdateUserPresence(userID, req.IsOnline); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "presence updated", "is_online": req.IsOnline})
}

// UpdateUserProfile updates a user's profile fields
func (h *AdminHandler) UpdateUserProfile(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		DisplayName *string  `json:"display_name,omitempty"`
		Age         *int     `json:"age,omitempty"`
		Bio         *string  `json:"bio,omitempty"`
		City        *string  `json:"city,omitempty"`
		State       *string  `json:"state,omitempty"`
		Latitude    *float64 `json:"latitude,omitempty"`
		Longitude   *float64 `json:"longitude,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.adminService.UpdateUserProfile(userID, req.DisplayName, req.Age, req.Bio, req.City, req.State, req.Latitude, req.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

// SendMessageAsUser sends a message from one user to another
func (h *AdminHandler) SendMessageAsUser(c *gin.Context) {
	var req struct {
		SenderID    string `json:"sender_id" binding:"required"`
		RecipientID string `json:"recipient_id" binding:"required"`
		Content     string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderID, err := uuid.Parse(req.SenderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sender ID"})
		return
	}

	recipientID, err := uuid.Parse(req.RecipientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipient ID"})
		return
	}

	if err := h.adminService.SendMessageAsUser(senderID, recipientID, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message sent"})
}

