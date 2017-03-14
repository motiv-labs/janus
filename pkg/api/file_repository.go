package api

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"sync"

	log "github.com/Sirupsen/logrus"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	sync.Mutex
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
			log.WithField("path", filePath).Info("Loading API definition from file")
			appConfigBody, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.WithError(err).WithField("path", filePath).Error("Couldn't load the api definition file")
				return nil, err
			}

			definition := repo.parseDefinition(appConfigBody)
			err = repo.Add(definition)
			if err != nil {
				log.WithError(err).Error("Can't add the definition to the repository")
				return nil, err
			}
		}
	}

	return repo, nil
}

// FindAll fetches all the api definitions available
func (r *FileSystemRepository) FindAll() ([]*Definition, error) {
	var definitions []*Definition
	for _, definition := range r.definitions {
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// FindBySlug find an api definition by slug
func (r *FileSystemRepository) FindBySlug(slug string) (*Definition, error) {
	definition, ok := r.definitions[slug]
	if false == ok {
		return nil, ErrAPIDefinitionNotFound
	}

	return definition, nil
}

// FindByListenPath searches an existing Proxy definition by its listen_path
func (r *FileSystemRepository) FindByListenPath(path string) (*Definition, error) {
	for _, definition := range r.definitions {
		if definition.Proxy.ListenPath == path {
			return definition, nil
		}
	}

	return nil, ErrAPIDefinitionNotFound
}

// Add adds an api definition to the repository
func (r *FileSystemRepository) Add(definition *Definition) error {
	r.Lock()
	defer r.Unlock()

	r.definitions[definition.Slug] = definition

	return nil
}

// Remove removes an api definition from the repository
func (r *FileSystemRepository) Remove(slug string) error {
	r.Lock()
	defer r.Unlock()

	delete(r.definitions, slug)

	return nil
}

func (r *FileSystemRepository) parseDefinition(apiDef []byte) *Definition {
	appConfig := &Definition{}
	if err := json.Unmarshal(apiDef, appConfig); err != nil {
		log.WithError(err).Error("[RPC] --> Couldn't unmarshal api configuration")
	}

	return appConfig
}
