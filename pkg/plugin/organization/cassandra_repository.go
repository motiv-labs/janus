package organization

import (
	"encoding/json"
	"github.com/hellofresh/janus/cassandra/wrapper"
	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	log "github.com/sirupsen/logrus"
)

// Repository represents an user repository
type Repository interface {
	FindAll() ([]*Organization, error)
	FindByUsername(username string) (*Organization, error)
	FindOrganization(organization string) (*OrganizationConfig, error)
	Add(organization *Organization) error
	AddOrganization(organization *OrganizationConfig) error
	Remove(username string) error
}

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session wrapper.Holder
	hash    encrypt.Hash
}

// NewCassandraRepository constructs CassandraRepository
func NewCassandraRepository(session wrapper.Holder) (*CassandraRepository, error) {
	log.Debugf("getting new organization cassandra repo")
	return &CassandraRepository{session: session}, nil
}

// FindAll fetches all the basic user definitions available
func (r *CassandraRepository) FindAll() ([]*Organization, error) {
	log.Debugf("finding all users servers")

	var results []*Organization

	iter := r.session.GetSession().Query("SELECT username, organization, password FROM organization").Iter()

	var username string
	var comp string
	var pass string

	err := iter.ScanAndClose(func() bool {
		var organization Organization
		organization.Username = username
		organization.Organization = comp
		organization.Password = pass
		results = append(results, &organization)
		return true
	}, &username, &comp, &pass)

	if err != nil {
		log.Errorf("error getting all organization users: %v", err)
	}
	return results, err

}

// FindByUsername find an user by username
// returns ErrUserNotFound when a user is not found.
func (r *CassandraRepository) FindByUsername(username string) (*Organization, error) {
	log.Debugf("finding: %s", username)

	var organization Organization

	err := r.session.GetSession().Query(
		"SELECT username, organization, password "+
			"FROM organization "+
			"WHERE username = ?",
		username).Scan(&organization.Username, &organization.Organization, &organization.Password)

	if err != nil {
		if err.Error() == "not found" {
			log.Debugf("organization not found")
			err = ErrUserNotFound
		}
		log.Errorf("error selecting organization user %s: %v", username, err)
	} else {
		log.Debugf("successfully found organization user %s", username)
	}

	return &organization, err
}

// FindOrganization find an organization by organization name
// returns ErrUserNotFound when a user is not found.
func (r *CassandraRepository) FindOrganization(organization string) (*OrganizationConfig, error) {
	log.Debugf("finding: %s", organization)

	var organizationConfig OrganizationConfig
	var bOrgConfig []byte

	err := r.session.GetSession().Query(
		"SELECT organization, priority, content_per_day, config "+
			"FROM organization_config "+
			"WHERE organization = ?",
		organization).Scan(
		&organizationConfig.Organization,
		&organizationConfig.Priority,
		&organizationConfig.ContentPerDay,
		&bOrgConfig)

	if err != nil {
		if err.Error() == "not found" {
			log.Debugf("organization not found")
			err = ErrUserNotFound
		}
		log.Errorf("error selecting organization %s: %v", organization, err)
		return &organizationConfig, err
	} else {
		log.Debugf("successfully found organization %s", organization)
	}

	err = json.Unmarshal(bOrgConfig, &organizationConfig.Config)
	if err != nil {
		log.Errorf("error unmarshalling config: %v", err)
	}

	return &organizationConfig, err
}

// Add adds an user to the repository
func (r *CassandraRepository) Add(organization *Organization) error {
	log.Debugf("adding: %s", organization.Username)

	hash, err := r.hash.Generate(organization.Password)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return err
	}

	err = r.session.GetSession().Query(
		"UPDATE organization "+
			"SET organization = ?, "+
			"password = ? "+
			"WHERE username = ?",
		organization.Organization, hash, organization.Username).Exec()

	if err != nil {
		log.Errorf("error saving organization user %s: %v", organization.Username, err)
	} else {
		log.Debugf("successfully saved organization user %s", organization.Username)
	}

	return err
}

// AddOrganization adds an organization to the repository
func (r *CassandraRepository) AddOrganization(organization *OrganizationConfig) error {
	log.Debugf("adding: %s", organization.Organization)

	bOrgConfig, err := json.Marshal(organization.Config)
	if err != nil {
		log.Errorf("error marshaling config %s: %v", organization.Config, err)
	}

	err = r.session.GetSession().Query(
		"UPDATE organization_config "+
			"SET priority = ?, "+
			"content_per_day = ?, "+
			"config = ? "+
			"WHERE organization = ?",
		organization.Priority, organization.ContentPerDay, bOrgConfig, organization.Organization).Exec()

	if err != nil {
		log.Errorf("error saving organization %s: %v", organization.Organization, err)
	} else {
		log.Debugf("successfully saved organization organization %s", organization.Organization)
	}

	return err
}

// Remove an user from the repository
func (r *CassandraRepository) Remove(username string) error {
	log.Debugf("removing: %s", username)

	err := r.session.GetSession().Query(
		"DELETE FROM organization WHERE username = ?", username).Exec()

	if err != nil {
		log.Errorf("error removing organization user %s: %v", username, err)
	} else {
		log.Debugf("successfully removed organization user %s", username)
	}

	return err
}
