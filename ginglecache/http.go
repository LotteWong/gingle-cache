package ginglecache

import (
	"fmt"
	"ginglecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_ginglecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self        string
	base        string
	mux         sync.Mutex // TODO: 分布式要，本地不要
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		base: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, values ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, values...))
}

func (p *HTTPPool) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, p.base) {
		panic("HTTPPool serving unexpected path: " + req.URL.Path)
	}

	p.Log("%s %s", req.Method, req.URL.Path)

	// TODO: 约定/<basepath>/<groupname>/<key>
	parts := strings.SplitN(req.URL.Path[len(p.base):], "/", 2)
	if len(parts) != 2 {
		http.Error(rw, "403 Bad Request", http.StatusBadRequest)
		return
	}

	name := parts[0]
	key := parts[1]

	group := GetGroup(name)
	if group == nil {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Write(view.Bytes())
}

type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

func (p *HTTPPool) ConfPeers(peers ...string) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Set(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			base: peer + p.base,
		}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)      // TODO: 用一致性哈希算节点的string，用string包装链接去发起请求充当对应客户端
		return p.httpGetters[peer], true // TODO: 每个 HTTPPool 既是服务器又是客户端
	}

	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

type httpGetter struct {
	base string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	url := fmt.Sprintf(
		"%v%v/%v",
		h.base,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO: 状态
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	// TODO: 内容
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
