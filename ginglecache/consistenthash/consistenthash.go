package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	keys []int // TODO: 真实节点
	dict map[int]string // TODO: 虚拟节点
	replicas int	
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash: hash,
		dict: make(map[int]string),
		replicas: replicas,
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return	m
}

// TODO: 此处的key为server
func (m *Map) Set(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.dict[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// TODO: 此处的key为client
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.dict[m.keys[idx%len(m.keys)]] // TODO: 环状结构要取余数
}
