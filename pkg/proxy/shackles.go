package proxy

import (
	"net/http"
	"net/http/httputil"

	stats "github.com/hellofresh/stats-go"
)

const (
	statsSectionRoundTrip = "round"
)

// InLink interface for inbound request plugins
type InLink func(req *http.Request) (*http.Request, error)

// OutLink interface for outbound request plugins
type OutLink func(req *http.Request, res *http.Response) (*http.Response, error)

// InChain typed array for inbound plugin sequence
type InChain []InLink

// OutChain typed array for outbound plugin sequence
type OutChain []OutLink

// Shackles construct holding plugin sequences
type Shackles struct {
	statsClient stats.Client
	inbound     InChain
	outbound    OutChain
}

// NewInChain variadic constructor for inbound plugin sequence
func NewInChain(in ...InLink) InChain {
	return append(([]InLink)(nil), in...)
}

// NewOutChain variadic constructor for outbound plugin sequence
func NewOutChain(out ...OutLink) OutChain {
	return append(([]OutLink)(nil), out...)
}

// RoundTrip provides the Transport.RoundTrip function to handle requests the proxy receives
func (s *Shackles) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	timing := s.statsClient.BuildTimeTracker().Start()

	// apply inbound request plugins (if any)
	req, err = s.applyInboundLinks(req)
	if err != nil {
		s.statsClient.SetHTTPRequestSection(statsSectionRoundTrip).
			TrackRequest(req, timing, false).
			ResetHTTPRequestSection()
		return nil, err
	}

	// use default RoundTrip function handle the actual request/response
	resp, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
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
func (s *Shackles) applyOutboundLinks(req *http.Request, resp *http.Response) (mod *http.Response, err error) {
	mod = resp

	for o := range s.outbound {
		mod, err = s.outbound[o](req, mod)
		if err != nil {
			return nil, err
		}
	}

	return mod, nil
}

// applies any inbound request plugins to the given request
func (s *Shackles) applyInboundLinks(req *http.Request) (mod *http.Request, err error) {
	mod = req

	for i := range s.inbound {
		mod, err = s.inbound[i](mod)
		if err != nil {
			return nil, err
		}
	}

	return mod, nil
}
