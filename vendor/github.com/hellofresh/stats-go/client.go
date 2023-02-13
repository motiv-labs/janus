package stats

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/hellofresh/stats-go/client"
)

const (
	// StatsD is a dsn scheme value for statsd client
	statsD = "statsd"
	// Log is a dsn scheme value for log client
	log = "log"
	// Memory is a dsn scheme value for memory client
	memory = "memory"
	// Noop is a dsn scheme value for noop client
	noop = "noop"
)

// ErrUnknownClient is an error returned when trying to create stats client of unknown type
var ErrUnknownClient = errors.New("unknown stats client type")

// NewClient creates and builds new stats client instance by given dsn
func NewClient(dsn string) (client.Client, error) {
	dsnURL, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	// do not care about parse error, as default value is set to false that is fine for us
	unicode, _ := strconv.ParseBool(dsnURL.Query().Get("unicode"))

	switch dsnURL.Scheme {
	case statsD:
		return client.NewStatsD(dsnURL.Host, strings.Trim(dsnURL.Path, "/"), unicode)
	case log:
		return client.NewLog(unicode), nil
	case memory:
		return client.NewMemory(unicode), nil
	case noop:
		return client.NewNoop(), nil
	}

	return nil, ErrUnknownClient
}
