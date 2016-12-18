package oauth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/session"
)

type OAuthAwareTransport struct {
	http.RoundTripper
	OauthManager *OAuthManager
}

func (t *OAuthAwareTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
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
			t.OauthManager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)
		}

		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return resp, err
}
