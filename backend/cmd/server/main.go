package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"heyspoilme/internal/config"
	"heyspoilme/internal/database"
	"heyspoilme/internal/handlers"
	"heyspoilme/internal/middleware"
	"heyspoilme/internal/repository"
	"heyspoilme/internal/services"
	"heyspoilme/internal/websocket"
	"heyspoilme/pkg/auth"
	"heyspoilme/pkg/email"
	"heyspoilme/pkg/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Debug: Log S3 configuration
	log.Printf("[Config] S3 Configuration:")
	log.Printf("  AWS_ACCESS_KEY_ID: %s", maskSecret(cfg.AWSAccessKeyID))
	log.Printf("  AWS_SECRET_ACCESS_KEY: %s", maskSecret(cfg.AWSSecretAccessKey))
	log.Printf("  AWS_REGION: %s", cfg.AWSRegion)
	log.Printf("  S3_BUCKET: %s", cfg.S3Bucket)
	log.Printf("  S3_ENDPOINT: %s", cfg.S3Endpoint)
	log.Printf("  S3_BASE_URL: %s", cfg.S3BaseURL)

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize S3 storage
	s3Client, err := storage.NewS3Client(cfg.AWSRegion, cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.S3Bucket, cfg.S3BaseURL, cfg.S3Endpoint)
	if err != nil {
		log.Printf("Warning: S3 client not initialized: %v", err)
	}

	// Initialize Google OAuth
	googleAuth := auth.NewGoogleAuth(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)

	// Initialize ZeptoMail client
	var emailClient *email.ZeptoMailClient
	if cfg.ZeptoMailAPIKey != "" {
		var err error
		emailClient, err = email.NewZeptoMailClient(cfg.ZeptoMailAPIKey, cfg.ZeptoMailFromEmail, cfg.ZeptoMailFromName, cfg.FrontendURL)
		if err != nil {
			log.Printf("Warning: Email client not initialized: %v", err)
		} else {
			log.Printf("[Config] ZeptoMail client initialized")
		}
	} else {
		log.Printf("Warning: ZEPTOMAIL_API_KEY not set - email verification disabled")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	presenceRepo := repository.NewPresenceRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	cityRepo := repository.NewCityRepository(db)
	featureFlagRepo := repository.NewFeatureFlagRepository(db)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Start background notification job
	notificationJob := services.NewNotificationJobService(messageRepo, emailClient)
	go notificationJob.Start()

	// Start background ranking job (updates profile scores every 15 minutes)
	rankingService := services.NewRankingService(db)
	go rankingService.Start()

	// Initialize feature flag service first (used by other services)
	featureFlagService := services.NewFeatureFlagService(featureFlagRepo)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, emailClient)
	profileService := services.NewProfileService(profileRepo, userRepo)
	chatService := services.NewChatService(messageRepo, profileRepo, userRepo, hub, featureFlagService)
	likeService := services.NewLikeService(likeRepo, notificationRepo, profileRepo, hub, featureFlagService)
	notificationService := services.NewNotificationService(notificationRepo)
	presenceService := services.NewPresenceService(presenceRepo, hub)
	accountService := services.NewAccountService(userRepo, profileRepo, messageRepo, likeRepo, notificationRepo, presenceRepo, s3Client)
	verificationService := services.NewVerificationService(verificationRepo, profileRepo)
	adminService := services.NewAdminService(adminRepo, s3Client)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, profileService, accountService, googleAuth, cfg.FrontendURL)
	profileHandler := handlers.NewProfileHandler(profileService)
	chatHandler := handlers.NewChatHandler(chatService)
	uploadHandler := handlers.NewUploadHandler(s3Client)
	likeHandler := handlers.NewLikeHandler(likeService, profileService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	verificationHandler := handlers.NewVerificationHandler(verificationService)
	wsHandler := handlers.NewWebSocketHandler(hub, authService, presenceService)
	adminHandler := handlers.NewAdminHandler(adminService, featureFlagService, cfg.AdminCode1, cfg.AdminCode2)
	cityHandler := handlers.NewCityHandler(cityRepo)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	verificationMiddleware := middleware.NewVerificationMiddleware(userRepo, featureFlagService)

	// Setup Gin
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL, "http://localhost:5173", "http://localhost:3000", "http://localhost:3003"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "timestamp": time.Now().UTC().Format(time.RFC3339)})
	})

	// Auth routes
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/signup", authHandler.Signup)
		authRoutes.POST("/signin", authHandler.Signin)
		authRoutes.GET("/google", authHandler.GoogleLogin)
		authRoutes.GET("/google/callback", authHandler.GoogleCallback)
		authRoutes.POST("/logout", authHandler.Logout)
		authRoutes.GET("/me", authMiddleware.RequireAuth(), authHandler.GetCurrentUser)
		authRoutes.GET("/verify-email", authHandler.VerifyEmail)
		authRoutes.POST("/resend-verification", authMiddleware.RequireAuth(), authHandler.ResendVerificationEmail)
		authRoutes.DELETE("/account", authMiddleware.RequireAuth(), authHandler.DeleteAccount)
	}

	// Public routes for autocomplete (no auth required)
	r.GET("/api/cities/search", cityHandler.SearchCities)

	// Public feature flags endpoint (returns restrictions status)
	r.GET("/api/feature-flags", adminHandler.GetPublicFeatureFlags)

	// Protected routes (basic auth only)
	api := r.Group("/api")
	api.Use(authMiddleware.RequireAuth())
	{
		// Profile routes - basic viewing allowed without verification
		api.POST("/profile", profileHandler.CreateProfile)
		api.GET("/profile", profileHandler.GetMyProfile)
		api.PUT("/profile", profileHandler.UpdateProfile)
		api.GET("/profiles/:id", profileHandler.GetProfile)

		// Like routes - viewing allowed without verification
		api.GET("/likes/received", likeHandler.GetReceivedLikes)
		api.GET("/likes/given", likeHandler.GetGivenLikes)

		// Chat routes - viewing allowed without verification
		api.GET("/conversations", chatHandler.GetConversations)
		api.GET("/conversations/:id/messages", chatHandler.GetMessages)
		api.GET("/messages/unread-count", chatHandler.GetUnreadCount)
		api.GET("/inbox", chatHandler.GetInbox) // Inbox with locked message support

		// Notification routes
		api.GET("/notifications", notificationHandler.GetNotifications)
		api.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		api.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
		api.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)

		// Identity verification routes
		api.GET("/verification/code", verificationHandler.GenerateCode)
		api.POST("/verification/submit", verificationHandler.SubmitVerification)
		api.GET("/verification/status", verificationHandler.GetStatus)
	}

	// Protected routes requiring email verification
	verifiedAPI := r.Group("/api")
	verifiedAPI.Use(authMiddleware.RequireAuth(), verificationMiddleware.RequireEmailVerified())
	{
		// Profile browsing with filters requires verification
		verifiedAPI.GET("/profiles", profileHandler.ListProfiles)

		// Upload routes require verification
		verifiedAPI.POST("/upload/presigned-url", uploadHandler.GetPresignedURL)
		verifiedAPI.POST("/upload/chat-image-url", uploadHandler.GetChatImagePresignedURL)
		verifiedAPI.POST("/profile/images", profileHandler.AddProfileImage)
		verifiedAPI.DELETE("/profile/images/:imageId", profileHandler.DeleteProfileImage)

		// Like actions require verification
		verifiedAPI.POST("/profiles/:id/like", likeHandler.LikeProfile)
		verifiedAPI.DELETE("/profiles/:id/like", likeHandler.UnlikeProfile)

		// Messaging requires verification
		verifiedAPI.POST("/conversations", chatHandler.CreateConversation)
		verifiedAPI.POST("/conversations/:id/messages", chatHandler.SendMessage)
	}

	// WebSocket route
	r.GET("/ws", wsHandler.HandleWebSocket)

	// Admin routes (protected by secret URL path)
	adminRoutes := r.Group("/admin/:code1/:code2")
	adminRoutes.Use(adminHandler.ValidateAdminAccess())
	{
		adminRoutes.GET("/stats", adminHandler.GetStats)
		adminRoutes.GET("/users", adminHandler.ListUsers)
		adminRoutes.GET("/users/:userId", adminHandler.GetUser)
		adminRoutes.DELETE("/users/:userId", adminHandler.DeleteUser)
		adminRoutes.PUT("/users/:userId/wealth-status", adminHandler.UpdateUserWealthStatus)
		adminRoutes.GET("/messages", adminHandler.ListMessages)
		adminRoutes.GET("/verifications", adminHandler.ListVerificationRequests)
		adminRoutes.POST("/verifications/:requestId/approve", adminHandler.ApproveVerification)
		adminRoutes.POST("/verifications/:requestId/reject", adminHandler.RejectVerification)
		adminRoutes.GET("/images", adminHandler.ListAllImages)
		adminRoutes.DELETE("/images/:imageId", adminHandler.DeleteProfileImage)
		adminRoutes.GET("/feature-flags", adminHandler.GetFeatureFlags)
		adminRoutes.PUT("/feature-flags/:key", adminHandler.UpdateFeatureFlag)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func maskSecret(s string) string {
	if s == "" {
		return "(not set)"
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****"
}

