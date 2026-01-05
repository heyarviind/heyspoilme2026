package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/pkg/storage"
)

type UploadHandler struct {
	s3Client *storage.S3Client
}

func NewUploadHandler(s3Client *storage.S3Client) *UploadHandler {
	return &UploadHandler{
		s3Client: s3Client,
	}
}

func (h *UploadHandler) GetPresignedURL(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

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
		log.Println("[Upload] ERROR: S3 client is nil - storage not configured")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "image upload not configured"})
		return
	}

	fileKey := fmt.Sprintf("profiles/%s/%s%s", userID.String(), uuid.New().String(), req.FileExt)
	log.Printf("[Upload] Generating presigned URL for key: %s, contentType: %s", fileKey, req.ContentType)

	uploadURL, err := h.s3Client.GetPresignedUploadURL(fileKey, req.ContentType)
	if err != nil {
		log.Printf("[Upload] ERROR generating presigned URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate presigned URL"})
		return
	}
	log.Printf("[Upload] Successfully generated presigned URL for key: %s", fileKey)

	publicURL := h.s3Client.GetPublicURL(fileKey)

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"s3_key":     fileKey,
		"public_url": publicURL,
	})
}

func (h *UploadHandler) GetChatImagePresignedURL(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		ContentType    string `json:"content_type" binding:"required"`
		FileExt        string `json:"file_ext" binding:"required"`
		ConversationID string `json:"conversation_id" binding:"required"`
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
		log.Println("[Upload] ERROR: S3 client is nil - storage not configured")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "image upload not configured"})
		return
	}

	// Store chat images in a separate folder with conversation context
	fileKey := fmt.Sprintf("chat/%s/%s/%s%s", req.ConversationID, userID.String(), uuid.New().String(), req.FileExt)
	log.Printf("[Upload] Generating chat image presigned URL for key: %s, contentType: %s", fileKey, req.ContentType)

	uploadURL, err := h.s3Client.GetPresignedUploadURL(fileKey, req.ContentType)
	if err != nil {
		log.Printf("[Upload] ERROR generating presigned URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate presigned URL"})
		return
	}
	log.Printf("[Upload] Successfully generated chat image presigned URL for key: %s", fileKey)

	publicURL := h.s3Client.GetPublicURL(fileKey)

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"s3_key":     fileKey,
		"public_url": publicURL,
	})
}
