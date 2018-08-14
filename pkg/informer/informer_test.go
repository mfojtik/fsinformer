package informer

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInformerIsDirectory(t *testing.T) {
	baseDir, _ := ioutil.TempDir("", "test")
	defer os.RemoveAll(baseDir)
	_, err := NewFileInformer(4*time.Second, baseDir)
	if err != ErrIsDirectory {
		t.Fatalf("unexpected error: %v (expected ErrIsDirectory)", err)
	}
}

func TestInformerBasic(t *testing.T) {
	baseDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.RemoveAll(baseDir)
	fooFilePath := filepath.Join(baseDir, "test_foo")

	informer, err := NewFileInformer(4*time.Second, fooFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		isTestFooObserved = make(chan struct{})
		isTestFooDeleted  = make(chan struct{})
		isTestFooUpdated  = make(chan struct{})
	)

	informer.AddEventHandler(FileEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			f := obj.(File)
			defer close(isTestFooObserved)
			if f.Name() != fooFilePath {
				t.Errorf("expected 'test_foo', got: %v", f.Name())
			}
			if f.Stat() == nil {
				t.Errorf("stat for 'test_foo' should not be nil")
			}
			if string(f.Content()) != "foo" {
				t.Errorf("expected 'test_foo' with 'foo' content, got %q", string(f.Content()))
			}
		},
		UpdateFunc: func(old, obj interface{}) {
			f := obj.(File)
			oldFile := old.(File)
			defer close(isTestFooUpdated)
			if f.Name() != fooFilePath {
				t.Errorf("expected 'test_foo', got: %v", f.Name())
			}
			if f.Stat() == nil {
				t.Errorf("stat for 'test_foo' should not be nil")
			}
			if string(f.Content()) != "updated foo" {
				t.Errorf("expected 'test_foo' with 'updated foo' content, got %q", string(f.Content()))
			}
			if string(oldFile.Content()) != "foo" {
				t.Errorf("expected old file to be 'foo', got: %q", string(oldFile.Content()))
			}
		},
		DeleteFunc: func(obj interface{}) {
			f := obj.(File)
			defer close(isTestFooDeleted)
			if f.Name() != fooFilePath {
				t.Errorf("expected 'test_foo', got: %v", f.Name())
			}
			if f.Stat() == nil {
				t.Errorf("stat for 'test_foo' should not be nil")
			}
			if string(f.Content()) != "updated foo" {
				t.Errorf("expected 'test_foo' with 'updated foo' content, got %q", string(f.Content()))
			}
		},
	})

	stopCh := make(chan struct{})

	informer.Run(stopCh)

	if err := ioutil.WriteFile(fooFilePath, []byte("foo"), 0644); err != nil {
		t.Fatalf("unable to write file: %v", err)
	}

	// Wait for test-foo
	select {
	case <-isTestFooObserved:
		break
	case <-time.After(4 * time.Second):
		t.Fatalf("timeout while waiting for test foo observed")
	}

	if err := ioutil.WriteFile(fooFilePath, []byte("updated foo"), 0644); err != nil {
		t.Fatalf("unable to write file: %v", err)
	}

	// Wait for test-foo is updated
	select {
	case <-isTestFooUpdated:
		break
	case <-time.After(4 * time.Second):
		t.Fatalf("timeout while waiting for test foo to be deleted")
	}

	if err := os.Remove(fooFilePath); err != nil {
		t.Fatalf("unable to deleted test_foo: %v", err)
	}

	// Wait for test-foo is deleted
	select {
	case <-isTestFooDeleted:
		break
	case <-time.After(4 * time.Second):
		t.Fatalf("timeout while waiting for test foo to be deleted")
	}
}
