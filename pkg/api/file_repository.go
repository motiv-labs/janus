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
	*InMemoryRepository
	sync.RWMutex
}

// NewFileSystemRepository creates a mongo country repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := &FileSystemRepository{InMemoryRepository: NewInMemoryRepository()}

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

func (r *FileSystemRepository) parseDefinition(apiDef []byte) *Definition {
	appConfig := NewDefinition()
	if err := json.Unmarshal(apiDef, appConfig); err != nil {
		log.WithError(err).Error("[RPC] --> Couldn't unmarshal api configuration")
	}

	return appConfig
}
