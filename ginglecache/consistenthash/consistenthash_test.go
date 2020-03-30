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

	// Nodes: 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Set("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for key, value := range testCases {
		if hash.Get(key) != value {
			t.Errorf("Asking for %s, should have yielded %s", key, value)
		}
	}

	// Nodes: 2, 4, 6, 8, 12, 14, 16, 18, 22, 24, 26, 28
	hash.Set("8")

	testCases["27"] = "8" // partially invalid not globally invalid

	for key, value := range testCases {
		if hash.Get(key) != value {
			t.Errorf("Asking for %s, should have yielded %s", key, value)
		}
	}
}
