package util

import (
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/utils/strings/slices"
)

var edgeNodes = []string{"uq7j5k991-01", "uq7g5w631-01", "uq7p7x251-01"}
var MasterNodeName = "uq7g5t611-01"

func IsNodeOnEdge(nodeKubernetesObject *v1.Node) bool {
	n := nodeKubernetesObject.GetName()
	return slices.Contains(edgeNodes, n)
}

func IsMasterNode(object *v1.Node) bool {
	nodeName := object.GetName()
	return strings.Trim(strings.ToLower(nodeName), " ") == strings.Trim(strings.ToLower(MasterNodeName), " ")
}
