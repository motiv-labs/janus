package store

import (
	"sync"

	"github.com/ulule/limiter"
)

// InMemoryStore is the redis store.
type InMemoryStore struct {
	sync.Mutex
	data map[string]string
}

// NewInMemoryStore returns an instance of memory store.
func NewInMemoryStore() *InMemoryStore {
	return NewInMemoryStoreWithOptions()
}

// NewInMemoryStoreWithOptions returns an instance of memory store with custom options.
func NewInMemoryStoreWithOptions() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]string),
	}
}

func (s *InMemoryStore) Exists(key string) (bool, error) {
	s.Lock()
	defer s.Unlock()

	return s.exists(key)
}

func (s *InMemoryStore) Get(key string) (string, error) {
	s.Lock()
	defer s.Unlock()

	return s.get(key)
}

func (s *InMemoryStore) Set(key string, value string, expire int64) error {
	s.Lock()
	defer s.Unlock()

	return s.set(key, value, expire)
}

func (s *InMemoryStore) ToLimiterStore(prefix string) (limiter.Store, error) {
	return limiter.NewMemoryStore(), nil
}

func (s *InMemoryStore) exists(key string) (bool, error) {
	_, exists := s.data[key]
	return exists, nil
}

func (s *InMemoryStore) get(key string) (string, error) {
	return s.data[key], nil
}

func (s *InMemoryStore) set(key string, value string, expire int64) error {
	s.data[key] = value
	return nil
}
