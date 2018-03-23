package opentracing

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/hellofresh/gcloud-opentracing"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

const (
	gcloudTracing = "googleCloud"
	jaegerTracing = "jaeger"
)

type noopCloser struct{}

func (n noopCloser) Close() error { return nil }

// Build a tracer based on the configuration provided
func Build(config config.Tracing) (opentracing.Tracer, io.Closer, error) {
	switch config.Provider {
	case gcloudTracing:
		log.Debug("Using google cloud platform (stackdriver trace) as tracing system")
		return buildGCloud(config.GoogleCloudTracing)
	case jaegerTracing:
		log.Debug("Using Jaeger as tracing system")
		return buildJaeger(config.ServiceName, config.JaegerTracing)
	default:
		log.Debug("No tracer selected")
		return &opentracing.NoopTracer{}, noopCloser{}, nil
	}
}

// FromContext creates a span from a context that contains a parent span
func FromContext(ctx context.Context, name string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, name)
	return span
}

// ToContext sets a span to a context
func ToContext(r *http.Request, span opentracing.Span) *http.Request {
	return r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
}

func buildGCloud(config config.GoogleCloudTracing) (opentracing.Tracer, io.Closer, error) {
	tracer, err := gcloudtracer.NewTracer(
		context.Background(),
		gcloudtracer.WithLogger(log.StandardLogger()),
		gcloudtracer.WithProject(config.ProjectID),
		gcloudtracer.WithJWTCredentials(gcloudtracer.JWTCredentials{
			Email:        config.Email,
			PrivateKey:   []byte(config.PrivateKey),
			PrivateKeyID: config.PrivateKeyID,
		}),
	)

	return tracer, noopCloser{}, err
}

func buildJaeger(componentName string, c config.JaegerTracing) (opentracing.Tracer, io.Closer, error) {
	bufferFLushInterval, err := time.ParseDuration(c.BufferFlushInterval)
	if err != nil {
		return nil, noopCloser{}, errors.Wrap(err, "could not parse buffer flush interval for jaeger")
	}

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  c.SamplingType,
			Param: c.SamplingParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            c.LogSpans,
			BufferFlushInterval: bufferFLushInterval,
			LocalAgentHostPort:  c.SamplingServerURL,
			QueueSize:           c.QueueSize,
		},
	}

	return cfg.New(
		componentName,
		jaegercfg.Logger(jaegerLoggerAdapter{log.StandardLogger()}),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
