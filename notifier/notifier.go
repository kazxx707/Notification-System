package notifier

// Notifier defines the interface for notification senders (Strategy Pattern)
type Notifier interface {
	Send(userID int64, itemID int64) error
}

// NotificationSenderFactory creates a notifier based on channel type (Factory Pattern)
func NotificationSenderFactory(channel string) Notifier {
	switch channel {
	case "email":
		return &EmailNotifier{}
	case "sms":
		return &SMSNotifier{}
	case "push":
		return &PushNotifier{}
	default:
		// Return a no-op notifier for unknown channels
		return &NoOpNotifier{}
	}
}

// NoOpNotifier handles disabled or unknown channels
type NoOpNotifier struct{}

func (n *NoOpNotifier) Send(userID int64, itemID int64) error {
	// Silently skip unknown channels
	return nil
}
