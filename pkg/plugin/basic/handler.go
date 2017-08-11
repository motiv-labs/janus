package basic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/render"
	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
)

// Controller is the api rest controller
type Controller struct {
	repo Repository
}

// NewController creates a new instance of Controller
func NewController(repo Repository) *Controller {
	return &Controller{repo}
}

// Get is the find all handler
func (c *Controller) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.FromContext(r.Context(), "datastore.user.FindAll")
		data, err := c.repo.FindAll()
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		render.JSON(w, http.StatusOK, data)
	}
}

// GetBy is the find by handler
func (c *Controller) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := router.URLParam(r, "username")
		span := opentracing.FromContext(r.Context(), "datastore.user.FindByUsername")
		data, err := c.repo.FindByUsername(username)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		render.JSON(w, http.StatusOK, data)
	}
}

// PutBy is the update handler
func (c *Controller) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		username := router.URLParam(r, "username")
		span := opentracing.FromContext(r.Context(), "datastore.user.FindByUsername")
		user, err := c.repo.FindByUsername(username)
		span.Finish()

		if user == nil {
			errors.Handler(w, ErrUserNotFound)
			return
		}

		if err != nil {
			errors.Handler(w, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			errors.Handler(w, err)
			return
		}

		span = opentracing.FromContext(r.Context(), "datastore.user.Add")
		err = c.repo.Add(user)
		span.Finish()

		if err != nil {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Post is the create handler
func (c *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &User{}

		err := json.NewDecoder(r.Body).Decode(user)
		if nil != err {
			errors.Handler(w, err)
			return
		}

		span := opentracing.FromContext(r.Context(), "datastore.users.Exists")
		_, err = c.repo.FindByUsername(user.Username)
		span.Finish()

		if err != ErrUserNotFound {
			log.WithError(err).Warn("An error occurrend when looking for an user")
			errors.Handler(w, ErrUserExists)
			return
		}

		span = opentracing.FromContext(r.Context(), "datastore.users.Add")
		err = c.repo.Add(user)
		span.Finish()

		if err != nil {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/credentials/basic_auth/%s", user.Username))
		w.WriteHeader(http.StatusCreated)
	}
}

// DeleteBy is the delete handler
func (c *Controller) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := router.URLParam(r, "username")

		span := opentracing.FromContext(r.Context(), "datastore.users.Remove")
		err := c.repo.Remove(username)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
