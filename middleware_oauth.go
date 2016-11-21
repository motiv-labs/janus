package main

import (
	"encoding/json"

	"errors"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type OAuthMiddleware struct {
	oauthManager *OAuthManager
	oauthSpec    *OAuthSpec
}

func (m *OAuthMiddleware) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	var newSession SessionState
	newSession.OAuthServerID = m.oauthSpec.ID

	log.WithFields(log.Fields{
		"req": req,
	}).Info("Getting body")

	data, exists := c.Get("body")

	if false == exists {
		return errors.New("Body from the proxy doesn't exists"), http.StatusInternalServerError
	}

	body := data.([]byte)

	if marshalErr := json.Unmarshal(body, &newSession); marshalErr != nil {
		return marshalErr, http.StatusInternalServerError
	}

	m.oauthManager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)

	return nil, http.StatusOK
}
