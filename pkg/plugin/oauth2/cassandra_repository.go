package oauth2

import (
	"encoding/json"
	cass "github.com/hellofresh/janus/cassandra"
	cassmod "github.com/motiv-labs/cassandra"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	session cassmod.Holder
}

func NewCassandraRepository(dsn string) (*CassandraRepository, error) {
	log.Debugf("getting new oauth cassandra repo")
	span := opentracing.StartSpan("NewCassandraRepository")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	// parse the dsn string for the cluster host, system key space, app key space and connection timeout.
	log.Infof("dsn is %s", dsn)
	clusterHost, systemKeyspace, appKeyspace, connectionTimeout := parseDSN(dsn)
	if clusterHost == "" {
		clusterHost = cass.ClusterHostName
	}
	if systemKeyspace == "" {
		systemKeyspace = cass.SystemKeyspace
	}
	if appKeyspace == "" {
		appKeyspace = cass.AppKeyspace
	}
	if connectionTimeout == 0 {
		connectionTimeout = cass.Timeout
	}

	// Wait for Cassandra to start, setup Cassandra keyspace if required
	cassmod.Initialize(cass.ClusterHostName, cass.SystemKeyspace, cass.AppKeyspace, cass.Timeout*time.Second, span)

	// Getting a cassandra connection initializer
	initializer := cassmod.New(cass.ClusterHostName, cass.AppKeyspace, span)

	// Starting a new cassandra session
	sessionHolder, err := initializer.NewSession(span)
	if err != nil {
		panic(err)
	}
	// set oauth cassandra repo session
	cass.SetSessionHolder(sessionHolder)

	return &CassandraRepository{
		session: sessionHolder,
	}, nil

}

// FindAll fetches all the OAuth Servers available
func (r *CassandraRepository) FindAll() ([]*OAuth, error) {
	span := opentracing.StartSpan("FindAll")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("finding all oauth servers")

	var results []*OAuth

	iter := r.session.GetSession(span).Query(span,
		"SELECT name, oauth FROM oauth").Iter(span)

	var savedDef string
	var oauth *OAuth

	for iter.Scan(span, &savedDef) {
		err := json.Unmarshal([]byte(savedDef), &oauth)
		if err != nil {
			log.Errorf("error trying to unmarshal oauth json: %v", err)
			return nil, err
		}
		results = append(results, oauth)
	}

	err := iter.Close(span)
	if err != nil {
		log.Errorf("error getting all oauths: %v", err)
	}
	return results, err
}

// FindByName find an OAuth Server by name
func (r *CassandraRepository) FindByName(name string) (*OAuth, error) {
	span := opentracing.StartSpan("add")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("finding: %s", name)

	var oauth *OAuth

	err := r.session.GetSession(span).Query(span,
		"SELECT oauth = ? " +
			"FROM oauth" +
			"WHERE name = ?",
		oauth, name).Exec(span)

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
	span := opentracing.StartSpan("add")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("adding: %s", oauth.Name)

	log.Infof("oauth is: %v", *oauth)

	saveOauth, err := json.Marshal(oauth)
	if err != nil {
		log.Errorf("error marshaling oauth: %v", err)
		return err
	}
	err = r.session.GetSession(span).Query(span,
		"UPDATE oauth " +
			"SET oauth = ? " +
			"WHERE name = ?",
		saveOauth, oauth.Name).Exec(span)

	if err != nil {
		log.Errorf("error saving oauth %s: %v", oauth.Name, err)
	} else {
		log.Debugf("successfully saved oauth %s", oauth.Name)
	}

	return err
}

// Save saves OAuth Server to the repository
func (r *CassandraRepository) Save(oauth *OAuth) error {
	span := opentracing.StartSpan("add")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("adding: %s", oauth.Name)

	log.Infof("oauth is: %v", *oauth)

	saveOauth, err := json.Marshal(oauth)
	if err != nil {
		log.Errorf("error marshaling oauth: %v", err)
		return err
	}
	err = r.session.GetSession(span).Query(span,
		"UPDATE oauth " +
			"SET oauth = ? " +
			"WHERE name = ?",
		saveOauth, oauth.Name).Exec(span)

	if err != nil {
		log.Errorf("error saving oauth %s: %v", oauth.Name, err)
	} else {
		log.Debugf("successfully saved oauth %s", oauth.Name)
	}

	return err
}

// Remove removes an OAuth Server from the repository
func (r *CassandraRepository) Remove(name string) error {
	span := opentracing.StartSpan("remove")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("removing: %s", name)

	err := r.session.GetSession(span).Query(span,
		"DELETE FROM oauth WHERE name = ?", name).Exec(span)

	if err != nil {
		log.Errorf("error removing oauth %s: %v", name, err)
	} else {
		log.Debugf("successfully removed oauth %s", name)
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

