package bucket

import (
	"strings"

	"github.com/fiam/gounidecode/unidecode"
)

const (
	totalBucket = "total"

	suffixStatusOk   = "ok"
	suffixStatusFail = "fail"

	prefixUnicode = "-u-"

	// SectionRequest is default section name for tracking HTTP requests
	SectionRequest = "request"

	// MetricEmptyPlaceholder is a string placeholder for empty (unset) sections of operation
	MetricEmptyPlaceholder = "-"
	// MetricIDPlaceholder is a string placeholder for ID section of operation if any
	MetricIDPlaceholder = "-id-"
)

var operationsStatus = map[bool]string{true: suffixStatusOk, false: suffixStatusFail}

// Bucket is an interface for building metric names for operations
type Bucket interface {
	// Metric builds simple metric name in the form "<section>.<operation-0>.<operation-1>.<operation-2>"
	Metric() string

	// MetricWithSuffix builds metric name with success suffix in the form "<section>-ok|fail.<operation-0>.<operation-1>.<operation-2>"
	MetricWithSuffix() string

	// MetricTotal builds simple total metric name in the form total.<section>"
	MetricTotal() string

	// MetricTotalWithSuffix builds total metric name with success suffix in the form total-ok|fail.<section>"
	MetricTotalWithSuffix() string
}

// SanitizeMetricName modifies metric name to work well with statsd
func SanitizeMetricName(metric string, uniDecode bool) string {
	if metric == "" {
		return MetricEmptyPlaceholder
	}

	if uniDecode {
		asciiMetric := unidecode.Unidecode(metric)
		if asciiMetric != metric {
			metric = prefixUnicode + asciiMetric
		}
	}

	return strings.Replace(
		// Double underscores
		strings.Replace(metric, "_", "__", -1),
		// and replace dots with single underscore
		".",
		"_",
		-1,
	)
}
