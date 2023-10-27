package util

import (
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
