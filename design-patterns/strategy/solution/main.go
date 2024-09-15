package main

import "fmt"

type Notifier interface {
	Send(message string)
}

type EmailNotifier struct{}

func (e *EmailNotifier) Send(message string) {
	// send email
	fmt.Printf("Sending email: %s (Sender: Email)\n", message)
}

type SMSNotifier struct{}

func (s *SMSNotifier) Send(message string) {
	// send sms
	fmt.Printf("Sending sms: %s (Sender: SMS)\n", message)
}

type PushNotifier struct {
	notifier Notifier
}

func (p *PushNotifier) SendMessage(message string) {
	p.notifier.Send(message)
}

func main() {
	emailNotifier := &EmailNotifier{}
	emailNotification := &PushNotifier{notifier: emailNotifier}
	emailNotification.SendMessage("Hello, welcome to the world of Go!")

	smsNotifier := &SMSNotifier{}
	smsNotification := &PushNotifier{notifier: smsNotifier}
	smsNotification.SendMessage("Hello, welcome to the world of Go!")
}
