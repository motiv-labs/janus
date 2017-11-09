package checker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/hellofresh/health-go"
	"github.com/hellofresh/janus/pkg/api"
	log "github.com/sirupsen/logrus"
)

// NewOverviewHandler creates instance of all status checks handler
func NewOverviewHandler(repo api.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		definitions, err := repo.FindValidAPIHealthChecks()
		if err != nil {
			log.WithError(err).Error("Error fetching API definitions for health check registering")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		log.WithField("len", len(definitions)).Debug("Loading health check endpoints")
		health.Reset()

		for _, definition := range definitions {
			if definition.HealthCheck.URL != "" {
				log.WithField("name", definition.Name).Debug("Registering health check")
				health.Register(health.Config{
					Name:      definition.Name,
					Timeout:   time.Second * time.Duration(definition.HealthCheck.Timeout),
					SkipOnErr: true,
					Check:     check(definition),
				})
			}
		}

		health.HandlerFunc(w, r)
	}
}

// NewStatusHandler creates instance of single proxy status check handler
func NewStatusHandler(repo api.Repository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		definitions, err := repo.FindValidAPIHealthChecks()
		if err != nil {
			log.WithError(err).Error("Error fetching API definitions for health check registering")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		name := chi.URLParam(r, "name")
		for _, definition := range definitions {
			if name == definition.Name {
				resp, err := doStatusRequest(definition, false)
				if err != nil {
					log.WithField("name", name).WithError(err).Error("Error requesting service health status")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}

				body, err := ioutil.ReadAll(resp.Body)
				if closeErr := resp.Body.Close(); closeErr != nil {
					log.WithField("name", name).WithError(closeErr).Error("Error closing health status body")
				}

				if err != nil {
					log.WithField("name", name).WithError(err).Error("Error reading health status body")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}

				w.WriteHeader(resp.StatusCode)
				w.Write(body)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Definition name is not found"))
	}
}

func doStatusRequest(definition *api.Definition, closeBody bool) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, definition.HealthCheck.URL, nil)
	if err != nil {
		log.WithError(err).Error("Creating the request for the health check failed")
		return nil, err
	}

	// Inform to close the connection after the transaction is complete
	req.Header.Set("Connection", "close")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Making the request for the health check failed")
		return resp, err
	}

	if closeBody {
		defer resp.Body.Close()
	}

	return resp, err
}

func check(definition *api.Definition) func() error {
	return func() error {
		resp, err := doStatusRequest(definition, true)
		if err != nil {
			return fmt.Errorf("%s health check endpoint %s is unreachable", definition.Name, definition.HealthCheck.URL)
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("%s is not available at the moment", definition.Name)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("%s is partially available at the moment", definition.Name)
		}

		return nil
	}
}
