package api

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

// InMemoryRepository represents a in memory repository
type InMemoryRepository struct {
	sync.RWMutex
	definitions map[string]*Definition
}

// NewInMemoryRepository creates a in memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{definitions: make(map[string]*Definition)}
}

// Close terminates the session.  It's a runtime error to use a session
// after it has been closed.
func (r *InMemoryRepository) Close() error {
	return nil
}

// Watch watches for changes on the database
func (r *InMemoryRepository) Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged) {

}

// FindAll fetches all the api definitions available
func (r *InMemoryRepository) FindAll() ([]*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	var definitions []*Definition
	for _, definition := range r.definitions {
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func (r *InMemoryRepository) findByName(name string) (*Definition, error) {
	definition, ok := r.definitions[name]
	if false == ok {
		return nil, ErrAPIDefinitionNotFound
	}

	return definition, nil
}

// Add adds an api definition to the repository
func (r *InMemoryRepository) add(definition *Definition) error {
	r.Lock()
	defer r.Unlock()

	isValid, err := definition.Validate()
	if false == isValid && err != nil {
		log.WithError(err).Error("Validation errors")
		return err
	}

	r.definitions[definition.Name] = definition

	return nil
}

// Remove removes an api definition from the repository
func (r *InMemoryRepository) remove(name string) error {
	r.Lock()
	defer r.Unlock()

	if _, err := r.findByName(name); err != nil {
		return err
	}

	delete(r.definitions, name)

	return nil
}
