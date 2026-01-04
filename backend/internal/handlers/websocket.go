package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gorillaws "github.com/gorilla/websocket"

	"heyspoilme/internal/services"
	"heyspoilme/internal/websocket"
)

type WebSocketHandler struct {
	hub             *websocket.Hub
	authService     *services.AuthService
	presenceService *services.PresenceService
	upgrader        gorillaws.Upgrader
}

func NewWebSocketHandler(hub *websocket.Hub, authService *services.AuthService, presenceService *services.PresenceService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:             hub,
		authService:     authService,
		presenceService: presenceService,
		upgrader: gorillaws.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get token from query parameter for WebSocket
	token := c.Query("token")
	if token == "" {
		// Also try Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID := claims.UserID

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket.NewClient(h.hub, conn, userID)
	h.hub.Register(client)

	h.presenceService.SetOnline(userID)

	go func() {
		client.WritePump()
	}()

	go func() {
		client.ReadPump()
		h.presenceService.SetOffline(userID)
	}()
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := websocket.NewClient(h.hub, conn, userID)
	h.hub.Register(client)

	h.presenceService.SetOnline(userID)

	go func() {
		client.WritePump()
	}()

	go func() {
		client.ReadPump()
		h.presenceService.SetOffline(userID)
	}()
}
