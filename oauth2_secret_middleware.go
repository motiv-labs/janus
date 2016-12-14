package janus

import (
	"net/http"

	"fmt"

	"encoding/base64"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/hellofresh/janus/errors"
)

type Oauth2Secret struct {
	oauthSpec *OAuthSpec
}

func (m *Oauth2Secret) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	log.Debug("Starting Oauth2Secret middleware")

	if "" != req.Header.Get("Authorization") {
		log.Debug("Authorization is set, proxying")
		return nil, http.StatusOK
	}

	clientID := req.URL.Query().Get("client_id")
	if "" == clientID {
		log.Debug("ClientID not set, proxying")
		return nil, http.StatusOK
	}

	clientSecret, exists := m.oauthSpec.Secrets[clientID]
	if false == exists {
		err := errors.ErrClientIdNotFound
		return err, err.Code
	}

	m.ChangeRequest(req, clientID, clientSecret)
	return nil, http.StatusOK
}

func (m *Oauth2Secret) ChangeRequest(req *http.Request, clientID, clientSecret string) {
	log.Debug("Modifying request")
	authString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authString))
}
