package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jordan-wright/email"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"sync"
)

const (
	smtpAuthAddress   = "smtp.zoho.com"
	smtpServerAddress = "smtp.zoho.com:587"
)

type EmailTemplate string

const (
	Welcome EmailTemplate = "welcome_email.html"
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
		to string,
		template EmailTemplate,
		variablesTemplate []string,
		cc []string,
	) error
}

func newEmailService(config Config) EmailSender {
	auth := smtp.PlainAuth("", config.fromEmailAddress, config.fromEmailPassword, smtpAuthAddress)
	return &emailService{fromEmail: config.fromEmailAddress, auth: auth}
}

func (es *emailService) SendEmail(to string, template EmailTemplate, variablesTemplate, cc []string) error {
	e := email.NewEmail()
	e.From = es.fromEmail
	e.To = []string{to}
	e.Cc = cc
	e.Subject = "Te damos la bienvenida a Olga :)"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return err
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: aws.String("prod-olga-backend-assets"),
		Key:    aws.String(string(template)),
	}

	resp, err := client.GetObject(context.TODO(), input)
	if err != nil {
		fmt.Println("Error fetching template from S3:", err)
		return err
	}

	htmlContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading index.html:", err)
		return err
	}

	body := strings.ReplaceAll(string(htmlContent), "{{name}}", variablesTemplate[0])
	e.HTML = []byte(body)
	err = e.Send(smtpServerAddress, es.auth)
	if err != nil {
		fmt.Println("Sending email from: ", es.fromEmail, " to: ", to, " with subject: ", e.Subject, "got error: ", err)
		return err
	}
	return nil
}

func NewDefaultEmailService() EmailSender {
	emailAddrB64 := os.Getenv("EMAIL_SENDER_ADDRESS")
	emailPassB64 := os.Getenv("EMAIL_SENDER_PASSWORD")
	if emailAddrB64 == "" || emailPassB64 == "" {
		panic("env variables must be set")
	}
	emailAddr, err := base64.StdEncoding.DecodeString(emailAddrB64)
	if err != nil {
		panic("failed to decode env variable")
	}
	emailPass, err := base64.StdEncoding.DecodeString(emailPassB64)
	if err != nil {
		panic("failed to decode env variable")
	}
	emailAddr = []byte(strings.TrimSuffix(string(emailAddr), "\n"))
	config := Config{
		fromEmailAddress:  string(emailAddr),
		fromEmailPassword: string(emailPass),
	}

	once.Do(func() {
		es = newEmailService(config).(*emailService)
	})

	return es
}
