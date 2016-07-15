package main

// SessionState objects represent a current API session, mainly used for rate limiting.
type SessionState struct {
	ExpiresIn     int64                       `json:"expires_in"`
	OauthClientID string                      `json:"oauth_client_id"`
	AccessToken   string                      `json:"client_id"`
	TokenType     string                      `json:"token_type"`
}
