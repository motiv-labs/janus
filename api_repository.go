package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// APISpecRepository defines the behaviour of a country repository
type APISpecRepository interface {
	FindAll() ([]APIDefinition, error)
	FindByID(id string) (APIDefinition, error)
	Add(app *APIDefinition) error
	Remove(id string) error
}

// MongoAPISpecRepository represents a mongodb repository
type MongoAPISpecRepository struct {
	coll *mgo.Collection
}

// NewMongoAppRepository creates a mongo country repo
func NewMongoAppRepository(db *mgo.Database) (*MongoAPISpecRepository, error) {
	coll := db.C("api_specs")

	return &MongoAPISpecRepository{coll}, nil
}

// FindAll fetches all the countries available
func (r *MongoAPISpecRepository) FindAll() ([]APIDefinition, error) {
	result := []APIDefinition{}

	err := r.coll.Find(nil).All(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// FindByID find a country by the iso2code provided
func (r *MongoAPISpecRepository) FindByID(id string) (*APIDefinition, error) {
	result := &APIDefinition{}

	if false == bson.IsObjectIdHex(id) {
		return result, ErrInvalidID
	}

	err := r.coll.FindId(bson.ObjectIdHex(id)).One(result)

	return result, err
}

// Add adds a country to the repository
func (r *MongoAPISpecRepository) Add(definition *APIDefinition) error {
	var id bson.ObjectId

	if len(definition.ID) == 0 {
		id = bson.NewObjectId()
		definition.CreatedAt = time.Now()
	} else {
		id = definition.ID
		definition.UpdatedAt = time.Now()
	}

	definition.ID = id

	isValid, err := govalidator.ValidateStruct(definition)
	if false == isValid && err != nil {
		fields := log.Fields{
			"errors": err.Error(),
		}
		log.WithFields(fields).Error("Validation errors")
		return err
	}

	_, err = r.coll.UpsertId(id, definition)
	if err != nil {
		log.Errorf("There was an error adding the resource %s", id)
		return err
	}

	log.Infof("Resource %s added", id)
	return nil
}

// Remove removes a country from the repository
func (r *MongoAPISpecRepository) Remove(id string) error {
	if false == bson.IsObjectIdHex(id) {
		return ErrInvalidID
	}

	err := r.coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Errorf("There was an error removing the resource %s", id)
		return err
	}

	log.Infof("Resource %s removed", id)
	return nil
}
