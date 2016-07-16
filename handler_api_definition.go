package main

import (
	"github.com/kataras/iris"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type AppsAPI struct {
	*iris.Context
	proxyRegister *ProxyRegister
}

// GET /apps
func (u AppsAPI) Get() {
	repo := u.getRepository()
	data, err := repo.FindAll()

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	u.JSON(iris.StatusOK, data)
}

// GET /apps/:param1 which its value passed to the id argument
func (u AppsAPI) GetBy(id string) {
	repo := u.getRepository()
	data, err := repo.FindByID(id)

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	u.JSON(iris.StatusOK, data)
}

// PUT /apps/:id
func (u AppsAPI) PutBy(id string) {
	repo := u.getRepository()
	apiSpec, err := repo.FindByID(id)

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	err = u.ReadJSON(&apiSpec)

	if err != nil {
		log.Errorf("Error when reading json: %s", err.Error())
	}

	repo.Add(&apiSpec)

	u.JSON(iris.StatusCreated, apiSpec)
}

// POST /apps
func (u AppsAPI) Post() {
	repo := u.getRepository()
	apiSpec := &APIDefinition{}
	err := u.ReadJSON(apiSpec)

	if err != nil {
		log.Errorf("Error when reading json: %s", err.Error())
	}

	repo.Add(apiSpec)

	u.JSON(iris.StatusCreated, apiSpec)
}

// DELETE /apps/:param1
func (u AppsAPI) DeleteBy(id string) {
	repo := u.getRepository()
	err := repo.Remove(id)

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	u.Response.SetStatusCode(iris.StatusNoContent)
}

// GetRepository gets the repository for the handlers
func (u AppsAPI) getRepository() *MongoAPISpecRepository {
	db := u.Context.Get("db").(*mgo.Database)
	repo, err := NewMongoAppRepository(db)

	if err != nil {
		log.Panic(err)
	}

	return repo
}
