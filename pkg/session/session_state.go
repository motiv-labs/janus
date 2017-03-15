package session

// State objects represent a current API session, mainly used for access tokens.
type State struct {
	OAuthServer   string `json:"oauth_server"`
	ExpiresIn     int64  `json:"expires_in"`
	OauthClientID string `json:"oauth_client_id"`
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
}
