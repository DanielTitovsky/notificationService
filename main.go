package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"

	app_client_kafka "github.com/DanielTitovsky/notificationService/internal/app/client/kafka"
	app_repository_postgres "github.com/DanielTitovsky/notificationService/internal/app/repository/postgres"
	app_transport_http_server "github.com/DanielTitovsky/notificationService/internal/app/transport/http/server"
	notification_postgres "github.com/DanielTitovsky/notificationService/internal/feature/notification/repository/postgres"
	notification_services "github.com/DanielTitovsky/notificationService/internal/feature/notification/services"
	notification_transport "github.com/DanielTitovsky/notificationService/internal/feature/notification/transport"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer cancel()

	//Создание HTTP SERVER
	httpConfig := app_transport_http_server.NewConfig()
	httpServer := app_transport_http_server.NewServerMux(httpConfig)
	router := app_transport_http_server.NewRouter("v01", httpServer.ServerMux)

	//Создание Postgres DB
	pgConfig := app_repository_postgres.NewConfig()
	pgPool := app_repository_postgres.NewPool(ctx, pgConfig)

	//Регитрация KAFKA CLIENT
	kafkaConfig := app_client_kafka.NewConfig(
		"localhost:9092",
		"auth_registered",
		"AuthUser",
	)
	kafkaClient, err := app_client_kafka.NewKafkaClient(kafkaConfig)

	if err != nil {
		log.Fatal("Failed to create kafka client: %w", err)
	}

	//Notification
	notificationRepository := notification_postgres.NewNotificationRepository(pgPool.Pool)
	notificationService := notification_services.NewNotificationService(&notificationRepository)
	notificationHandler := notification_transport.NewNotificationHandler(&notificationService)

	errCh := make(chan error, 2)

	var wg sync.WaitGroup

	wg.Add(2)

	fmt.Print("1")

	go func() {
		defer wg.Done()
		errCh <- kafkaClient.ListenAuthEvents(ctx, notificationService.RegisterEventService)
	}()

	router.RegisterRouter(notificationHandler.Register()...)

	fmt.Print("2")

	go func() {
		defer wg.Done()
		errCh <- httpServer.Run(ctx)
	}()

	fmt.Print("3")

	select {
	case err := <-errCh:
		log.Printf("Server error: %v", err)
		cancel()
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	}

	wg.Wait()
}

// Написать SMTP-обёртку и проверить отправку одного письма на собственный адрес.
// Подключить её к сервису и отправить одну задачу из БД.
// Добавить фоновый цикл, который периодически проверяет ожидающие задачи.
// Добавить ограниченные повторные попытки и окончательный статус ошибки.
