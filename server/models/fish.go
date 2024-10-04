package models

import (
	"time"

	"github.com/google/uuid"
)

type Fish struct {
	ID         uuid.UUID `json:"id"`
	AquariumID uuid.UUID `json:"aquarium_id"`
	Filename   string    `json:"filename"`
	Name       string    `json:"name"`
	Approved   bool      `json:"approved"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ApprovedAt *time.Time `json:"approved_at"`
}
