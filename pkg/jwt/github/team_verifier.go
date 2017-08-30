package github

import (
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Team represents a github team within the organization
type Team struct {
	Name         string
	Organization string
}

// TeamVerifier checks if the current user belongs any of the defined teams
type TeamVerifier struct {
	teams        []Team
	gitHubClient Client
}

// NewTeamVerifier creates a new instance of TeamVerifier
func NewTeamVerifier(teams []Team, gitHubClient Client) *TeamVerifier {
	return &TeamVerifier{
		teams:        teams,
		gitHubClient: gitHubClient,
	}
}

// Verify makes a check and return a boolean if the check was successful or not
func (verifier TeamVerifier) Verify(r *http.Request, httpClient *http.Client) (bool, error) {
	usersOrgTeams, err := verifier.gitHubClient.Teams(httpClient)
	if err != nil {
		return false, errors.Wrap(err, "failed to get teams")
	}

	for _, team := range verifier.teams {
		if teams, ok := usersOrgTeams[team.Organization]; ok {
			for _, teamUserBelongsTo := range teams {
				if teamUserBelongsTo == team.Name {
					return true, nil
				}
			}
		}
	}

	log.WithFields(log.Fields{
		"have": usersOrgTeams,
		"want": verifier.teams,
	}).Debug("not in teams")

	return false, nil
}
