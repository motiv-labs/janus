package api

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/pkg/errors"
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
	FindByName(name string) (*Definition, error)
	FindByListenPath(path string) (*Definition, error)
	Exists(def *Definition) (bool, error)
	Add(app *Definition) error
	Remove(name string) error
	FindValidAPIHealthChecks() ([]*Definition, error)
	Watch(ctx context.Context) <-chan ConfigrationChanged
}

func exists(r Repository, def *Definition) (bool, error) {
	_, err := r.FindByName(def.Name)
	if nil != err && err != ErrAPIDefinitionNotFound {
		return false, err
	} else if err != ErrAPIDefinitionNotFound {
		return true, ErrAPINameExists
	}

	_, err = r.FindByListenPath(def.Proxy.ListenPath)
	if nil != err && err != ErrAPIDefinitionNotFound {
		return false, err
	} else if err != ErrAPIDefinitionNotFound {
		return true, ErrAPIListenPathExists
	}

	return false, nil
}

// BuildRepository creates a repository instance that will depend on your given DSN
func BuildRepository(dsn string, refreshTime time.Duration) (Repository, error) {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing the DSN")
	}

	switch dsnURL.Scheme {
	case mongodb:
		log.Debug("MongoDB configuration chosen")
		return NewMongoAppRepository(dsn, refreshTime)
	case file:
		log.Debug("File system based configuration chosen")
		apiPath := fmt.Sprintf("%s/apis", dsnURL.Path)

		log.WithField("api_path", apiPath).Debug("Trying to load configuration files")
		repo, err := NewFileSystemRepository(apiPath)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a file system repository")
		}
		return repo, nil
	default:
		return nil, errors.New("The selected scheme is not supported to load API definitions")
	}
}
