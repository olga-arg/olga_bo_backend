package services

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"os"
	"sync"
)

const (
	smtpAuthAddress   = "smtp.zoho.com"
	smtpServerAddress = "smtp.zoho.com:587"
)

type Config struct {
	fromEmailAddress  string
	fromEmailPassword string
}

type emailService struct {
	fromEmail string
	auth      smtp.Auth
}

var (
	es   *emailService
	once sync.Once
)

type EmailSender interface {
	SendEmail(
		subject string,
		body string,
		to []string,
		cc []string,
	) error
}

func newEmailService(config Config) EmailSender {
	auth := smtp.PlainAuth("", config.fromEmailPassword, config.fromEmailPassword, smtpAuthAddress)
	return &emailService{fromEmail: config.fromEmailAddress, auth: auth}
}

func (es *emailService) SendEmail(subject, body string, to, cc []string) error {
	e := email.NewEmail()
	e.From = es.fromEmail
	e.To = to
	e.Cc = cc
	e.Subject = subject
	e.Text = []byte(body)
	err := e.Send(smtpServerAddress, es.auth)
	if err != nil {
		return err
	}
	return nil
}

func NewDefaultEmailService() EmailSender {
	config := Config{
		fromEmailAddress:  os.Getenv("EMAIL_SENDER_ADDRESS"),
		fromEmailPassword: os.Getenv("EMAIL_SENDER_PASSWORD"),
	}
	if config.fromEmailAddress == "" || config.fromEmailPassword == "" {
		panic("env variables must be set")
	}

	once.Do(func() {
		es = newEmailService(config).(*emailService)
	})

	return es
}
