package app_client_smtp

import (
	"context"
	"crypto/tls"
	"fmt"

	usecases "github.com/DanielTitovsky/notificationService/internal/app/useCases"
	"github.com/wneessen/go-mail"
)

type SMTPClient struct {
	config Config
	cl     *mail.Client
}

func NewSMTPClient(config Config) (SMTPClient, error) {
	client, err := mail.NewClient(
		config.Host,
		mail.WithPort(config.Port),
		mail.WithSMTPAuth(mail.SMTPAuthType(config.AuthPlan)),
		mail.WithUsername(config.UserName),
		mail.WithPassword(config.Password),
		mail.WithTLSConfig(&tls.Config{
			ServerName: config.Host,
			MinVersion: tls.VersionTLS12,
		}),
	)

	if err != nil {
		return SMTPClient{}, fmt.Errorf("Failed to create SMTP client: %w", err)
	}

	return SMTPClient{
		config: config,
		cl:     client,
	}, nil
}

func (s *SMTPClient) SendMessage(ctx context.Context, to string, messageData usecases.Message) error {
	message := mail.NewMsg()

	if err := message.From(s.config.UserName); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}

	if err := message.To(to); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}

	message.Subject(messageData.Subject)
	message.SetBodyString(mail.TypeTextPlain, messageData.Body)

	if err := s.cl.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
