package main

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"time"
	log "github.com/Sirupsen/logrus"
)

// AppRepository defines the behaviour of a country repository
type APISpecRepository interface {
	FindAll() ([]*APIDefinition, error)
	FindByID(id string) (*APIDefinition, error)
	Add(app *APIDefinition) error
	Remove(id string) error
}

// MongoAppRepository represents a mongodb repository
type MongoAPISpecRepository struct {
	coll *mgo.Collection
}

// NewMongoCountryRepository creates a mongo country repo
func NewMongoAppRepository(db *mgo.Database) (*MongoAPISpecRepository, error) {
	coll := db.C("api_specs")

	return &MongoAPISpecRepository{coll}, nil
}

// FindAll fetches all the countries available
func (r MongoAPISpecRepository) FindAll() ([]*APIDefinition, error) {
	result := []*APIDefinition{}
	err := r.coll.Find(nil).All(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

// FindByID find a country by the iso2code provided
func (r MongoAPISpecRepository) FindByID(id string) (*APIDefinition, error) {
	result := &APIDefinition{}
	err := r.coll.FindId(bson.ObjectIdHex(id)).One(result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// Add adds a country to the repository
func (r MongoAPISpecRepository) Add(apiSpec *APIDefinition) error {
	var id bson.ObjectId

	if len(apiSpec.ID) == 0 {
		id = bson.NewObjectId()
		apiSpec.CreatedAt = time.Now()
	} else {
		id = apiSpec.ID
		apiSpec.UpdatedAt = time.Now()
	}

	_, err := r.coll.UpsertId(id, apiSpec)

	if err != nil {
		log.Errorf("There was an error adding the resource %s", id)
		return err
	}

	apiSpec.ID = id
	log.Infof("Resource %s added", id)

	return nil
}

// Remove removes a country from the repository
func (r MongoAPISpecRepository) Remove(id string) error {
	err := r.coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Errorf("There was an error removing the resource %s", id)
		return err
	}

	log.Infof("Resource %s removed", id)
	return nil
}