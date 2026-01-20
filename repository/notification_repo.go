package repository

import (
	"database/sql"
	"fmt"

	"github.com/kavya/noti/models"
)

// NotificationRepo handles notification database operations
type NotificationRepo struct {
	db *sql.DB
}

// NewNotificationRepo creates a new notification repository
func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{db: db}
}

// CreateNotification creates a new notification within a transaction
func (r *NotificationRepo) CreateNotification(tx *sql.Tx, notif models.Notification) error {
	query := `
		INSERT INTO notifications (user_id, item_id, channel, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := tx.QueryRow(query, notif.UserID, notif.ItemID, notif.Channel, notif.Status).
		Scan(&notif.ID, &notif.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}
