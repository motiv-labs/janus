package oauth

import (
	"strings"

	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/store"
)

const (
	// Storage enables you to store the tokens in a cache (This way you don't need to validate the token against
	// the auth provider on every request)
	Storage ManagerType = iota
	// JWT provides a way to check the `exp` field on the JWT and make sure the token is still valid. This is
	// probably the most versatile way to check for tokens, since it doesn't require any storage or extra calls in
	// each request.
	JWT
	// Auth strategy makes sure to validade the provided token on every request against the athentication provider.
	Auth
)

// ParseType takes a string type and returns the Manager type constant.
func ParseType(lvl string) (ManagerType, error) {
	switch strings.ToLower(lvl) {
	case "storage":
		return Storage, nil
	case "jwt":
		return JWT, nil
	case "auth":
		return Auth, nil
	}

	var m ManagerType
	return m, ErrUnknownStrategy
}

// ManagerType type
type ManagerType uint8

// Manager holds the methods to handle tokens
type Manager interface {
	Set(accessToken string, session session.SessionState, resetTTLTo int64) error
	Remove(accessToken string) error
	IsKeyAuthorised(accessToken string) (session.SessionState, bool)
}

type ManagerFactory struct {
	Storage store.Store
	Secret  string
}

func NewManagerFactory(storage store.Store, secret string) *ManagerFactory {
	return &ManagerFactory{storage, secret}
}

func (f *ManagerFactory) Build(t ManagerType) (Manager, error) {
	switch t {
	case Storage:
		return &StorageTokenManager{Storage: f.Storage}, nil
	case JWT:
		if f.Secret == "" {
			return nil, ErrJWTSecretMissing
		}
		return &JWTManager{Secret: f.Secret}, nil
	case Auth:
		// TODO: Create an Auth Manager that always validated tokens against an auth provider
	}

	return nil, ErrUnknownManager
}
