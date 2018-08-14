package informer

import (
	"time"

	"github.com/mfojtik/fsinformer/pkg/cache"
	"github.com/mfojtik/fsinformer/pkg/types"
)

func NewFileInformer(resyncPeriod time.Duration, paths ...string) (types.FileInformer, error) {
	store := cache.NewStore()
	if err := AddFiles(store, nil, paths...); err != nil {
		return nil, err
	}
	return &fsHandler{paths: paths, store: store, resyncPeriod: resyncPeriod}, nil
}
