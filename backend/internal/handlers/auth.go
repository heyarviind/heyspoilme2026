package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"heyspoilme/internal/models"
	"heyspoilme/internal/services"
	"heyspoilme/pkg/auth"
)

type AuthHandler struct {
	authService    *services.AuthService
	profileService *services.ProfileService
	accountService *services.AccountService
	googleAuth     *auth.GoogleAuth
	frontendURL    string
}

func NewAuthHandler(authService *services.AuthService, profileService *services.ProfileService, accountService *services.AccountService, googleAuth *auth.GoogleAuth, frontendURL string) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		profileService: profileService,
		accountService: accountService,
		googleAuth:     googleAuth,
		frontendURL:    frontendURL,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Signup(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token":  token,
		"user":   user,
		"is_new": true,
	})
}

func (h *AuthHandler) Signin(c *gin.Context) {
	var req models.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Signin(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	profile, _ := h.profileService.GetProfile(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user":    user,
		"profile": profile,
		"is_new":  false,
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	url := h.googleAuth.GetAuthURL()
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?error=missing_code")
		return
	}

	token, err := h.googleAuth.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?error=exchange_failed")
		return
	}

	userInfo, err := h.googleAuth.GetUserInfo(context.Background(), token)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?error=userinfo_failed")
		return
	}

	user, isNew, err := h.authService.FindOrCreateUser(userInfo.ID, userInfo.Email)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?error=user_creation_failed")
		return
	}

	jwtToken, err := h.authService.GenerateToken(user)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"/auth/error?error=token_generation_failed")
		return
	}

	redirectURL := h.frontendURL + "/auth/callback?token=" + jwtToken
	if isNew {
		redirectURL += "&is_new=true"
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	user, err := h.authService.GetUserByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	profile, _ := h.profileService.GetProfile(userID)

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"profile": profile,
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification token required"})
		return
	}

	if err := h.authService.VerifyEmail(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.authService.ResendVerificationEmail(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification email sent"})
}

func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	if err := h.accountService.DeleteAccount(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
