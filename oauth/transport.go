package oauth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/session"
)

type AwareTransport struct {
	http.RoundTripper
	manager *Manager
}

func NewAwareTransport(roundTripper http.RoundTripper, manager *Manager) *AwareTransport {
	return &AwareTransport{roundTripper, manager}
}

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
