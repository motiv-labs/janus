package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/proxy/balancer"
)

type oAuthResponse struct {
	Active bool `json:"active"`
}

// IntrospectionManager is responsible for using OAuth2 Introspection definition to
// validate tokens from an authentication provider
type IntrospectionManager struct {
	balancer balancer.Balancer
	urls     proxy.Targets
	settings *IntrospectionSettings
}

// NewIntrospectionManager creates a new instance of Introspection
func NewIntrospectionManager(def *proxy.Definition, settings *IntrospectionSettings) (*IntrospectionManager, error) {
	bb, err := balancer.New(def.Upstreams.Balancing)
	if err != nil {
		return nil, fmt.Errorf("could not create a bb: %w", err)
	}

	return &IntrospectionManager{bb, def.Upstreams.Targets, settings}, nil
}

// IsKeyAuthorized checks if the access token is valid
func (o *IntrospectionManager) IsKeyAuthorized(ctx context.Context, accessToken string) bool {
	resp, err := o.doStatusRequest(accessToken)
	defer resp.Body.Close()

	if err != nil {
		log.WithError(err).
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

func (o *IntrospectionManager) doStatusRequest(accessToken string) (*http.Response, error) {
	upstream, err := o.balancer.Elect(o.urls.ToBalancerTargets())
	if err != nil {
		return nil, fmt.Errorf("could not elect one upstream: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, upstream.Target, nil)
	if err != nil {
		log.WithError(err).Error("Creating the request for the health check failed")
		return nil, err
	}

	if o.settings.UseAuthHeader {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", o.settings.AuthHeaderType, accessToken))
	} else if o.settings.UseCustomHeader {
		req.Header.Add(o.settings.HeaderName, accessToken)
	} else {
		req.Form = make(url.Values)
		req.Form.Add(o.settings.ParamName, accessToken)
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
