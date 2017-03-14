package api

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/request"
	"github.com/hellofresh/janus/pkg/response"
	"github.com/hellofresh/janus/pkg/router"
	"gopkg.in/mgo.v2"
)

// Controller is the api rest controller
type Controller struct {
	repo Repository
}

// NewController creates a new instance of Controller
func NewController(repo Repository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := c.repo.FindAll()
		if err != nil {
			panic(err.Error())
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (c *Controller) GetBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.FromContext(r.Context()).ByName("name")

		data, err := c.repo.FindByName(name)
		if data == nil {
			panic(ErrAPIDefinitionNotFound)
		}

		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (c *Controller) PutBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		name := router.FromContext(r.Context()).ByName("name")
		definition, err := c.repo.FindByName(name)
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

		err = c.repo.Add(definition)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		response.JSON(w, http.StatusOK, nil)
	}
}

func (c *Controller) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		definition := NewDefinition()

		err := request.BindJSON(r, definition)
		if nil != err {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		def, err := c.repo.FindByListenPath(definition.Proxy.ListenPath)
		if nil != err && err != mgo.ErrNotFound {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		if def != nil {
			panic(errors.ErrProxyExists)
		}

		err = c.repo.Add(definition)
		if nil != err {
			panic(errors.New(http.StatusBadRequest, err.Error()))
		}

		w.Header().Add("Location", "")
		response.JSON(w, http.StatusCreated, nil)
	}
}

func (c *Controller) DeleteBy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := router.FromContext(r.Context()).ByName("name")

		err := c.repo.Remove(name)
		if err != nil {
			panic(errors.New(http.StatusInternalServerError, err.Error()))
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}
