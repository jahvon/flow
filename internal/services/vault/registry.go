package vault

import "sync"

type Registry struct {
	providers map[string]Provider
	mu        sync.RWMutex
}

func (r *Registry) Register(name string, provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

func (r *Registry) Get(name string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, exists := r.providers[name]
	return provider, exists
}
