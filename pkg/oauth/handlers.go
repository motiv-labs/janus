package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
)

// Controller is the api rest controller
type Controller struct {
	repo      Repository
	publisher store.Publisher
}

// NewController creates a new instance of Controller
func NewController(repo Repository, publisher store.Publisher) *Controller {
	return &Controller{repo, publisher}
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

		if data.Name == "" {
			panic(ErrOauthServerNotFound)
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
		oauth, err := c.repo.FindByName(name)
		span.Finish()

		if oauth.Name == "" {
			panic(ErrOauthServerNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		err = request.BindJSON(r, oauth)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		span = opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(oauth)
		c.dispatch(oauth)
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
		var oauth OAuth

		err := request.BindJSON(r, &oauth)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		span := opentracing.FromContext(r.Context(), "datastore.Add")
		err = c.repo.Add(&oauth)
		c.dispatch(&oauth)
		span.Finish()

		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

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

		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (c *Controller) dispatch(server *OAuth) {
	if c.publisher != nil {
		raw, _ := json.Marshal(server)
		c.publisher.Publish("oauth_updates", raw)
	}
}
