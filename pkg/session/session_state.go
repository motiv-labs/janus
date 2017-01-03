package session

import "gopkg.in/mgo.v2/bson"

// SessionState objects represent a current API session, mainly used for access tokens.
type SessionState struct {
	OAuthServerID bson.ObjectId `json:"server_id"`
	ExpiresIn     int64         `json:"expires_in"`
	OauthClientID string        `json:"oauth_client_id"`
	AccessToken   string        `json:"access_token"`
	TokenType     string        `json:"token_type"`
}
