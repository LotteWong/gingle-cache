package main

import (
	"fmt"
	"ginglecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	_ = ginglecache.NewGroup("scores", 2<<10, ginglecache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if value, ok := db[key]; ok {
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	addr := "localhost:8080"
	pool := ginglecache.NewHTTPPool(addr)
	log.Println("ginglecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, pool))
}
