package health

import (
	"net/http"

	"github.com/containous/traefik/log"
	"github.com/hellofresh/janus/pkg/api"
)

const (
	// Healthy represents that a service is fully operational
	Healthy Status = "healthy"
	// PartiallyHealthy represents that a service can still reply for requests,
	// but it's not fully operational
	PartiallyHealthy = "partially_healthy"
	// Unhealthy represents that a service is down
	Unhealthy = "unhealthy"
)

// Status represents the health status of a service
type Status string

// Checker represents a list of definitions to be checked
type Checker struct {
	Definitions []*api.Definition
	Response    chan *Response
	Err         chan error
}

// Response represents a sigle check response of a service
type Response struct {
	Name    string `json:"name"`
	Check   Check  `json:"check"`
	message map[string]string
}

// Check represents one check of a service
type Check struct {
	Status  Status `json:"status"`
	Subject string `json:"subject"`
}

// New creates a new instance of Checker
func New(definitions []*api.Definition) *Checker {
	return &Checker{definitions, make(chan *Response), make(chan error)}
}

// Check runs checks on all API definitions
func (c *Checker) Check() []*Response {
	urls := make(map[string]string)
	for _, definition := range c.Definitions {
		urls[definition.Name] = definition.HealthCheck.URL
	}

	for name, url := range urls {
		go doRequest(name, url, c)
	}

	var checks []*Response

	for range urls {
		select {
		case resp := <-c.Response:
			checks = append(checks, resp)
		case err := <-c.Err:
			log.WithError(err).Error("Check went wrong")
		}
	}

	return checks
}

func doRequest(name string, url string, checker *Checker) {
	res, err := http.Get(url)
	if err != nil {
		checker.Err <- err
	}

	response := &Response{Name: name, Check: Check{}}
	if res.StatusCode > http.StatusInternalServerError {
		response.Check.Status = Unhealthy
	} else if res.StatusCode > http.StatusBadRequest {
		response.Check.Status = PartiallyHealthy
	} else {
		response.Check.Status = Healthy
	}

	checker.Response <- response
}
