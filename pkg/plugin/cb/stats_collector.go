package cb

import (
	"github.com/hellofresh/stats-go/bucket"
	"github.com/hellofresh/stats-go/timer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/hellofresh/stats-go/client"
)

// StatsCollector fulfills the metricCollector interface allowing users to ship circuit
// stats to a metric backend. To use users must call InitializeStatsCollector before
// circuits are started.
type StatsCollector struct {
	client                  client.Client
	section                 string
	name                    string
	circuitOpenPrefix       string
	attemptsPrefix          string
	errorsPrefix            string
	successesPrefix         string
	failuresPrefix          string
	rejectsPrefix           string
	shortCircuitsPrefix     string
	timeoutsPrefix          string
	fallbackSuccessesPrefix string
	fallbackFailuresPrefix  string
	canceledPrefix          string
	deadlinePrefix          string
	totalDurationPrefix     string
	runDurationPrefix       string
}

// NewCollectorRegistry returns a function to be registerd with metricCollector.Registry.Register(NewCollectorRegistry).
func NewCollectorRegistry(client client.Client) func(string) metricCollector.MetricCollector {
	return func(name string) metricCollector.MetricCollector {
		c, err := NewStatsCollector(name, client)
		if err != nil {
			log.WithError(err).Error("could not initialize the stats collector")
		}

		return c
	}
}

// NewStatsCollector creates a collector for a specific circuit. The
// prefix given to this circuit will be {config.Prefix}.{circuit_name}.{metric}.
// Circuits with "/" in their names will have them replaced with ".".
func NewStatsCollector(name string, client client.Client) (*StatsCollector, error) {
	if client == nil {
		return nil, errors.New("metrics client must be initialized before circuits are created")
	}

	return &StatsCollector{
		client:                  client,
		name:                    name,
		circuitOpenPrefix:       "circuitOpen",
		attemptsPrefix:          "attempts",
		errorsPrefix:            "errors",
		successesPrefix:         "successes",
		failuresPrefix:          "failures",
		rejectsPrefix:           "rejects",
		shortCircuitsPrefix:     "shortCircuits",
		timeoutsPrefix:          "timeouts",
		fallbackSuccessesPrefix: "fallbackSuccesses",
		fallbackFailuresPrefix:  "fallbackFailures",
		canceledPrefix:          "contextCanceled",
		deadlinePrefix:          "contextDeadlineExceeded",
		totalDurationPrefix:     "totalDuration",
		runDurationPrefix:       "runDuration",
	}, nil
}

// Update metrics
func (g *StatsCollector) Update(r metricCollector.MetricResult) {
	if r.Successes > 0 {
		g.client.TrackState(g.section, bucket.MetricOperation{g.name, g.circuitOpenPrefix}, 0)
	} else if r.ShortCircuits > 0 {
		g.client.TrackState(g.section, bucket.MetricOperation{g.name, g.circuitOpenPrefix}, 1)
	}

	g.incrementCounterMetric(g.attemptsPrefix, r.Attempts)
	g.incrementCounterMetric(g.errorsPrefix, r.Errors)
	g.incrementCounterMetric(g.successesPrefix, r.Successes)
	g.incrementCounterMetric(g.failuresPrefix, r.Failures)
	g.incrementCounterMetric(g.rejectsPrefix, r.Rejects)
	g.incrementCounterMetric(g.shortCircuitsPrefix, r.ShortCircuits)
	g.incrementCounterMetric(g.timeoutsPrefix, r.Timeouts)
	g.incrementCounterMetric(g.fallbackSuccessesPrefix, r.FallbackSuccesses)
	g.incrementCounterMetric(g.fallbackFailuresPrefix, r.FallbackFailures)
	g.incrementCounterMetric(g.canceledPrefix, r.ContextCanceled)
	g.incrementCounterMetric(g.deadlinePrefix, r.ContextDeadlineExceeded)

	g.client.TrackOperation(
		g.section,
		bucket.MetricOperation{g.name, g.totalDurationPrefix},
		timer.NewDuration(r.TotalDuration),
		r.Successes > 0,
	)

	g.client.TrackOperation(
		g.section,
		bucket.MetricOperation{g.name, g.runDurationPrefix},
		timer.NewDuration(r.RunDuration),
		r.Successes > 0,
	)
}

// Reset is a noop operation in this collector.
func (g *StatsCollector) Reset() {}

func (g *StatsCollector) incrementCounterMetric(prefix string, i float64) {
	if i == 0 {
		return
	}
	g.client.TrackMetricN(g.section, bucket.MetricOperation{g.name, prefix}, int(i))
}
