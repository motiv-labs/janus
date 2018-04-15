package api

import (
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName string = "api_specs"
)

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	// TODO: we need to expose this so the plugins can use the same session. We should abstract the session and provide
	// the plugins with a simple interface to search, insert, update and remove data from whatever backend implementation
	Session     *mgo.Session
	refreshTime time.Duration
}

// NewMongoAppRepository creates a mongo API definition repo
func NewMongoAppRepository(dsn string, refreshTime time.Duration) (*MongoRepository, error) {
	log.WithField("dsn", dsn).Debug("Trying to connect to MongoDB...")
	session, err := mgo.Dial(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to mongodb")
	}

	log.Debug("Connected to MongoDB")
	session.SetMode(mgo.Monotonic, true)

	return &MongoRepository{Session: session, refreshTime: refreshTime}, nil
}

// Close terminates the session.  It's a runtime error to use a session
// after it has been closed.
func (r *MongoRepository) Close() error {
	r.Session.Close()
	return nil
}

// Listen watches for changes on the configuration
func (r *MongoRepository) Listen(ctx context.Context, cfgChan <-chan ConfigurationMessage) {
	go func() {
		log.Debug("Listening for changes on the provider...")
		for {
			select {
			case cfg := <-cfgChan:
				switch cfg.Operation {
				case AddedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not add the configuration on the provider")
					}
				case UpdatedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not update the configuration on the provider")
					}
				case RemovedOperation:
					err := r.remove(cfg.Configuration.Name)
					if err != nil {
						log.WithError(err).Error("Could not remove the configuration from the provider")
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Watch watches for changes on the database
func (r *MongoRepository) Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged) {
	t := time.NewTicker(r.refreshTime)
	go func(refreshTicker *time.Ticker) {
		defer refreshTicker.Stop()
		log.Debug("Watching Provider...")
		for {
			select {
			case <-refreshTicker.C:
				defs, err := r.FindAll()
				if err != nil {
					log.WithError(err).Error("Failed to get configurations on watch")
					continue
				}

				cfgChan <- ConfigurationChanged{
					Configurations: &Configuration{Definitions: defs},
				}
			case <-ctx.Done():
				return
			}
		}
	}(t)
}

// FindAll fetches all the API definitions available
func (r *MongoRepository) FindAll() ([]*Definition, error) {
	result := []*Definition{}
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Add adds an API definition to the repository
func (r *MongoRepository) add(definition *Definition) error {
	session, coll := r.getSession()
	defer session.Close()

	isValid, err := definition.Validate()
	if false == isValid && err != nil {
		log.WithError(err).Error("Validation errors")
		return err
	}

	_, err = coll.Upsert(bson.M{"name": definition.Name}, definition)
	if err != nil {
		log.WithField("name", definition.Name).Error("There was an error adding the resource")
		return err
	}

	log.WithField("name", definition.Name).Debug("Resource added")
	return nil
}

// Remove removes an API definition from the repository
func (r *MongoRepository) remove(name string) error {
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Remove(bson.M{"name": name})
	if err != nil {
		if err == mgo.ErrNotFound {
			return ErrAPIDefinitionNotFound
		}
		log.WithField("name", name).Error("There was an error removing the resource")
		return err
	}

	log.WithField("name", name).Debug("Resource removed")
	return nil
}

func (r *MongoRepository) getSession() (*mgo.Session, *mgo.Collection) {
	session := r.Session.Copy()
	coll := session.DB("").C(collectionName)

	return session, coll
}
