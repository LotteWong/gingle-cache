package lru

import "container/list"

// Value is the general type of value
type Value interface {
	Len() int
}

// Cache is the essential data structure of LRU
type Cache struct {
	maxBytes  int64                    // capacity of cache
	curBytes  int64                    // length of cache
	dll       *list.List               // double linked list relating key and value
	dict      map[string]*list.Element // hash table relating key and list node
	OnEvicted func(entry)              // callback when evicted
}

// entry is the element of double linked list
type entry struct {
	key   string
	value Value
}

// New returns a instance of Cache
func New(maxBytes int64, onEvicted func(entry)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		curBytes:  0,
		dll:       list.New(),
		dict:      make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Len makes Cache implement the Value interface
func (c *Cache) Len() int {
	return c.dll.Len()
}

// Get defines how Cache get a element
func (c *Cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.dict[key]; ok {
		kvpair := elem.Value.(*entry) // double linked list: get value on key
		c.dll.MoveToFront(elem)       // double linked list: move the element to tail
		return kvpair.value, true
	}

	return
}

// Set defines how Cache set a element
func (c *Cache) Set(key string, value Value) {
	elem, ok := c.dict[key]

	if ok { // If found it, modify the element
		kvpair := elem.Value.(*entry)
		kvpair.value = value                                         // double linked list: set value on key
		c.dll.MoveToFront(elem)                                      // double linked list: move the element to tail
		c.curBytes += int64(value.Len()) - int64(kvpair.value.Len()) // update the length of value
	} else { // If not found, insert the element
		kvpair := &entry{key: key, value: value}           // double linked list: new an entry
		elem = c.dll.PushFront(kvpair)                     // double linked list: push the element to tail
		c.dict[key] = elem                                 // hash table: append a new relation
		c.curBytes += int64(len(key)) + int64(value.Len()) // update the length of key and value
	}

	for c.maxBytes != 0 && c.curBytes > c.maxBytes {
		c.Remove() // check replacement strategy
	}
}

// Remove defines how Cache replace expired elements
func (c *Cache) Remove() {
	elem := c.dll.Back() // double linked list: get the element of front
	if elem != nil {
		kvpair := elem.Value.(*entry)
		c.dll.Remove(elem)                                               // double linked list: remove the elelment of front
		delete(c.dict, kvpair.key)                                       // hash table: delete the expired relation
		c.curBytes -= int64(len(kvpair.key)) + int64(kvpair.value.Len()) // update the length of key and value
		if c.OnEvicted != nil {                                          // call the evicted callback
			c.OnEvicted(entry{
				key:   kvpair.key,
				value: kvpair.value,
			})
		}
	}
}
