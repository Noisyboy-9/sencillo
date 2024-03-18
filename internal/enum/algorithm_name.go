package enum

import "errors"

type AlgorithmName string

const (
	RandomScheduler                  AlgorithmName = "random"
	SmallestFittingEdgeNodeScheduler AlgorithmName = "smallest-fitting-edge-node"
)

func ParseAlgorithmName(name string) (AlgorithmName, error) {
	switch name {
	case string(RandomScheduler):
		return RandomScheduler, nil
	case string(SmallestFittingEdgeNodeScheduler):
		return SmallestFittingEdgeNodeScheduler, nil
	default:
		return "", errors.New("unknown scheduler algorithm")
	}
}
