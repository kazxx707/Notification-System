package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kavya/noti/models"
	"github.com/kavya/noti/notifier"
	"github.com/kavya/noti/repository"
)

// NotificationService handles business logic for notifications
type NotificationService struct {
	subscriptionRepo *repository.SubscriptionRepo
	notificationRepo *repository.NotificationRepo
	db               *sql.DB
	jsonFileMutex    sync.Mutex
	jsonFilePath     string
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	subscriptionRepo *repository.SubscriptionRepo,
	notificationRepo *repository.NotificationRepo,
	db *sql.DB,
	jsonFilePath string,
) *NotificationService {
	return &NotificationService{
		subscriptionRepo: subscriptionRepo,
		notificationRepo: notificationRepo,
		db:               db,
		jsonFilePath:     jsonFilePath,
	}
}

// Subscribe creates a subscription for a user to receive notifications
func (s *NotificationService) Subscribe(userID, itemID int64, channels []string) error {
	return s.subscriptionRepo.CreateSubscription(userID, itemID, channels)
}

// ProcessRestock handles restock event and sends notifications
// Uses transaction to ensure atomicity: notification insert + subscription update
func (s *NotificationService) ProcessRestock(itemID int64) error {
	// Get all pending subscriptions for this item
	subscriptions, err := s.subscriptionRepo.GetPendingSubscriptions(itemID)
	if err != nil {
		return fmt.Errorf("failed to get pending subscriptions: %w", err)
	}

	// Process each subscription
	for _, sub := range subscriptions {
		err := s.sendNotificationsForSubscription(sub)
		if err != nil {
			// Log error but continue with other subscriptions
			// One channel failure shouldn't block others
			fmt.Printf("Error processing subscription %d: %v\n", sub.ID, err)
		}
	}

	return nil
}

// sendNotificationsForSubscription sends notifications via all channels for one subscription
// Uses transaction to ensure atomicity (Transaction Pattern)
func (s *NotificationService) sendNotificationsForSubscription(sub models.Subscription) error {
	// Start transaction - ensures notification is stored before marking subscription as NOTIFIED
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Track if at least one notification succeeded
	hasSuccess := false

	// Send notification via each channel
	for _, channel := range sub.Channels {
		notifierInstance := notifier.NotificationSenderFactory(channel)
		if _, ok := notifierInstance.(*notifier.NoOpNotifier); ok {
			// Skip disabled channels
			continue
		}

		// Attempt to send notification
		err := notifierInstance.Send(sub.UserID, sub.ItemID)
		status := "SUCCESS"
		if err != nil {
			status = "FAILED"
			// Continue with other channels even if one fails
		} else {
			hasSuccess = true
		}

		// Create notification record
		notif := models.Notification{
			UserID:  sub.UserID,
			ItemID:  sub.ItemID,
			Channel: channel,
			Status:  status,
		}

		err = s.notificationRepo.CreateNotification(tx, notif)
		if err != nil {
			return fmt.Errorf("failed to create notification record: %w", err)
		}

		// Append to JSON file
		err = s.appendToJSONFile(notif)
		if err != nil {
			// Log but don't fail - DB is source of truth
			fmt.Printf("Warning: failed to write to JSON file: %v\n", err)
		}
	}

	// Only mark subscription as NOTIFIED if at least one notification succeeded
	// This ensures we retry if all channels failed
	if hasSuccess {
		err = s.subscriptionRepo.MarkAsNotified(tx, sub.ID)
		if err != nil {
			return fmt.Errorf("failed to mark subscription as notified: %w", err)
		}
	}

	// Commit transaction - this atomically saves notifications and updates subscription
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// appendToJSONFile appends notification to JSON file
func (s *NotificationService) appendToJSONFile(notif models.Notification) error {
	s.jsonFileMutex.Lock()
	defer s.jsonFileMutex.Unlock()

	// Read existing notifications
	var notifications []models.Notification
	data, err := os.ReadFile(s.jsonFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	if len(data) > 0 {
		err = json.Unmarshal(data, &notifications)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	// Add new notification with timestamp
	notif.CreatedAt = time.Now()
	notifications = append(notifications, notif)

	// Write back to file
	newData, err := json.MarshalIndent(notifications, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(s.jsonFilePath, newData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
