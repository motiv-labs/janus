package basic

import (
	"sync"

	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	log "github.com/sirupsen/logrus"
)

// InMemoryRepository represents a in memory repository
type InMemoryRepository struct {
	sync.RWMutex
	users map[string]*User
	hash  encrypt.Hash
}

// NewInMemoryRepository creates a in memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{users: make(map[string]*User)}
}

// FindAll fetches all the users available
func (r *InMemoryRepository) FindAll() ([]*User, error) {
	r.RLock()
	defer r.RUnlock()

	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// FindByUsername find an user by username
func (r *InMemoryRepository) FindByUsername(username string) (*User, error) {
	r.RLock()
	defer r.RUnlock()
	return r.findByUsername(username)
}

// Add adds an user to the repository
func (r *InMemoryRepository) Add(user *User) error {
	r.Lock()
	defer r.Unlock()

	hash, err := r.hash.Generate(user.Password)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return err
	}
	user.Password = hash

	r.users[user.Username] = user

	return nil
}

// Remove removes an user from the repository
func (r *InMemoryRepository) Remove(username string) error {
	r.Lock()
	defer r.Unlock()

	if _, err := r.findByUsername(username); err == ErrUserNotFound {
		return err
	}

	delete(r.users, username)

	return nil
}

func (r *InMemoryRepository) findByUsername(username string) (*User, error) {
	user, ok := r.users[username]
	if false == ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}
