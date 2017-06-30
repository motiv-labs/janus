package checker

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hellofresh/health-go"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/router"
	log "github.com/sirupsen/logrus"
)

// Register registers the health checks for valid API definitions
func Register(r router.Router, repo api.Repository) error {
	definitions, err := repo.FindValidAPIHealthChecks()
	if err != nil {
		log.WithError(err).Error("Error fetching API definitions for health check registering")
		return err
	}

	for _, definition := range definitions {
		log.Debugf("%s health check registered", definition.Name)
		health.Register(health.Config{
			Name:      definition.Name,
			Timeout:   time.Second * time.Duration(definition.HealthCheck.Timeout),
			SkipOnErr: true,
			Check:     check(definition),
		})
	}

	r.GET("/status", health.HandlerFunc)
	return nil
}

func check(definition *api.Definition) func() error {
	return func() error {
		req, err := http.NewRequest(http.MethodGet, definition.HealthCheck.URL, nil)
		if err != nil {
			log.WithError(err).Error("Creating the request for the health check failed")
			return err
		}

		// Inform to close the connection after the transaction is complete
		req.Header.Set("Connection", "close")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.WithError(err).Error("Making the request for the health check failed")
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= http.StatusInternalServerError {
			return fmt.Errorf("%s is not available at the moment", definition.Name)
		}

		return nil
	}
}
