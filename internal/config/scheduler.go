package config

type scheduler struct {
	Name      string `json:"name,omitempty" default:"name"`
	Namespace string `json:"namespace,omitempty" default:"namespace"`
	Algorithm string `json:"algorithm,omitempty" default:"algorithm"`
}

var Scheduler *scheduler
