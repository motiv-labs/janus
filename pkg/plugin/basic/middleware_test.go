package basic

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/hellofresh/janus/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizedAccess(t *testing.T) {
	mw := NewBasicAuth(setupRepo())

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Basic " + basicAuth("test", "test"),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestInvalidBasicHeader(t *testing.T) {
	mw := NewBasicAuth(setupRepo())

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Basic wrong",
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestUnauthorizedAccess(t *testing.T) {
	mw := NewBasicAuth(setupRepo())

	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Basic " + basicAuth("wrong", "wrong"),
		},
		mw(http.HandlerFunc(test.Ping)),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func setupRepo() Repository {
	repo := NewInMemoryRepository()
	repo.Add(&User{Username: "test", Password: "test"})

	return repo
}
