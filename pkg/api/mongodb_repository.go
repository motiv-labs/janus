package api

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName string = "api_specs"
)

// Repository defines the behaviour of a country repository
type Repository interface {
	FindAll() ([]*Definition, error)
	FindByName(name string) (*Definition, error)
	Exists(def *Definition) (bool, error)
	Add(app *Definition) error
	Remove(name string) error
	FindValidAPIHealthChecks() ([]*Definition, error)
}

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
	var result = NewDefinition()
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(bson.M{"name": name}).One(&result)
	if err != nil {
		return nil, ErrAPIDefinitionNotFound
	}

	return result, err
}

// Exists searches an existing API definition by its listen_path
func (r *MongoRepository) Exists(def *Definition) (bool, error) {
	session, coll := r.getSession()
	defer session.Close()

	count, err := coll.Find(bson.M{"name": def.Name}).Count()
	if err != nil {
		return false, err
	}

	if count >= 1 {
		return true, ErrAPINameExists
	}

	count, err = coll.Find(bson.M{"proxy.listen_path": def.Proxy.ListenPath}).Count()
	if err != nil {
		return false, err
	}

	if count >= 1 {
		return true, ErrAPIListenPathExists
	}

	return false, nil
}

// Add adds an API definition to the repository
func (r *MongoRepository) Add(definition *Definition) error {
	session, coll := r.getSession()
	defer session.Close()

	isValid, err := definition.Validate()
	if false == isValid && err != nil {
		fields := log.Fields{
			"errors": err.Error(),
		}
		log.WithFields(fields).Error("Validation errors")
		return err
	}

	_, err = coll.Upsert(bson.M{"name": definition.Name}, definition)
	if err != nil {
		log.Errorf("There was an error adding the resource %s", definition.Name)
		return err
	}

	log.Debugf("Resource %s added", definition.Name)
	return nil
}

// Remove removes an API definition from the repository
func (r *MongoRepository) Remove(name string) error {
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Remove(bson.M{"name": name})
	if err != nil {
		log.Errorf("There was an error removing the resource %s", name)
		return err
	}

	log.Debugf("Resource %s removed", name)
	return nil
}

// FindValidAPIHealthChecks retreives all apis that has health check configured
func (r *MongoRepository) FindValidAPIHealthChecks() ([]*Definition, error) {
	result := []*Definition{}
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(bson.M{"health_check.url": bson.M{"$exists": true}}).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *MongoRepository) getSession() (*mgo.Session, *mgo.Collection) {
	session := r.session.Copy()
	coll := session.DB("").C(collectionName)

	return session, coll
}
