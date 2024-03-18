package enum

import "errors"

type AlgorithmName string

const (
	RandomScheduler AlgorithmName = "random"
)

func ParseAlgorithmName(name string) (AlgorithmName, error) {
	switch name {
	case string(RandomScheduler):
		return RandomScheduler, nil
	default:
		return "", errors.New("unknown scheduler algorithm")
	}
}
