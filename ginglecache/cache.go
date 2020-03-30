package ginglecache

import (
	"ginglecache/lru"
	"sync"
)

// cache is equipped with mutex to guarantee thread safty
type cache struct {
	mux sync.Mutex
	lru *lru.Cache
	cap int64
}

// get defines how cache thread-safely get a element
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.lru == nil {
		return
	}

	if value, ok := c.lru.Get(key); ok {
		return value.(ByteView), ok
	}

	return
}

// set defines how cache thread-safely set a element
func (c *cache) set(key string, value ByteView) {
	c.mux.Lock()
	defer c.mux.Unlock()

	// Delay initailization
	if c.lru == nil {
		c.lru = lru.New(c.cap, nil)
	}

	c.lru.Set(key, value)
}
