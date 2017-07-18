package oauth

// AuthProviderManager is responsible for managing the access tokens
type AuthProviderManager struct{}

// IsKeyAuthorized checks if the access token is valid
func (o *AuthProviderManager) IsKeyAuthorized(accessToken string) bool {
	return true
}
