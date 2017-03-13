package api

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	definitions map[string]*Definition
}

// NewFileSystemRepository creates a mongo country repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := &FileSystemRepository{}

	// Grab json files from directory
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			filePath := filepath.Join(dir, f.Name())
			log.Info("Loading API Specification from ", filePath)
			appConfigBody, err := ioutil.ReadFile(filePath)
			definition := repo.parseDefinition(appConfigBody)
			if err != nil {
				log.Error("Couldn't load app configuration file: ", err)
				return nil, err
			}

			err = repo.Add(definition)
			if err != nil {
				log.WithError(err).Error("Can't add the definition to the repository")
				return nil, err
			}

		}
	}

	return repo, nil
}

// FindAll fetches all the countries available
func (r *FileSystemRepository) FindAll() ([]*Definition, error) {
	var definitions []*Definition
	for _, definition := range r.definitions {
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// FindBySlug find a country by the iso2code provided
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

// Add adds a country to the repository
func (r *FileSystemRepository) Add(definition *Definition) error {
	r.definitions[definition.Slug] = definition

	return nil
}

// Remove removes a country from the repository
func (r *FileSystemRepository) Remove(slug string) error {
	delete(r.definitions, slug)
	return nil
}

func (r *FileSystemRepository) parseDefinition(apiDef []byte) *Definition {
	appConfig := &Definition{}
	if err := json.Unmarshal(apiDef, appConfig); err != nil {
		log.Error("[RPC] --> Couldn't unmarshal api configuration: ", err)
	}

	return appConfig
}
