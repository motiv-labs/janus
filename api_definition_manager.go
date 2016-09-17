package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

// APIDefinitionLoader will load an Api definition from a storage system. It has two methods LoadDefinitionsFromMongo()
// and LoadDefinitions(), each will pull api specifications from different locations.
type APIDefinitionLoader struct{}

func (a *APIDefinitionLoader) LoadDefinitions(dir string) {

}

func (a *APIDefinitionLoader) LoadDefinitionsFromDatastore(dbSession *mgo.Session) []*APISpec {
	repo, err := NewMongoAppRepository(dbSession.DB(""))

	if err != nil {
		log.Panic(err)
	}

	definitions, err := repo.FindAll()

	if err != nil {
		log.Panic(err)
	}

	var APISpecs = []*APISpec{}
	for _, definition := range definitions {
		newAppSpec := a.MakeSpec(definition)
		APISpecs = append(APISpecs, &newAppSpec)
	}

	return APISpecs
}

func (a *APIDefinitionLoader) MakeSpec(definition APIDefinition) APISpec {
	newAppSpec := APISpec{}
	newAppSpec.APIDefinition = definition

	return newAppSpec
}
