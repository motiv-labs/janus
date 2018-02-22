package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/hellofresh/janus/pkg/router"
	stats "github.com/hellofresh/stats-go"
	"github.com/pkg/errors"
)

const (
	statsSectionRoundTrip = "round"

	// DefaultIdleConnsPerHost the default value set for http.Transport.MaxIdleConnsPerHost.
	DefaultIdleConnsPerHost = 64

	// DefaultCloseIdleConnsPeriod the default period at which the idle connections are forcibly
	// closed.
	DefaultCloseIdleConnsPeriod = 20 * time.Second
)

// Params initialization options.
type Params struct {
	// StatsClient defines the stats client for tracing
	StatsClient stats.Client

	// When set, the proxy will skip the TLS verification on outgoing requests.
	InsecureSkipVerify bool

	// Same as net/http.Transport.MaxIdleConnsPerHost, but the default
	// is 64. This value supports scenarios with relatively few remote
	// hosts. When the routing table contains different hosts in the
	// range of hundreds, it is recommended to set this options to a
	// lower value.
	IdleConnectionsPerHost int

	// Defines the time period of how often the idle connections are
	// forcibly closed. The default is 12 seconds. When set to less than
	// 0, the proxy doesn't force closing the idle connections.
	CloseIdleConnsPeriod time.Duration

	// The Flush interval for copying upgraded connections
	FlushInterval time.Duration

	Outbound OutChain

	Breaker Breaker
}

// Breaker is the circuit breaker identifier configuration
type Breaker struct {
	ID string
}

// OutLink interface for outbound request plugins
type OutLink func(req *http.Request, res *http.Response) (*http.Response, error)

// OutChain typed array for outbound plugin sequence
type OutChain []OutLink

// InChain typed array for inbound plugin sequence, normally this is a middleware chain
type InChain []router.Constructor

// Transport construct holding plugin sequences
type Transport struct {
	statsClient stats.Client
	outbound    OutChain
	breaker     Breaker
}

// NewTransport creates a new instance of Transport
func NewTransport(statsClient stats.Client, outbound OutChain, breaker Breaker) *Transport {
	return &Transport{statsClient: statsClient, outbound: outbound, breaker: breaker}
}

// NewTransportWithParams creates a new instance of Transport with the given params
func NewTransportWithParams(o Params) *Transport {
	if o.IdleConnectionsPerHost <= 0 {
		o.IdleConnectionsPerHost = DefaultIdleConnsPerHost
	}

	if o.CloseIdleConnsPeriod == 0 {
		o.CloseIdleConnsPeriod = DefaultCloseIdleConnsPeriod
	}

	tr := http.DefaultTransport.(*http.Transport)
	tr.Proxy = http.ProxyFromEnvironment
	tr.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	tr.MaxIdleConns = 100
	tr.IdleConnTimeout = 90 * time.Second
	tr.TLSHandshakeTimeout = 10 * time.Second
	tr.ExpectContinueTimeout = 1 * time.Second
	tr.MaxIdleConnsPerHost = o.IdleConnectionsPerHost
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: o.InsecureSkipVerify}

	if o.CloseIdleConnsPeriod > 0 {
		go func() {
			for range time.After(o.CloseIdleConnsPeriod) {
				tr.CloseIdleConnections()
			}
		}()
	}

	return NewTransport(o.StatsClient, o.Outbound, o.Breaker)
}

// NewInChain variadic constructor for inbound plugin sequence
func NewInChain(in ...router.Constructor) InChain {
	return append(([]router.Constructor)(nil), in...)
}

// NewOutChain variadic constructor for outbound plugin sequence
func NewOutChain(out ...OutLink) OutChain {
	return append(([]OutLink)(nil), out...)
}

// RoundTrip provides the Transport.RoundTrip function to handle requests the proxy receives
func (s *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	timing := s.statsClient.BuildTimer().Start()

	err = hystrix.Do(s.breaker.ID, func() error {
		resp, err = http.DefaultTransport.RoundTrip(req)
		if err != nil {
			return err
		}

		// treat 500 and above as errors for the sake of the circuit breaker
		if resp.StatusCode >= http.StatusInternalServerError {
			return errors.New("internal server error")
		}

		return nil
	}, nil)

	if err != nil && resp == nil {
		s.statsClient.SetHTTPRequestSection(statsSectionRoundTrip).
			TrackRequest(req, timing, false).
			ResetHTTPRequestSection()
		return nil, err
	}

	// block until the entire body has been read
	_, err = httputil.DumpResponse(resp, true)
	if err != nil {
		s.statsClient.SetHTTPRequestSection(statsSectionRoundTrip).
			TrackRequest(req, timing, false).
			ResetHTTPRequestSection()
		return nil, err
	}

	// apply outbound response plugins (if any)
	resp, err = s.applyOutboundLinks(req, resp)
	if err != nil {
		s.statsClient.SetHTTPRequestSection(statsSectionRoundTrip).
			TrackRequest(req, timing, false).
			ResetHTTPRequestSection()
		return nil, err
	}

	statusCodeSuccess := resp.StatusCode < http.StatusInternalServerError
	s.statsClient.SetHTTPRequestSection(statsSectionRoundTrip).
		TrackRequest(req, timing, statusCodeSuccess).
		ResetHTTPRequestSection()

	// pass response back to client
	return resp, nil
}

// applies any outbound response plugins to the given response
func (s *Transport) applyOutboundLinks(req *http.Request, resp *http.Response) (mod *http.Response, err error) {
	mod = resp

	for o := range s.outbound {
		mod, err = s.outbound[o](req, mod)
		if err != nil {
			return nil, err
		}
	}

	return mod, nil
}
