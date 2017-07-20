package oauth

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type oAuthResponse struct {
	Active bool `json:"active"`
}

// IntrospectionManager is responsible for using OAuth2 Introspection definition to
// validate tokens from an authentication provider
type IntrospectionManager struct {
	URL string
}

// NewIntrospectionManager creates a new instance of Introspection
func NewIntrospectionManager(url string) (*IntrospectionManager, error) {
	if url == "" {
		return nil, ErrInvalidIntrospectionURL
	}

	return &IntrospectionManager{url}, nil
}

// IsKeyAuthorized checks if the access token is valid
func (o *IntrospectionManager) IsKeyAuthorized(accessToken string) bool {
	resp, err := doStatusRequest(o.URL)
	defer resp.Body.Close()

	if err != nil {
		log.WithField("url", o.URL).
			WithError(err).
			Error("Error making a request to the authentication provider")
	}

	if resp.StatusCode != http.StatusOK {
		log.Info("The token check was invalid")
		return false
	}

	var oauthResp oAuthResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&oauthResp)
	if err != nil {
		return false
	}

	return oauthResp.Active
}

func doStatusRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.WithError(err).Error("Creating the request for the health check failed")
		return nil, err
	}

	// Inform to close the connection after the transaction is complete
	req.Header.Set("Connection", "close")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Making the request for the health check failed")
		return resp, err
	}

	return resp, err
}
