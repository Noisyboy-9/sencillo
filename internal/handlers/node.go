package handlers

type NodeEventHandler struct {
}

func (n NodeEventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	//TODO implement me
	panic("implement me")
}

func (n NodeEventHandler) OnUpdate(oldObj, newObj interface{}) {
	//TODO implement me
	panic("implement me")
}

func (n NodeEventHandler) OnDelete(obj interface{}) {
	//TODO implement me
	panic("implement me")
}
