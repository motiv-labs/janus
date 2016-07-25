package main

// Enums for keys to be stored in a session context - this is how gorilla expects
// these to be implemented and is lifted pretty much from docs
const (
	SessionData     = "session_data"
	AuthHeaderValue = "auth_header"
)

// SessionState objects represent a current API session, mainly used for rate limiting.
type SessionState struct {
	ExpiresIn     int64  `json:"expires_in"`
	OauthClientID string `json:"oauth_client_id"`
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
}
