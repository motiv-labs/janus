package api

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName string = "api_specs"
)

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	session *mgo.Session
}

// NewMongoAppRepository creates a mongo API definition repo
func NewMongoAppRepository(session *mgo.Session) (*MongoRepository, error) {
	return &MongoRepository{session}, nil
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

// FindByName find an API definition by name
func (r *MongoRepository) FindByName(name string) (*Definition, error) {
	return r.findOneByQuery(bson.M{"name": name})
}

// FindByListenPath find an API definition by proxy listen path
func (r *MongoRepository) FindByListenPath(path string) (*Definition, error) {
	return r.findOneByQuery(bson.M{"proxy.listen_path": path})
}

func (r *MongoRepository) findOneByQuery(query interface{}) (*Definition, error) {
	var result = NewDefinition()
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(query).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrAPIDefinitionNotFound
		}
		return nil, err
	}

	return result, err
}

// Exists searches an existing API definition by its listen_path
func (r *MongoRepository) Exists(def *Definition) (bool, error) {
	return exists(r, def)
}

// Add adds an API definition to the repository
func (r *MongoRepository) Add(definition *Definition) error {
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
func (r *MongoRepository) Remove(name string) error {
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

// FindValidAPIHealthChecks retrieves all active apis that has health check configured
func (r *MongoRepository) FindValidAPIHealthChecks() ([]*Definition, error) {
	session, coll := r.getSession()
	defer session.Close()

	query := bson.M{
		"active": true,
		"health_check.url": bson.M{
			"$exists": true,
		},
	}

	result := []*Definition{}
	if err := coll.Find(query).All(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *MongoRepository) getSession() (*mgo.Session, *mgo.Collection) {
	session := r.session.Copy()
	coll := session.DB("").C(collectionName)

	return session, coll
}
