package notifier

import "fmt"

// PushNotifier handles push notifications (Strategy Pattern implementation)
type PushNotifier struct{}

func (p *PushNotifier) Send(userID int64, itemID int64) error {
	// Mock push notification - in real system, this would call push service
	fmt.Printf("[MOCK] Sending push notification to user %d for item %d\n", userID, itemID)
	return nil
}
