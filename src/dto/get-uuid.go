package dto

import "github.com/google/uuid"

type GetUUIDDto struct {
	Uuid uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}
