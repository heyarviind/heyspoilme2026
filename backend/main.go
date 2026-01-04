package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var discordWebhookURL string

type EmailSubscription struct {
	Email  string `json:"email" binding:"required"`
	Gender string `json:"gender" binding:"required"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type DiscordWebhook struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Color       int           `json:"color,omitempty"`
	Fields      []EmbedField  `json:"fields,omitempty"`
	Timestamp   string        `json:"timestamp,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func main() {
	// Discord webhook URL
	discordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	if discordWebhookURL == "" {
		discordWebhookURL = "https://discord.com/api/webhooks/1457381901642895444/ubTk45MMMQ5XZBxYXeN44ibk6j7MtQSVcqPwPK5fZeRxjjkViwQYb82JtjmwD1uqxqe1"
	}

	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5433/heyspoilme?sslmode=disable"
	}

	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for database... attempt %d/30", i+1)
		time.Sleep(time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS email_subscriptions (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			gender VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Setup Gin
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	r.GET("/health", healthCheck)
	r.POST("/api/subscribe", subscribe)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func subscribe(c *gin.Context) {
	var input EmailSubscription

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email and gender are required",
		})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid email format",
		})
		return
	}

	// Validate gender
	if input.Gender != "man" && input.Gender != "woman" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid gender selection",
		})
		return
	}

	// Insert into database
	_, err := db.Exec("INSERT INTO email_subscriptions (email, gender) VALUES ($1, $2) ON CONFLICT (email) DO UPDATE SET gender = $2", input.Email, input.Gender)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save email",
		})
		return
	}

	// Send Discord notification (async)
	go sendDiscordNotification(input.Email, input.Gender)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Thank you! We'll notify you when we launch.",
	})
}

func sendDiscordNotification(email, gender string) {
	if discordWebhookURL == "" {
		return
	}

	genderLabel := "Man"
	if gender == "woman" {
		genderLabel = "Woman"
	}

	webhook := DiscordWebhook{
		Embeds: []DiscordEmbed{
			{
				Title:       "New Waitlist Signup",
				Description: "Someone just joined the HeySpoilMe waitlist!",
				Color:       0x00FF00, // Green
				Fields: []EmbedField{
					{
						Name:   "Email",
						Value:  email,
						Inline: true,
					},
					{
						Name:   "Gender",
						Value:  genderLabel,
						Inline: true,
					},
				},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	jsonData, err := json.Marshal(webhook)
	if err != nil {
		log.Printf("Discord webhook marshal error: %v", err)
		return
	}

	resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Discord webhook error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("Discord webhook returned status: %d", resp.StatusCode)
	}
}
