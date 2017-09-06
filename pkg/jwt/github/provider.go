package github

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt/provider"
)

func init() {
	provider.Register("github", &Provider{})
}

// Provider abstracts the authentication for github
type Provider struct {
	provider.Verifier
}

// Build acts like the constructor for a provider
func (gp *Provider) Build(config config.Credentials) provider.Provider {
	client := NewClient()

	return &Provider{
		Verifier: provider.NewVerifierBasket(
			NewTeamVerifier(teamConfigsToTeam(config.Github.Teams), client),
			NewOrganizationVerifier(config.Github.Organizations, client),
		),
	}
}

// GetClaims returns a JWT Map Claim
func (gp *Provider) GetClaims(httpClient *http.Client) (jwt.MapClaims, error) {
	client := NewClient()

	user, err := client.CurrentUser(httpClient)
	if err != nil {
		return nil, err
	}

	return jwt.MapClaims{
		"sub": *user.Login,
	}, nil
}

func teamConfigsToTeam(configTeams map[string]string) []Team {
	teams := []Team{}
	for org, team := range configTeams {
		teams = append(teams, Team{
			Name:         team,
			Organization: org,
		})
	}
	return teams
}
