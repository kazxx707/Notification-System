package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kavya/noti/models"
)

// SubscriptionRepo handles subscription database operations
type SubscriptionRepo struct {
	db *sql.DB
}

// NewSubscriptionRepo creates a new subscription repository
func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

// CreateSubscription creates a new subscription
// If user already subscribed, updates channels and resets status to PENDING
func (r *SubscriptionRepo) CreateSubscription(userID, itemID int64, channels []string) error {
	channelsJSON, err := json.Marshal(channels)
	if err != nil {
		return fmt.Errorf("failed to marshal channels: %w", err)
	}

	query := `
		INSERT INTO subscriptions (user_id, item_id, channels, status)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, item_id)
		DO UPDATE SET channels = $3, status = $4
	`

	_, err = r.db.Exec(query, userID, itemID, string(channelsJSON), models.StatusPending)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

// GetPendingSubscriptions returns all PENDING subscriptions for an item
func (r *SubscriptionRepo) GetPendingSubscriptions(itemID int64) ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, item_id, channels, status
		FROM subscriptions
		WHERE item_id = $1 AND status = $2
	`

	rows, err := r.db.Query(query, itemID, models.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		var channelsJSON string

		err := rows.Scan(&sub.ID, &sub.UserID, &sub.ItemID, &channelsJSON, &sub.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		err = json.Unmarshal([]byte(channelsJSON), &sub.Channels)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal channels: %w", err)
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

// MarkAsNotified marks a subscription as NOTIFIED (used in transaction)
func (r *SubscriptionRepo) MarkAsNotified(tx *sql.Tx, subscriptionID int64) error {
	query := `UPDATE subscriptions SET status = $1 WHERE id = $2`
	_, err := tx.Exec(query, models.StatusNotified, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to mark subscription as notified: %w", err)
	}
	return nil
}
