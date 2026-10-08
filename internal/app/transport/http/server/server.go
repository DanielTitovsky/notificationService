package app_transport_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	app_middleware "github.com/DanielTitovsky/notificationService/internal/app/middleware"
)

type ServerMux struct {
	Config     Config
	ServerMux  *http.ServeMux
	Middleware []app_middleware.Middleware
}

func NewServerMux(config Config, middleware ...app_middleware.Middleware) ServerMux {
	return ServerMux{
		Config:     config,
		ServerMux:  http.NewServeMux(),
		Middleware: middleware,
	}
}

func (s *ServerMux) Run(ctx context.Context) error {
	mux := app_middleware.MiddlewareChain(s.ServerMux, s.Middleware...)

	server := http.Server{
		Addr:    s.Config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		ch <- server.ListenAndServe()
	}()

	select {
	case err := <-ch:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and Serve: %w", err)
		}
	case <-ctx.Done():
		fmt.Println("shutdown initiated")
		shutgownContext, cancel := context.WithTimeout(
			context.Background(),
			s.Config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutgownContext); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown http server %w", err)
		}

		fmt.Println("shutdown completed successfully")
	}
}
