package email_service

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAdress = "smtp.olga.lat"
	smtpServerAdress = "smtp.olga.lat:587"
)

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
	name string
	fromEmailAdress string
	fromEmailPassword string
}

func NewEmailSender(name string, fromEmailAdress string, fromEmailPassword string)
	return &emailSender{
		name: name,
		fromEmailAdress: fromEmailAdress,
		fromEmailPassword: fromEmailPassword,
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
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAdress)
	e.Subject = subject
	e.HTML = []byte(body)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, file := range attachFile {
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", file, err)
	}

	smtpAuth := stmp.PlainAuth("", sender.fromEmailAdress, sender.fromEmailPassword, smtpAuthAdress)
	return e.Send(smtpServerAdress, smtpAuth)
}

