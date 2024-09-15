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

type TelegramNotifier struct{}

func (t *TelegramNotifier) Send(message string) {
	// send sms
	fmt.Printf("Sending telegram: %s (Sender: Telegram)\n", message)
}

type NotifierDecorator struct {
	notifier Notifier
	core     *NotifierDecorator // point to another decorator
}

func (nd NotifierDecorator) Send(message string) {
	nd.notifier.Send(message)

	if nd.core != nil {
		nd.core.Send(message)
	}
}

func (nd NotifierDecorator) Decorate(notifier Notifier) NotifierDecorator {
	return NotifierDecorator{
		core:     &nd,
		notifier: notifier,
	}
}

func NewNotifierDecorator(notifier Notifier) NotifierDecorator {
	return NotifierDecorator{
		notifier: notifier,
	}
}

func main() {
	notifier := NewNotifierDecorator(&EmailNotifier{}).
		Decorate(&SMSNotifier{}).
		Decorate(&TelegramNotifier{})

	notifier.Send("Hello, welcome to the world of Go!")
}
