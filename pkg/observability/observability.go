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

// PromExporter is the prometheus exporter containing HTTP handler for "/metrics"
var PromExporter *prometheus.Exporter

// Tags
var (
	KeyHostname, _               = tag.NewKey("hostname")
	KeyJWTValidationErrorType, _ = tag.NewKey("error")
)

// Metrics
var (
	MRequestsByHostname         = stats.Int64("opencensus.io/http/proxy/request_total_by_host", "Number of proxied requests by target hostname", dimensionless)
	MJWTManagerValidationErrors = stats.Int64("opencensus.io/plugin/jwt_manager/validation_error_total", "Number of validation errors by error type", dimensionless)
	MOAuth2MissingHeader        = stats.Int64("opencensus.io/plugin/oauth2/missing_header_total", "Number of failed oauth2 authentication due to missing header", dimensionless)
	MOAuth2MalformedHeader      = stats.Int64("opencensus.io/plugin/oauth2/malformed_header_total", "Number of failed oauth2 authentication due to malformed bearer header", dimensionless)
	MOAuth2Authorized           = stats.Int64("opencensus.io/plugin/oauth2/authorized_request_total", "Number of successful and authorized oauth2 authentication", dimensionless)
	MOAuth2Unauthorized         = stats.Int64("opencensus.io/plugin/oauth2/unauthorized_request_total", "Number of successful but unauthorized oauth2 authentication", dimensionless)
)

// AllViews aggregates the metrics
var AllViews = []*view.View{
	{
		Name:        "opencensus.io/http/proxy/request_total_by_host",
		TagKeys:     []tag.Key{KeyHostname},
		Measure:     MRequestsByHostname,
		Aggregation: view.Count(),
	},
	{
		Name:        "opencensus.io/plugin/jwt_manager/validation_error_total",
		TagKeys:     []tag.Key{KeyJWTValidationErrorType},
		Measure:     MJWTManagerValidationErrors,
		Aggregation: view.Count(),
	},
	{
		Name:        "opencensus.io/plugin/oauth2/missing_header_total",
		Measure:     MOAuth2MissingHeader,
		Aggregation: view.Count(),
	},
	{
		Name:        "opencensus.io/plugin/oauth2/malformed_header_total",
		Measure:     MOAuth2MalformedHeader,
		Aggregation: view.Count(),
	},
	{
		Name:        "opencensus.io/plugin/oauth2/authorized_request_total",
		Measure:     MOAuth2Authorized,
		Aggregation: view.Count(),
	},
	{
		Name:        "opencensus.io/plugin/oauth2/unauthorized_request_total",
		Measure:     MOAuth2Unauthorized,
		Aggregation: view.Count(),
	},
}
