package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/stats"
	"github.com/hellofresh/janus/pkg/store"
)

// AwareTransport is a RoundTripper that is aware of the access tokens that come back from the
// authentication server. After retrieving the token, we delagete the storage of it for the
// oauth manager
type AwareTransport struct {
	statsClient *stats.StatsClient
	storage     store.Store
	repo        Repository
}

// NewAwareTransport creates a new instance of AwareTransport
func NewAwareTransport(statsClient *stats.StatsClient, storage store.Store, repo Repository) *AwareTransport {
	return &AwareTransport{statsClient, storage, repo}
}

// GetRoundTripper returns initialized RoundTripper insnace
func (at *AwareTransport) GetRoundTripper(roundTripper http.RoundTripper) http.RoundTripper {
	return &RoundTripper{roundTripper, at.statsClient, at.storage, at.repo}
}

type RoundTripper struct {
	RoundTripper http.RoundTripper
	statsClient  *stats.StatsClient
	storage      store.Store
	repo         Repository
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
		var newSession session.State

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
				tokenURL := url.URL{Scheme: req.URL.Scheme, Host: req.URL.Host, Path: req.URL.Path}
				log.WithField("token_url", tokenURL.String()).Debug("Looking for OAuth provider who issued the token")
				manager, oAuthServer, err := t.getManager(tokenURL)
				if err != nil {
					log.WithError(err).Error("Failed to find OAuth server by token URL", err)
				} else {
					newSession.OAuthServerSlug = oAuthServer.Slug
					log.Debug("Setting body in the oauth storage")
					manager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)
				}
			}
		}

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return resp, err
}

func (t *RoundTripper) getManager(url url.URL) (Manager, *OAuth, error) {
	oauthServer, err := t.repo.FindByTokenURL(url)
	if nil != err {
		return nil, nil, err
	}

	managerType, err := ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, nil, err
	}

	manager, err := NewManagerFactory(t.storage, oauthServer.TokenStrategy.Settings).Build(managerType)

	return manager, oauthServer, err
}
