package model

import "github.com/google/uuid"

type ServiceResource struct {
	BaseModel

	ApplicationID uuid.UUID `json:"application_id" db:"application_id"`
	Name          string    `json:"name" db:"name"`
	Internet      Internet  `json:"internet" db:"internet"`
	Ports         []Port    `json:"ports" db:"ports"`
	Status        string    `json:"status" db:"status"`
}

func (ServiceResource) CollectionName() string { return "services" }
