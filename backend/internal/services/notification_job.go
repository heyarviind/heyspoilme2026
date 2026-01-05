package services

import (
	"log"
	"time"

	"heyspoilme/internal/repository"
	"heyspoilme/pkg/email"
)

type NotificationJobService struct {
	messageRepo *repository.MessageRepository
	emailClient *email.ZeptoMailClient
	interval    time.Duration
	unreadDelay time.Duration
	stopChan    chan struct{}
}

func NewNotificationJobService(messageRepo *repository.MessageRepository, emailClient *email.ZeptoMailClient) *NotificationJobService {
	return &NotificationJobService{
		messageRepo: messageRepo,
		emailClient: emailClient,
		interval:    1 * time.Minute,  // Check every minute
		unreadDelay: 5 * time.Minute,  // Send notification after 5 minutes unread
		stopChan:    make(chan struct{}),
	}
}

// Start begins the background job that checks for unread messages
func (s *NotificationJobService) Start() {
	if s.emailClient == nil {
		log.Printf("[NotificationJob] Email client not configured, skipping notification job")
		return
	}

	log.Printf("[NotificationJob] Starting unread message notification job (checking every %v, notifying after %v unread)", s.interval, s.unreadDelay)
	
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run immediately on start
	s.processUnreadMessages()

	for {
		select {
		case <-ticker.C:
			s.processUnreadMessages()
		case <-s.stopChan:
			log.Printf("[NotificationJob] Stopping notification job")
			return
		}
	}
}

// Stop stops the background job
func (s *NotificationJobService) Stop() {
	close(s.stopChan)
}

func (s *NotificationJobService) processUnreadMessages() {
	notifications, err := s.messageRepo.GetUnreadMessagesNeedingNotification(s.unreadDelay)
	if err != nil {
		log.Printf("[NotificationJob] Error fetching unread messages: %v", err)
		return
	}

	if len(notifications) == 0 {
		return
	}

	log.Printf("[NotificationJob] Found %d messages needing notification", len(notifications))

	for _, n := range notifications {
		err := s.emailClient.SendNewMessageNotification(n.RecipientEmail, n.SenderName, n.Content)
		if err != nil {
			log.Printf("[NotificationJob] Failed to send notification email to %s: %v", n.RecipientEmail, err)
			continue
		}

		// Mark the notification as sent
		if err := s.messageRepo.MarkNotificationEmailSent(n.MessageID); err != nil {
			log.Printf("[NotificationJob] Failed to mark notification as sent for message %s: %v", n.MessageID, err)
		}

		log.Printf("[NotificationJob] Sent notification email to %s for message from %s", n.RecipientEmail, n.SenderName)
	}
}


