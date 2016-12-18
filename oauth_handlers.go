package janus

import (
	"net/http"

	"github.com/hellofresh/janus/errors"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/request"
	"github.com/hellofresh/janus/response"
	"github.com/hellofresh/janus/router"
	"gopkg.in/mgo.v2"
)

type OAuthAPI struct{}

func (u *OAuthAPI) Get() http.HandlerFunc {
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
func (u *OAuthAPI) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		data, err := repo.FindByID(id)
		if data.ID == "" {
			panic(errors.New(http.StatusNotFound, "OAuth server not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (u *OAuthAPI) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))
		definition, err := repo.FindByID(id)
		if definition.ID == "" {
			panic(errors.New(http.StatusNotFound, "OAuth server not found"))
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		repo.Add(definition)
		response.JSON(w, http.StatusOK, definition)
	}
}

func (u *OAuthAPI) Post() http.HandlerFunc {
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

		response.JSON(w, http.StatusCreated, oauth)
	}
}

func (u *OAuthAPI) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := router.FromContext(r.Context()).ByName("id")
		repo := u.getRepository(u.getDatabase(r))

		err := repo.Remove(id)
		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (u *OAuthAPI) getDatabase(r *http.Request) *mgo.Database {
	db := r.Context().Value(middleware.ContextKeyDatabase)

	if nil == db {
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
