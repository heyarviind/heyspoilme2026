package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/services"
)

type LikeHandler struct {
	likeService    *services.LikeService
	profileService *services.ProfileService
}

func NewLikeHandler(likeService *services.LikeService, profileService *services.ProfileService) *LikeHandler {
	return &LikeHandler{
		likeService:    likeService,
		profileService: profileService,
	}
}

func (h *LikeHandler) LikeProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	likedIDStr := c.Param("id")
	likedID, err := uuid.Parse(likedIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile ID"})
		return
	}

	if userID == likedID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot like your own profile"})
		return
	}

	profile, _ := h.profileService.GetProfile(userID)
	images, _ := h.profileService.GetProfileImages(userID)

	var likerName string
	var likerImage string
	if profile != nil {
		likerName = profile.Bio
	}
	if len(images) > 0 {
		likerImage = images[0].URL
	}

	like, err := h.likeService.LikeProfile(userID, likedID, likerName, likerImage)
	if err != nil {
		if errors.Is(err, services.ErrNotVerified) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":           "person_verification_required",
				"message":         "You must complete identity verification to like profiles",
				"person_verified": false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if like == nil {
		c.JSON(http.StatusOK, gin.H{"message": "already liked"})
		return
	}

	c.JSON(http.StatusCreated, like)
}

func (h *LikeHandler) UnlikeProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	likedIDStr := c.Param("id")
	likedID, err := uuid.Parse(likedIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile ID"})
		return
	}

	if err := h.likeService.UnlikeProfile(userID, likedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

func (h *LikeHandler) GetReceivedLikes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	limit := 20
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

	likes, total, err := h.likeService.GetReceivedLikes(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type likeWithProfile struct {
		Like    interface{} `json:"like"`
		Profile interface{} `json:"profile"`
	}

	var results []likeWithProfile
	for _, like := range likes {
		profile, _ := h.profileService.GetProfileWithDetails(like.LikerID, userID)
		results = append(results, likeWithProfile{
			Like:    like,
			Profile: profile,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"likes": results,
		"total": total,
	})
}

func (h *LikeHandler) GetGivenLikes(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	limit := 20
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

	likes, total, err := h.likeService.GetGivenLikes(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type likeWithProfile struct {
		Like    interface{} `json:"like"`
		Profile interface{} `json:"profile"`
	}

	var results []likeWithProfile
	for _, like := range likes {
		profile, _ := h.profileService.GetProfileWithDetails(like.LikedID, userID)
		results = append(results, likeWithProfile{
			Like:    like,
			Profile: profile,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"likes": results,
		"total": total,
	})
}
