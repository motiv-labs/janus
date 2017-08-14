package api

import (
	"fmt"
	"net/url"

	mgo "gopkg.in/mgo.v2"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	mongodb = "mongodb"
	file    = "file"
)

// Repository defines the behavior of a proxy specs repository
type Repository interface {
	FindAll() ([]*Definition, error)
	FindByName(name string) (*Definition, error)
	FindByListenPath(path string) (*Definition, error)
	Exists(def *Definition) (bool, error)
	Add(app *Definition) error
	Remove(name string) error
	FindValidAPIHealthChecks() ([]*Definition, error)
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
func BuildRepository(dsn string, session *mgo.Session) (Repository, error) {
	var repo Repository
	dsnURL, err := url.Parse(dsn)
	switch dsnURL.Scheme {
	case mongodb:
		repo, err = NewMongoAppRepository(session)
		if err != nil {
			return nil, errors.Wrap(err, "Could not create a mongodb repository for api definitions")
		}
	case file:
		log.Debug("File system based configuration chosen")
		apiPath := fmt.Sprintf("%s/apis", dsnURL.Path)

		log.WithField("api_path", apiPath).Debug("Trying to load configuration files")
		repo, err = NewFileSystemRepository(apiPath)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a file system repository")
		}
	default:
		return nil, errors.New("The selected scheme is not supported to load API definitions")
	}

	return repo, nil
}
