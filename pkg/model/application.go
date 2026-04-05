package model

import "github.com/google/uuid"

type Application struct {
	BaseModel

	ProjectID        uuid.UUID   `json:"project_id" db:"project_id"`
	Name             string      `json:"name" db:"name"`
	RepoAddress      string      `json:"repo_address" db:"repo_address"`
	ActiveManifestID *uuid.UUID  `json:"active_manifest_id,omitempty" db:"active_manifest_id"`
	Replica          *int32      `json:"replica,omitempty" db:"replica"`
	Type             ReleaseType `json:"type" db:"type"`
	Status           string      `json:"status" db:"status"`
}

func (Application) CollectionName() string { return "applications" }
