package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellofresh/janus/errors"
	"gopkg.in/mgo.v2"
)

type AppsAPI struct {
	apiManager *APIManager
}

// GET /apps
func (u *AppsAPI) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))

		data, err := repo.FindAll()
		if err != nil {
			panic(err.Error())
		}

		c.JSON(http.StatusOK, data)
	}
}

// GetBy gets an application by its id
func (u *AppsAPI) GetBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			panic(errors.New(http.StatusNotFound, "Application not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		c.JSON(http.StatusOK, data)
	}
}

// PUT /apps/:id
func (u *AppsAPI) PutBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))
		definition, err := repo.FindByID(id)
		if definition.ID == "" {
			panic(errors.New(http.StatusNotFound, "Application not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = c.Bind(definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		repo.Add(definition)
		u.apiManager.Load()

		c.JSON(http.StatusOK, definition)
	}
}

// POST /apps
func (u *AppsAPI) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))
		definition := &APIDefinition{}

		err := c.Bind(definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = repo.Add(definition)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		u.apiManager.Load()
		c.JSON(http.StatusCreated, definition)
	}
}

// DELETE /apps/:param1
func (u *AppsAPI) DeleteBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))

		err := repo.Remove(id)
		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		u.apiManager.Load()
		c.String(http.StatusNoContent, "")
	}
}

func (u *AppsAPI) getDatabase(c *gin.Context) *mgo.Database {
	db, exists := c.Get("db")

	if false == exists {
		panic(errors.New(http.StatusInternalServerError, "DB context was not set for this request"))
	}

	return db.(*mgo.Database)
}

// GetRepository gets the repository for the handlers
func (u *AppsAPI) getRepository(db *mgo.Database) *MongoAPISpecRepository {
	repo, err := NewMongoAppRepository(db)
	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	return repo
}
