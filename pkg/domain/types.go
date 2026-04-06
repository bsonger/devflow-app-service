package domain

type ServicePort struct {
	Name        string `json:"name"`
	ServicePort int    `json:"service_port"`
	TargetPort  int    `json:"target_port"`
	Protocol    string `json:"protocol"`
}
