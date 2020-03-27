package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	curBytes  int64
	dll       *list.List
	dict      map[string]*list.Element
	OnEvicted func(entry)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(entry)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		dll:       list.New(),
		dict:      make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Len() int {
	return c.dll.Len()
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.dict[key]; ok {
		kvpair := elem.Value.(*entry) // TODO: 返回查找的值
		c.dll.MoveToFront(elem)       // TODO: 移动到队列尾
		return kvpair.value, true
	}

	return
}

func (c *Cache) Set(key string, value Value) {
	elem, ok := c.dict[key]

	// TODO: 修改
	if ok {
		kvpair := elem.Value.(*entry)
		c.dll.MoveToFront(elem)
		c.curBytes += int64(value.Len()) - int64(kvpair.value.Len())
		kvpair.value = value
		// TODO: 新增
	} else {
		elem = c.dll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.curBytes += int64(len(key)) + int64(value.Len())
		c.dict[key] = elem
	}

	for c.maxBytes != 0 && c.curBytes > c.maxBytes {
		c.Remove()
	}
}

func (c *Cache) Remove() {
	elem := c.dll.Back()
	if elem != nil {
		kvpair := elem.Value.(*entry)
		c.dll.Remove(elem)                                               // TODO: 链表出队
		delete(c.dict, kvpair.key)                                       // TODO: 字典删除
		c.curBytes -= int64(len(kvpair.key)) + int64(kvpair.value.Len()) // TODO: 更新字节
		if c.OnEvicted != nil {                                          // TODO: 触发回调
			c.OnEvicted(entry{
				key:   kvpair.key,
				value: kvpair.value,
			})
		}
	}
}
