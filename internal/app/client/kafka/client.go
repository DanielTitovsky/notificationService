package app_client_kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	app_domain "github.com/DanielTitovsky/notificationService/internal/app/domain"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaClient struct {
	Config
	cl *kgo.Client
}

func NewKafkaClient(config Config) (KafkaClient, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(config.addr),
		kgo.ConsumerGroup(config.consumerGroup),
		kgo.ConsumeTopics(config.consumerTopics),
	)

	if err != nil {
		return KafkaClient{}, fmt.Errorf("Failed to connect kafka: %w", err)
	}

	return KafkaClient{
		Config: config,
		cl:     cl,
	}, nil
}

type RegisterEventService func(ctx context.Context, events app_domain.RegisterEvent) error

func (kc *KafkaClient) ListenAuthEvents(ctx context.Context, eventService RegisterEventService) {
	for ctx.Err() == nil {
		fetches := kc.cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}

		eventsData := make([]app_domain.RegisterEvent, 0, 5)

		iter := fetches.RecordIter()

		for !iter.Done() {
			record := iter.Next()

			var event app_domain.RegisterEvent

			err := json.Unmarshal(record.Value, &event)

			if err != nil {
				eventsData = append(eventsData, app_domain.RegisterEvent{
					Id:    event.Id,
					Email: event.Email,
					Error: err,
				})
			} else {
				eventsData = append(eventsData, app_domain.RegisterEvent{
					Id:    event.Id,
					Email: event.Email,
					Error: nil,
				})
			}
		}

		for _, event := range eventsData {
			if event.Error != nil {
				log.Printf("failed to parse event: %v", event.Error)
				continue
			}

			err := eventService(ctx, event)

			if err != nil {
				log.Printf("failed to notify: %v", err)
			}
		}
	}
}
