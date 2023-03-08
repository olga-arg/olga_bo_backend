package services

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress   = "smtp.zoho.com"
	smtpServerAddress = "smtp.zoho.com:587"
)

type Config struct {
	EmailSenderName     string `env:"EMAIL_SENDER_NAME"`
	fromEmailAddress    string `env:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `env:"EMAIL_SENDER_PASSWORD"`
}

type EmailSender interface {
	SendEmail(
		subject string,
		body string,
		to []string,
		cc []string,
		bcc []string,
		attachFile []string,
	) error
}

type emailSender struct {
	fromEmailAddress  string
	fromEmailPassword string
}

func NewEmailSender(fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &emailSender{
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *emailSender) SendEmail(
	subject string,
	body string,
	to []string,
	cc []string,
	bcc []string,
	attachFile []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf(sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(body)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
