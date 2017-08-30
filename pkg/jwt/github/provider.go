package github

import (
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt/provider"
)

func init() {
	provider.Register("github", Provider{})
}

// Provider abstracts the authentication for github
type Provider struct {
	provider.Verifier
}

// Build acts like the constructor for a provider
func (gp GithubProvider) Build(config config.Credentials) provider.Provider {
	client := NewClient()

	return &GithubProvider{
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
