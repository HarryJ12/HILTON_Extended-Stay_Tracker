package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func sendEmail(subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	to := os.Getenv("MANAGER_EMAIL")

	if host == "" || port == "" || from == "" || pass == "" || to == "" {
		return fmt.Errorf("smtp or manager env vars not set")
	}

	auth := smtp.PlainAuth("", from, pass, host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
			body + "\r\n",
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
}
