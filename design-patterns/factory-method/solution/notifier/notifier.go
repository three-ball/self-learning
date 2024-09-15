package notifier

import "fmt"

type Notifier interface {
	Send(message string)
}

type emailNotifier struct{}

func (e *emailNotifier) Send(message string) {
	// send email
	fmt.Printf("Sending email: %s (Sender: Email)\n", message)
}

type sMSNotifier struct{}

func (s *sMSNotifier) Send(message string) {
	// send sms
	fmt.Printf("Sending sms: %s (Sender: SMS)\n", message)
}

func NewNotifier(notificationType string) Notifier {
	if notificationType == "email" {
		return &emailNotifier{}
	}
	if notificationType == "sms" {
		return &sMSNotifier{}
	}
	return nil
}
