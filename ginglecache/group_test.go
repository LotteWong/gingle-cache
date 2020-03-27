package ginglecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestLocalGetter(t *testing.T) {
	var f LocalGetter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if value, err := f.Get("key"); err != nil || !reflect.DeepEqual(value, expect) {
		t.Errorf("LocalGetter callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	group := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if value, ok := db[key]; ok {
			if _, ok := loadCounts[key]; ok {
				loadCounts[key]++
			}
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for key, value := range db {
		// TODO: 本地读取
		if view, err := group.Get(key); err != nil || view.String() != value {
			t.Fatalf("failed to get value of `Tom`")
		}

		// TODO: 缓存读取
		if _, err := group.Get(key); err != nil || loadCounts[key] > 1 {
			t.Fatalf("cache `%s` miss", key)
		}
	}

	// TODO: 不存在的查询
	if view, err := group.Get("unknown"); err == nil {
		t.Fatalf("the value of `unknown` should be empty, but `%s` got", view)
	}
}
