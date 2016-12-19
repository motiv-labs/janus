package oauth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/session"
)

// AwareTransport is a RoundTripper that is aware of the access tokens that come back from the
// authentication server. After retriving the token, we delagete the storage of it for the
// oauth manager
type AwareTransport struct {
	http.RoundTripper
	manager *Manager
}

// NewAwareTransport creates a new instace of AwareTransport
func NewAwareTransport(roundTripper http.RoundTripper, manager *Manager) *AwareTransport {
	return &AwareTransport{roundTripper, manager}
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
func (t *AwareTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if nil != err {
		return resp, err
	}

	if resp.StatusCode < 300 && resp.Body != nil {
		var newSession session.SessionState

		//This is useful for the middlewares
		var bodyBytes []byte

		defer resp.Body.Close()
		bodyBytes, _ = ioutil.ReadAll(resp.Body)

		// Use the content
		log.Info("Setting body")

		if marshalErr := json.Unmarshal(bodyBytes, &newSession); marshalErr == nil {
			t.manager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)
		}

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return resp, err
}
