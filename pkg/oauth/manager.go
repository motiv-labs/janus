package oauth

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/store"
)

// Manager is responsible for managing the access tokens
type Manager struct {
	Storage store.Store
}

// KeyExists checks if the given access token exits in the storage
func (o *Manager) KeyExists(accessToken string) (bool, error) {
	log.Debugf("Searching for key %s", accessToken)
	return o.Storage.Exists(accessToken)
}

// Set a new access token and its session to the storage
func (o *Manager) Set(accessToken string, session session.SessionState, resetTTLTo int64) error {
	value, err := json.Marshal(session)
	if err != nil {
		return err
	}

	log.Debugf("Storing key %s for %d seconds", accessToken, resetTTLTo)
	go o.Storage.Set(accessToken, string(value), resetTTLTo)

	return nil
}

// IsKeyAuthorised checks if the access token is valid
func (o *Manager) IsKeyAuthorised(accessToken string) (session.SessionState, bool) {
	var newSession session.SessionState
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
