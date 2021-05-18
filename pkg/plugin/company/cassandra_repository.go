package company

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
	log "github.com/sirupsen/logrus"
)

// Repository represents an user repository
type Repository interface {
	FindAll() ([]*Company, error)
	FindByUsername(username string) (*Company, error)
	Add(company *Company) error
	Remove(username string) error
}

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session wrapper.Holder
}

func NewCassandraRepository(session wrapper.Holder) (*CassandraRepository, error) {
	log.Debugf("getting new company cassandra repo")
	return &CassandraRepository{session: session}, nil
}

// FindAll fetches all the basic user definitions available
func (r *CassandraRepository) FindAll() ([]*Company, error) {
	log.Debugf("finding all users servers")

	var results []*Company

	iter := r.session.GetSession().Query("SELECT username, company FROM company").Iter()

	var username string
	var comp string

	err := iter.ScanAndClose(func() bool {
		var company Company
		company.Username = username
		company.Company = comp
		results = append(results, &company)
		return true
	}, &username, &comp)

	if err != nil {
		log.Errorf("error getting all company users: %v", err)
	}
	return results, err

}

// FindByUsername find an user by username
// returns ErrUserNotFound when a user is not found.
func (r *CassandraRepository) FindByUsername(username string) (*Company, error) {
	log.Debugf("finding: %s", username)

	var company Company

	err := r.session.GetSession().Query(
		"SELECT username, company " +
			"FROM company " +
			"WHERE username = ?",
		username).Scan(&company.Username, &company.Company)

	if err.Error() == "not found"{
		log.Debugf("company not found")
		err = ErrUserNotFound
	} else if err != nil {
		log.Errorf("error selecting company user %s: %v", username, err)
	} else {
		log.Debugf("successfully found company user %s", username)
	}

	return &company, err
}

// Add adds an user to the repository
func (r *CassandraRepository) Add(company *Company) error {
	log.Debugf("adding: %s", company.Username)

	err := r.session.GetSession().Query(
		"UPDATE company " +
			"SET company = ? " +
			"WHERE username = ?",
		company.Company, company.Username).Exec()

	if err != nil {
		log.Errorf("error saving company user %s: %v", company.Username, err)
	} else {
		log.Debugf("successfully saved company user %s", company.Username)
	}

	return err
}

// Remove an user from the repository
func (r *CassandraRepository) Remove(username string) error {
	log.Debugf("removing: %s", username)

	err := r.session.GetSession().Query(
		"DELETE FROM user WHERE username = ?", username).Exec()

	if err != nil {
		log.Errorf("error removing company user %s: %v", username, err)
	} else {
		log.Debugf("successfully removed company user %s", username)
	}

	return err
}
