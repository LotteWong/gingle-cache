package ginglecache

import (
	"fmt"
	"ginglecache/singleflight"
	"log"
	"sync"
)

// Group includes cache, getters, loader and identity
type Group struct {
	name string // cache content classification id

	getter LocalGetter         // helper to get data from local
	picker RemotePicker        // helper to get data from remote
	loader *singleflight.Calls // no-breakdown strategy

	mainCache cache // thread-safe cache
}

var (
	mux    sync.RWMutex              // global read-write lock
	groups = make(map[string]*Group) // global namespace map
)

// NewGroup returns a new instance of Group
func NewGroup(name string, cap int64, getter LocalGetter) *Group {
	// getter cannot be nil, picker can be nil
	if getter == nil {
		panic("Group cannot have nil getter")
	}

	// Write operation, use write lock
	mux.Lock()
	defer mux.Unlock()

	group := &Group{
		name:      name,
		getter:    getter,
		loader:    &singleflight.Calls{},
		mainCache: cache{cap: cap},
	}

	groups[name] = group
	return group
}

// GetGroup returns a existed instance of Group
func GetGroup(name string) *Group {
	// Read Operation, use read lock
	mux.RLock()
	defer mux.RUnlock()

	group := groups[name]
	return group
}

// Get defines how cache get a element with getters and singleflight
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	// If cache hit, return
	if value, ok := g.mainCache.get(key); ok {
		log.Println("[GingleCache] hit ")
		return value, nil
	}

	// If cache missed, load
	return g.load(key)
}

// load defines how to get data from local or remote
func (g *Group) load(key string) (ByteView, error) {
	// Apply singleflight strategy to avoid cache breakdown
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.picker != nil {
			if peer, ok := g.picker.PickPeer(key); ok {
				// First try to get data from remote
				value, err := g.getRemotely(peer, key)
				if err == nil {
					return value, nil
				}
				log.Println("[GingleCache] failed to get from peer:", err)
			}
		}

		// Second try to get data from local
		return g.getLocally(key)
	})

	if err != nil {
		return ByteView{}, err
	}
	return view.(ByteView), nil
}

// getLocally defines how to get data from local
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)

	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		bytes: cloneByteView(bytes),
	}

	// Write back local data to main cache
	g.populateCache(key, value)

	return value, nil
}

// getRemotely defines how to get data from remote
func (g *Group) getRemotely(getter RemoteGetter, key string) (ByteView, error) {
	bytes, err := getter.Get(g.name, key)

	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{
		bytes: cloneByteView(bytes),
	}

	return value, nil
}

// populateCache helps write back local data to main cache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.set(key, value)
}

// RegisterPicker assigns picker to group
func (g *Group) RegisterPicker(picker RemotePicker) {
	// picker can only be initialized once
	if g.picker != nil {
		panic("RegisterPeers called more than once")
	}

	g.picker = picker
}
