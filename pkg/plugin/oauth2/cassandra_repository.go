package oauth2

import (
	"encoding/json"
	"github.com/hellofresh/janus/cassandra/wrapper"
	log "github.com/sirupsen/logrus"
)

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session wrapper.Holder
}

func NewCassandraRepository(session wrapper.Holder) (*CassandraRepository, error) {
	return &CassandraRepository{session: session}, nil

}

// FindAll fetches all the OAuth Servers available
func (r *CassandraRepository) FindAll() ([]*OAuth, error) {
	log.Debugf("finding all oauth servers")

	var results []*OAuth

	iter := r.session.GetSession().Query("SELECT name, oauth FROM oauth").Iter()

	var savedOauth string

	err := iter.ScanAndClose(func() bool {
		var oauth *OAuth
		err := json.Unmarshal([]byte(savedOauth), &oauth)
		if err != nil {
			log.Errorf("error trying to unmarshal oauth json: %v", err)
			return false
		}
		results = append(results, oauth)
		return true
	}, &savedOauth)

	if err != nil {
		log.Errorf("error getting all oauths: %v", err)
	}
	return results, err
}

// FindByName find an OAuth Server by name
func (r *CassandraRepository) FindByName(name string) (*OAuth, error) {
	log.Debugf("finding: %s", name)

	var savedOauth string
	var oauth *OAuth

	err := r.session.GetSession().Query(
		"SELECT oauth " +
			"FROM oauth " +
			"WHERE name = ?",
		name).Scan(&savedOauth)

	err = json.Unmarshal([]byte(savedOauth), &oauth)

	if err != nil {
		log.Errorf("error selecting oauth %s: %v", name, err)
	} else {
		log.Debugf("successfully found oauth %s", name)
	}

	return oauth, err
}

// Add add a new OAuth Server to the repository
// Add is the same as Save because Cassandra only upserts and I didn't want to write an existence checker
func (r *CassandraRepository) Add(oauth *OAuth) error {
	log.Debugf("adding: %s", oauth.Name)

	saveOauth, err := json.Marshal(oauth)
	if err != nil {
		log.Errorf("error marshaling oauth: %v", err)
		return err
	}
	err = r.session.GetSession().Query(
		"UPDATE oauth " +
			"SET oauth = ? " +
			"WHERE name = ?",
		saveOauth, oauth.Name).Exec()

	if err != nil {
		log.Errorf("error saving oauth %s: %v", oauth.Name, err)
	} else {
		log.Debugf("successfully saved oauth %s", oauth.Name)
	}

	return err
}

// Save saves OAuth Server to the repository
func (r *CassandraRepository) Save(oauth *OAuth) error {
	log.Debugf("adding: %s", oauth.Name)

	saveOauth, err := json.Marshal(oauth)
	if err != nil {
		log.Errorf("error marshaling oauth: %v", err)
		return err
	}
	err = r.session.GetSession().Query(
		"UPDATE oauth " +
			"SET oauth = ? " +
			"WHERE name = ?",
		saveOauth, oauth.Name).Exec()

	if err != nil {
		log.Errorf("error saving oauth %s: %v", oauth.Name, err)
	} else {
		log.Debugf("successfully saved oauth %s", oauth.Name)
	}

	return err
}

// Remove removes an OAuth Server from the repository
func (r *CassandraRepository) Remove(name string) error {
	log.Debugf("removing: %s", name)

	err := r.session.GetSession().Query(
		"DELETE FROM oauth WHERE name = ?", name).Exec()

	if err != nil {
		log.Errorf("error removing oauth %s: %v", name, err)
	} else {
		log.Debugf("successfully removed oauth %s", name)
	}

	return err
}
