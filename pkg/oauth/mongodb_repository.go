package oauth

import (
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName string = "oauth_servers"
)

// Repository defines the behaviour of a authentication repo
type Repository interface {
	FindAll() ([]*OAuth, error)
	FindBySlug(slug string) (*OAuth, error)
	FindByTokenURL(url url.URL) (*OAuth, error)
	Add(oauth *OAuth) error
	Remove(id string) error
}

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	session *mgo.Session
}

// NewMongoRepository creates a mongo country repo
func NewMongoRepository(session *mgo.Session) (*MongoRepository, error) {
	return &MongoRepository{session}, nil
}

// FindAll fetches all the countries available
func (r *MongoRepository) FindAll() ([]*OAuth, error) {
	var result []*OAuth
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindBySlug find an oauth server by slug
func (r *MongoRepository) FindBySlug(slug string) (*OAuth, error) {
	var result *OAuth
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(bson.M{"slug": slug}).One(&result)

	return result, err
}

// Add adds a country to the repository
func (r *MongoRepository) Add(oauth *OAuth) error {
	session, coll := r.getSession()
	defer session.Close()

	isValid, err := govalidator.ValidateStruct(oauth)
	if false == isValid && err != nil {
		fields := log.Fields{
			"errors": err.Error(),
		}
		log.WithFields(fields).Error("Validation errors")
		return err
	}

	_, err = coll.Upsert(bson.M{"slug": oauth.Slug}, oauth)
	if err != nil {
		log.Errorf("There was an error adding the resource %s", oauth.Name)
		return err
	}

	log.Debugf("Resource %s added", oauth.Name)
	return nil
}

// Remove removes a country from the repository
func (r *MongoRepository) Remove(slug string) error {
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

// FindByTokenURL returns OAuth server records with corresponding token url
func (r *MongoRepository) FindByTokenURL(url url.URL) (*OAuth, error) {
	var result *OAuth
	session, coll := r.getSession()
	defer session.Close()

	err := coll.Find(bson.M{"oauth_endpoints.token.target_url": url.String()}).One(&result)

	return result, err
}

func (r *MongoRepository) getSession() (*mgo.Session, *mgo.Collection) {
	session := r.session.Copy()
	coll := session.DB("").C(collectionName)

	return session, coll
}
