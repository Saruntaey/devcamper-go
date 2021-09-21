package utils

import (
	"net/smtp"
	"os"
)

func SendMail(to string, subj string, body string) bool {
	// Set up authentication information.
	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	msg := []byte("From: " + os.Getenv("FROM_EMAIL") + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subj + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, os.Getenv("FROM_EMAIL"), []string{to}, msg)

	return err == nil
}
