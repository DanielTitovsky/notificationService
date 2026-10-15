package notification_postgres

import (
	"context"
	"fmt"
	"time"
)

func (r *NotificationRepository) LockPendingRegEvents(ctx context.Context, limit int) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)

	defer cancel()

	sql := `
	WITH selected_events AS (
    	SELECT id
    	FROM public.events
    	WHERE status = 'pending'
    	ORDER BY created_at, id
    	LIMIT $1
    	FOR UPDATE SKIP LOCKED
	)
	UPDATE public.events AS e
	SET status = 'processing'
	FROM selected_events AS selected
	WHERE e.id = selected.id
	RETURNING
    	e.id,
    	e."type",
    	e.payload,
    	e.status,
    	e.created_at;
	`

	rows, err := r.pool.Query(ctx, sql, limit)

	if err != nil {
		return fmt.Errorf("Failed to save event: %w", err)
	}

	return nil
}
