package api

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/response"
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
			panic(err.Error())
		}

		response.JSON(w, http.StatusOK, data)
	}
}

// GetBy is the find by handler
func (c *Controller) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.URLParam(r, "name")
		span := opentracing.FromContext(r.Context(), "datastore.FindByName")
		data, err := c.repo.FindByName(name)
		span.Finish()

		if data == nil {
			panic(ErrAPIDefinitionNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusOK, data)
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
			panic(ErrAPIDefinitionNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(definition)
		c.dispatch(notifier.NoticeAPIUpdated)
		span.Finish()

		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		response.JSON(w, http.StatusOK, nil)
	}
}

// Post is the create handler
func (c *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		definition := NewDefinition()

		err := request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		span := opentracing.FromContext(r.Context(), "datastore.Exists")
		exists, err := c.repo.Exists(definition)
		span.Finish()

		if nil != err && err != mgo.ErrNotFound {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		if exists {
			panic(err)
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(definition)
		c.dispatch(notifier.NoticeAPIAdded)
		span.Finish()

		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		w.Header().Add("Location", "")
		response.JSON(w, http.StatusCreated, nil)
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
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		c.dispatch(notifier.NoticeAPIRemoved)
		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (c *Controller) dispatch(cmd notifier.NotificationCommand) {
	if c.notifier != nil {
		c.notifier.Notify(notifier.Notification{Command: cmd})
	}
}
