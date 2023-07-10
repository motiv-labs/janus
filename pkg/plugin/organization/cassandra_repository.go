package organization

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	log "github.com/sirupsen/logrus"
)

// Repository represents an user repository
type Repository interface {
	FindAll() ([]*Organization, error)
	FindByUsername(username string) (*Organization, error)
	Add(organization *Organization) error
	Remove(username string) error
}

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session wrapper.Holder
	hash    encrypt.Hash
}

func NewCassandraRepository(session wrapper.Holder) (*CassandraRepository, error) {
	log.Debugf("getting new organization cassandra repo")
	return &CassandraRepository{session: session}, nil
}

// FindAll fetches all the basic user definitions available
func (r *CassandraRepository) FindAll() ([]*Organization, error) {
	log.Debugf("finding all users servers")

	var results []*Organization

	iter := r.session.GetSession().Query("SELECT username, organization, password, priority, content_per_day FROM organization").Iter()

	var username string
	var comp string
	var pass string
	var priority int
	var contentPerDay int

	err := iter.ScanAndClose(func() bool {
		var organization Organization
		organization.Username = username
		organization.Organization = comp
		organization.Password = pass
		organization.Priority = priority
		organization.ContentPerDay = contentPerDay
		results = append(results, &organization)
		return true
	}, &username, &comp, &pass, &priority, &contentPerDay)

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
		"SELECT username, organization, password, priority, content_per_day "+
			"FROM organization "+
			"WHERE username = ?",
		username).Scan(&organization.Username, &organization.Organization, &organization.Password, &organization.Priority, &organization.ContentPerDay)

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
			"password = ?, "+
			"priority = ?, "+
			"content_per_day = ? "+
			"WHERE username = ?",
		organization.Organization, hash, organization.Priority, organization.ContentPerDay, organization.Username).Exec()

	if err != nil {
		log.Errorf("error saving organization user %s: %v", organization.Username, err)
	} else {
		log.Debugf("successfully saved organization user %s", organization.Username)
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
