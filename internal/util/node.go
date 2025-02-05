package util

import (
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/utils/strings/slices"

	"github.com/noisyboy-9/sencillo/internal/config"
)

func IsNodeOnEdge(nodeKubernetesObject *v1.Node) bool {
	n := nodeKubernetesObject.GetName()
	return slices.Contains(config.Cluster.EdgeNodes, n)
}

func IsMasterNode(object *v1.Node) bool {
	nodeName := object.GetName()

	return strings.Trim(strings.ToLower(nodeName), " ") ==
		strings.Trim(strings.ToLower(config.Cluster.MasterNode), " ")
}
