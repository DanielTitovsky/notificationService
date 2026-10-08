package app_transport_http_server

import (
	"net/http"

	app_middleware "github.com/DanielTitovsky/notificationService/internal/app/middleware"
)

type Router struct {
	Version   string
	ServerMux *http.ServeMux
}

func NewRouter(version string, server *http.ServeMux) Router {
	return Router{
		Version:   version,
		ServerMux: server,
	}
}

func (r *Router) RegisterRouter(routers ...Route) {
	for _, route := range routers {
		handler := http.Handler(route.Handler)

		handler = app_middleware.MiddlewareChain(handler, route.Midleware...)

		fullPath := route.Method + " /api/" + r.Version + route.Path

		r.ServerMux.Handle(fullPath, handler)
	}
}
