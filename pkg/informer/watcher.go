package informer

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type fsHandler struct {
	handlerFuncs []FileEventHandler
	store        Store

	// mutex is needed to avoid race between relist and watcher
	mutex sync.Mutex

	// resyncPeriod is the time we should perform list of the path
	resyncPeriod time.Duration

	watcher *fsnotify.Watcher

	paths     []string
	isStarted bool
}

func (f *fsHandler) AddEventHandler(handler FileEventHandlerFuncs) {
	if f.isStarted {
		panic("cannot add handler funcs when started")
	}
	f.handlerFuncs = append(f.handlerFuncs, handler)
}

func (f *fsHandler) Run(stopCh <-chan struct{}) {
	var err error
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("unable to create new watcher: %v", err)
	}
	f.isStarted = true
	go f.runFileSystemRelist(stopCh)
	go f.runFileSystemWatch(stopCh)
}

func (f *fsHandler) HasSynced() bool {
	return f.isStarted
}

func (f *fsHandler) runFileSystemRelist(stopCh <-chan struct{}) {
	// Perform the initial sweep and register all existing files into watch
	f.relist()
	// Periodically re-list the on-disk files and store to synchronize the cache to match reality.
	ticker := time.NewTicker(f.resyncPeriod)
	for {
		select {
		case <-ticker.C:
			f.relist()
		case <-stopCh:
			ticker.Stop()
			return
		}
	}
}

func (f *fsHandler) relist() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	// Refresh the store from on-disk to match the reality
	// In case a path was specified to non-existing file, this will check if the file exists now
	// and register it into filesystem watcher (and run OnAdd() handlers).
	postAddFunc := func(item File) error {
		if err := f.watcher.Add(item.Name()); err != nil {
			return err
		}
		return nil
	}
	if err := AddFiles(f.store, postAddFunc, f.paths...); err != nil {
		log.Printf("error adding file: %v", err)
	}

	// Relist the store periodically and execute the OnAdd() handlers for all items periodically.
	var wg sync.WaitGroup
	for _, item := range f.store.List() {
		obj, err := NewFile(item.(File).Name())
		if err != nil {
			log.Printf("error creating file: %v", err)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			f.handleCreate(obj)
		}()
	}
	wg.Wait()
}

func (f *fsHandler) runFileSystemWatch(stopCh <-chan struct{}) {
	defer f.watcher.Close()
	for {
		select {
		case event := <-f.watcher.Events:
			f.mutex.Lock()
			item, err := NewFile(event.Name)
			if os.IsNotExist(err) {
				obj, exists, err := f.store.GetByKey(event.Name)
				if err != nil {
					log.Printf("unable to get %q from store: %v", event.Name, err)
					continue
				}
				if !exists {
					log.Printf("file %q does not exist in store", event.Name)
					continue
				}
				item = obj.(File)
			} else if err != nil {
				log.Printf("error gathering file information: %v", err)
				continue
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				go f.handleCreate(item)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				go f.handleWrite(item)
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				go f.handleDelete(item)
			}
			f.mutex.Unlock()
		case <-stopCh:
			break
		case err := <-f.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (f *fsHandler) handleCreate(item File) {
	if err := f.store.Add(item); err != nil {
		log.Printf("error adding %#+v to store: %v", item, err)
		return
	}
	for _, h := range f.handlerFuncs {
		h.OnAdd(item)
	}
}

func (f *fsHandler) handleWrite(item File) {
	oldItem, exists, err := f.store.Get(item)
	if err != nil {
		log.Printf("unable to get item from store: %v", err)
	}
	if !exists {
		log.Printf("item update called, but does not exists in the store: %v", err)
		return
	}
	// No content update (in some cases, the Update() is registered when the FS first create
	// the empty file and then writes the content to it. It might be specific to OSX...
	if oldItem.(File).ContentSum256() == item.ContentSum256() {
		return
	}
	if err := f.store.Update(item); err != nil {
		log.Printf("error adding %#+v to store: %v", item, err)
		return
	}
	for _, h := range f.handlerFuncs {
		h.OnUpdate(oldItem, item)
	}
}

func (f *fsHandler) handleDelete(item File) {
	if err := f.store.Delete(item); err != nil {
		log.Printf("error adding %#+v to store: %v", item, err)
		return
	}
	for _, h := range f.handlerFuncs {
		h.OnDelete(item)
	}
}
