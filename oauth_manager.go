package main

import (
	"gopkg.in/redis.v3"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type OAuthManager struct {
	storage *redis.Client
}

func NewOAuthManager(client *redis.Client) *OAuthManager {
	return &OAuthManager{client}
}

func (o OAuthManager) CheckSessionAndIdentityForValidKey(accessToken string) (SessionState, bool) {
	var newSession SessionState

	//Checks if the key is present on the cache and if it didn't expire yet
	if o.storage.Exists(accessToken).Val() {
		jsonKeyVal := o.storage.Get(accessToken).String()

		if marshalErr := json.Unmarshal([]byte(jsonKeyVal), &newSession); marshalErr != nil {
			log.Error("Couldn't unmarshal session object")
			log.Error(marshalErr)
			return newSession, false
		}

		return newSession, true
	}


	//if its not in the cache
	//make a request to /info on the auth service
	//if the key is invalid return nil and false
	//otherwise save the key on the cache
	//and return the session and true

	return newSession, true
}
