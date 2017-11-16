package api

import (
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

// FindValidAPIHealthChecks retrieves all apis that has health check configured
func (r *InMemoryRepository) FindValidAPIHealthChecks() ([]*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	var definitions []*Definition
	for _, definition := range r.definitions {
		isValid := definition.Active && definition.HealthCheck.URL != ""
		if isValid {
			definitions = append(definitions, definition)
		}
	}

	return definitions, nil
}

// FindByName find an api definition by name
func (r *InMemoryRepository) FindByName(name string) (*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	return r.findByName(name)
}

func (r *InMemoryRepository) findByName(name string) (*Definition, error) {
	definition, ok := r.definitions[name]
	if false == ok {
		return nil, ErrAPIDefinitionNotFound
	}

	return definition, nil
}

// FindByListenPath find an API definition by proxy listen path
func (r *InMemoryRepository) FindByListenPath(path string) (*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	for _, definition := range r.definitions {
		if definition.Proxy.ListenPath == path {
			return definition, nil
		}
	}

	return nil, ErrAPIDefinitionNotFound
}

// Exists searches an existing Proxy definition by its listen_path
func (r *InMemoryRepository) Exists(def *Definition) (bool, error) {
	return exists(r, def)
}

// Add adds an api definition to the repository
func (r *InMemoryRepository) Add(definition *Definition) error {
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
func (r *InMemoryRepository) Remove(name string) error {
	r.Lock()
	defer r.Unlock()

	if _, err := r.findByName(name); err != nil {
		return err
	}

	delete(r.definitions, name)

	return nil
}
