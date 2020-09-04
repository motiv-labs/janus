package oauth2

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "oauth_servers"

	mongoQueryTimeout = 10 * time.Second
)

// Repository defines the behavior of a OAuth Server repo
type Repository interface {
	FindAll() ([]*OAuth, error)
	FindByName(name string) (*OAuth, error)
	Add(oauth *OAuth) error
	Save(oauth *OAuth) error
	Remove(id string) error
}

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a mongodb OAuth Server repo
func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{db.Collection(collectionName)}, nil
}

// FindAll fetches all the OAuth Servers available
func (r *MongoRepository) FindAll() ([]*OAuth, error) {
	var result []*OAuth

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		o := new(OAuth)
		if err := cur.Decode(o); err != nil {
			return nil, err
		}

		result = append(result, o)
	}

	return result, cur.Err()
}

// FindByName find an OAuth Server by name
func (r *MongoRepository) FindByName(name string) (*OAuth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	result := NewOAuth()
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, ErrOauthServerNotFound
	}

	return result, err
}

// Add add a new OAuth Server to the repository
func (r *MongoRepository) Add(oauth *OAuth) error {
	isValid, err := govalidator.ValidateStruct(oauth)
	if !isValid && err != nil {
		log.WithField("errors", err.Error()).Error("Validation errors")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	_, err = r.collection.InsertOne(ctx, oauth)
	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrOauthServerNameExists
		}
		log.WithField("name", oauth.Name).WithError(err).Error("There was an error persisting the resource")
		return err
	}

	log.WithField("name", oauth.Name).Debug("Resource persisted")
	return nil
}

// Save saves OAuth Server to the repository
func (r *MongoRepository) Save(oauth *OAuth) error {
	isValid, err := govalidator.ValidateStruct(oauth)
	if !isValid && err != nil {
		log.WithField("errors", err.Error()).Error("Validation errors")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"name": oauth.Name},
		bson.M{"$set": oauth},
		options.FindOneAndUpdate().SetUpsert(true),
	).Err(); err != nil {
		log.WithField("name", oauth.Name).WithError(err).Error("There was an error adding the resource")
		return err
	}

	log.WithField("name", oauth.Name).Debug("Resource added")
	return nil
}

// Remove removes an OAuth Server from the repository
func (r *MongoRepository) Remove(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		log.WithField("name", name).Error("There was an error removing the resource")
		return err
	}

	log.WithField("name", name).Debug("Resource removed")
	return nil
}

func isDuplicateKeyError(err error) bool {
	// TODO: maybe there is (or will be) a better way of checking duplicate key error
	// this one is based on https://github.com/mongodb/mongo-go-driver/blob/master/mongo/integration/collection_test.go#L54-L65
	we, ok := err.(mongo.WriteException)
	if !ok {
		return false
	}

	return len(we.WriteErrors) > 0 && we.WriteErrors[0].Code == 11000
}
