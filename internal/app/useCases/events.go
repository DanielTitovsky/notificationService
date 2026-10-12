package usecases

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventsService struct {
	UserId uuid.UUID
	Email  string
}

type EventsRepository struct {
	Type    string
	Payload json.RawMessage
	Status  string
}
