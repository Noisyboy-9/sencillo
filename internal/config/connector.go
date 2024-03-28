package config

type connector struct {
	Mode           string `json:"mode,omitempty" yaml:"mode"`
	MasterURL      string `json:"masterURL,omitempty" yaml:"masterURL"`
	KubeConfigPath string `json:"kubeConfigPath,omitempty" yaml:"kubeConfigPath"`
}

var Connector *connector
