package notifier

import "fmt"

// EmailNotifier handles email notifications (Strategy Pattern implementation)
type EmailNotifier struct{}

func (e *EmailNotifier) Send(userID int64, itemID int64) error {
	// Mock email sending - in real system, this would call email service
	fmt.Printf("[MOCK] Sending email notification to user %d for item %d\n", userID, itemID)
	return nil
}
