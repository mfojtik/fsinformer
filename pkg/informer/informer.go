package informer

func NewFileInformer(baseDir string, paths ...string) (FileInformer, error) {
	store := NewStore()
	if err := AddFiles(store, baseDir, paths...); err != nil {
		return nil, err
	}
	return &fsHandler{baseDir: baseDir, paths: paths, store: store}, nil
}
