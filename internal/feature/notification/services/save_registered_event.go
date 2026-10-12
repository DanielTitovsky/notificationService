package notification_services

import (
	"context"
	"encoding/json"
	"fmt"

	app_domain "github.com/DanielTitovsky/notificationService/internal/app/domain"
	usecases "github.com/DanielTitovsky/notificationService/internal/app/useCases"
)

func (s *NotificationService) RegisterEventService(ctx context.Context, event app_domain.RegisterEvent) error {

	eventPayload, err := json.Marshal([]byte(event.Email))

	err = s.repository.SaveRegisterdEvent(ctx, usecases.EventsRepository{
		Type:    "register",
		Payload: eventPayload,
		Status:  "Pending",
	})

	if err != nil {
		return fmt.Errorf("Failed to save evetn: %w", err)
	}

	return nil
}
