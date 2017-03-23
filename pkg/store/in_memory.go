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

// Exists checks if a key exists in the store
func (s *InMemoryStore) Exists(key string) (bool, error) {
	s.Lock()
	defer s.Unlock()

	return s.exists(key)
}

// Get retreives a value from the store
func (s *InMemoryStore) Get(key string) (string, error) {
	s.Lock()
	defer s.Unlock()

	return s.get(key)
}

// Remove a value from the store
func (s *InMemoryStore) Remove(key string) error {
	s.Lock()
	defer s.Unlock()

	return s.remove(key)
}

// Set a value in the store
func (s *InMemoryStore) Set(key string, value string, expire int64) error {
	s.Lock()
	defer s.Unlock()

	return s.set(key, value, expire)
}

// ToLimiterStore converts a storage into a limmiter compliant storage
func (s *InMemoryStore) ToLimiterStore(prefix string) (limiter.Store, error) {
	return limiter.NewMemoryStore(), nil
}

func (s *InMemoryStore) exists(key string) (bool, error) {
	_, exists := s.data[key]
	return exists, nil
}

func (s *InMemoryStore) remove(key string) error {
	delete(s.data, key)
	return nil
}

func (s *InMemoryStore) get(key string) (string, error) {
	return s.data[key], nil
}

func (s *InMemoryStore) set(key string, value string, expire int64) error {
	s.data[key] = value
	return nil
}
