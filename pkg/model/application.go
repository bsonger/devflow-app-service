package model

import "github.com/google/uuid"

type Application struct {
	BaseModel

	ProjectID        uuid.UUID         `json:"project_id" db:"project_id"`
	Name             string            `json:"name" db:"name"`
	RepoAddress      string            `json:"repo_address" db:"repo_address"`
	ActiveManifestID *uuid.UUID        `json:"active_manifest_id,omitempty" db:"active_manifest_id"`
	Labels           map[string]string `json:"labels,omitempty" db:"labels"`
}

func (Application) CollectionName() string { return "applications" }
