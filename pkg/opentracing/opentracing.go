package opentracing

import (
	"context"

	"net/http"

	log "github.com/Sirupsen/logrus"
	gcloudtracer "github.com/hellofresh/gcloud-opentracing"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/opentracing/appdash"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// CtxSpanID is used to store the SpanID in a request's context
	CtxSpanID = 0
)

// Build a tracer based on the configuration provided
func Build(config config.Tracing) (opentracing.Tracer, error) {
	if config.IsGoogleCloudEnabled() {
		tracer, err := gcloudtracer.NewTracer(
			context.Background(),
			gcloudtracer.WithLogger(log.StandardLogger()),
			gcloudtracer.WithProject(config.GoogleCloudTracing.ProjectID),
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
	} else if config.IsAppdashEnabled() {
		server, err := appdash.NewServer(config.AppdashTracing.DSN, config.AppdashTracing.URL)
		if err != nil {
			return nil, err
		}

		server.Listen()
		return server.GetTracer(), nil
	}

	return nil, nil
}

func FromContext(ctx context.Context, name string) opentracing.Span {
	parentSpan := ctx.Value(CtxSpanID).(opentracing.Span)
	return opentracing.StartSpan(name, opentracing.ChildOf(parentSpan.Context()))
}

func ToContext(r *http.Request, span opentracing.Span) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), CtxSpanID, span))
}
