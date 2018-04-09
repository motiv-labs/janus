package oauth2

import (
	"sync"
)

// InMemoryRepository represents a in memory repository
type InMemoryRepository struct {
	sync.RWMutex
	servers map[string]*OAuth
}

// NewInMemoryRepository creates a in memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{servers: make(map[string]*OAuth)}
}

// FindAll fetches all the OAuth Servers available
func (r *InMemoryRepository) FindAll() ([]*OAuth, error) {
	r.RLock()
	defer r.RUnlock()

	var servers []*OAuth
	for _, server := range r.servers {
		servers = append(servers, server)
	}

	return servers, nil
}

// FindByName find an OAuth Server by name
func (r *InMemoryRepository) FindByName(name string) (*OAuth, error) {
	r.RLock()
	defer r.RUnlock()

	return r.findByName(name)
}

// Add add a new OAuth Server to the repository
func (r *InMemoryRepository) Add(server *OAuth) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.servers[server.Name]; ok {
		return ErrOauthServerNameExists
	}

	r.servers[server.Name] = server

	return nil
}

// Save saves a OAuth Server to the repository
func (r *InMemoryRepository) Save(server *OAuth) error {
	r.Lock()
	defer r.Unlock()

	r.servers[server.Name] = server

	return nil
}

// Remove removes an OAuth Server from the repository
func (r *InMemoryRepository) Remove(name string) error {
	r.Lock()
	defer r.Unlock()

	if _, err := r.findByName(name); err != nil {
		return err
	}

	delete(r.servers, name)

	return nil
}

func (r *InMemoryRepository) findByName(name string) (*OAuth, error) {
	server, ok := r.servers[name]
	if false == ok {
		return nil, ErrOauthServerNotFound
	}

	return server, nil
}
