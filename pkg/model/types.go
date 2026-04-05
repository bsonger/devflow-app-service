package model

type ReleaseType string
type Internet string

const (
	Normal    ReleaseType = "normal"
	Canary    ReleaseType = "canary"
	BlueGreen ReleaseType = "blue-green"
)

const (
	Internal Internet = "internal"
	External Internet = "external"
)

type Port struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	TargetPort int    `json:"target_port"`
}
