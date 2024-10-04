package models

import (
	"time"

	"github.com/google/uuid"
)

type Aquarium struct {
	ID uuid.UUID `json:"id"`

	NeedApproval bool `json:"need_approval"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
