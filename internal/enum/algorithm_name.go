package enum

type AlgorithmName string

const (
	RandomScheduler                  AlgorithmName = "random"
	SmallestFittingEdgeNodeScheduler AlgorithmName = "smallest-fitting-edge-node"
	BiggestFittingEdgeNodeScheduler  AlgorithmName = "biggest-fitting-edge-node"

	CloudFirstNodeScheduler AlgorithmName = "cloud-first"
)
