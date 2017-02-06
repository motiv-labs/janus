package stats

import (
	"net/http"

	statsd "gopkg.in/alexcesaro/statsd.v2"
)

type StatsClient struct {
	client *statsd.Client
}

func NewStatsClient(client *statsd.Client) *StatsClient {
	return &StatsClient{client}
}

func (f *StatsClient) BuildTimeTracker() *TimeTracker {
	return NewTimeTracker(f.client)
}

func (f *StatsClient) TrackRequest(r *http.Request, tt *TimeTracker, success bool) {
	b := RequestBucket(r)
	i := NewIncrementer(f.client)

	tt.Finish(b)
	i.Increment(b)
	i.Increment(totalRequestBucket)

	i.Increment(TotalRequestsWithSuffixBucket(success))
	i.Increment(RequestsWithSuffixBucket(r, success))
}

func (f *StatsClient) TrackRoundTrip(r *http.Request, tt *TimeTracker, success bool) {
	b := RoundTripBucket(r, success)
	i := NewIncrementer(f.client)

	tt.Finish(b)
	i.Increment(b)
	i.Increment(totalRoundTripBucket)
	i.Increment(RoundTripSuffixBucket(success))
}
