package util

import (
	"errors"
	"fmt"

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

func GetNodeResourceSum(node *model.Node) (sum *resource.Quantity, err error) {
	cpu, ok := node.Cores().AsInt64()
	if !ok {
		return nil, errors.New(fmt.Sprintf("error in converting node: %s core count to int64", node.Name()))
	}

	memory, ok := node.Memory().AsInt64()
	if !ok {
		return nil, errors.New(fmt.Sprintf("error in converting node: %s memory to int64", node.Name()))
	}

	return resource.NewQuantity(cpu+memory, node.Cores().Format), nil
}
