package notification_postgres

import (
	"context"
	"fmt"
	"time"

	usecases "github.com/DanielTitovsky/notificationService/internal/app/useCases"
)

func (r *NotificationRepository) SaveRegisterdEvent(ctx context.Context, event usecases.EventsRepository) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	sql := `
		INSERT INTO events (type, payload,status)
		VALUES ($1, $2, $3)
	`

	_, err := r.pool.Exec(ctx, sql,
		event.Type,
		event.Payload,
		event.Status,
	)

	if err != nil {
		return fmt.Errorf("Failed to save event: %w", err)
	}

	return nil
}
