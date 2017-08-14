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

// Handler is the api rest handlers
type Handler struct {
	repo Repository
}

// NewHandler creates a new instance of Handler
func NewHandler(repo Repository) *Handler {
	return &Handler{repo}
}

// Index is the find all handler
func (c *Handler) Index() http.HandlerFunc {
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

// Show is the find by handler
func (c *Handler) Show() http.HandlerFunc {
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

// Update is the update handler
func (c *Handler) Update() http.HandlerFunc {
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

// Create is the create handler
func (c *Handler) Create() http.HandlerFunc {
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

// Delete is the delete handler
func (c *Handler) Delete() http.HandlerFunc {
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
