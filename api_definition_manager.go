package main

import (
	"gopkg.in/mgo.v2"
	"github.com/hellofresh/api-gateway/storage"
	log "github.com/Sirupsen/logrus"
)

// APIDefinitionLoader will load an Api definition from a storage system. It has two methods LoadDefinitionsFromMongo()
// and LoadDefinitions(), each will pull api specifications from different locations.
type APIDefinitionLoader struct{}

func (a *APIDefinitionLoader) LoadDefinitions(dir string) {

}

func (a *APIDefinitionLoader) LoadDefinitionsFromDatastore(dbSession *mgo.Session, dbConfig  storage.Database) []*APIDefinition {
	repo, err := NewMongoAppRepository(dbSession.DB(dbConfig.Name))

	if err != nil {
		log.Panic(err)
	}

	apiSpecs, err := repo.FindAll()

	if err != nil {
		log.Panic(err)
	}

	return apiSpecs;
}
