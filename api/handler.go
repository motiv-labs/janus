package api

import (
	"net/http"

	"github.com/hellofresh/janus/errors"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/request"
	"github.com/hellofresh/janus/response"
	"github.com/hellofresh/janus/router"
	"gopkg.in/mgo.v2"
)

type API struct {
	Loader *Loader
}

// GET /apps
func (u *API) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := u.getRepository(u.getDatabase(r))

		data, err := repo.FindAll()
		if err != nil {
			panic(err.Error())
		}

		response.JSON(w, http.StatusOK, data)
	}
}

// GetBy gets an application by its id
func (u *API) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			panic(errors.New(http.StatusNotFound, "Application not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusOK, data)
	}
}

// PUT /apps/:id
func (u *API) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))
		definition, err := repo.FindByID(id)
		if definition.ID == "" {
			panic(errors.New(http.StatusNotFound, "Application not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		repo.Add(definition)
		u.Loader.Load()

		response.JSON(w, http.StatusOK, definition)
	}
}

// POST /apps
func (u *API) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := u.getRepository(u.getDatabase(r))
		definition := NewDefinition()

		err := request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = repo.Add(definition)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		u.Loader.Load()
		response.JSON(w, http.StatusOK, definition)
	}
}

// DELETE /apps/:param1
func (u *API) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		err := repo.Remove(id)
		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		u.Loader.Load()
		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (u *API) getDatabase(r *http.Request) *mgo.Database {
	db := r.Context().Value(middleware.ContextKeyDatabase)

	if nil == db {
		panic(errors.New(http.StatusInternalServerError, "DB context was not set for this request"))
	}

	return db.(*mgo.Database)
}

// GetRepository gets the repository for the handlers
func (u *API) getRepository(db *mgo.Database) *MongoAPISpecRepository {
	repo, err := NewMongoAppRepository(db)
	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	return repo
}
