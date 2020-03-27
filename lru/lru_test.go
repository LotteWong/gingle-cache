package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGetAndSet(t *testing.T) {
	lru := New(int64(0), nil)

	lru.Set("key1", String("value1"))

	if value, ok := lru.Get("key1"); !ok || string(value.(String)) != "value1" {
		t.Fatalf("cache hit key1 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)

	lru.Set(k1, String(v1))
	lru.Set(k2, String(v2))
	lru.Set(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("LRU algorithm failed")
	}
}

func TestOnEvicted(t *testing.T) {
	expiredKeys := make([]string, 0)
	callback := func(item entry) {
		expiredKeys = append(expiredKeys, item.key)
	}

	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), callback)

	lru.Set(k1, String(v1))
	lru.Set(k2, String(v2))
	lru.Set(k3, String(v3))

	expectKeys := []string{"key1"}

	if !reflect.DeepEqual(expiredKeys, expectKeys) {
		t.Fatalf("OnEvicted callback failed")
	}
}
