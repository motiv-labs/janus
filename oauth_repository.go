package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// OAuthRepository defines the behaviour of a authentication repo
type OAuthRepository interface {
	FindAll() ([]*OAuth, error)
	FindByID(id string) (*OAuth, error)
	Add(oauth *OAuth) error
	Remove(id string) error
}

// MongoOAuthRepository represents a mongodb repository
type MongoOAuthRepository struct {
	coll *mgo.Collection
}

// NewMongoOAuthRepository creates a mongo country repo
func NewMongoOAuthRepository(db *mgo.Database) (*MongoOAuthRepository, error) {
	return &MongoOAuthRepository{db.C("oauth_servers")}, nil
}

// FindAll fetches all the countries available
func (r *MongoOAuthRepository) FindAll() ([]*OAuth, error) {
	var result []*OAuth

	err := r.coll.Find(nil).All(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

// FindByID find a country by the iso2code provided
func (r *MongoOAuthRepository) FindByID(id string) (*OAuth, error) {
	var result *OAuth

	if false == bson.IsObjectIdHex(id) {
		return result, ErrInvalidID
	}

	err := r.coll.FindId(bson.ObjectIdHex(id)).One(result)

	return result, err
}

// Add adds a country to the repository
func (r *MongoOAuthRepository) Add(oauth *OAuth) error {
	var id bson.ObjectId

	if len(oauth.ID) == 0 {
		id = bson.NewObjectId()
		oauth.CreatedAt = time.Now()
	} else {
		id = oauth.ID
		oauth.UpdatedAt = time.Now()
	}

	oauth.ID = id

	isValid, err := govalidator.ValidateStruct(oauth)
	if false == isValid && err != nil {
		fields := log.Fields{
			"errors": err.Error(),
		}
		log.WithFields(fields).Error("Validation errors")
		return err
	}

	_, err = r.coll.UpsertId(id, oauth)
	if err != nil {
		log.Errorf("There was an error adding the resource %s", id)
		return err
	}

	log.Debugf("Resource %s added", id)
	return nil
}

// Remove removes a country from the repository
func (r *MongoOAuthRepository) Remove(id string) error {
	if false == bson.IsObjectIdHex(id) {
		return ErrInvalidID
	}

	err := r.coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Errorf("There was an error removing the resource %s", id)
		return err
	}

	log.Debugf("Resource %s removed", id)
	return nil
}
