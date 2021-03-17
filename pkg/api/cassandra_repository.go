package api

import (
	cass "github.com/hellofresh/janus/cassandra"
	cassmod "github.com/motiv-labs/cassandra"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// CassandraRepository represents a cassandra repository
type CassandraRepository struct {
	//TODO: we need to expose this so the plugins can use the same session. We should abstract mongo DB and provide
	// the plugins with a simple interface to search, insert, update and remove data from whatever backend implementation
	session cassmod.Holder
	refreshTime time.Duration
}

func NewCassandraRepository(dsn string, refreshTime time.Duration) (*CassandraRepository, error) {
	log.Debugf("getting new cassandra repo")
	//span := opentracing.StartSpan("NewCassandraRepository")
	//defer span.Finish()
	//span.SetTag("Interface", "CassandraRepository")

	// parse the dsn string for the cluster host, system key space, app key space and connection timeout.
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
	cassmod.Initialize(cass.ClusterHostName, cass.SystemKeyspace, cass.AppKeyspace, cass.Timeout*time.Second, nil)

	// Getting a cassandra connection initializer
	initializer := cassmod.New(cass.ClusterHostName, cass.AppKeyspace, nil)

	// Starting a new cassandra session
	sessionHolder, err := initializer.NewSession(nil)
	if err != nil {
		panic(err)
	}
	// Global session for Janus
	cass.SetSessionHolder(sessionHolder)

	return &CassandraRepository{
		session: sessionHolder,
		refreshTime: refreshTime,
	}, nil

}

func (r *CassandraRepository) Close() error {
	//span := opentracing.StartSpan("Close")
	//defer span.Finish()
	//span.SetTag("Interface", "CassandraRepository")
	// Close the session
	r.session.CloseSession(nil)
	return nil
}

// FindAll fetches all the API definitions available
func (r *CassandraRepository) FindAll() ([]*Definition, error) {
	//span := opentracing.StartSpan("FindAll")
	//defer span.Finish()
	//span.SetTag("Interface", "CassandraRepository")

	var results []*Definition
	// todo fill in the select with the actual column names
	iter := r.session.GetSession(nil).Query(nil,
		"SELECT name, definition FROM api_definition").Iter(nil)

	var savedDef string
	var definition *Definition

	for iter.Scan(nil, definition) {
		err := definition.UnmarshalJSON([]byte(savedDef))
		if err != nil {
			log.Errorf("error trying to unmarshal definition json: %v", err)
			return nil, err
		}
		results = append(results, definition)
	}

	err := iter.Close(nil)
	if err != nil {
		log.Errorf("error getting all definitions: %v", err)
	}
	return results, err
}

// Add adds an API definition to the repository
func (r *CassandraRepository) add(definition *Definition) error {
	//span := opentracing.StartSpan("add")
	//defer span.Finish()
	//span.SetTag("Interface", "CassandraRepository")

	isValid, err := definition.Validate()
	if false == isValid && err != nil {
		log.WithError(err).Error("Validation errors")
		return err
	}

	// todo I might need to marshal the definition before saving it.
	err = r.session.GetSession(nil).Query(nil,
		"UPDATE api_definition " +
		"SET name = ?, " +
		"definition = ? " +
		"WHERE name = ?",
		definition.Name, definition).Exec(nil)

	if err != nil {
		log.Errorf("error saving definition %s: %v", definition.Name, err)
	} else {
		log.Debugf("successfully saved definition %s", definition.Name)
	}

	return err
}

// Remove removes an API definition from the repository
func (r *CassandraRepository) remove(name string) error {
	//span := opentracing.StartSpan("remove")
	//defer span.Finish()
	//span.SetTag("Interface", "CassandraRepository")

	// todo I might need to marshal the definition before saving it.
	err := r.session.GetSession(nil).Query(nil,
		"DELETE FROM api_definition WHERE name = ?", name).Exec(nil)

	if err != nil {
		log.Errorf("error saving definition %s: %v", name, err)
	} else {
		log.Debugf("successfully saved definition %s", name)
	}

	return err
}

func parseDSN(dsn string) (clusterHost string, systemKeyspace string, appKeyspace string, connectionTimeout int) {
	trimDSN := strings.TrimSpace(dsn)
	if len(trimDSN) == 0 {
		return "", "", "", 0
	}
	// split each `:`
	splitDSN := strings.Split(trimDSN, ":")
	// list of info
	for i, v := range splitDSN {
		switch i {
		case 0:
			clusterHost = v
		case 1:
			systemKeyspace = v
		case 2:
			appKeyspace = v
		case 3:
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
