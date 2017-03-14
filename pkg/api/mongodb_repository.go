package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName string = "api_specs"
)

// APISpecRepository defines the behaviour of a country repository
type APISpecRepository interface {
	FindAll() ([]*Definition, error)
	FindBySlug(slug string) (*Definition, error)
	FindByListenPath(path string) (*Definition, error)
	Add(app *Definition) error
	Remove(slug string) error
}

// MongoAPISpecRepository represents a mongodb repository
type MongoAPISpecRepository struct {
	session *mgo.Session
}

// NewMongoAppRepository creates a mongo API definition repo
func NewMongoAppRepository(session *mgo.Session) (*MongoAPISpecRepository, error) {
	return &MongoAPISpecRepository{session}, nil
}

// FindAll fetches all the API definitions available
func (r *MongoAPISpecRepository) FindAll() ([]*Definition, error) {
	result := []*Definition{}
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindBySlug find an API definition by its slug
func (r *MongoAPISpecRepository) FindBySlug(slug string) (*Definition, error) {
	var result *Definition
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(bson.M{"slug": slug}).One(&result)
	return result, err
}

// FindByListenPath searches an existing API definition by its listen_path
func (r *MongoAPISpecRepository) FindByListenPath(path string) (*Definition, error) {
	var result *Definition
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(bson.M{"proxy.listen_path": path}).One(&result)

	return result, err
}

// Add adds an API definition to the repository
func (r *MongoAPISpecRepository) Add(definition *Definition) error {
	session, coll := r.getSession()
	defer session.Close()

	isValid, err := govalidator.ValidateStruct(definition)
	if false == isValid && err != nil {
		fields := log.Fields{
			"errors": err.Error(),
		}
		log.WithFields(fields).Error("Validation errors")
		return err
	}

	_, err = coll.Upsert(bson.M{"slug": definition.Slug}, definition)
	if err != nil {
		log.Errorf("There was an error adding the resource %s", definition.Slug)
		return err
	}

	log.Debugf("Resource %s added", definition.Slug)
	return nil
}

// Remove removes an API definition from the repository
func (r *MongoAPISpecRepository) Remove(slug string) error {
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Remove(bson.M{"slug": slug})
	if err != nil {
		log.Errorf("There was an error removing the resource %s", slug)
		return err
	}

	log.Debugf("Resource %s removed", slug)
	return nil
}

func (r *MongoAPISpecRepository) getSession() (*mgo.Session, *mgo.Collection) {
	session := r.session.Copy()
	coll := session.DB("").C(collectionName)

	return session, coll
}
