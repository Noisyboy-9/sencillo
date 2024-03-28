package config

type scheduler struct {
	Name                string `json:"name"`
	Namespace           string `json:"namespace"`
	Algorithm           string `json:"algorithm"`
	EdgeAnnotationKey   string `json:"edge_annotation_key"`
	EdgeAnnotationValue string `json:"edge_annotation_value"`
}

var Scheduler *scheduler
