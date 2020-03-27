package ginglecache

import (
	"gingle-cache/lru"
	"sync"
)

// TODO: 封装更底层的lru
// TODO: Value更加具体，是只读的ByteView保证并发性

type cache struct {
	mux sync.Mutex // TODO: 读写锁的升级
	lru *lru.Cache
	cap int64
}

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

func (c *cache) set(key string, value ByteView) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cap, nil) // TODO: 初始化封装性不统一
	}

	c.lru.Set(key, value)
}
