package github

import (
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// OrganizationVerifier checks if the current user belongs any of the defined organizations
type OrganizationVerifier struct {
	organizations []string
	gitHubClient  Client
}

// NewOrganizationVerifier creates a new instance of OrganizationVerifier
func NewOrganizationVerifier(organizations []string, gitHubClient Client) *OrganizationVerifier {
	return &OrganizationVerifier{
		organizations: organizations,
		gitHubClient:  gitHubClient,
	}
}

// Verify makes a check and return a boolean if the check was successful or not
func (verifier OrganizationVerifier) Verify(r *http.Request) (bool, error) {
	accessToken, err := extractAccessToken(r)
	if err != nil {
		return false, errors.Wrap(err, "failed to extract access token")
	}

	orgs, err := verifier.gitHubClient.Organizations(getClient(accessToken))
	if err != nil {
		return false, errors.Wrap(err, "failed to get organizations")
	}

	for _, name := range orgs {
		for _, authorizedOrg := range verifier.organizations {
			if name == authorizedOrg {
				return true, nil
			}
		}
	}

	log.WithFields(log.Fields{
		"have": orgs,
		"want": verifier.organizations,
	}).Debug("not in the organizations")

	return false, nil
}
