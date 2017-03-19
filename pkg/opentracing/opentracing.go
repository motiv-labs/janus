package opentracing

import (
	"context"

	log "github.com/Sirupsen/logrus"
	gcloudtracer "github.com/hellofresh/gcloud-opentracing"
	"github.com/hellofresh/janus/pkg/config"
	opentracing "github.com/opentracing/opentracing-go"
)

// Build a tracer based on the configuration provided
func Build(config config.Tracing) (opentracing.Tracer, error) {
	if config.IsGoogleCloudEnabled() {
		tracer, err := gcloudtracer.NewTracer(
			context.Background(),
			gcloudtracer.WithLogger(log.StandardLogger()),
			gcloudtracer.WithProject(config.ProjectID),
			gcloudtracer.WithJWTCredentials(gcloudtracer.JWTCredentials{
				Email:        config.GoogleCloudTracing.Email,
				PrivateKey:   []byte(config.GoogleCloudTracing.PrivateKey),
				PrivateKeyID: config.GoogleCloudTracing.PrivateKeyID,
			}),
		)
		if err != nil {
			return nil, err
		}

		return tracer, nil
	}

	return nil, nil
}
