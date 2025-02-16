package email

import (
	"os"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Dialer *gomail.Dialer
	From   string
}

func NewEmailSender() *EmailSender {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587 // Porta padrão, altere conforme necessário
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	dialer := gomail.NewDialer(
		smtpHost,
		smtpPort,
		smtpUsername,
		smtpPassword,
	)

	return &EmailSender{
		Dialer: dialer,
		From:   "fideliplusapp@gmail.com",
	}
}

func (s *EmailSender) SendEmail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return s.Dialer.DialAndSend(msg)
}
