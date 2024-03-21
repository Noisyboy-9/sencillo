package util

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

func GetNodeResourceSum(node *model.Node) *resource.Quantity {
	sum := new(resource.Quantity)
	sum.Add(*node.Cores())
	sum.Add(*node.Memory())
	return sum
}
