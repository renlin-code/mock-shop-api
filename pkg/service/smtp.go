package service

import (
	"net/smtp"
	"os"

	"github.com/spf13/viper"
)

func sendMail(to []string, subject, body string) error {
	smtpServer := viper.GetString("smtp.server")
	smtpPort := viper.GetString("smtp.port")
	sender := viper.GetString("smtp.sender")
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
