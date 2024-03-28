package config

import "time"

type scheduler struct {
	Name                string        `json:"name"`
	InformerSyncPeriod  time.Duration `json:"informerSyncPeriod"`
	Namespace           string        `json:"namespace"`
	Algorithm           string        `json:"algorithm"`
	EdgeAnnotationKey   string        `json:"edgeAnnotationKey"`
	EdgeAnnotationValue string        `json:"edgeAnnotationValue"`
}

var Scheduler *scheduler
