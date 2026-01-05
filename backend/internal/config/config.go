package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Server
	Port string

	// Database
	DatabaseURL string

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// JWT
	JWTSecret string

	// S3 / Cloudflare R2
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	S3Bucket           string
	S3BaseURL          string
	S3Endpoint         string

	// Discord (optional)
	DiscordWebhookURL string

	// Frontend URL
	FrontendURL string

	// ZeptoMail
	ZeptoMailAPIKey    string
	ZeptoMailFromEmail string
	ZeptoMailFromName  string

	// Admin
	AdminCode1 string
	AdminCode2 string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5433/heyspoilme?sslmode=disable"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/google/callback"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", "auto"),
		S3Bucket:           getEnv("S3_BUCKET", "heyspoilme"),
		S3BaseURL:          getEnv("S3_BASE_URL", ""),
		S3Endpoint:         getEnv("S3_ENDPOINT", ""),
		DiscordWebhookURL:  getEnv("DISCORD_WEBHOOK_URL", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3003"),
		ZeptoMailAPIKey:    getEnv("ZEPTOMAIL_API_KEY", ""),
		ZeptoMailFromEmail: getEnv("ZEPTOMAIL_FROM_EMAIL", "noreply@heyspoilme.com"),
		ZeptoMailFromName:  getEnv("ZEPTOMAIL_FROM_NAME", "HeySpoilMe"),
		AdminCode1:         getEnv("ADMIN_CODE_1", "super"),
		AdminCode2:         getEnv("ADMIN_CODE_2", "secret"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
