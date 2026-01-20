package notifier

import "fmt"

// SMSNotifier handles SMS notifications (Strategy Pattern implementation)
type SMSNotifier struct{}

func (s *SMSNotifier) Send(userID int64, itemID int64) error {
	// Mock SMS sending - in real system, this would call SMS service
	fmt.Printf("[MOCK] Sending SMS notification to user %d for item %d\n", userID, itemID)
	return nil
}
