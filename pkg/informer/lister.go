package informer

import "os"

// AddFiles add all on-disk files into store
func AddFiles(store Store, baseDir string, paths ...string) error {
	for _, path := range paths {
		f, err := NewFile(baseDir, path)
		// Files might exist later, lets ignore non-existing files for now
		if err != nil && err != os.ErrNotExist {
			return err
		}
		if err := store.Add(f); err != nil {
			return err
		}
	}
	return nil
}
