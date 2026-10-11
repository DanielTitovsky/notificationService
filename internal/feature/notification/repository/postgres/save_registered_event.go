package notification_postgres

import (
	"context"

	app_domain "github.com/DanielTitovsky/notificationService/internal/app/domain"
)

func (r *NotificationRepository) SaveRegisterdEvent(ctx context.Context, event app_domain.RegisterEvent) error {
	return nil
}
