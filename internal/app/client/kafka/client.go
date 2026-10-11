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

func (kc *KafkaClient) ListenAuthEvents(ctx context.Context, eventService RegisterEventService) error {
	var errs []kgo.FetchError

	for ctx.Err() == nil {
		fetches := kc.cl.PollFetches(ctx)

		errs = fetches.Errors()

		if len(errs) > 0 {
			break
		}

		iter := fetches.RecordIter()

		for !iter.Done() {
			record := iter.Next()

			var event app_domain.RegisterEvent

			err := json.Unmarshal(record.Value, &event)

			if err != nil {
				log.Printf("failed to parse event: %v", event.Error)
			} else {
				err = eventService(ctx, event)

				if err != nil {
					log.Printf("failed to notify: %v", err)
				}
			}
		}
	}

	return fmt.Errorf("Failed to fetches event: %w", fmt.Sprint(errs))
}
