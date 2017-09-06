package github

import (
	"net/http"

	"github.com/pkg/errors"
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
func (v *OrganizationVerifier) Verify(r *http.Request, httpClient *http.Client) (bool, error) {
	orgs, err := v.gitHubClient.Organizations(httpClient)
	if err != nil {
		return false, errors.Wrap(err, "failed to get organizations")
	}

	for _, name := range orgs {
		for _, authorizedOrg := range v.organizations {
			if name == authorizedOrg {
				return true, nil
			}
		}
	}

	return false, errors.New("you are not part of the allowed organizations")
}
