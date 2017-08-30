package github

import (
	"context"
	"net/http"
	"strings"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt/provider"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func init() {
	provider.Register("github", Provider{})
}

// Provider abstracts the authentication for github
type Provider struct {
	provider.Verifier
}

// Build acts like the constructor for a provider
func (gp Provider) Build(config config.Credentials) provider.Provider {
	client := NewClient()

	return &Provider{
		Verifier: provider.NewVerifierBasket(
			NewTeamVerifier(teamConfigsToTeam(config.Github.Teams), client),
			NewOrganizationVerifier(config.Github.Organizations, client),
		),
	}
}

func teamConfigsToTeam(configTeams []config.GitHubTeamConfig) []Team {
	teams := []Team{}
	for _, team := range configTeams {
		teams = append(teams, Team{
			Name:         team.TeamName,
			Organization: team.OrganizationName,
		})
	}
	return teams
}

func extractAccessToken(r *http.Request) (string, error) {
	// We're using OAuth, start checking for access keys
	authHeaderValue := r.Header.Get("Authorization")
	parts := strings.Split(authHeaderValue, " ")
	if len(parts) < 2 {
		return "", errors.New("attempted access with malformed header, no auth header found")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("bearer token malformed")
	}

	return parts[1], nil
}

func getClient(token string) *http.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	return oauth2.NewClient(ctx, ts)
}
