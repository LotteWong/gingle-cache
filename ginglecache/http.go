package ginglecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/ginglecache/"

type HTTPPool struct {
	self string
	base string
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
