package services

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type Config struct {
	EmailSenderName     string `env:"EMAIL_SENDER_NAME"`
	fromEmailAddress    string `env:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword string `env:"EMAIL_SENDER_PASSWORD"`
}

func TestSendEmail(t *testing.T) {
	sender := NewEmailSender("", "")

	subject := "Test email"
	body := "This is a test email"
	to := []string{""}

	err := sender.SendEmail(subject, body, to, nil, nil, nil)
	require.NoError(t, err)
}
