package ginglecache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// LocalGetter helps get data from local when cache misses (interface)
type LocalGetter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc helps get data from local when cache misses (function)
type GetterFunc func(key string) ([]byte, error)

// Get makes GetterFunc implement LocalGetter
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// RemotePicker picks specific server by using consistent hash (interface)
type RemotePicker interface {
	PickPeer(key string) (RemoteGetter, bool)
}

// RemoteGetter helps get data from remote when cache misses (interface)
type RemoteGetter interface {
	Get(group string, key string) ([]byte, error)
}

// httpGetter helps get data from remote when cache misses (struct)
type httpGetter struct {
	base string
}

// Get makes httpGetter implement RemoteGetter
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// Safely construct query pattern
	url := fmt.Sprintf(
		"%v%v/%v",
		h.base,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	// Client (self) sends request to server (peer)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
