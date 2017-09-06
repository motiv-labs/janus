package basic

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	r          *http.Request
	httpClient *http.Client
)

func TestPasswordVerifier(t *testing.T) {
	tests := []struct {
		scenario string
		function func(*testing.T, *PasswordVerifier)
	}{
		{
			scenario: "when credentials are sent as form parameters",
			function: testSendFormParamsCredentials,
		},
		{
			scenario: "when credentials are sent as application/json",
			function: testSendJSONCredentials,
		},
		{
			scenario: "when basic header is sent",
			function: testBasicHeaderCredentials,
		},
		{
			scenario: "when credentials are sent as application/json;charset=UTF-8",
			function: testSendJSONWithCharsetCredentials,
		},
		{
			scenario: "when invalid credentials are given we should get an error",
			function: testInvalidCredentialsGiven,
		},
		{
			scenario: "when no credentials are given we should get an error",
			function: testNoCredentialsGiven,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			t.Parallel()
			verifier := NewPasswordVerifier([]*user{
				{Username: "user1", Password: "test"},
				{Username: "user2", Password: "test"},
			})

			test.function(t, verifier)
		})
	}
}

func testSendFormParamsCredentials(t *testing.T, v *PasswordVerifier) {
	r := httptest.NewRequest("GET", "/", nil)
	r.ParseForm()
	r.Form.Add("username", "user1")
	r.Form.Add("password", "test")

	result, err := v.Verify(r, httpClient)

	require.NoError(t, err)
	assert.True(t, result)
}

func testSendJSONCredentials(t *testing.T, v *PasswordVerifier) {
	r := httptest.NewRequest("GET", "/", strings.NewReader(`{"username": "user1", "password": "test"}`))
	r.Header.Add("Content-Type", "application/json")
	result, err := v.Verify(r, httpClient)

	require.NoError(t, err)
	assert.True(t, result)
}

func testBasicHeaderCredentials(t *testing.T, v *PasswordVerifier) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetBasicAuth("user1", "test")
	result, err := v.Verify(r, httpClient)

	require.NoError(t, err)
	assert.True(t, result)
}

func testSendJSONWithCharsetCredentials(t *testing.T, v *PasswordVerifier) {
	r := httptest.NewRequest("GET", "/", strings.NewReader(`{"username": "user1", "password": "test"}`))
	r.Header.Add("Content-Type", "application/json;charset=UTF-8")
	result, err := v.Verify(r, httpClient)

	require.NoError(t, err)
	assert.True(t, result)
}

func testInvalidCredentialsGiven(t *testing.T, v *PasswordVerifier) {
	r := httptest.NewRequest("GET", "/", nil)
	r.ParseForm()
	r.Form.Add("username", "user1")
	r.Form.Add("password", "wrong")

	result, err := v.Verify(r, httpClient)

	require.Error(t, err)
	assert.False(t, result)
}

func testNoCredentialsGiven(t *testing.T, v *PasswordVerifier) {
	r := httptest.NewRequest("GET", "/", nil)
	result, err := v.Verify(r, httpClient)

	require.Error(t, err)
	assert.False(t, result)
}
