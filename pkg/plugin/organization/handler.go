package organization

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/render"
	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
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
		_, span := trace.StartSpan(r.Context(), "repo.FindAll")
		data, err := c.repo.FindAll()
		span.End()

		if err != nil {
			errors.Handler(w, r, err)
			return
		}

		render.JSON(w, http.StatusOK, data)
	}
}

// Show is the find by handler
func (c *Handler) Show() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := router.URLParam(r, "username")
		_, span := trace.StartSpan(r.Context(), "repo.Show")
		data, err := c.repo.FindByUsername(username)
		span.End()

		if err != nil {
			errors.Handler(w, r, err)
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
		_, span := trace.StartSpan(r.Context(), "repo.FindByUsername")
		user, err := c.repo.FindByUsername(username)
		span.End()

		if user == nil {
			errors.Handler(w, r, ErrUserNotFound)
			return
		}

		if err != nil {
			errors.Handler(w, r, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			errors.Handler(w, r, err)
			return
		}

		_, span = trace.StartSpan(r.Context(), "repo.Add")
		err = c.repo.Add(user)
		span.End()

		if err != nil {
			errors.Handler(w, r, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Create is the create handler
func (c *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organization := &Organization{}

		err := json.NewDecoder(r.Body).Decode(organization)
		if nil != err {
			errors.Handler(w, r, err)
			return
		}

		if organization.Organization == "" {
			err = errors.New(http.StatusBadRequest, "Invalid request body")
			log.WithError(err).Error("No organization provided. Must provide an organization.")
			errors.Handler(w, r, err)
			return
		}

		_, span := trace.StartSpan(r.Context(), "repo.FindByUsername")
		_, err = c.repo.FindByUsername(organization.Username)
		span.End()

		if err != ErrUserNotFound {
			log.WithError(err).Warn("An error occurred when looking for an organization")
			errors.Handler(w, r, ErrUserExists)
			return
		}

		_, span = trace.StartSpan(r.Context(), "repo.Add")
		err = c.repo.Add(organization)
		span.End()

		if err != nil {
			errors.Handler(w, r, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/credentials/organization/%s", organization.Username))
		w.WriteHeader(http.StatusCreated)
	}
}

// Delete is the delete handler
func (c *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := router.URLParam(r, "username")

		_, span := trace.StartSpan(r.Context(), "repo.Remove")
		err := c.repo.Remove(username)
		span.End()

		if err != nil {
			errors.Handler(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
