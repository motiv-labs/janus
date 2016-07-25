package main

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
)

type OAuthHandler struct {
	spec *APISpec
}

func (h OAuthHandler) Serve(c *iris.Context) {
	var newSession SessionState

	body := string(c.Response.Body())

	if marshalErr := json.Unmarshal([]byte(body), &newSession); marshalErr != nil {
		log.Error("Couldn't unmarshal session object")
		log.Panic(marshalErr)
	}

	h.spec.OAuthManager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)
}
