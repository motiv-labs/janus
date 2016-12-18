package oauth

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/session"
	"gopkg.in/redis.v3"
)

type Manager struct {
	Storage *redis.Client
}

func (o *Manager) KeyExists(accessToken string) bool {
	log.Debugf("Searching for key %s", accessToken)
	return o.Storage.Exists(accessToken).Val()
}

func (o *Manager) Set(accessToken string, session session.SessionState, resetTTLTo int64) error {
	value, _ := json.Marshal(session)
	expireDuration := time.Duration(resetTTLTo) * time.Second

	log.Debugf("Storing key %s for %d seconds", accessToken, expireDuration)
	go o.Storage.Set(accessToken, string(value), expireDuration)

	return nil
}

func (o *Manager) IsKeyAuthorised(accessToken string) (session.SessionState, bool) {
	var newSession session.SessionState
	jsonKeyVal := o.Storage.Get(accessToken).Val()

	if marshalErr := json.Unmarshal([]byte(jsonKeyVal), &newSession); marshalErr != nil {
		log.Errorf("Couldn't unmarshal session object: %s", marshalErr.Error())
		return newSession, false
	}

	return newSession, true
}
