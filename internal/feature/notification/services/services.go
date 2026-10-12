package notification_services

import (
	"context"

	usecases "github.com/DanielTitovsky/notificationService/internal/app/useCases"
)

type NotificationService struct {
	repository NotificationRepository
}

type NotificationRepository interface {
	SaveRegisterdEvent(ctx context.Context, event usecases.EventsRepository) error
}

func NewNotificationService(repository NotificationRepository) NotificationService {
	return NotificationService{
		repository: repository,
	}
}
