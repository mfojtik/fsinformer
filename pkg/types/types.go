package types

type FileEventHandler interface {
	OnAdd(obj interface{})
	OnUpdate(oldObj, newObj interface{})
	OnDelete(obj interface{})
}

type FileEventHandlerFuncs struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}

func (r FileEventHandlerFuncs) OnAdd(obj interface{}) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

func (r FileEventHandlerFuncs) OnUpdate(oldObj, newObj interface{}) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

func (r FileEventHandlerFuncs) OnDelete(obj interface{}) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
}

type FileInformer interface {
	AddEventHandler(handler FileEventHandlerFuncs)
	Run(stopCh <-chan struct{})
	HasSynced() bool
}
