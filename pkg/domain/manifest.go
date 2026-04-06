package domain

import "github.com/google/uuid"

type Manifest struct {
	BaseModel
	ApplicationID uuid.UUID `json:"application_id" db:"application_id"`
	Name          string    `json:"name" db:"name"`
}

func (Manifest) CollectionName() string { return "manifests" }
