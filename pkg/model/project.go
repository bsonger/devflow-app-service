package model

type ProjectStatus string

const (
	ProjectActive   ProjectStatus = "active"
	ProjectArchived ProjectStatus = "archived"
)

type Project struct {
	BaseModel

	Name        string            `json:"name" db:"name"`
	Key         string            `json:"key" db:"key"`
	Description string            `json:"description,omitempty" db:"description"`
	Namespace   string            `json:"namespace,omitempty" db:"namespace"`
	Owner       string            `json:"owner,omitempty" db:"owner"`
	Labels      map[string]string `json:"labels,omitempty" db:"labels"`
	Status      ProjectStatus     `json:"status" db:"status"`
}

func (p *Project) ApplyDefaults() {
	if p.Status == "" {
		p.Status = ProjectActive
	}
	if p.Namespace == "" {
		p.Namespace = p.Name
	}
}

func (Project) CollectionName() string { return "projects" }
