package basic

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	log "github.com/sirupsen/logrus"
)

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session wrapper.Holder
	hash encrypt.Hash
}

func NewCassandraRepository(session wrapper.Holder) (*CassandraRepository, error) {
	log.Debugf("getting new basic cassandra repo")
	return &CassandraRepository{session: session, hash: encrypt.Hash{}}, nil

}

// FindAll fetches all the basic user definitions available
func (r *CassandraRepository) FindAll() ([]*User, error) {
	log.Infof("finding all users servers")

	var results []*User

	iter := r.session.GetSession().Query("SELECT username, password FROM user").Iter()

	var username string
	var password string

	for iter.Scan(&username, &password) {
		var user User
		user.Username = username
		user.Password = password
		results = append(results, &user)
	}

	err := iter.Close()
	if err != nil {
		log.Errorf("error getting all oauths: %v", err)
	}
	return results, err

}

// FindByUsername find an user by username
// returns ErrUserNotFound when a user is not found.
func (r *CassandraRepository) FindByUsername(username string) (*User, error) {
	log.Infof("finding: %s", username)

	var user User

	iter := r.session.GetSession().Query(
		"SELECT username, password " +
			"FROM user " +
			"WHERE username = ?",
		username).Iter()

	iter.Scan(&user.Username, &user.Password)
	err := iter.Close()

	if err != nil {
		log.Errorf("error selecting user %s: %v", username, err)
	} else {
		log.Debugf("successfully found user %s", username)
		err = ErrUserNotFound
	}

	return &user, err
}

// Add adds an user to the repository
func (r *CassandraRepository) Add(user *User) error {
	log.Infof("adding: %s", user.Username)

	hash, err := r.hash.Generate(user.Password)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return err
	}

	err = r.session.GetSession().Query(
		"UPDATE user " +
			"SET password = ? " +
			"WHERE username = ?",
		hash, user.Username).Exec()

	if err != nil {
		log.Errorf("error saving user %s: %v", user.Username, err)
	} else {
		log.Debugf("successfully saved user %s", user.Username)
	}

	return err
}

// Remove an user from the repository
func (r *CassandraRepository) Remove(username string) error {
	log.Infof("removing: %s", username)

	err := r.session.GetSession().Query(
		"DELETE FROM user WHERE username = ?", username).Exec()

	if err != nil {
		log.Errorf("error removing user %s: %v", username, err)
	} else {
		log.Debugf("successfully removed user %s", username)
	}

	return err
}
