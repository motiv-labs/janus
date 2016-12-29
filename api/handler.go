package api

import (
	"net/http"

	"github.com/hellofresh/janus/errors"
	"github.com/hellofresh/janus/loader"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/request"
	"github.com/hellofresh/janus/response"
	"github.com/hellofresh/janus/router"
	"gopkg.in/mgo.v2"
)

// Controller is the api rest controller
type Controller struct {
	changeTracker *loader.Tracker
}

// NewController creates a new instance of Controller
func NewController(changeTracker *loader.Tracker) *Controller {
	return &Controller{changeTracker}
}

// GET /apps
func (u *Controller) Get() http.HandlerFunc {
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
func (u *Controller) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			panic(ErrAPIDefinitionNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusOK, data)
	}
}

// PUT /apps/:id
func (u *Controller) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))
		definition, err := repo.FindByID(id)
		if definition.ID == "" {
			panic(ErrAPIDefinitionNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = repo.Add(definition)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}
		u.changeTracker.Change()
		response.JSON(w, http.StatusOK, definition)
	}
}

// Post corresponds to POST /apis to create a new Proxy definition
func (u *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := u.getRepository(u.getDatabase(r))
		definition := NewDefinition()

		err := request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		def, err := repo.FindByListenPath(definition.Proxy.ListenPath)
		if nil != err && err != mgo.ErrNotFound {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		if def != nil {
			panic(errors.ErrProxyExists)
		}
		err = repo.Add(definition)
		if nil != err {

			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		u.changeTracker.Change()
		response.JSON(w, http.StatusOK, definition)
	}
}

// DELETE /apps/:param1
func (u *Controller) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		err := repo.Remove(id)
		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		u.changeTracker.Change()
		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (u *Controller) getDatabase(r *http.Request) *mgo.Database {
	db := r.Context().Value(middleware.ContextKeyDatabase)

	if nil == db {
		panic(ErrDBContextNotSet)
	}

	return db.(*mgo.Database)
}

// GetRepository gets the repository for the handlers
func (u *Controller) getRepository(db *mgo.Database) *MongoAPISpecRepository {
	repo, err := NewMongoAppRepository(db)
	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	return repo
}
