你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
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

- [ ] HTTP2 Support
- [ ] Protobuf Support
- [ ] More Replacement Strategy
- [ ] Distributed Unique ID Generation
- [ ] Distributed Lock Manager
- [ ] Goroutine and Connect Pool
