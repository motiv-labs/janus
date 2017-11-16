package oauth2

import (
	"net/url"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionName string = "oauth_servers"
)

// Repository defines the behavior of a OAuth Server repo
type Repository interface {
	FindAll() ([]*OAuth, error)
	FindByName(name string) (*OAuth, error)
	FindByTokenURL(url url.URL) (*OAuth, error)
	Add(oauth *OAuth) error
	Save(oauth *OAuth) error
	Remove(id string) error
}

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	session *mgo.Session
}

// NewMongoRepository creates a mongodb OAuth Server repo
func NewMongoRepository(session *mgo.Session) (*MongoRepository, error) {
	return &MongoRepository{session}, nil
}

// FindAll fetches all the OAuth Servers available
func (r *MongoRepository) FindAll() ([]*OAuth, error) {
	session, coll := r.getSession()
	defer session.Close()

	var result []*OAuth
	err := coll.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByName find an OAuth Server by name
func (r *MongoRepository) FindByName(name string) (*OAuth, error) {
	session, coll := r.getSession()
	defer session.Close()

	var result *OAuth
	if err := coll.Find(bson.M{"name": name}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrOauthServerNotFound
		}
		return nil, err
	}

	return result, nil
}

// Add add a new OAuth Server to the repository
func (r *MongoRepository) Add(oauth *OAuth) error {
	session, coll := r.getSession()
	defer session.Close()

	isValid, err := govalidator.ValidateStruct(oauth)
	if !isValid && err != nil {
		log.WithField("errors", err.Error()).Error("Validation errors")
		return err
	}

	if err = coll.Insert(oauth); err != nil {
		log.WithField("name", oauth.Name).
			WithError(err).
			Error("There was an error persisting the resource")
		if mgo.IsDup(err) {
			return ErrOauthServerNameExists
		}
		return err
	}

	log.WithField("name", oauth.Name).Debug("Resource persisted")
	return nil
}

// Save saves OAuth Server to the repository
func (r *MongoRepository) Save(oauth *OAuth) error {
	session, coll := r.getSession()
	defer session.Close()

	isValid, err := govalidator.ValidateStruct(oauth)
	if !isValid && err != nil {
		log.WithField("errors", err.Error()).Error("Validation errors")
		return err
	}

	_, err = coll.Upsert(bson.M{"name": oauth.Name}, oauth)
	if err != nil {
		log.WithField("name", oauth.Name).
			WithError(err).
			Error("There was an error adding the resource")
		return err
	}

	log.WithField("name", oauth.Name).Debug("Resource added")
	return nil
}

// Remove removes an OAuth Server from the repository
func (r *MongoRepository) Remove(name string) error {
	session, coll := r.getSession()
	defer session.Close()

	if err := coll.Remove(bson.M{"name": name}); err != nil {
		log.WithField("name", name).
			WithError(err).
			Error("There was an error removing the resource")
		return err
	}

	log.WithField("name", name).Debug("Resource removed")
	return nil
}

// FindByTokenURL returns OAuth Server records with corresponding token url
func (r *MongoRepository) FindByTokenURL(url url.URL) (*OAuth, error) {
	session, coll := r.getSession()
	defer session.Close()

	query := bson.M{
		"oauth_endpoints.token.upstreams.targets.target": url.String(),
	}

	var result *OAuth
	err := coll.Find(query).One(&result)

	return result, err
}

func (r *MongoRepository) getSession() (*mgo.Session, *mgo.Collection) {
	session := r.session.Copy()
	coll := session.DB("").C(collectionName)

	return session, coll
}
