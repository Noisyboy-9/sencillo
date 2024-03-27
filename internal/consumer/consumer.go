package consumer

import (
	"context"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/handlers"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"time"
)

type Consumer struct {
	State                 *model.ClusterState
	NodeHandlerRegisterer cache.ResourceEventHandlerRegistration
	PodHandlerRegisterer  cache.ResourceEventHandlerRegistration
}

var C *Consumer

func Start() {
	C = new(Consumer)
	C.State = model.NewClusterState()

	nodesResource := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "nodes"}
	podsResource := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "pods"}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(connector.C.DynamicConfig(), time.Minute, config.Scheduler.Namespace, nil)
	nodeInformer := factory.ForResource(nodesResource).Informer()
	podInformer := factory.ForResource(podsResource).Informer()

	var err error
	C.NodeHandlerRegisterer, err = nodeInformer.AddEventHandler(handlers.NodeEventHandler{})
	if err != nil {
		log.App.WithError(err).Panic("error in registering node informer event handlers ")
	}

	C.PodHandlerRegisterer, err = podInformer.AddEventHandler(handlers.PodEventHandler{})
	if err != nil {
		log.App.WithError(err).Panic("error in registering node informer event handlers ")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go nodeInformer.Run(ctx.Done())
	go podInformer.Run(ctx.Done())

	isNodesSynced := cache.WaitForCacheSync(ctx.Done(), nodeInformer.HasSynced)
	isPodsSynced := cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

	C.State.SetIsPodsSynced(isPodsSynced)
	C.State.SetIsNodesSynced(isNodesSynced)
	<-ctx.Done()
}
