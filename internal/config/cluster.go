package config

type cluster struct {
	EdgeNodes  []string `yaml:"edgeNodes"`
	MasterNode string   `yaml:"masterNode"`
}

var Cluster *cluster
