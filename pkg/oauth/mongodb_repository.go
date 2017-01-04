package oauth

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/hellofresh/janus/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Repository defines the behaviour of a authentication repo
type Repository interface {
	FindAll() ([]*OAuth, error)
	FindByID(id string) (*OAuth, error)
	Add(oauth *OAuth) error
	Remove(id string) error
}

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	coll *mgo.Collection
}

// NewMongoRepository creates a mongo country repo
func NewMongoRepository(db *mgo.Database) (*MongoRepository, error) {
	return &MongoRepository{db.C("oauth_servers")}, nil
}

// FindAll fetches all the countries available
func (r *MongoRepository) FindAll() ([]*OAuth, error) {
	var result []*OAuth

	err := r.coll.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByID find a country by the iso2code provided
func (r *MongoRepository) FindByID(id string) (*OAuth, error) {
	var result *OAuth

	if false == bson.IsObjectIdHex(id) {
		return nil, errors.ErrInvalidID
	}

	err := r.coll.FindId(bson.ObjectIdHex(id)).One(&result)

	return result, err
}

// Add adds a country to the repository
func (r *MongoRepository) Add(oauth *OAuth) error {
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
func (r *MongoRepository) Remove(id string) error {
	if false == bson.IsObjectIdHex(id) {
		return errors.ErrInvalidID
	}

	err := r.coll.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		log.Errorf("There was an error removing the resource %s", id)
		return err
	}

	log.Debugf("Resource %s removed", id)
	return nil
}
