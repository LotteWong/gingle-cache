package main

import (
	"flag"
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

func createGroup() *ginglecache.Group {
	return ginglecache.NewGroup("scores", 2<<10, ginglecache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if value, ok := db[key]; ok {
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, group *ginglecache.Group) {
	picker := ginglecache.NewHTTPPool(addr)
	picker.ConfPeers(addrs...)
	group.RegisterPicker(picker)

	log.Println("ginglecache server is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], picker))
}

func startCacheClient(addr string, group *ginglecache.Group) {
	http.Handle("/api", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		key := req.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/octet-stream")
		rw.Write(view.Bytes())
	}))

	log.Println("ginglecache client is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nil))
}

func main() {
	var server int
	var client bool
	flag.IntVar(&server, "server", 8000, "ginglecache server port")
	flag.BoolVar(&client, "client", false, "whether start a client or not")
	flag.Parse()

	clientAddr := "http://localhost:8080"
	serverAddrs := map[int]string{
		8000: "http://localhost:8000",
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, addr := range serverAddrs {
		addrs = append(addrs, addr)
	}

	group := createGroup()
	if client {
		go startCacheClient(clientAddr, group)
	}
	startCacheServer(serverAddrs[server], []string(addrs), group)
}
