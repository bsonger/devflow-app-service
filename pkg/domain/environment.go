package domain

type Environment struct {
	BaseModel

	Name        string      `json:"name" db:"name"`
	Cluster     string      `json:"cluster" db:"cluster"`
	Description string      `json:"description,omitempty" db:"description"`
	Labels      []LabelItem `json:"labels,omitempty" db:"labels"`
}

func (Environment) CollectionName() string { return "environments" }
