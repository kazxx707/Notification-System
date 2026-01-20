package models

// Subscription represents a user's subscription to item restock notifications
type Subscription struct {
	ID       int64    `json:"id" db:"id"`
	UserID   int64    `json:"user_id" db:"user_id"`
	ItemID   int64    `json:"item_id" db:"item_id"`
	Channels []string `json:"channels" db:"channels"` // JSON array stored as text
	Status   string   `json:"status" db:"status"`     // PENDING or NOTIFIED
}

// Status constants for subscription
const (
	StatusPending = "PENDING"
	StatusNotified = "NOTIFIED"
)
