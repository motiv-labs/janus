package oauth

import (
	"strings"

	"github.com/hellofresh/janus/pkg/jwt"
	"github.com/hellofresh/janus/pkg/session"
	"github.com/hellofresh/janus/pkg/store"
	log "github.com/sirupsen/logrus"
)

const (
	// Storage enables you to store the tokens in a cache (This way you don't need to validate the token against
	// the auth provider on every request)
	Storage ManagerType = iota
	// JWT provides a way to check the `exp` field on the JWT and make sure the token is still valid. This is
	// probably the most versatile way to check for tokens, since it doesn't require any storage or extra calls in
	// each request.
	JWT
	// Auth strategy makes sure to validate the provided token on every request against the athentication provider.
	Auth
)

var typesMap = map[string]ManagerType{
	"storage": Storage,
	"jwt":     JWT,
	"auth":    Auth,
}

// ParseType takes a string type and returns the Manager type constant.
func ParseType(lvl string) (ManagerType, error) {
	m, ok := typesMap[strings.ToLower(lvl)]
	if !ok {
		var m ManagerType
		return m, ErrUnknownStrategy
	}
	return m, nil
}

// ManagerType type
type ManagerType uint8

// Manager holds the methods to handle tokens
type Manager interface {
	Set(accessToken string, session session.State, resetTTLTo int64) error
	Remove(accessToken string) error
	IsKeyAuthorised(accessToken string) (session.State, bool)
}

// ManagerFactory is used for creating a new manager
type ManagerFactory struct {
	Storage  store.Store
	settings TokenStrategySettings
}

// NewManagerFactory creates a new instance of ManagerFactory
func NewManagerFactory(storage store.Store, settings TokenStrategySettings) *ManagerFactory {
	return &ManagerFactory{storage, settings}
}

// Build creates a manager based on the type
func (f *ManagerFactory) Build(t ManagerType) (Manager, error) {
	// FIXME: make it nicer with BiMap - GetByType, GetByName
	typesMapReversed := make(map[ManagerType]string, len(typesMap))
	for k, v := range typesMap {
		typesMapReversed[v] = k
	}

	log.WithField("name", typesMapReversed[t]).
		Debug("Building token strategy")

	switch t {
	case Storage:
		return &StorageTokenManager{Storage: f.Storage}, nil
	case JWT:
		value, err := f.settings.GetJWTSecret()
		if nil != err {
			return nil, err
		}

		return NewJWTManager(jwt.NewParser(jwt.NewConfig("HS256", value))), nil
	case Auth:
		// TODO: Create an Auth Manager that always validated tokens against an auth provider
	}

	return nil, ErrUnknownManager
}
