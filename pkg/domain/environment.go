package domain

type Environment struct {
	BaseModel

	Name        string            `json:"name" db:"name"`
	Cluster     string            `json:"cluster" db:"cluster"`
	Namespace   string            `json:"namespace" db:"namespace"`
	Description string            `json:"description,omitempty" db:"description"`
	Labels      map[string]string `json:"labels,omitempty" db:"labels"`
}

func (Environment) CollectionName() string { return "environments" }
