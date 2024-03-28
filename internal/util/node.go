package util

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"k8s.io/api/core/v1"
)

func IsNodeOnEdge(nodeKubernetesObject *v1.Node) bool {
	return nodeKubernetesObject.GetAnnotations()[config.Scheduler.EdgeAnnotationKey] == config.Scheduler.EdgeAnnotationValue
}
