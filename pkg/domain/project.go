package domain

type Project struct {
	BaseModel

	Name        string            `json:"name" db:"name"`
	Description string            `json:"description,omitempty" db:"description"`
	Labels      map[string]string `json:"labels,omitempty" db:"labels"`
}

func (Project) CollectionName() string { return "projects" }
