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
	"heyspoilme/pkg/storage"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize S3 storage
	s3Client, err := storage.NewS3Client(cfg.AWSRegion, cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.S3Bucket)
	if err != nil {
		log.Printf("Warning: S3 client not initialized: %v", err)
	}

	// Initialize Google OAuth
	googleAuth := auth.NewGoogleAuth(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	profileRepo := repository.NewProfileRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	presenceRepo := repository.NewPresenceRepository(db)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	profileService := services.NewProfileService(profileRepo, userRepo)
	chatService := services.NewChatService(messageRepo, profileRepo, hub)
	likeService := services.NewLikeService(likeRepo, notificationRepo, hub)
	notificationService := services.NewNotificationService(notificationRepo)
	presenceService := services.NewPresenceService(presenceRepo, hub)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, profileService, googleAuth, cfg.FrontendURL)
	profileHandler := handlers.NewProfileHandler(profileService)
	chatHandler := handlers.NewChatHandler(chatService)
	uploadHandler := handlers.NewUploadHandler(s3Client)
	likeHandler := handlers.NewLikeHandler(likeService, profileService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	wsHandler := handlers.NewWebSocketHandler(hub, authService, presenceService)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

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
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(authMiddleware.RequireAuth())
	{
		// Profile routes
		api.POST("/profile", profileHandler.CreateProfile)
		api.GET("/profile", profileHandler.GetMyProfile)
		api.PUT("/profile", profileHandler.UpdateProfile)
		api.GET("/profiles", profileHandler.ListProfiles)
		api.GET("/profiles/:id", profileHandler.GetProfile)

		// Upload routes
		api.POST("/upload/presigned-url", uploadHandler.GetPresignedURL)
		api.POST("/profile/images", profileHandler.AddProfileImage)
		api.DELETE("/profile/images/:imageId", profileHandler.DeleteProfileImage)

		// Like routes
		api.POST("/profiles/:id/like", likeHandler.LikeProfile)
		api.DELETE("/profiles/:id/like", likeHandler.UnlikeProfile)
		api.GET("/likes/received", likeHandler.GetReceivedLikes)
		api.GET("/likes/given", likeHandler.GetGivenLikes)

		// Chat routes
		api.GET("/conversations", chatHandler.GetConversations)
		api.POST("/conversations", chatHandler.CreateConversation)
		api.GET("/conversations/:id/messages", chatHandler.GetMessages)
		api.POST("/conversations/:id/messages", chatHandler.SendMessage)

		// Notification routes
		api.GET("/notifications", notificationHandler.GetNotifications)
		api.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		api.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
		api.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
	}

	// WebSocket route
	r.GET("/ws", wsHandler.HandleWebSocket)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

