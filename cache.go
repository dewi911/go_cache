package cache

import "sync"

type Cache struct {
	storage map[string]interface{}
	mu      sync.RWMutex
}

func New() *Cache {
	s := make(map[string]interface{})
	return &Cache{
		storage: s,
	}
}

func (c *Cache) Set(k string, v interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[k] = v
}

func (c *Cache) Del(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.storage, k)

}

func (c *Cache) Get(k string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.storage[k]
}
