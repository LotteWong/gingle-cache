package consistenthash

import (
	"strconv"
	"testing"
)

func TestConsistentHash(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// TODO: 节点：2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Set("6", "4", "2")

	testCases := map[string]string{
		"2":  "2", // TODO: 正好命中真实节点
		"11": "2", // TODO: 顺序命中虚拟节点
		"23": "4", // TODO: 前后同距离顺时针
		"27": "2", // TODO: 首尾相接处顺时针
	}

	for key, value := range testCases {
		if hash.Get(key) != value {
			t.Errorf("Asking for %s, should have yielded %s", key, value)
		}
	}

	// TODO: 节点：2, 4, 6, 8, 12, 14, 16, 18, 22, 24, 26, 28
	hash.Set("8")

	testCases["27"] = "8" // TODO: 部分失效非全部失效

	for key, value := range testCases {
		if hash.Get(key) != value {
			t.Errorf("Asking for %s, should have yielded %s", key, value)
		}
	}
}
