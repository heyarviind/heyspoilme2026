package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"heyspoilme/internal/models"
	"heyspoilme/internal/repository"
	"heyspoilme/pkg/email"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	jwtSecret   string
	emailClient *email.ZeptoMailClient
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, emailClient *email.ZeptoMailClient) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		emailClient: emailClient,
	}
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

func (s *AuthService) FindOrCreateUser(googleID, email string) (*models.User, bool, error) {
	return s.userRepo.FindOrCreate(googleID, email)
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *AuthService) Signup(userEmail, password string) (*models.User, error) {
	existing, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		return nil, errors.New("failed to generate verification token")
	}
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	user, err := s.userRepo.CreateWithPassword(userEmail, string(hashedPassword), token, tokenExpiresAt)
	if err != nil {
		return nil, err
	}

	// Send verification email asynchronously
	if s.emailClient != nil {
		go func() {
			if err := s.emailClient.SendVerificationEmail(userEmail, token); err != nil {
				// Log error but don't fail signup
				println("Failed to send verification email:", err.Error())
			}
		}()
	}

	return user, nil
}

func (s *AuthService) VerifyEmail(token string) error {
	user, err := s.userRepo.FindByVerificationToken(token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid or expired verification token")
	}

	// Check if token is expired
	if user.VerificationTokenExpiresAt.Valid && user.VerificationTokenExpiresAt.Time.Before(time.Now()) {
		return errors.New("verification token has expired")
	}

	return s.userRepo.VerifyEmail(user.ID)
}

func (s *AuthService) ResendVerificationEmail(userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if user.EmailVerified {
		return errors.New("email already verified")
	}

	// Generate new verification token
	token, err := generateVerificationToken()
	if err != nil {
		return errors.New("failed to generate verification token")
	}
	tokenExpiresAt := time.Now().Add(24 * time.Hour)

	if err := s.userRepo.UpdateVerificationToken(userID, token, tokenExpiresAt); err != nil {
		return err
	}

	if s.emailClient == nil {
		return errors.New("email service not configured")
	}

	return s.emailClient.SendVerificationEmail(user.Email, token)
}

func (s *AuthService) Signin(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !user.PasswordHash.Valid {
		return nil, errors.New("please sign in with Google")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "heyspoilme",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
