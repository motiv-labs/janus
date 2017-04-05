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
		filePath := filepath.Join(dir, f.Name())
		definition := new(Definition)

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
	var definitions []*Definition
	for _, definition := range r.definitions {
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

// FindByName find an api definition by name
func (r *FileSystemRepository) FindByName(name string) (*Definition, error) {
	definition, ok := r.definitions[name]
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
