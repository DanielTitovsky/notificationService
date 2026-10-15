package notification_services

import (
	"context"
	"encoding/json"
	"fmt"

	usecases "github.com/DanielTitovsky/notificationService/internal/app/useCases"
)

func (s *NotificationService) SendPendingEvents(ctx context.Context, event usecases.SendRegEventService) error {
	var payload usecases.Message

	err := json.Unmarshal(event.Payload, &payload)

	if err != nil {
		return fmt.Errorf("Failed to unmarshal event payload: %w", err)
	}

	err = s.SMTPClient.SendMessage(ctx, event.Email, usecases.Message{
		Subject: payload.Subject,
		Body:    payload.Body,
	})

	if err != nil {
		fmt.Print("Failed to create SMTP client: %w", err)
	}

	return nil
}
