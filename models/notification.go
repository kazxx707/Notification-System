package models

import "time"

// Notification represents a sent notification
type Notification struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ItemID    int64     `json:"item_id" db:"item_id"`
	Channel   string    `json:"channel" db:"channel"` // email, sms, push
	Status    string    `json:"status" db:"status"`   // SUCCESS or FAILED
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
