package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mfojtik/fsinformer/pkg/informer"
)

func setup() []string {
	baseDir, _ := ioutil.TempDir("", "sample")

	// Write some existing files...
	ioutil.WriteFile(filepath.Join(baseDir, "existing_sample.yaml"), []byte("sample"), os.ModePerm)

	// Return file paths
	return []string{
		filepath.Join(baseDir, "existing_sample.yaml"),
		filepath.Join(baseDir, "future_sample.yaml"),
	}
}

func main() {
	paths := setup()

	i, err := informer.NewFileInformer(3*time.Second, paths...)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Register handlers:
	i.AddEventHandler(informer.FileEventHandlerFuncs{
		// AddFunc is called when the file is added to the store (observed).
		// AddFunc is also called when the resync happens (every 3 seconds in this example)
		AddFunc: func(item interface{}) {
			f := item.(informer.File)
			log.Printf("OnAdd called for %q (content: %s)", f.Name(), string(f.Content()))
		},
		// UpdateFunc is called when the file content change on disk.
		// The first argument represents old (stored) version of the file, second argument is updated file.
		UpdateFunc: func(oldItem, newItem interface{}) {
			newFile := newItem.(informer.File)
			oldFile := oldItem.(informer.File)
			log.Printf("OnUpdate called for %q (old content: %s, new content: %s)", newFile.Name(),
				string(oldFile.Content()), string(newFile.Content()))
		},
		// DeleteFunc is called when the file was removed from the filesystem.
		DeleteFunc: func(item interface{}) {
			deletedFile := item.(informer.File)
			log.Printf("OnDelete called for %q", deletedFile.Name())
		},
	})

	// Run the informer until we send it stop signal
	stopChan := make(chan struct{})
	defer close(stopChan)

	i.Run(stopChan)

	// Wait until informer is fully started
	for {
		if i.HasSynced() {
			break
		}
	}

	// Now create the second sample file. We registered it when we started the informer, so on every relist, the informer
	// will attempt to add it into store when it is available.
	// TODO: In future we can watch the directory which will avoid waiting for resync
	time.Sleep(5 * time.Second)
	log.Printf("Creating %s file ...", paths[1])
	ioutil.WriteFile(paths[1], []byte("future sample"), os.ModePerm)

	// Now update the content of the future file
	time.Sleep(5 * time.Second)
	log.Printf("Updating %s file ...", paths[1])
	ioutil.WriteFile(paths[1], []byte("updated future sample"), os.ModePerm)

	// Finally, delete it existing file
	time.Sleep(5 * time.Second)
	log.Printf("Deleting %s file ...", paths[0])
	os.Remove(paths[0])

	// "Th-Th-The, Th-Th-The, Th-Th... That's all, folks!"...
	time.Sleep(5 * time.Second)
}
