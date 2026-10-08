package app_transport_http_server

import (
	"net/http"

	app_middleware "github.com/DanielTitovsky/notificationService/internal/app/middleware"
)

type Route struct {
	Path      string
	Method    string
	Handler   http.HandlerFunc
	Midleware []app_middleware.Middleware
}

func NewRoute(path string, method string, handler http.HandlerFunc, middleware ...app_middleware.Middleware) Route {
	return Route{
		Path:      path,
		Method:    method,
		Handler:   handler,
		Midleware: middleware,
	}
}
