package ginglecache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mux    sync.RWMutex              // TODO: 全局读写锁
	groups = make(map[string]*Group) // TODO: 全局命名空间表
)

func NewGroup(name string, cap int64, getter Getter) *Group {
	if getter == nil {
		panic("Group cannot have nil getter")
	}

	mux.Lock()
	defer mux.Unlock()

	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cap: cap},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mux.RLock()
	defer mux.RUnlock()

	group := groups[name]
	return group
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// TODO: 命中，从缓存中取出数据
	if value, ok := g.mainCache.get(key); ok {
		log.Println("[GingleCache] hit ")
		return value, nil
	}

	// TODO: 错失，使用策略取出数据
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		bytes: cloneByteView(bytes),
	}

	g.populate(key, value) // TODO: 将数据写回缓存中

	return value, nil
}

func (g *Group) populate(key string, value ByteView) {
	g.mainCache.set(key, value)
}
