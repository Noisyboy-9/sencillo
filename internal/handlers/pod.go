package handlers

type PodEventHandler struct {
}

func (p PodEventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	//TODO implement me
	panic("implement me")
}

func (p PodEventHandler) OnUpdate(oldObj, newObj interface{}) {
	//TODO implement me
	panic("implement me")
}

func (p PodEventHandler) OnDelete(obj interface{}) {
	//TODO implement me
	panic("implement me")
}
