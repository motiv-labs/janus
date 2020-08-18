package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// FileSystemRepository represents a mongodb repository
type FileSystemRepository struct {
	*InMemoryRepository
	watcher *fsnotify.Watcher
}

// Type used for JSON.Unmarshaller
type definitionList struct {
	defs []*Definition
}

// NewFileSystemRepository creates a mongo country repo
func NewFileSystemRepository(dir string) (*FileSystemRepository, error) {
	repo := FileSystemRepository{InMemoryRepository: NewInMemoryRepository()}

	// Grab json files from directory
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create a file system watcher: %w", err)
	}

	repo.watcher = watcher

	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			filePath := filepath.Join(dir, f.Name())
			logger := log.WithField("path", filePath)

			appConfigBody, err := ioutil.ReadFile(filePath)
			if err != nil {
				logger.WithError(err).Error("Couldn't load the api definition file")
				return nil, err
			}

			err = repo.watcher.Add(filePath)
			if err != nil {
				logger.WithError(err).Error("Couldn't load the api definition file")
				return nil, err
			}

			definition := repo.parseDefinition(appConfigBody)
			for _, v := range definition.defs {
				if err = repo.add(v); err != nil {
					logger.WithField("name", v.Name).WithError(err).Error("Failed during add definition to the repository")
					return nil, err
				}
			}
		}
	}

	return &repo, nil
}

// Close terminates the session.  It's a runtime error to use a session
// after it has been closed.
func (r *FileSystemRepository) Close() error {
	return r.watcher.Close()
}

// Watch watches for changes on the database
func (r *FileSystemRepository) Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged) {
	go func() {
		for {
			select {
			case event := <-r.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {

					body, err := ioutil.ReadFile(event.Name)
					if err != nil {
						log.WithError(err).Error("Couldn't load the api definition file")
						continue
					}
					cfgChan <- ConfigurationChanged{
						Configurations: &Configuration{Definitions: r.parseDefinition(body).defs},
					}
				}
			case err := <-r.watcher.Errors:
				log.WithError(err).Error("error received from file system notify")
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (r *FileSystemRepository) parseDefinition(apiDef []byte) definitionList {
	appConfigs := definitionList{}

	// Try unmarshalling as if json is an unnamed Array of multiple definitions
	if err := json.Unmarshal(apiDef, &appConfigs); err != nil {
		// Try unmarshalling as if json is a single Definition
		appConfigs.defs = append(appConfigs.defs, NewDefinition())
		if err := json.Unmarshal(apiDef, &appConfigs.defs[0]); err != nil {
			log.WithError(err).Error("[RPC] --> Couldn't unmarshal api configuration")
		}
	}

	return appConfigs
}

func (d *definitionList) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &d.defs)
}
