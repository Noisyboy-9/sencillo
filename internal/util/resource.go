package util

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/noisyboy-9/sencillo/internal/model"
)

func RequiredCpuSum(containers []v1.Container) *resource.Quantity {
	result := resource.NewQuantity(0, resource.DecimalExponent)
	for _, container := range containers {
		result.Add(*container.Resources.Limits.Cpu())
	}
	return result
}
func RequiredMemorySum(containers []v1.Container) *resource.Quantity {
	result := resource.NewQuantity(0, resource.DecimalExponent)
	for _, container := range containers {
		result.Add(*container.Resources.Limits.Memory())
	}
	return result
}

func GetNodeResourceSum(node model.Node) *resource.Quantity {
	sum := new(resource.Quantity)
	sum.Add(*node.RemainingCores)
	sum.Add(*node.RemainingMemory)
	return sum
}

func FindSmallestNode(nodes []model.Node) model.Node {
	smallestNode := nodes[0]
	smallestNodeResources := GetNodeResourceSum(smallestNode)

	for _, node := range nodes {
		resourceSum := GetNodeResourceSum(node)

		if smallestNodeResources.Cmp(*resourceSum) == 1 {
			smallestNode = node
			smallestNodeResources = resourceSum
		}
	}

	return smallestNode
}

func FindLargestNode(nodes []model.Node) model.Node {
	biggestNode := nodes[0]
	biggestNodeResources := GetNodeResourceSum(biggestNode)

	for _, node := range nodes {
		resourceSum := GetNodeResourceSum(node)

		if biggestNodeResources.Cmp(*resourceSum) == -1 {
			biggestNode = node
			biggestNodeResources = resourceSum
		}
	}

	return biggestNode
}
