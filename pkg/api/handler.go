package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/render"
	"github.com/hellofresh/janus/pkg/router"
)

// Controller is the api rest controller
type Controller struct {
	repo     Repository
	notifier notifier.Notifier
}

// NewController creates a new instance of Controller
func NewController(repo Repository, notifier notifier.Notifier) *Controller {
	return &Controller{repo, notifier}
}

// Get is the find all handler
func (c *Controller) Get() http.HandlerFunc {
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
func (c *Controller) GetBy() http.HandlerFunc {
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
func (c *Controller) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		name := router.URLParam(r, "name")
		span := opentracing.FromContext(r.Context(), "datastore.FindByName")
		definition, err := c.repo.FindByName(name)
		span.Finish()

		if definition == nil {
			errors.Handler(w, ErrAPIDefinitionNotFound)
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

		if err != nil && err != ErrAPIDefinitionNotFound {
			errors.Handler(w, err)
			return
		}

		if nil != existingPathDefinition && existingPathDefinition.Name != definition.Name {
			errors.Handler(w, ErrAPIListenPathExists)
			return
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(definition)
		span.Finish()

		if err != nil {
			errors.Handler(w, errors.New(http.StatusBadRequest, err.Error()))
			return
		}

		c.dispatch(notifier.NoticeAPIUpdated)
		w.WriteHeader(http.StatusOK)
	}
}

// Post is the create handler
func (c *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		definition := NewDefinition()

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
		c.dispatch(notifier.NoticeAPIAdded)
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
func (c *Controller) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")

		span := opentracing.FromContext(r.Context(), "datastore.Remove")
		err := c.repo.Remove(name)
		span.Finish()

		if err != nil {
			errors.Handler(w, err)
			return
		}

		c.dispatch(notifier.NoticeAPIRemoved)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *Controller) dispatch(cmd notifier.NotificationCommand) {
	if c.notifier != nil {
		c.notifier.Notify(notifier.Notification{Command: cmd})
	}
}
