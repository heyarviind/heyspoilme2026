package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/services"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.profileService.CreateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	profile, err := h.profileService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	images, _ := h.profileService.GetProfileImages(userID)

	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
		"images":  images,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.profileService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	profileUserIDStr := c.Param("id")
	profileUserID, err := uuid.Parse(profileUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile ID"})
		return
	}

	profile, err := h.profileService.GetProfileWithDetails(profileUserID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) ListProfiles(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	query := &models.ListProfilesQuery{
		Gender:      c.Query("gender"),
		City:        c.Query("city"),
		State:       c.Query("state"),
		Page:        1,
		Limit:       20,
		MaxDistance: 0,
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			query.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			query.Limit = l
		}
	}
	if minAge := c.Query("min_age"); minAge != "" {
		if a, err := strconv.Atoi(minAge); err == nil {
			query.MinAge = a
		}
	}
	if maxAge := c.Query("max_age"); maxAge != "" {
		if a, err := strconv.Atoi(maxAge); err == nil {
			query.MaxAge = a
		}
	}
	if maxDistance := c.Query("max_distance"); maxDistance != "" {
		if d, err := strconv.ParseFloat(maxDistance, 64); err == nil {
			query.MaxDistance = d
		}
	}
	if onlineOnly := c.Query("online_only"); onlineOnly == "true" || onlineOnly == "1" {
		query.OnlineOnly = true
	}

	profiles, total, err := h.profileService.ListProfiles(userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"total":    total,
		"page":     query.Page,
		"limit":    query.Limit,
	})
}

func (h *ProfileHandler) AddProfileImage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req struct {
		S3Key     string `json:"s3_key" binding:"required"`
		URL       string `json:"url" binding:"required"`
		IsPrimary bool   `json:"is_primary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := h.profileService.AddProfileImage(userID, req.S3Key, req.URL, req.IsPrimary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, image)
}

func (h *ProfileHandler) DeleteProfileImage(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	imageIDStr := c.Param("imageId")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image ID"})
		return
	}

	if err := h.profileService.DeleteProfileImage(imageID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image deleted"})
}
