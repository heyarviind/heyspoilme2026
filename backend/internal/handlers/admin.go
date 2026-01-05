package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/services"
)

type AdminHandler struct {
	adminService *services.AdminService
	adminCode1   string
	adminCode2   string
}

func NewAdminHandler(adminService *services.AdminService, adminCode1, adminCode2 string) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		adminCode1:   adminCode1,
		adminCode2:   adminCode2,
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

