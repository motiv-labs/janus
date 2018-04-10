package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/render"
	"github.com/hellofresh/janus/pkg/router"
)

// APIHandler is the api rest controller
type APIHandler struct {
	repo api.Repository
}

// NewAPIHandler creates a new instance of Controller
func NewAPIHandler(repo api.Repository) *APIHandler {
	return &APIHandler{repo}
}

// Get is the find all handler
func (c *APIHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := opentracing.FromContext(r.Context(), "datastore.FindAll")
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
func (c *APIHandler) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")
		span := opentracing.FromContext(r.Context(), "datastore.FindByName")
		data, err := c.repo.FindByName(name)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		render.JSON(w, http.StatusOK, data)
	}
}

// PutBy is the update handler
func (c *APIHandler) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		name := router.URLParam(r, "name")
		span := opentracing.FromContext(r.Context(), "datastore.FindByName")
		definition, err := c.repo.FindByName(name)
		span.Finish()

		if definition == nil {
			errors.Handler(w, api.ErrAPIDefinitionNotFound)
			return
		}

		if err != nil {
			errors.Handler(w, err)
			return
		}

		err = json.NewDecoder(r.Body).Decode(definition)
		if err != nil {
			errors.Handler(w, err)
			return
		}

		// avoid situation when trying to update existing definition with new path
		// that is already registered with another name
		span = opentracing.FromContext(r.Context(), "datastore.FindByListenPath")
		existingPathDefinition, err := c.repo.FindByListenPath(definition.Proxy.ListenPath)
		span.Finish()

		if err != nil && err != api.ErrAPIDefinitionNotFound {
			errors.Handler(w, err)
			return
		}

		if nil != existingPathDefinition && existingPathDefinition.Name != definition.Name {
			errors.Handler(w, api.ErrAPIListenPathExists)
			return
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(definition)
		span.Finish()

		if err != nil {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// Post is the create handler
func (c *APIHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		definition := api.NewDefinition()

		err := json.NewDecoder(r.Body).Decode(definition)
		if nil != err {
			errors.Handler(w, err)
			return
		}

		span := opentracing.FromContext(r.Context(), "datastore.Exists")
		exists, err := c.repo.Exists(definition)
		span.Finish()

		if err != nil || exists {
			errors.Handler(w, err)
			return
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(definition)
		span.Finish()

		if err != nil {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/apis/%s", definition.Name))
		w.WriteHeader(http.StatusCreated)
	}
}

// DeleteBy is the delete handler
func (c *APIHandler) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")

		span := opentracing.FromContext(r.Context(), "datastore.Remove")
		err := c.repo.Remove(name)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
