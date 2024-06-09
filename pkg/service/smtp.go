package service

import (
	"net/smtp"
	"os"
)

func sendMail(to []string, subject, body string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	sender := os.Getenv("SMTP_SENDER")
	password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", sender, password, smtpServer)
	smtpAddr := smtpServer + ":" + smtpPort

	return smtp.SendMail(smtpAddr, auth, sender, to, []byte(
		"Subject: "+subject+"\r\n"+
			"MIME-version: 1.0;\r\n"+
			"Content-Type: text/plain; charset=\"UTF-8\";\r\n"+
			"\r\n"+
			body,
	))
}
