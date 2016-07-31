package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"errors"
)

type OAuthMiddleware struct {
	*Middleware
}

func (m OAuthMiddleware) ProcessRequest(req *http.Request, c *gin.Context) (error, int) {
	var newSession SessionState

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

	m.Spec.OAuthManager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)

	return nil, http.StatusOK
}
