package guds

import (
	"sync"
)

type Router struct {
	sync.RWMutex
	Callbacks map[string]func(msg Data, c chan Reply)
}

func (h *Router) Register(api string, callback func(msg Data, c chan Reply)) {
	h.Lock()
	defer h.Unlock()
	h.Callbacks[api] = callback
}

func (h *Router) Get(api string) func(msg Data, c chan Reply) {
	h.RLock()
	defer h.RUnlock()
	return h.Callbacks[api]
}

func (h *Router) Remove(api string) {
	h.Lock()
	defer h.Unlock()
	delete(h.Callbacks, api)
}

func (h *Router) RemoveAll() {
	h.Callbacks = map[string]func(msg Data, c chan Reply){}
}
