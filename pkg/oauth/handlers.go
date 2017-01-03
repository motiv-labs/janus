package oauth

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/loader"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/janus/pkg/router"
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

// GetBy gets an oauth server by its id
func (u *Controller) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			panic(ErrOauthServerNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (u *Controller) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))
		oauth, err := repo.FindByID(id)
		if oauth.ID == "" {
			panic(ErrOauthServerNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = request.BindJSON(r, oauth)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = repo.Add(oauth)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		u.changeTracker.Change()
		response.JSON(w, http.StatusOK, oauth)
	}
}

func (u *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := u.getRepository(u.getDatabase(r))
		var oauth OAuth

		err := request.BindJSON(r, &oauth)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = repo.Add(&oauth)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		u.changeTracker.Change()
		response.JSON(w, http.StatusCreated, oauth)
	}
}

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
func (u *Controller) getRepository(db *mgo.Database) Repository {
	repo, err := NewMongoRepository(db)
	if err != nil {
		panic(errors.New(http.StatusInternalServerError, err.Error()))
	}

	return repo
}
