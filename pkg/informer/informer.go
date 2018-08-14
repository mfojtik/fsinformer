package informer

import "time"

func NewFileInformer(resyncPeriod time.Duration, paths ...string) (FileInformer, error) {
	store := NewStore()
	if err := AddFiles(store, nil, paths...); err != nil {
		return nil, err
	}
	return &fsHandler{paths: paths, store: store, resyncPeriod: resyncPeriod}, nil
}
