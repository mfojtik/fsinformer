package informer

import (
	"os"
)

// AddFiles add all on-disk files into store
func AddFiles(store Store, postAddFunc func(item File) error, paths ...string) error {
	for _, path := range paths {
		f, err := NewFile(path)
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
