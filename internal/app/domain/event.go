package app_domain

import "github.com/google/uuid"

type RegisterEvent struct {
	UserId uuid.UUID `json:"id"`
	Email  string    `json:"email" validate:"required,email" db:"email"`
}
