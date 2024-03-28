package consumer

import (
	"context"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/handlers"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/scheduler"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"os"
	"os/signal"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	factory := informers.NewSharedInformerFactory(connector.C.Client(), config.Scheduler.InformerSyncPeriod)
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
		State:        C.State,
		PodScheduler: scheduler.S,
	})
	if err != nil {
		log.App.WithError(err).Panic("error in registering node informer event handlers ")
	}

	go nodeInformer.Run(ctx.Done())
	go podInformer.Run(ctx.Done())

	C.State.SetIsNodesSynced(cache.WaitForCacheSync(ctx.Done(), nodeInformer.HasSynced))
	<-ctx.Done()
}
