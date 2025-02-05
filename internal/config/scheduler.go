package config

import "time"

type scheduler struct {
	Name               string        `json:"name"`
	InformerSyncPeriod time.Duration `json:"informerSyncPeriod"`
	Namespace          string        `json:"namespace"`
	Algorithm          string        `json:"algorithm"`
}

var Scheduler *scheduler
