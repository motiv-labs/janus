package api

import (
	"io/ioutil"
	"path/filepath"

	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
		filePath := filepath.Join(dir, f.Name())
		definition := NewDefinition()

		v := viper.New()
		v.SetConfigFile(filePath)
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.WithFields(log.Fields{"name": e.Name, "op": e.Op.String()}).Debug("API configuration changed, reloading...")
			if err := v.Unmarshal(definition); err != nil {
				log.WithError(err).Error("Can't unmarshal the API configuration")
			}
		})

		if err := v.ReadInConfig(); err != nil {
			log.WithError(err).Error("Couldn't load the api definition file")
			return nil, err
		}

		if err := v.Unmarshal(definition); err != nil {
			return nil, err
		}

		if err = repo.Add(definition); err != nil {
			log.WithError(err).Error("Can't add the definition to the repository")
			return nil, err
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

	definition, ok := r.definitions[name]
	if false == ok {
		return nil, ErrAPIDefinitionNotFound
	}

	return definition, nil
}

// Exists searches an existing Proxy definition by its listen_path
func (r *FileSystemRepository) Exists(def *Definition) (bool, error) {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.definitions[def.Name]
	if ok {
		return true, ErrAPINameExists
	}

	for _, definition := range r.definitions {
		if definition.Proxy.ListenPath == def.Proxy.ListenPath {
			return true, ErrAPIListenPathExists
		}
	}

	return false, nil
}

// Add adds an api definition to the repository
func (r *FileSystemRepository) Add(definition *Definition) error {
	r.Lock()
	defer r.Unlock()

	r.definitions[definition.Name] = definition

	return nil
}

// Remove removes an api definition from the repository
func (r *FileSystemRepository) Remove(name string) error {
	r.Lock()
	defer r.Unlock()

	delete(r.definitions, name)

	return nil
}
