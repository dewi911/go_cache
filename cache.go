package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, k interface{}) (interface{}, error)
	Set(ctx context.Context, k, v interface{}, d time.Duration) error
	Del(ctx context.Context, k interface{}) error
	Flush(ctx context.Context) error
}

type cache struct {
	mu      sync.RWMutex
	values  map[interface{}]value
	ticker  *time.Ticker
	maxSize int
	metrics metrics
}

type value struct {
	data    interface{}
	expires time.Time
}

type metrics struct {
	hits   int64
	misses int64
}

type Config struct {
	CleanupInterval time.Duration
	MaxSize         int
}

func New(cfg Config) Cache {
	c := &cache{
		values:  make(map[interface{}]value),
		ticker:  time.NewTicker(cfg.CleanupInterval),
		maxSize: cfg.MaxSize,
	}
	go c.cleanupLoop()
	return c
}

func (c *cache) cleanupLoop() {
	for _ = range c.ticker.C {
		c.cleanup()
	}
}

func (c *cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, v := range c.values {
		if !v.expires.IsZero() && now.After(v.expires) {
			delete(c.values, k)
		}
	}
}

func (c *cache) Get(ctx context.Context, k interface{}) (interface{}, error) {
	c.mu.RLock()
	v, ok := c.values[k]
	c.mu.RUnlock()

	if !ok {
		c.metrics.misses++
		return nil, fmt.Errorf("key not found")
	}

	if !v.expires.IsZero() && time.Now().After(v.expires) {
		c.mu.Lock()
		delete(c.values, k)
		c.mu.Unlock()
		c.metrics.misses++
		return nil, fmt.Errorf("key expired")
	}

	c.metrics.hits++
	return v.data, nil
}

func (c *cache) Set(ctx context.Context, k, v interface{}, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.maxSize > 0 && len(c.values) >= c.maxSize {
		return fmt.Errorf("cache is full")
	}

	var expires time.Time
	if d > 0 {
		expires = time.Now().Add(d)
	}

	c.values[k] = value{
		data:    v,
		expires: expires,
	}

	return nil
}

func (c *cache) Del(ctx context.Context, k interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.values[k]; ok {
		return fmt.Errorf("key not found")
	}

	delete(c.values, k)
	return nil
}

func (c *cache) Flush(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values = make(map[interface{}]value)
	return nil
}
