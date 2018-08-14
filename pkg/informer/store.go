package informer

import (
	"fmt"
	"sync"
)

type Store interface {
	Add(obj interface{}) error
	Update(obj interface{}) error
	Delete(obj interface{}) error
	List() []interface{}
	ListKeys() []string

	Get(obj interface{}) (item interface{}, exists bool, err error)
	GetByKey(key string) (item interface{}, exists bool, err error)
	Replace([]interface{}, string) error

	Resync() error
}

type syncMapStore struct {
	mutex sync.Mutex
	store sync.Map
}

func NewStore() Store {
	return &syncMapStore{
		store: sync.Map{},
	}
}

func (c *syncMapStore) Add(obj interface{}) error {
	f, ok := obj.(File)
	if !ok {
		return fmt.Errorf("%#+v is not a file", obj)
	}
	c.store.Store(f.Name(), f)
	return nil
}

func (c *syncMapStore) Update(obj interface{}) error {
	_, exists, err := c.Get(obj)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%#+v does not exists", obj)
	}
	return c.Add(obj)
}

func (c *syncMapStore) Delete(obj interface{}) error {
	_, exists, err := c.Get(obj)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%#+v does not exists", obj)
	}
	c.store.Delete(obj.(File).Name())
	return nil
}

func (c *syncMapStore) List() []interface{} {
	var items []interface{}
	c.store.Range(func(_, value interface{}) bool {
		items = append(items, value)
		return true
	})
	return items
}

func (c *syncMapStore) ListKeys() []string {
	var keys []string
	c.store.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func (c *syncMapStore) Get(obj interface{}) (interface{}, bool, error) {
	f, ok := obj.(File)
	if !ok {
		return nil, false, fmt.Errorf("%#+v is not a file", obj)
	}
	return c.GetByKey(f.Name())
}

func (c *syncMapStore) GetByKey(key string) (interface{}, bool, error) {
	item, exists := c.store.Load(key)
	return item, exists, nil
}

func (c *syncMapStore) Replace(items []interface{}, _ string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store = sync.Map{}
	for _, item := range items {
		if err := c.Add(item); err != nil {
			return err
		}
	}
	return nil
}

func (c *syncMapStore) Resync() error {
	return nil
}
