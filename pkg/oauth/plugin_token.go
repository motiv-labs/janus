package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/store"
)

// TokenPlugin represents an outbound plugin for handling oauth tokens
type TokenPlugin struct {
	storage store.Store
	repo    Repository
}

// NewTokenPlugin creates a new instance of TokenPlugin
func NewTokenPlugin(storage store.Store, repo Repository) *TokenPlugin {
	return &TokenPlugin{storage, repo}
}

// Out is the entry point for a plugin
func (t *TokenPlugin) Out(req *http.Request, res *http.Response) (*http.Response, error) {
	if res.StatusCode < http.StatusMultipleChoices && res.Body != nil {
		var newSession session.State

		//This is useful for the middlewares
		var bodyBytes []byte

		defer func(body io.Closer) {
			err := body.Close()
			if err != nil {
				log.Error(err)
			}
		}(res.Body)
		bodyBytes, _ = ioutil.ReadAll(res.Body)

		if marshalErr := json.Unmarshal(bodyBytes, &newSession); marshalErr == nil {
			if newSession.AccessToken != "" {
				tokenURL := url.URL{Scheme: req.URL.Scheme, Host: req.URL.Host, Path: req.URL.Path}
				log.WithField("token_url", tokenURL.String()).Debug("Looking for OAuth provider who issued the token")
				manager, oAuthServer, err := t.getManager(tokenURL)
				if err != nil {
					log.WithError(err).Error("Failed to find OAuth server by token URL")
				} else {
					newSession.OAuthServer = oAuthServer.Name
					log.Debug("Setting body in the oauth storage")
					manager.Set(newSession.AccessToken, newSession, newSession.ExpiresIn)
				}
			}
		}

		// Restore the io.ReadCloser to its original state
		res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return res, nil
}

func (t *TokenPlugin) getManager(url url.URL) (Manager, *OAuth, error) {
	oauthServer, err := t.repo.FindByTokenURL(url)
	if nil != err {
		return nil, nil, err
	}

	managerType, err := ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, nil, err
	}

	manager, err := NewManagerFactory(t.storage, oauthServer.TokenStrategy.Settings).Build(managerType)

	return manager, oauthServer, err
}
