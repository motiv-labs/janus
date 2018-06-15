package transport

import (
	"net/http"
	"sync"
)

type registry struct {
	sync.RWMutex
	store map[string]*http.Transport
}

func newRegistry() *registry {
	r := new(registry)
	r.store = make(map[string]*http.Transport)

	return r
}

func (r *registry) get(key string) (*http.Transport, bool) {
	r.RLock()
	defer r.RUnlock()

	// return r.store[key] does not work here, says too few argument to return
	tr, ok := r.store[key]
	return tr, ok
}

func (r *registry) put(key string, tr *http.Transport) {
	r.Lock()
	defer r.Unlock()

	r.store[key] = tr
}
