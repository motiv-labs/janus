package store

import (
	"sync"

	"github.com/hellofresh/janus/pkg/notifier"
)

// InMemoryStore is the redis store.
type InMemoryStore struct {
	sync.Mutex
	data   map[string]string
	client chan struct{}
}

// NewInMemoryStore returns an instance of memory store.
func NewInMemoryStore() *InMemoryStore {
	return NewInMemoryStoreWithOptions()
}

// NewInMemoryStoreWithOptions returns an instance of memory store with custom options.
func NewInMemoryStoreWithOptions() *InMemoryStore {
	return &InMemoryStore{
		data:   make(map[string]string),
		client: make(chan struct{}),
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

// Publish publishes in memory
func (s *InMemoryStore) Publish(topic string, data []byte) error {
	s.client <- struct{}{}
	return nil
}

// Subscribe subscribes to messages in memory
func (s *InMemoryStore) Subscribe(channel string, callback func(notifier.Notification)) error {
	for range s.client {
		notification := notifier.Notification{}
		callback(notification)
	}
	return nil
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
