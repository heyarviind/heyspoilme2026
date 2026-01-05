package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/internal/services"
)

type VerificationMiddleware struct {
	userRepo           *repository.UserRepository
	featureFlagService *services.FeatureFlagService
}

func NewVerificationMiddleware(userRepo *repository.UserRepository, featureFlagService *services.FeatureFlagService) *VerificationMiddleware {
	return &VerificationMiddleware{
		userRepo:           userRepo,
		featureFlagService: featureFlagService,
	}
}

// RequireEmailVerified returns a middleware that checks if the user's email is verified.
// This should be used AFTER RequireAuth() middleware.
// If restrictions are disabled via feature flag, this check is skipped.
func (m *VerificationMiddleware) RequireEmailVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip verification check if restrictions are disabled
		if !m.featureFlagService.IsEnabled(models.FlagRestrictionsEnabled) {
			c.Next()
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		user, err := m.userRepo.FindByID(userID.(uuid.UUID))
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		if !user.EmailVerified {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "email_not_verified",
				"message":        "Please verify your email to access this feature",
				"email_verified": false,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}



