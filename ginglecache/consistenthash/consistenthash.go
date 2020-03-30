package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash is custom hash function to convert data to id
type Hash func(data []byte) uint32

// Map is the essential data structure of consistent hash
type Map struct {
	hash     Hash           // hash function
	keys     []int          // hash ring storing all nodes
	dict     map[int]string // hash table mapping virtual nodes to actual nodes
	replicas int            // virtual node count for one actual node
}

// New returns a instance of Map
func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		keys:     make([]int, 0),
		dict:     make(map[int]string),
		replicas: replicas,
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // set default hash function
	}
	return m
}

// Get defines how consistent hash select corresponding server
func (m *Map) Get(key string) string {
	// key comes from client input

	if len(m.keys) == 0 {
		return ""
	}

	// hash is related to data
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash // hash ring: clockwise choose the first id greater than hash
	})

	return m.dict[m.keys[idx%len(m.keys)]] // hash dict: mod to retrieve the actual node
}

// Set defines how consistent hash append new servers
func (m *Map) Set(keys ...string) {
	// keys come from server config

	for _, key := range keys {
		// Generate virtual node id
		for i := 0; i < m.replicas; i++ {
			// hash is related to virtual node index and actual node id
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))

			m.keys = append(m.keys, hash) // hash ring: insert id into ring
			m.dict[hash] = key            // hash table: map virtual node (int) to actual node (string)
		}
	}

	// Sort the hash ring
	sort.Ints(m.keys)
}
