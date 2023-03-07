package services

import (
	"testing"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/require"
)

type Config struct {
	EmailSenderName string `env:"EMAIL_SENDER_NAME"`
	EmailSenderAdress string `env:"EMAIL_SENDER_ADRESS"`
	EmailSenderPassword string `env:"EMAIL_SENDER_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

func TestSendEmail(t *testing.T) {
	cfg, err := LoadConfig()
	require.NoError(t, err)

	sender := NewEmailSender(cfg.EmailSenderName, cfg.EmailSenderAdress, cfg.EmailSenderPassword)

	subject := "Test email"
	body := "This is a test email"
	to := []string{"ignacio.nahuel.ramos@gmail.com"}

	err = sender.SendEmail(subject, body, to, nil, nil, nil)
	require.NoError(t, err)
}

