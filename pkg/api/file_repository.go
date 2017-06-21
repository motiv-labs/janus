package api

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"sync"

	log "github.com/sirupsen/logrus"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	sync.RWMutex
	definitions map[string]*Definition
}

// NewFileSystemRepository creates a mongo country repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := &FileSystemRepository{definitions: make(map[string]*Definition)}

	// Grab json files from directory
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		return nil, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			filePath := filepath.Join(dir, f.Name())
			appConfigBody, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.WithError(err).WithField("path", filePath).Error("Couldn't load the api definition file")
				return nil, err
			}

			definition := repo.parseDefinition(appConfigBody)
			if err = repo.Add(definition); err != nil {
				log.WithError(err).Error("Can't add the definition to the repository")
				return nil, err
			}
		}
	}

	return repo, nil
}

// FindAll fetches all the api definitions available
func (r *FileSystemRepository) FindAll() ([]*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	var definitions []*Definition
	for _, definition := range r.definitions {
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// FindValidAPIHealthChecks retrieves all apis that has health check configured
func (r *FileSystemRepository) FindValidAPIHealthChecks() ([]*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	var definitions []*Definition
	for _, definition := range r.definitions {
		if definition.HealthCheck.URL != "" {
			definitions = append(definitions, definition)
		}
	}

	return definitions, nil
}

// FindByName find an api definition by name
func (r *FileSystemRepository) FindByName(name string) (*Definition, error) {
	r.RLock()
	defer r.RUnlock()

	return r.findByName(name)
}

func (r *FileSystemRepository) findByName(name string) (*Definition, error) {
	definition, ok := r.definitions[name]
	if false == ok {
		return nil, ErrAPIDefinitionNotFound
	}

	return definition, nil
}

// FindByListenPath find an API definition by proxy listen path
func (r *FileSystemRepository) FindByListenPath(path string) (*Definition, error) {
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
func (r *FileSystemRepository) Exists(def *Definition) (bool, error) {
	return exists(r, def)
}

// Add adds an api definition to the repository
func (r *FileSystemRepository) Add(definition *Definition) error {
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
func (r *FileSystemRepository) Remove(name string) error {
	r.Lock()
	defer r.Unlock()

	if _, err := r.findByName(name); err != nil {
		return err
	}

	delete(r.definitions, name)

	return nil
}

func (r *FileSystemRepository) parseDefinition(apiDef []byte) *Definition {
	appConfig := NewDefinition()
	if err := json.Unmarshal(apiDef, appConfig); err != nil {
		log.WithError(err).Error("[RPC] --> Couldn't unmarshal api configuration")
	}

	return appConfig
}
