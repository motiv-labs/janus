package observability

import (
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

const (
	by            = "By"
	ms            = "ms"
	dimensionless = "1"
)

// Known exporters
const (
	AzureMonitor = "azure_monitor"
	Datadog      = "datadog"
	Jaeger       = "jaeger"
	Prometheus   = "prometheus"
	Stackdriver  = "stackdriver"
	Zipkin       = "zipkin"
)

// PrometheusExporter is the prometheus exporter containing HTTP handler for "/metrics"
var PrometheusExporter *prometheus.Exporter

// Tags
var (
	KeyJWTValidationErrorType, _ = tag.NewKey("error")
)

// Metrics
var (
	MJWTManagerValidationErrors = stats.Int64("plugin_jwt_manager_validation_error_total", "Number of validation errors by error type", dimensionless)
	MOAuth2MissingHeader        = stats.Int64("plugin_oauth2_missing_header_total", "Number of failed oauth2 authentication due to missing header", dimensionless)
	MOAuth2MalformedHeader      = stats.Int64("plugin_oauth2_malformed_header_total", "Number of failed oauth2 authentication due to malformed bearer header", dimensionless)
	MOAuth2Authorized           = stats.Int64("plugin_oauth2_authorized_request_total", "Number of successful and authorized oauth2 authentication", dimensionless)
	MOAuth2Unauthorized         = stats.Int64("plugin_oauth2_unauthorized_request_total", "Number of successful but unauthorized oauth2 authentication", dimensionless)
)

// AllViews aggregates the metrics
var AllViews = []*view.View{
	{
		Name:        "plugin_jwt_manager_validation_error_total",
		TagKeys:     []tag.Key{KeyJWTValidationErrorType},
		Measure:     MJWTManagerValidationErrors,
		Aggregation: view.Count(),
	},
	{
		Name:        "plugin_oauth2_missing_header_total",
		Measure:     MOAuth2MissingHeader,
		Aggregation: view.Count(),
	},
	{
		Name:        "plugin_oauth2_malformed_header_total",
		Measure:     MOAuth2MalformedHeader,
		Aggregation: view.Count(),
	},
	{
		Name:        "plugin_oauth2_authorized_request_total",
		Measure:     MOAuth2Authorized,
		Aggregation: view.Count(),
	},
	{
		Name:        "plugin_oauth2_unauthorized_request_total",
		Measure:     MOAuth2Unauthorized,
		Aggregation: view.Count(),
	},
}
