package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	mongodb = "mongodb"
	file    = "file"
)

// Repository defines the behavior of a proxy specs repository
type Repository interface {
	io.Closer
	FindAll() ([]*Definition, error)
}

// Watcher defines how a provider should watch for changes on configurations
type Watcher interface {
	Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged)
}

// Listener defines how a provider should listen for changes on configurations
type Listener interface {
	Listen(ctx context.Context, cfgChan <-chan ConfigurationMessage)
}

// BuildRepository creates a repository instance that will depend on your given DSN
func BuildRepository(dsn string, refreshTime time.Duration) (Repository, error) {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing the DSN: %w", err)
	}

	switch dsnURL.Scheme {
	case mongodb:
		log.Debug("MongoDB configuration chosen")
		return NewMongoAppRepository(dsn, refreshTime)
	case file:
		log.Debug("File system based configuration chosen")
		apiPath := fmt.Sprintf("%s/apis", dsnURL.Path)

		log.WithField("path", apiPath).Debug("Trying to load API configuration files")
		repo, err := NewFileSystemRepository(apiPath)
		if err != nil {
			return nil, fmt.Errorf("could not create a file system repository: %w", err)
		}
		return repo, nil
	default:
		return nil, errors.New("selected scheme is not supported to load API definitions")
	}
}
