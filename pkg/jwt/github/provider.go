package github

import (
	"net/http"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/jwt/provider"
	"github.com/pkg/errors"
)

func init() {
	provider.Register("github", &Provider{})
}

// Provider abstracts the authentication for github
type Provider struct {
	provider.Verifier

	teams  []Team
	config config.Credentials
}

// Build acts like the constructor for a provider
func (gp *Provider) Build(config config.Credentials) provider.Provider {
	client := NewClient()
	teams := gp.teamConfigsToTeam(config.Github.Teams)

	return &Provider{
		Verifier: provider.NewVerifierBasket(
			NewTeamVerifier(teams, client),
			NewOrganizationVerifier(config.Github.Organizations, client),
		),
		teams:  teams,
		config: config,
	}
}

// GetClaims returns a JWT Map Claim
func (gp *Provider) GetClaims(httpClient *http.Client) (jwt.MapClaims, error) {
	client := NewClient()

	var (
		wg            sync.WaitGroup
		user          *github.User
		usersOrgTeams OrganizationTeams
		err           error
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		res, err := client.CurrentUser(httpClient)
		if err != nil {
			err = append(errs, errors.Wrap(err, "failed to get github users"))
			return
		}
		user = res
	}()

	go func() {
		defer wg.Done()
		res, err := client.Teams(httpClient)
		if err != nil {
			err = append(errs, errors.Wrap(err, "failed to get github teams"))
			return
		}
		usersOrgTeams = res
	}()

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return jwt.MapClaims{
		"sub":      *user.Login,
		"is_admin": gp.isAdmin(usersOrgTeams),
	}, nil
}

func (gp *Provider) teamConfigsToTeam(configTeams map[string]string) []Team {
	teams := []Team{}
	for org, team := range configTeams {
		teams = append(teams, Team{
			Name:         team,
			Organization: org,
		})
	}
	return teams
}

func (gp *Provider) isAdmin(usersOrgTeams OrganizationTeams) bool {
	for _, team := range gp.teams {
		if teams, ok := usersOrgTeams[team.Organization]; ok {
			for _, teamUserBelongsTo := range teams {
				if teamUserBelongsTo == gp.config.JanusAdminTeam {
					return true
				}
			}
		}
	}
	return false
}
