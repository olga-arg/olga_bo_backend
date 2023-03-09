package services

import (
	"encoding/base64"
	"github.com/jordan-wright/email"
	"log"
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
		to string,
		cc []string,
	) error
}

func newEmailService(config Config) EmailSender {
	auth := smtp.PlainAuth("", config.fromEmailPassword, config.fromEmailPassword, smtpAuthAddress)
	return &emailService{fromEmail: config.fromEmailAddress, auth: auth}
}

func (es *emailService) SendEmail(subject, body, to string, cc []string) error {
	e := email.NewEmail()
	e.From = es.fromEmail
	log.Println("email to: ", to)
	log.Println("email to parsed: ", []string{to})
	//e.To = []string{to}
	e.To = []string{"vilavalentin@gmail.com"}
	e.Cc = cc
	e.Subject = subject
	e.Text = []byte(body)
	err := e.Send(smtpServerAddress, es.auth)
	log.Println("email sent: ", err)
	if err != nil {
		return err
	}
	return nil
}

func NewDefaultEmailService() EmailSender {
	emailAddrB64 := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailPassB64 := os.Getenv("EMAIL_SENDER_PASSWORD")
	emailAddr, err := base64.StdEncoding.DecodeString(emailAddrB64)
	if err != nil {
		panic("env variables must be set")
	}
	emailPass, err := base64.StdEncoding.DecodeString(emailPassB64)
	if err != nil {
		panic("env variables must be set")
	}
	log.Println("emailAddr: ", string(emailPass[0]), string(emailPass[11]))
	config := Config{
		fromEmailAddress:  string(emailAddr),
		fromEmailPassword: string(emailPass),
	}

	once.Do(func() {
		es = newEmailService(config).(*emailService)
	})

	return es
}
