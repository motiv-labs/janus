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
	repo := u.GetRepository()
	data, err := repo.FindAll()

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	u.JSON(iris.StatusOK, data)
}

// GET /apps/:param1 which its value passed to the id argument
func (u AppsAPI) GetBy(id string) {
	repo := u.GetRepository()
	data, err := repo.FindByID(id)

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	u.JSON(iris.StatusOK, data)
}

// PUT /apps/:id
func (u AppsAPI) Put(id string) {
	repo := u.GetRepository()
	app, err := repo.FindByID(id)

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	err = u.ReadJSON(app)

	if err != nil {
		log.Errorf("Error when reading json: %s", err.Error())
	}

	repo.Add(app)
	u.Response.SetStatusCode(iris.StatusOK)
}

// POST /apps
func (u AppsAPI) Post() {
	repo := u.GetRepository()
	definition := &APIDefinition{}
	err := u.ReadJSON(definition)

	if err != nil {
		log.Errorf("Error when reading json: %s", err.Error())
	}

	repo.Add(definition)
	u.proxyRegister.Register(definition.Proxy)

	u.JSON(iris.StatusCreated, definition)
}

// DELETE /apps/:param1
func (u AppsAPI) DeleteBy(id string) {
	repo := u.GetRepository()
	err := repo.Remove(id)

	if err != nil {
		log.Errorf(err.Error())
		u.JSON(iris.StatusInternalServerError, err.Error())
	}

	u.Response.SetStatusCode(iris.StatusNoContent)
}

// GetRepository gets the repository for the handlers
func (u AppsAPI) GetRepository() *MongoAPISpecRepository {
	db := u.Context.Get("db").(*mgo.Database)
	repo, err := NewMongoAppRepository(db)

	if err != nil {
		log.Panic(err)
	}

	return repo
}
