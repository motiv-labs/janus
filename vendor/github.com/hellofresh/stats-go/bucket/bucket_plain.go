package bucket

import (
	"strings"
)

// MetricOperation is a list of metric operations to use for metric
type MetricOperation [3]string

// Plain struct in an implementation of Bucket interface that produces metric names for given section and operation
type Plain struct {
	section   string
	operation string
	success   bool
}

// NewPlain builds and returns new Plain instance
func NewPlain(section string, operation MetricOperation, success, uniDecode bool) *Plain {
	operationSanitized := make([]string, cap(operation))
	for k, v := range operation {
		operationSanitized[k] = SanitizeMetricName(v, uniDecode)
	}
	return &Plain{SanitizeMetricName(section, uniDecode), strings.Join(operationSanitized, "."), success}
}

// Metric builds simple metric name in the form:
//  <section>.<operation-0>.<operation-1>.<operation-2>
func (b *Plain) Metric() string {
	return b.section + "." + b.operation
}

// MetricWithSuffix builds metric name with success suffix in the form:
//  <section>-ok|fail.<operation-0>.<operation-1>.<operation-2>
func (b *Plain) MetricWithSuffix() string {
	return b.section + "-" + operationsStatus[b.success] + "." + b.operation
}

// MetricTotal builds simple total metric name in the form:
//  total.<section>
func (b *Plain) MetricTotal() string {
	return totalBucket + "." + b.section
}

// MetricTotalWithSuffix builds total metric name with success suffix in the form
//  total-ok|fail.<section>
func (b *Plain) MetricTotalWithSuffix() string {
	return totalBucket + "." + b.section + "-" + operationsStatus[b.success]
}
