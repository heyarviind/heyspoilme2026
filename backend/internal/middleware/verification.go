package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/repository"
)

type VerificationMiddleware struct {
	userRepo *repository.UserRepository
}

func NewVerificationMiddleware(userRepo *repository.UserRepository) *VerificationMiddleware {
	return &VerificationMiddleware{userRepo: userRepo}
}

// RequireEmailVerified returns a middleware that checks if the user's email is verified.
// This should be used AFTER RequireAuth() middleware.
func (m *VerificationMiddleware) RequireEmailVerified() gin.HandlerFunc {
	return func(c *gin.Context) {
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


