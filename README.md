# gingle-cache

A simple groupcache-like distributed cache implemented by Golang.

---

## Features

- [x] Cache Eviction Algorithm (LRU)
- [x] Stand-alone Cache
- [x] Distributed Cache
- [x] HTTP Server and Client
- [X] Consistent Hash <u>(prevent cache avalanche)</u>
- [x] Singleflight <u>(prevent cache breakdown)</u>

## Quick Start

### StartCacheServer

```go
func startCacheServer(addr string, addrs []string, group *ginglecache.Group) {
  picker := ginglecache.NewHTTPPool(addr)
  picker.ConfPeers(addrs...)
  
  group.RegisterPicker(picker)
  
  log.Println("ginglecache server is running at", addr)
  log.Fatal(http.ListenAndServe(addr[7:], picker))
}
```

### StartCacheClient

```go
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
```

## TODOs

- [ ] Comments
- [ ] Blogs
- [ ] HTTP2 Support
- [ ] Protobuf Support
