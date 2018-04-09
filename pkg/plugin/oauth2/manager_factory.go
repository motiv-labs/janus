package oauth2

import (
	"context"
	"fmt"
	"strings"

	"github.com/hellofresh/janus/pkg/jwt"
	log "github.com/sirupsen/logrus"
)

const (
	// JWT provides a way to check the `exp` field on the JWT and make sure the token is still valid. This is
	// probably the most versatile way to check for tokens, since it doesn't require any storage or extra calls in
	// each request.
	JWT ManagerType = iota
	// Introspection strategy makes sure to validate the provided token on every request against the authentication provider.
	Introspection
)

var typesMap = map[string]ManagerType{
	"jwt":           JWT,
	"introspection": Introspection,
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
	IsKeyAuthorized(ctx context.Context, accessToken string) bool
}

// ManagerFactory is used for creating a new manager
type ManagerFactory struct {
	oAuthServer *OAuth
}

// NewManagerFactory creates a new instance of ManagerFactory
func NewManagerFactory(oAuthServer *OAuth) *ManagerFactory {
	return &ManagerFactory{oAuthServer}
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
	case JWT:
		signingMethods, err := f.oAuthServer.TokenStrategy.GetJWTSigningMethods()
		if nil != err {
			return nil, err
		}

		logEntry := log.WithField("leeway", f.oAuthServer.TokenStrategy.Leeway)
		for i, signingMethod := range signingMethods {
			logEntry = logEntry.WithField(fmt.Sprintf("alg_%d", i), signingMethod.Alg)
		}
		logEntry.Debug("Building JWT token parser")

		return NewJWTManager(jwt.NewParser(jwt.NewParserConfig(f.oAuthServer.TokenStrategy.Leeway, signingMethods...))), nil
	case Introspection:
		settings, err := f.oAuthServer.TokenStrategy.GetIntrospectionSettings()
		if nil != err {
			return nil, err
		}

		manager, err := NewIntrospectionManager(f.oAuthServer.Endpoints.Introspect, settings)
		if err != nil {
			return nil, err
		}

		return manager, nil
	}

	return nil, ErrUnknownManager
}
