package opentracing

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/hellofresh/gcloud-opentracing"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/opentracing/appdash"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

type noopCloser struct{}

func (n noopCloser) Close() error { return nil }

// Build a tracer based on the configuration provided
func Build(config config.Tracing) (opentracing.Tracer, io.Closer, error) {
	switch config.Tracer {
	case "gcloud":
		log.Debug("Using google cloud platform (stackdriver trace) as tracing system")
		tracer, err := buildGCloud(config.GoogleCloudTracing)
		return tracer, noopCloser{}, err
	case "appdash":
		tracer, err := buildAppdash(config.AppdashTracing)
		return tracer, noopCloser{}, err
	case "jaeger":
		return buildJaeger(config.JaegerTracing)
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

func buildGCloud(config config.GoogleCloudTracing) (opentracing.Tracer, error) {
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
	if err != nil {
		return nil, err
	}

	return tracer, nil
}

func buildAppdash(config config.AppdashTracing) (opentracing.Tracer, error) {
	server := appdash.NewServer(config.DSN, config.URL)

	appdashFields := log.WithFields(log.Fields{
		"appdash_dsn":    config.DSN,
		"appdash_web_ui": config.URL,
	})

	if config.URL != "" {
		appdashFields.Debug("Using local appdash server as tracing system")
		err := server.Listen()
		if err != nil {
			return nil, err
		}
	} else {
		appdashFields.Debug("Using remote appdash server as tracing system")
	}

	return server.GetTracer(), nil
}

func buildJaeger(c config.JaegerTracing) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  c.DSN,
		},
	}

	return cfg.New(
		"janus",
		jaegercfg.Logger(jaegerLoggerAdapter{log.StandardLogger()}),
		jaegercfg.Metrics(metrics.NullFactory),
	)
}
