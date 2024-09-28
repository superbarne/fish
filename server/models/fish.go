package models

import "github.com/google/uuid"

type Fish struct {
	ID         uuid.UUID `json:"id"`
	Filename   string    `json:"filename"`
	Name       string    `json:"name"`
	UploadTime string    `json:"upload_time"`
	Approved   bool      `json:"approved"`
}
