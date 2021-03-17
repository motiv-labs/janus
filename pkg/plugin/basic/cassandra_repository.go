package basic

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session wrapper.Holder
}

func NewCassandraRepository(session wrapper.Holder) (*CassandraRepository, error) {
	log.Debugf("getting new basic cassandra repo")
	return &CassandraRepository{session: session}, nil

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

	log.Infof("user is: %v", *user)

	err := r.session.GetSession().Query(
		"UPDATE user " +
			"SET password = ? " +
			"WHERE username = ?",
		user.Password, user.Username).Exec()

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

func parseDSN(dsn string) (clusterHost string, systemKeyspace string, appKeyspace string, connectionTimeout int) {
	trimDSN := strings.TrimSpace(dsn)
	log.Infof("trimDSN: %s", trimDSN)
	if len(trimDSN) == 0 {
		return "", "", "", 0
	}
	// split each `:`
	splitDSN := strings.Split(trimDSN, "/")
	// list of info
	for i, v := range splitDSN {
		log.Infof("splitDSN i is %d and v is %s", i, v)
		// start at 1 because the dsn path comes in with a leading /
		switch i {
		case 1:
			clusterHost = v
		case 2:
			systemKeyspace = v
		case 3:
			appKeyspace = v
		case 4:
			timeout, err := strconv.Atoi(v)
			if err != nil {
				log.Error("timeout is not an int")
				timeout = 0
			}
			connectionTimeout = timeout
		}
	}
	return clusterHost, systemKeyspace, appKeyspace, connectionTimeout
}
