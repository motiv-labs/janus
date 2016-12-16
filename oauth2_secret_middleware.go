package janus

import (
	"net/http"

	"fmt"

	"encoding/base64"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/errors"
)

// Oauth2SecretMiddleware prevents requests to an API from exceeding a specified rate limit.
type Oauth2SecretMiddleware struct {
	oauthSpec *OAuthSpec
}

// ProcessRequest is the middleware method.
func (m *Oauth2SecretMiddleware) ProcessRequest(rw http.ResponseWriter, req *http.Request) (int, error) {
	log.Debug("Starting Oauth2Secret middleware")

	if "" != req.Header.Get("Authorization") {
		log.Debug("Authorization is set, proxying")
		return http.StatusOK, nil
	}

	clientID := req.URL.Query().Get("client_id")
	if "" == clientID {
		log.Debug("ClientID not set, proxying")
		return http.StatusOK, nil
	}

	clientSecret, exists := m.oauthSpec.Secrets[clientID]
	if false == exists {
		err := errors.ErrClientIdNotFound
		return err.Code, err
	}

	m.ChangeRequest(req, clientID, clientSecret)
	return http.StatusOK, nil
}

// ChangeRequest modifies the request to add the Authorization headers.
func (m *Oauth2SecretMiddleware) ChangeRequest(req *http.Request, clientID, clientSecret string) {
	log.Debug("Modifying request")
	authString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authString))
}
