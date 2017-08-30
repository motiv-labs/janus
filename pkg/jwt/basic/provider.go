package basic

import (
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt/provider"
)

func init() {
	provider.Register("basic", Provider{})
}

// Provider abstracts the authentication for github
type Provider struct {
	provider.Verifier
}

// Build acts like the constructor for a provider
func (gp Provider) Build(config config.Credentials) provider.Provider {
	return &Provider{
		Verifier: provider.NewVerifierBasket(
			NewPasswordVerifier(userConfigToTeam(config.Basic.Users)),
		),
	}
}

func userConfigToTeam(configUser []config.BasicUsersConfig) []*User {
	users := []*User{}
	for _, u := range configUser {
		users = append(users, &User{
			Username: u.Username,
			Password: u.Password,
		})
	}
	return users
}
