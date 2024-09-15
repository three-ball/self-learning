package main

import "factory-method/solution/notifier"

func main() {
	emailNoti := notifier.NewNotifier("email")
	emailNoti.Send("Hello, welcome to the world of Go!")

	smsNoti := notifier.NewNotifier("sms")
	smsNoti.Send("Hello, welcome to the world of Go!")
}
