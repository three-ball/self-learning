package main

import "fmt"

type NotificationService struct {
	notificationType string
}

func (n *NotificationService) SendNotification(message string) {
	if n.notificationType == "email" {
		fmt.Printf("Sending email: %s (Sender: Email)", message)
	} else if n.notificationType == "sms" {
		fmt.Printf("Sending sms: %s (Sender: SMS)", message)
	} else {
		// send push notification
	}
}

func main() {
	emailNotification := NotificationService{notificationType: "email"}
	emailNotification.SendNotification("Hello, welcome to the world of Go!")

	smsNotification := NotificationService{notificationType: "sms"}
	smsNotification.SendNotification("Hello, welcome to the world of Go!")
}
