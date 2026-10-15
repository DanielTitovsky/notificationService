package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"

	app_client_kafka "github.com/DanielTitovsky/notificationService/internal/app/client/kafka"
	app_client_smtp "github.com/DanielTitovsky/notificationService/internal/app/client/smtp"
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

	//Регистрация SMTP Client
	SMTPConfig := app_client_smtp.NewConfig()
	SMTPClient, err := app_client_smtp.NewSMTPClient(SMTPConfig)

	if err != nil {
		log.Fatal("Failed to create SMTP client: %w", err)
	}

	//Notification
	notificationRepository := notification_postgres.NewNotificationRepository(pgPool.Pool)
	notificationService := notification_services.NewNotificationService(&notificationRepository, &SMTPClient)
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

//Продолжить с написания репозитория

// Забирать задачи из БД. Метод репозитория выбирает ограниченную порцию pending-событий и
// атомарно переводит их в processing, возвращая выбранные записи. Так два обработчика не заберут
//  одну задачу. Транзакция заканчивается до обращения к SMTP.

// Обрабатывать сохранённое событие. По полю type сервис определяет нужное действие,
//  извлекает данные из payload и формирует usecases.Message. Например, для регистрации
//  берёт email и готовит приветственное письмо.

// Отправлять порцию через ограниченный пул. SendPendingMessages распределяет задачи
//  между несколькими воркерами и ждёт их завершения. После успешной SMTP-отправки сохраняет sent.

// Запускать обработку периодически. Фоновый цикл вызывает SendPendingMessages раз
//  в N секунд. Следующий вызов начинается после завершения предыдущего; при отмене
// контекста цикл останавливается.

// Определить поведение при сбоях. Временную ошибку SMTP нужно повторять
//  с задержкой, окончательную — переводить в failed. Ещё потребуется возвращать
// в обработку зависшие processing-задачи после падения приложения. Для этого одного
//  status недостаточно: понадобятся данные о попытках и времени начала обработки.
