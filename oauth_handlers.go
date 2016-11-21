package janus

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellofresh/janus/errors"
	"gopkg.in/mgo.v2"
)

type OAuthAPI struct{}

func (u *OAuthAPI) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))

		data, err := repo.FindAll()
		if err != nil {
			panic(err.Error())
		}

		c.JSON(http.StatusOK, data)
	}
}

// GetBy gets an oauth server by its id
func (u *OAuthAPI) GetBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			panic(errors.New(http.StatusNotFound, "OAuth server not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		c.JSON(http.StatusOK, data)
	}
}

func (u *OAuthAPI) PutBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))
		definition, err := repo.FindByID(id)
		if definition.ID == "" {
			panic(errors.New(http.StatusNotFound, "OAuth server not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = c.Bind(definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		repo.Add(definition)
		c.JSON(http.StatusOK, definition)
	}
}

func (u *OAuthAPI) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := u.getRepository(u.getDatabase(c))
		var oauth OAuth

		err := c.Bind(&oauth)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = repo.Add(&oauth)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		c.JSON(http.StatusCreated, oauth)
	}
}

func (u *OAuthAPI) DeleteBy() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")
		repo := u.getRepository(u.getDatabase(c))

		err := repo.Remove(id)
		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		c.String(http.StatusNoContent, "")
	}
}

func (u *OAuthAPI) getDatabase(c *gin.Context) *mgo.Database {
	db, exists := c.Get("db")

	if false == exists {
		panic(errors.New(http.StatusInternalServerError, "DB context was not set for this request"))
	}

	return db.(*mgo.Database)
}

// GetRepository gets the repository for the handlers
func (u *OAuthAPI) getRepository(db *mgo.Database) *MongoOAuthRepository {
	repo, err := NewMongoOAuthRepository(db)
	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	return repo
}
