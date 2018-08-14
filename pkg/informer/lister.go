package informer

import (
	"os"

	"github.com/mfojtik/fsinformer/pkg/cache"
	"github.com/mfojtik/fsinformer/pkg/types"
)

// AddFiles add all on-disk files into store
func AddFiles(store cache.Store, postAddFunc func(item types.File) error, paths ...string) error {
	for _, path := range paths {
		f, err := types.NewFile(path)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
		if err := store.Add(f); err != nil {
			return err
		}
		if postAddFunc != nil {
			if err := postAddFunc(f); err != nil {
				return err
			}
		}
	}
	return nil
}
