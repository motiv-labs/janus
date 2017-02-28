package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/stats"
)

// AwareTransport is a RoundTripper that is aware of the access tokens that come back from the
// authentication server. After retrieving the token, we delagete the storage of it for the
// oauth manager
type AwareTransport struct {
	manager          Manager
	oAuthServersRepo *MongoRepository
	statsClient      *stats.StatsClient
}

// NewAwareTransport creates a new instance of AwareTransport
func NewAwareTransport(manager Manager, oAuthServersRepo *MongoRepository, statsClient *stats.StatsClient) *AwareTransport {
	return &AwareTransport{manager, oAuthServersRepo, statsClient}
}

// GetRoundTripper returns initialized RoundTripper insnace
func (at *AwareTransport) GetRoundTripper(roundTripper http.RoundTripper) http.RoundTripper {
	return &RoundTripper{roundTripper, at.manager, at.oAuthServersRepo, at.statsClient}
}

type RoundTripper struct {
	RoundTripper     http.RoundTripper
	manager          Manager
	oAuthServersRepo *MongoRepository
	statsClient      *stats.StatsClient
}

// RoundTrip executes a single HTTP transaction, returning
// a Response for the provided Request.
//
// RoundTrip should not attempt to interpret the response. In
// particular, RoundTrip must return err == nil if it obtained
// a response, regardless of the response's HTTP status code.
// A non-nil err should be reserved for failure to obtain a
// response. Similarly, RoundTrip should not attempt to
// handle higher-level protocol details such as redirects,
// authentication, or cookies.
//
// RoundTrip should not modify the request, except for
// consuming and closing the Request's Body.
//
// RoundTrip must always close the body, including on errors,
// but depending on the implementation may do so in a separate
// goroutine even after RoundTrip returns. This means that
// callers wanting to reuse the body for subsequent requests
// must arrange to wait for the Close call before doing so.
//
// The Request's URL and Header fields must be initialized.
func (t *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	timing := t.statsClient.BuildTimeTracker()
	timing.Start()
	resp, err = t.RoundTripper.RoundTrip(req)

	if nil != err {
		t.statsClient.TrackRoundTrip(req, timing, false)
		log.Error("Response from the server was an error", err)
		return resp, err
	}

	statusCodeSuccess := resp.StatusCode < http.StatusInternalServerError
	t.statsClient.TrackRoundTrip(req, timing, statusCodeSuccess)

	if resp.StatusCode < http.StatusMultipleChoices && resp.Body != nil {
		var newSession session.SessionState

		//This is useful for the middlewares
		var bodyBytes []byte

		defer func(body io.Closer) {
			err := body.Close()
			if err != nil {
				log.Error(err)
			}
		}(resp.Body)
		bodyBytes, _ = ioutil.ReadAll(resp.Body)

		if marshalErr := json.Unmarshal(bodyBytes, &newSession); marshalErr == nil {
			if newSession.AccessToken != "" {
				log.Debug("Setting body in the oauth storage")
				t.manager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)
			}
		}

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return resp, err
}
