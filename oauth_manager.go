package main

import (
	"gopkg.in/redis.v3"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"
)

type OAuthManager struct {
	storage *redis.Client
}

func (o OAuthManager) KeyExists(accessToken string) bool {
	log.Debugf("Searching for key %s", accessToken)
	return o.storage.Exists(accessToken).Val()
}

func (o OAuthManager) Set(accessToken string, session SessionState, resetTTLTo int64) error {
	value, _ := json.Marshal(session)
	expireDuration := time.Duration(resetTTLTo) * time.Second

	log.Debugf("Storing key %s for %d seconds", accessToken, expireDuration)
	go o.storage.Set(accessToken, string(value), expireDuration)
	return nil
}

func (o OAuthManager) IsKeyAuthorised(accessToken string) (SessionState, bool) {
	var newSession SessionState
	jsonKeyVal := o.storage.Get(accessToken).Val()

	if marshalErr := json.Unmarshal([]byte(jsonKeyVal), &newSession); marshalErr != nil {
		log.Error("Couldn't unmarshal session object")
		log.Error(marshalErr)
		return newSession, false
	}

	return newSession, true
}
