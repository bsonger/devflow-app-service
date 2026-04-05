package model

type Environment struct {
	BaseModel

	Key       string            `json:"key" db:"key"`
	Name      string            `json:"name" db:"name"`
	Cluster   string            `json:"cluster" db:"cluster"`
	Namespace string            `json:"namespace" db:"namespace"`
	Labels    map[string]string `json:"labels,omitempty" db:"labels"`
}

func (Environment) CollectionName() string { return "environments" }
