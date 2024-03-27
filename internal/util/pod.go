package util

import (
	"errors"
	"fmt"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

func FindPodByName(pods []*model.Pod, name string) (*model.Pod, error) {
	for _, pod := range pods {
		if pod.Name() == name {
			return pod, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("can't find pod with name: %s in list of pods", name))
}
