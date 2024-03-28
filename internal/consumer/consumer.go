package consumer

import (
	"context"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/handlers"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
	"time"
)

type consumer struct {
	State                 *model.ClusterState
	NodeHandlerRegisterer cache.ResourceEventHandlerRegistration
	PodHandlerRegisterer  cache.ResourceEventHandlerRegistration
}

var C *consumer

func Start() {
	C = new(consumer)
	C.State = model.NewClusterState()

	factory := informers.NewSharedInformerFactory(connector.C.Client(), time.Minute)
	nodeInformer := factory.Core().V1().Nodes().Informer()
	podInformer := factory.Core().V1().Pods().Informer()

	var err error
	C.NodeHandlerRegisterer, err = nodeInformer.AddEventHandler(handlers.NodeEventHandler{
		State: C.State,
	})
	if err != nil {
		log.App.WithError(err).Panic("error in registering node informer event handlers ")
	}

	C.PodHandlerRegisterer, err = podInformer.AddEventHandler(handlers.PodEventHandler{
		State: C.State,
	})
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
