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

	baseDir   string
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
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("unable to create new watcher: %v", err)
	}
	defer watcher.Close()
	for _, name := range f.paths {
		if err := watcher.Add(name); err != nil {
			log.Fatalf("unable to add %s: %v", name, err)
		}
	}
	f.isStarted = true
	go f.runFileSystemWatch(watcher, stopCh)
}

func (f *fsHandler) HasSynced() bool {
	return true
}

func (f *fsHandler) runFileSystemRelist(stopCh <-chan struct{}) {
	// Perform initial list
	f.relist()
	// Update the store periodically from disk
	ticker := time.NewTicker(f.resyncPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				f.relist()
			case <-stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

func (f *fsHandler) relist() {
	for _, item := range f.store.List() {
		f.mutex.Lock()
		obj, err := NewFile(f.baseDir, item.(File).Name())
		if err != nil && err == os.ErrNotExist {
			go f.handleDelete(obj)
			continue
		}
		go f.handleCreate(obj)
		f.mutex.Unlock()
	}
}

func (f *fsHandler) runFileSystemWatch(watcher *fsnotify.Watcher, stopCh <-chan struct{}) {
	for {
		select {
		case event := <-watcher.Events:
			f.mutex.Lock()
			item, err := NewFile(f.baseDir, event.Name)
			if err != nil {
				log.Println("error:", err)
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
		case err := <-watcher.Errors:
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
