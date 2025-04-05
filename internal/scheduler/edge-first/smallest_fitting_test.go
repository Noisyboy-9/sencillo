package edge_first

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/model"
)

func TestSmallestFittingEdgeNodeScheduler_Filter(t *testing.T) {
	cases := map[string]struct {
		clusterNodes []model.Node
		targetPod    model.Pod

		expectedCloudNodes []model.Node
		expectedEdgeNodes  []model.Node
	}{
		"master is not being selected": {
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
			},

			targetPod: model.Pod{
				ID:        "1",
				Name:      "new-pod",
				Namespace: "example-namespace",
				Cores:     resource.NewQuantity(1, resource.DecimalSI),
				Memory:    resource.NewQuantity(20*1024, resource.BinarySI),
			},

			expectedCloudNodes: []model.Node{},
			expectedEdgeNodes:  []model.Node{},
		},
		"no cloud nodes in the cluster": {
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			targetPod: model.Pod{
				ID:        "1",
				Name:      "new-pod",
				Namespace: "example-namespace",
				Cores:     resource.NewQuantity(1, resource.DecimalSI),
				Memory:    resource.NewQuantity(20*1024, resource.BinarySI),
			},
			expectedCloudNodes: []model.Node{},
			expectedEdgeNodes: []model.Node{
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
		},
		"no edge nodes in the cluster": {
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
			},
			targetPod: model.Pod{
				ID:        "1",
				Name:      "new-pod",
				Namespace: "example-namespace",
				Cores:     resource.NewQuantity(1, resource.DecimalSI),
				Memory:    resource.NewQuantity(20*1024, resource.BinarySI),
			},
			expectedCloudNodes: []model.Node{
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
			},
			expectedEdgeNodes: []model.Node{},
		},
		"all cloud nodes are filtered": {
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("0"),
					RemainingMemory:   resource.MustParse("0Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			targetPod: model.Pod{
				ID:        "1",
				Name:      "new-pod",
				Namespace: "example-namespace",
				Cores:     resource.NewQuantity(1, resource.DecimalSI),
				Memory:    resource.NewQuantity(20*1024, resource.BinarySI),
			},
			expectedCloudNodes: []model.Node{},
			expectedEdgeNodes: []model.Node{
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
		},
		"all edge nodes are filtered": {
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("500m"),
					RemainingMemory:   resource.MustParse("0.5Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			targetPod: model.Pod{
				ID:        "1",
				Name:      "new-pod",
				Namespace: "example-namespace",
				Cores:     resource.NewQuantity(1, resource.DecimalSI),
				Memory:    resource.NewQuantity(20*1024, resource.BinarySI),
			},
			expectedCloudNodes: []model.Node{
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
			},
			expectedEdgeNodes: []model.Node{},
		},
		"some edge and cloud nodes remain": {
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			targetPod: model.Pod{
				ID:        "1",
				Name:      "new-pod",
				Namespace: "example-namespace",
				Cores:     resource.NewQuantity(1, resource.DecimalSI),
				Memory:    resource.NewQuantity(20*1024, resource.BinarySI),
			},
			expectedCloudNodes: []model.Node{

				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
			},
			expectedEdgeNodes: []model.Node{
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			log.App = new(logrus.Logger)

			scheduler := SmallestFittingEdgeNodeScheduler{}
			edgeNodes, cloudNodes := scheduler.Filter(tc.targetPod, tc.clusterNodes)

			assert.ElementsMatch(t, tc.expectedEdgeNodes, edgeNodes)
			assert.ElementsMatch(t, tc.expectedCloudNodes, cloudNodes)
		})
	}
}

func TestSmallestFittingEdgeNodeScheduler_Run(t *testing.T) {
	cases := map[string]struct {
		targetPod    model.Pod
		clusterNodes []model.Node

		expectedNode model.Node
		expectError  bool
	}{
		"smallest edge node is prioritized over edge nodes": {
			targetPod: model.Pod{
				Cores:  resource.NewQuantity(1, resource.DecimalSI),
				Memory: resource.NewQuantity(512, resource.BinarySI),
			},
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("3"),
					RemainingMemory:   resource.MustParse("3Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("5"),
					RemainingMemory:   resource.MustParse("5Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
				{
					ID:                "4",
					Name:              "node-4",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("10Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			expectedNode: model.Node{
				ID:                "2",
				Name:              "node-2",
				AllocatableCores:  resource.MustParse("10"),
				AllocatableMemory: resource.MustParse("100Gi"),
				RemainingCores:    resource.MustParse("3"),
				RemainingMemory:   resource.MustParse("3Gi"),
				IsOnEdge:          true,
				IsMaster:          false,
			},
			expectError: false,
		},
		"no eligible cloud nodes found": {
			targetPod: model.Pod{
				Cores:  resource.NewQuantity(1, resource.DecimalSI),
				Memory: resource.NewQuantity(20*1024, resource.BinarySI),
			},
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("0"),
					AllocatableMemory: resource.MustParse("0Gi"),
					RemainingCores:    resource.MustParse("0"),
					RemainingMemory:   resource.MustParse("0Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			expectedNode: model.Node{
				ID:                "3",
				Name:              "node-3",
				AllocatableCores:  resource.MustParse("10"),
				AllocatableMemory: resource.MustParse("100Gi"),
				RemainingCores:    resource.MustParse("10"),
				RemainingMemory:   resource.MustParse("100Gi"),
				IsOnEdge:          true,
				IsMaster:          false,
			},
			expectError: false,
		},
		"no eligible nodes are found": {
			targetPod: model.Pod{
				Cores:  resource.NewQuantity(1, resource.DecimalSI),
				Memory: resource.NewQuantity(20*1024, resource.BinarySI),
			},
			clusterNodes: []model.Node{
				{
					ID:                "1",
					Name:              "node-1",
					AllocatableCores:  resource.MustParse("10"),
					AllocatableMemory: resource.MustParse("100Gi"),
					RemainingCores:    resource.MustParse("10"),
					RemainingMemory:   resource.MustParse("100Gi"),
					IsOnEdge:          false,
					IsMaster:          true,
				},
				{
					ID:                "2",
					Name:              "node-2",
					AllocatableCores:  resource.MustParse("0"),
					AllocatableMemory: resource.MustParse("0Gi"),
					RemainingCores:    resource.MustParse("0"),
					RemainingMemory:   resource.MustParse("0Gi"),
					IsOnEdge:          false,
					IsMaster:          false,
				},
				{
					ID:                "3",
					Name:              "node-3",
					AllocatableCores:  resource.MustParse("0"),
					AllocatableMemory: resource.MustParse("0Gi"),
					RemainingCores:    resource.MustParse("0"),
					RemainingMemory:   resource.MustParse("0Gi"),
					IsOnEdge:          true,
					IsMaster:          false,
				},
			},
			expectedNode: model.Node{},
			expectError:  true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			log.App = new(logrus.Logger)

			s := SmallestFittingEdgeNodeScheduler{}
			selectedNode, err := s.Run(tc.targetPod, tc.clusterNodes)
			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedNode, selectedNode)
		})
	}
}
