package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/hellofresh/health-go/v3"
	"github.com/hellofresh/janus/pkg/api"
	log "github.com/sirupsen/logrus"
)

// NewOverviewHandler creates instance of all status checks handler
func NewOverviewHandler(cfg *api.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defs := findValidAPIHealthChecks(cfg.Definitions)

		log.WithField("len", len(defs)).Debug("Loading health check endpoints")
		health.Reset()

		for _, def := range defs {
			log.WithField("name", def.Name).Debug("Registering health check")
			health.Register(health.Config{
				Name:      def.Name,
				Timeout:   time.Second * time.Duration(def.HealthCheck.Timeout),
				SkipOnErr: true,
				Check:     check(def),
			})
		}

		health.HandlerFunc(w, r)
	}
}

// NewStatusHandler creates instance of single proxy status check handler
func NewStatusHandler(cfgs *api.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defs := findValidAPIHealthChecks(cfgs.Definitions)

		name := chi.URLParam(r, "name")
		for _, def := range defs {
			if name == def.Name {
				resp, err := doStatusRequest(def, false)
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

func doStatusRequest(def *api.Definition, closeBody bool) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, def.HealthCheck.URL, nil)
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

func check(def *api.Definition) func() error {
	return func() error {
		resp, err := doStatusRequest(def, true)
		if err != nil {
			return fmt.Errorf("%s health check endpoint %s is unreachable", def.Name, def.HealthCheck.URL)
		}

		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("%s is not available at the moment", def.Name)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("%s is partially available at the moment", def.Name)
		}

		return nil
	}
}

func findValidAPIHealthChecks(defs []*api.Definition) []*api.Definition {
	var validDefs []*api.Definition

	for _, def := range defs {
		if def.Active && def.HealthCheck.URL != "" {
			validDefs = append(validDefs, def)
		}
	}

	return validDefs
}
