package github

import (
	"context"
	"net/http"

	"github.com/google/go-github/github"
)

// Client contains the methods that abstract an API
type Client interface {
	CurrentUser(*http.Client) (*github.User, error)
	Organizations(*http.Client) ([]string, error)
	Teams(*http.Client) (OrganizationTeams, error)
}

type client struct {
	baseURL string
}

// NewClient creates a new instance of client
func NewClient() Client {
	return &client{}
}

// OrganizationTeams is a map of organization names and teams
type OrganizationTeams map[string][]string

// CurrentUser retrieves the current authenticated user for an http client
func (c *client) CurrentUser(httpClient *http.Client) (*github.User, error) {
	client := github.NewClient(httpClient)

	currentUser, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		return nil, err
	}

	return currentUser, nil
}

// Teams retrieves the teams that the authenticated user belongs
func (c *client) Teams(httpClient *http.Client) (OrganizationTeams, error) {
	client := github.NewClient(httpClient)

	nextPage := 1
	organizationTeams := OrganizationTeams{}

	for nextPage != 0 {
		teams, resp, err := client.Teams.ListUserTeams(context.TODO(), &github.ListOptions{Page: nextPage})
		if err != nil {
			return nil, err
		}

		for _, team := range teams {
			organizationName := *team.Organization.Login

			if _, found := organizationTeams[organizationName]; !found {
				organizationTeams[organizationName] = []string{}
			}

			// We add both forms (slug and name) of team
			organizationTeams[organizationName] = append(organizationTeams[organizationName], *team.Name)
			organizationTeams[organizationName] = append(organizationTeams[organizationName], *team.Slug)
		}

		nextPage = resp.NextPage
	}

	return organizationTeams, nil
}

// Organizations retrieves the organizations that the authenticated user belongs
func (c *client) Organizations(httpClient *http.Client) ([]string, error) {
	client := github.NewClient(httpClient)

	nextPage := 1
	organizations := []string{}

	for nextPage != 0 {
		orgs, resp, err := client.Organizations.List(context.TODO(), "", &github.ListOptions{Page: nextPage})

		if err != nil {
			return nil, err
		}

		for _, org := range orgs {
			organizations = append(organizations, *org.Login)
		}

		nextPage = resp.NextPage
	}

	return organizations, nil
}
