package oauth

import (
	"encoding/base64"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// SecretMiddleware is used as a helper for client applications that don't want to send the client secret
// on the request. The applications should only send the `client_id` and this middleware will try to find
// the secret on it's configuration.
// If the secret is found then the middleware will build a valid `Authorization` header to be sent to the
// authentication provider.
// If the secret is not found then and error is returned to the client application.
type SecretMiddleware struct {
	oauth *Spec
}

// NewSecretMiddleware creates an instance of SecretMiddleware
func NewSecretMiddleware(oauth *Spec) *SecretMiddleware {
	return &SecretMiddleware{oauth}
}

// Handler is the middleware method.
func (m *SecretMiddleware) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Starting Oauth2Secret middleware")

		if "" != r.Header.Get("Authorization") {
			log.Debug("Authorization is set, proxying")
			handler.ServeHTTP(w, r)
			return
		}

		clientID := r.URL.Query().Get("client_id")
		if "" == clientID {
			log.Debug("ClientID not set, proxying")
			handler.ServeHTTP(w, r)
			return
		}

		clientSecret, exists := m.oauth.Secrets[clientID]
		if false == exists {
			panic(ErrClientIDNotFound)
		}

		m.changeRequest(r, clientID, clientSecret)
		handler.ServeHTTP(w, r)
	})
}

// changeRequest modifies the request to add the Authorization headers.
func (m *SecretMiddleware) changeRequest(req *http.Request, clientID, clientSecret string) {
	log.Debug("Modifying request")
	authString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authString))
}
