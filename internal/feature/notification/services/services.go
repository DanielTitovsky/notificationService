package notification_services

import (
	"context"

	usecases "github.com/DanielTitovsky/notificationService/internal/app/useCases"
)

type NotificationService struct {
	repository NotificationRepository
	SMTPClient NotificationSMTPClient
}

type NotificationRepository interface {
	SaveRegisterdEvent(ctx context.Context, event usecases.EventsRepository) error
}

type NotificationSMTPClient interface {
	SendMessage(ctx context.Context, to string, messageData usecases.Message) error
}

func NewNotificationService(repository NotificationRepository, SMTPClient NotificationSMTPClient) NotificationService {
	return NotificationService{
		repository: repository,
		SMTPClient: SMTPClient,
	}
}
