package services

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSendEmail(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	require.NoError(t, err)
	fromEmailAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	fromEmailPassword := os.Getenv("EMAIL_SENDER_PASSWORD")
	sender := NewEmailSender(fromEmailAddress, fromEmailPassword)
	subject := "Test email"
	body := "This is a test email"
	to := []string{""}

	err = sender.SendEmail(subject, body, to, nil, nil, nil)
	require.NoError(t, err)
}
