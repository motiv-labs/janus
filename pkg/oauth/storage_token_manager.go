package oauth

import (
	"encoding/json"

	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/store"
	log "github.com/sirupsen/logrus"
)

// StorageTokenManager is responsible for managing the access tokens
type StorageTokenManager struct {
	Storage store.Store
}

// Set a new access token and its session to the storage
func (o *StorageTokenManager) Set(accessToken string, session session.State, resetTTLTo int64) error {
	value, err := json.Marshal(session)
	if err != nil {
		return err
	}

	log.Debugf("Storing key %s for %d seconds", accessToken, resetTTLTo)
	go o.Storage.Set(accessToken, string(value), resetTTLTo)

	return nil
}

// Remove an access token from the storage
func (o *StorageTokenManager) Remove(accessToken string) error {
	log.WithField("token", accessToken).Debug("removing token from the storage")
	go o.Storage.Remove(accessToken)

	return nil
}

// IsKeyAuthorised checks if the access token is valid
func (o *StorageTokenManager) IsKeyAuthorised(accessToken string) (session.State, bool) {
	var newSession session.State

	exists, err := o.Storage.Exists(accessToken)
	if !exists || err != nil {
		log.WithError(err).Warn("Key not found in keystore")
		return newSession, false
	}

	jsonKeyVal, err := o.Storage.Get(accessToken)
	if nil != err {
		log.WithError(err).Error("Couldn't get the access token from storage")
		return newSession, false
	}

	if marshalErr := json.Unmarshal([]byte(jsonKeyVal), &newSession); marshalErr != nil {
		log.WithError(marshalErr).Error("Couldn't unmarshal session object")
		return newSession, false
	}

	return newSession, true
}
