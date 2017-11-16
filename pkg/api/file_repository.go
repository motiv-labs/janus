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

// Type used for JSON.Unmarshaller
type definitionList struct {
	defs []Definition
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
			for _, v := range definition.defs {
				if err = repo.Add(&v); err != nil {
					log.WithField("name", v.Name).WithError(err).Error("Failed during add definition to the repository")
					return nil, err
				}
			}
		}
	}

	return repo, nil
}

func (r *FileSystemRepository) parseDefinition(apiDef []byte) definitionList {
	appConfigs := definitionList{}

	//Try unmarshalling as if json is an unnamed Array of multiple definitions
	if err := json.Unmarshal(apiDef, &appConfigs); err != nil {
		//Try unmarshalling as if json is a single Definition
		appConfigs.defs = append(appConfigs.defs, *NewDefinition())
		if err := json.Unmarshal(apiDef, &appConfigs.defs[0]); err != nil {
			log.WithError(err).Error("[RPC] --> Couldn't unmarshal api configuration")
		}
	}

	return appConfigs
}

func (d *definitionList) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &d.defs)
}
