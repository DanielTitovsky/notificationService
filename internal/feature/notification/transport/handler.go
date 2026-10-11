package notification_transport

import (
	"context"

	app_domain "github.com/DanielTitovsky/notificationService/internal/app/domain"
	app_transport_http_server "github.com/DanielTitovsky/notificationService/internal/app/transport/http/server"
)

type NotificationHandler struct {
	service NotificationService
}

type NotificationService interface {
	RegisterEventService(ctx context.Context, events app_domain.RegisterEvent) error
}

func NewNotificationHandler(service NotificationService) NotificationHandler {
	return NotificationHandler{
		service: service,
	}
}

func (h *NotificationHandler) Register() []app_transport_http_server.Route {
	return []app_transport_http_server.Route{}
}
