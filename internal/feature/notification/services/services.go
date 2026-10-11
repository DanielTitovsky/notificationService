package notification_services

import (
	"context"

	app_domain "github.com/DanielTitovsky/notificationService/internal/app/domain"
)

type NotificationService struct {
	repository NotificationRepository
}

type NotificationRepository interface {
	SaveRegisterdEvent(ctx context.Context, event app_domain.RegisterEvent) error
}

func NewNotificationService(repository NotificationRepository) NotificationService {
	return NotificationService{
		repository: repository,
	}
}
