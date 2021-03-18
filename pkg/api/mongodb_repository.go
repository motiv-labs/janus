package api

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	collectionName = "api_specs"

	mongoConnTimeout  = 10 * time.Second
	mongoQueryTimeout = 10 * time.Second
)

// MongoRepository represents a mongodb repository
type MongoRepository struct {
	//TODO: we need to expose this so the plugins can use the same session. We should abstract mongo DB and provide
	// the plugins with a simple interface to search, insert, update and remove data from whatever backend implementation
	DB          *mongo.Database
	collection  *mongo.Collection
	client      *mongo.Client
	refreshTime time.Duration
}

// NewMongoAppRepository creates a mongo API definition repo
func NewMongoAppRepository(dsn string, refreshTime time.Duration) (*MongoRepository, error) {
	log.WithField("dsn", dsn).Debug("Trying to connect to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), mongoConnTimeout)
	defer cancel()

	connString, err := connstring.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse mongodb connection string: %w", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, fmt.Errorf("could not connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("could not ping mongodb after connect: %w", err)
	}

	mongoDB := client.Database(connString.Database)
	return &MongoRepository{
		DB:          mongoDB,
		client:      client,
		collection:  mongoDB.Collection(collectionName),
		refreshTime: refreshTime,
	}, nil
}

// Close terminates underlying mongo connection. It's a runtime error to use a session
// after it has been closed.
func (r *MongoRepository) Close() error {
	return r.client.Disconnect(context.TODO())
}

// Listen watches for changes on the configuration
func (r *MongoRepository) Listen(ctx context.Context, cfgChan <-chan ConfigurationMessage) {
	go func() {
		log.Debug("Listening for changes on the provider...")
		for {
			select {
			case cfg := <-cfgChan:
				switch cfg.Operation {
				case AddedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not add the configuration on the provider")
					}
				case UpdatedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not update the configuration on the provider")
					}
				case RemovedOperation:
					err := r.remove(cfg.Configuration.Name)
					if err != nil {
						log.WithError(err).Error("Could not remove the configuration from the provider")
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Watch watches for changes on the database
func (r *MongoRepository) Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged) {
	t := time.NewTicker(r.refreshTime)
	go func(refreshTicker *time.Ticker) {
		defer refreshTicker.Stop()
		log.Debug("Watching Provider...")
		for {
			select {
			case <-refreshTicker.C:
				defs, err := r.FindAll()
				if err != nil {
					log.WithError(err).Error("Failed to get configurations on watch")
					continue
				}

				cfgChan <- ConfigurationChanged{
					Configurations: &Configuration{Definitions: defs},
				}
			case <-ctx.Done():
				return
			}
		}
	}(t)
}

// FindAll fetches all the API definitions available
func (r *MongoRepository) FindAll() ([]*Definition, error) {
	var result []*Definition

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	// sort by name to have the same order all the time - for easier comparison
	cur, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		d := new(Definition)
		if err := cur.Decode(d); err != nil {
			return nil, err
		}

		result = append(result, d)
	}

	return result, cur.Err()
}

// Add adds an API definition to the repository
func (r *MongoRepository) add(definition *Definition) error {
	isValid, err := definition.Validate()
	if false == isValid && err != nil {
		log.WithError(err).Error("Validation errors")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"name": definition.Name},
		bson.M{"$set": definition},
		options.FindOneAndUpdate().SetUpsert(true),
	).Err(); err != nil {
		log.WithField("name", definition.Name).Error("There was an error adding the resource")
		return err
	}

	log.WithField("name", definition.Name).Debug("Resource added")
	return nil
}

// Remove removes an API definition from the repository
func (r *MongoRepository) remove(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	res, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		log.WithField("name", name).Error("There was an error removing the resource")
		return err
	}

	if res.DeletedCount < 1 {
		return ErrAPIDefinitionNotFound
	}

	log.WithField("name", name).Debug("Resource removed")
	return nil
}
