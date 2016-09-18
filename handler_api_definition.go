package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type AppsAPI struct {
	apiManager *APIManager
}

// GET /apps
func (u AppsAPI) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))
		data, err := repo.FindAll()

		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, data)
	}
}

// GET /apps/:param1 which its value passed to the id argument
func (u AppsAPI) GetBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			c.JSON(http.StatusNotFound, "Application not found")
			return
		}

		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

// PUT /apps/:id
func (u AppsAPI) PutBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))
		definition, err := repo.FindByID(id)
		if definition.ID == "" {
			c.JSON(http.StatusNotFound, "Application not found")
			return
		}

		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err = c.Bind(definition)
		if nil != err {
			log.Errorf("Error when reading json: %s", err.Error())
			return
		}

		repo.Add(definition)
		u.apiManager.Load()

		c.JSON(http.StatusOK, definition)
	}
}

// POST /apps
func (u AppsAPI) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))
		definition := &APIDefinition{}

		err := c.Bind(definition)
		if nil != err {
			log.Fatal("Error when reading json")
		}

		err = repo.Add(definition)
		if nil != err {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		u.apiManager.Load()
		c.JSON(http.StatusCreated, definition)
	}
}

// DELETE /apps/:param1
func (u AppsAPI) DeleteBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))

		err := repo.Remove(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		u.apiManager.Load()
		c.String(http.StatusNoContent, "")
	}
}

func (u AppsAPI) getDatabase(c *gin.Context) *mgo.Database {
	db, exists := c.Get("db")

	if false == exists {
		log.Error("DB context was not set for this request")
	}

	return db.(*mgo.Database)
}

// GetRepository gets the repository for the handlers
func (u AppsAPI) getRepository(db *mgo.Database) *MongoAPISpecRepository {
	repo, err := NewMongoAppRepository(db)

	if err != nil {
		log.Panic(err)
	}

	return repo
}
