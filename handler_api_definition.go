package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"github.com/gin-gonic/gin"
	"net/http"
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
		id := c.Param("id")
		repo := u.getRepository(u.getDatabase(c))
		data, err := repo.FindByID(id)

		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, data)
	}
}

// PUT /apps/:id
func (u AppsAPI) PutBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		repo := u.getRepository(u.getDatabase(c))
		apiSpec, err := repo.FindByID(id)

		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		err = c.BindJSON(&apiSpec)

		if err != nil {
			log.Errorf("Error when reading json: %s", err.Error())
		}

		repo.Add(&apiSpec)
		u.apiManager.Load()

		c.JSON(http.StatusCreated, apiSpec)
	}
}

// POST /apps
func (u AppsAPI) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))
		apiSpec := &APIDefinition{}
		err := c.BindJSON(apiSpec)

		if err != nil {
			log.Errorf("Error when reading json: %s", err.Error())
		}

		repo.Add(apiSpec)
		u.apiManager.Load()

		c.JSON(http.StatusCreated, apiSpec)
	}
}

// DELETE /apps/:param1
func (u AppsAPI) DeleteBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		repo := u.getRepository(u.getDatabase(c))

		err := repo.Remove(id)
		if err != nil {
			log.Errorf(err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		u.apiManager.Load()
		c.Status(http.StatusNoContent)
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
