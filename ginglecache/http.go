package ginglecache

import (
	"fmt"
	"ginglecache/consistenthash"
	"log"
	"net/http"
	"strings"
	"sync"
)

// Default peer configurations
const (
	defaultBasePath = "/_ginglecache/"
	defaultReplicas = 50
)

// HTTPPool is a router with peer information
type HTTPPool struct {
	self string // such as http://localhost:8080
	base string // such as /_ginglecache/

	mux sync.Mutex // guarantee thread safty

	peers       *consistenthash.Map    // consistent hash algorithm
	httpGetters map[string]*httpGetter // map relating key to handler
}

// NewHTTPPool returns a instance of HTTPPool
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self,
		base: defaultBasePath,
	}
}

// Log helps logs and formats procedure information
func (p *HTTPPool) Log(format string, values ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, values...))
}

// ServeHTTP defines how HTTPPool handle requests and responses
func (p *HTTPPool) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// We suppose pattern is /<basepath>/<groupname>/<key>
	if !strings.HasPrefix(req.URL.Path, p.base) {
		panic("HTTPPool serving unexpected path: " + req.URL.Path)
	}

	p.Log("%s %s", req.Method, req.URL.Path)

	// We suppose pattern is /<basepath>/<groupname>/<key>
	parts := strings.SplitN(req.URL.Path[len(p.base):], "/", 2)
	if len(parts) != 2 {
		http.Error(rw, "403 Bad Request", http.StatusBadRequest)
		return
	}

	name := parts[0]
	group := GetGroup(name)
	if group == nil {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}

	key := parts[1]
	view, err := group.Get(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Write(view.Bytes())
}

// ConfPeers receive params from http request to config peers (server)
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

// PickPeer receive param from http request to pick peer (client)
func (p *HTTPPool) PickPeer(key string) (RemoteGetter, bool) {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Do not pick empty peer or self peer
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}

	return nil, false
}

// Pre-compile to check whether HTTPPool implements RemotePicker
var _ RemotePicker = (*HTTPPool)(nil)
