package basic

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName string = "basic_auth"

	mongoQueryTimeout = 10 * time.Second
)

// User represents an user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Repository represents an user repository
type Repository interface {
	FindAll() ([]*User, error)
	FindByUsername(username string) (*User, error)
	Add(user *User) error
	Remove(username string) error
}

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository creates a mongo API definition repo
func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{db.Collection(collectionName)}, nil
}

// FindAll fetches all the API definitions available
func (r *MongoRepository) FindAll() ([]*User, error) {
	var result []*User

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "username", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		u := new(User)
		if err := cur.Decode(u); err != nil {
			return nil, err
		}

		result = append(result, u)
	}

	return result, cur.Err()
}

// FindByUsername find an user by username
func (r *MongoRepository) FindByUsername(username string) (*User, error) {
	return r.findOneByQuery(bson.M{"username": username})
}

func (r *MongoRepository) findOneByQuery(query interface{}) (*User, error) {
	var result User

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	err := r.collection.FindOne(ctx, query).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	}

	return &result, err
}

// Add adds an user to the repository
func (r *MongoRepository) Add(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"username": user.Username},
		bson.M{"$set": user},
		options.FindOneAndUpdate().SetUpsert(true),
	).Err(); err != nil {
		log.WithField("username", user.Username).Error("There was an error adding the user")
		return err
	}

	log.WithField("username", user.Username).Debug("User added")
	return nil
}

// Remove an user from the repository
func (r *MongoRepository) Remove(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	res, err := r.collection.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		log.WithField("username", username).Error("There was an error removing the user")
		return err
	}

	if res.DeletedCount < 1 {
		return ErrUserNotFound
	}

	log.WithField("username", username).Debug("User removed")
	return nil
}
