package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

var (
	recovery = middleware.NewRecovery(test.RecoveryHandler)
)

func TestValidKeyStorage(t *testing.T) {
	session := session.State{
		AccessToken: "123",
	}

	mw, err := createMiddlewareWithSession(session)
	assert.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", session.AccessToken),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestWrongAuthHeader(t *testing.T) {
	session := session.State{
		AccessToken: "123",
	}

	mw, err := createMiddlewareWithSession(session)
	assert.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Wrong %s", session.AccessToken),
		},
		recovery(mw(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestMissingAuthHeader(t *testing.T) {
	session := session.State{
		AccessToken: "123",
	}

	mw, err := createMiddlewareWithSession(session)
	assert.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		recovery(mw(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestMissingKeyStorage(t *testing.T) {
	session := session.State{
		AccessToken: "123",
	}

	mw, err := createMiddlewareWithSession(session)
	assert.NoError(t, err)

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", "1234"),
		},
		recovery(mw(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func createMiddlewareWithSession(session session.State) (func(http.Handler) http.Handler, error) {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	storage := store.NewInMemoryStore()
	storage.Set(session.AccessToken, string(sessionJSON), 0)

	manager, err := oauth.NewManagerFactory(storage, oauth.TokenStrategySettings{}).Build(oauth.Storage)
	if err != nil {
		return nil, err
	}

	return NewKeyExistsMiddleware(manager), nil
}
