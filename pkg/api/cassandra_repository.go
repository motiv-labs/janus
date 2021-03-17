package api

import (
	"context"
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
	//TODO: we need to expose this so the plugins can use the same session. We should abstract mongo DB and provide
	// the plugins with a simple interface to search, insert, update and remove data from whatever backend implementation
	session cassmod.Holder
	refreshTime time.Duration
}

func NewCassandraRepository(dsn string, refreshTime time.Duration) (*CassandraRepository, error) {
	log.Debugf("getting new api cassandra repo")
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
	// api cassandra repo session
	cass.SetSessionHolder(sessionHolder)

	return &CassandraRepository{
		session: sessionHolder,
		refreshTime: refreshTime,
	}, nil

}

func (r *CassandraRepository) Close() error {
	span := opentracing.StartSpan("Close")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")
	// Close the session
	r.session.CloseSession(span)
	return nil
}

// Listen watches for changes on the configuration
func (r *CassandraRepository) Listen(ctx context.Context, cfgChan <-chan ConfigurationMessage) {
	go func() {
		log.Debug("Listening for changes on the provider...")
		for {
			select {
			case cfg := <-cfgChan:
				switch cfg.Operation {
				case AddedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not add the configuration on the provider")
					}
				case UpdatedOperation:
					err := r.add(cfg.Configuration)
					if err != nil {
						log.WithError(err).Error("Could not update the configuration on the provider")
					}
				case RemovedOperation:
					err := r.remove(cfg.Configuration.Name)
					if err != nil {
						log.WithError(err).Error("Could not remove the configuration from the provider")
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Watch watches for changes on the database
func (r *CassandraRepository) Watch(ctx context.Context, cfgChan chan<- ConfigurationChanged) {
	t := time.NewTicker(r.refreshTime)
	go func(refreshTicker *time.Ticker) {
		defer refreshTicker.Stop()
		log.Debug("Watching Provider...")
		for {
			select {
			case <-refreshTicker.C:
				defs, err := r.FindAll()
				if err != nil {
					log.WithError(err).Error("Failed to get configurations on watch")
					continue
				}

				cfgChan <- ConfigurationChanged{
					Configurations: &Configuration{Definitions: defs},
				}
			case <-ctx.Done():
				return
			}
		}
	}(t)
}

// FindAll fetches all the API definitions available
func (r *CassandraRepository) FindAll() ([]*Definition, error) {
	span := opentracing.StartSpan("FindAll")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("finding all definitions")

	var results []*Definition

	iter := r.session.GetSession(span).Query(span,
		"SELECT definition FROM api_definition").Iter(span)

	var savedDef string
	var definition *Definition

	for iter.Scan(span, &savedDef) {
		err := json.Unmarshal([]byte(savedDef), &definition)
		if err != nil {
			log.Errorf("error trying to unmarshal definition json: %v", err)
			return nil, err
		}
		results = append(results, definition)
	}

	err := iter.Close(span)
	if err != nil {
		log.Errorf("error getting all definitions: %v", err)
	}
	return results, err
}

// Add adds an API definition to the repository
func (r *CassandraRepository) add(definition *Definition) error {
	span := opentracing.StartSpan("add")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("adding: %s", definition.Name)

	log.Infof("definition is: %v", *definition)

	isValid, err := definition.Validate()
	if false == isValid && err != nil {
		log.WithError(err).Error("Validation errors")
		return err
	}

	saveDef, err := json.Marshal(definition)
	if err != nil {
		log.Errorf("error marshaling oauth: %v", err)
		return err
	}

	err = r.session.GetSession(span).Query(span,
		"UPDATE api_definition " +
		"SET definition = ? " +
		"WHERE name = ?",
		saveDef, definition.Name).Exec(span)

	if err != nil {
		log.Errorf("error saving definition %s: %v", definition.Name, err)
	} else {
		log.Debugf("successfully saved definition %s", definition.Name)
	}

	return err
}

// Remove removes an API definition from the repository
func (r *CassandraRepository) remove(name string) error {
	span := opentracing.StartSpan("remove")
	defer span.Finish()
	span.SetTag("Interface", "CassandraRepository")

	log.Infof("removing: %s", name)

	err := r.session.GetSession(span).Query(span,
		"DELETE FROM api_definition WHERE name = ?", name).Exec(span)

	if err != nil {
		log.Errorf("error saving definition %s: %v", name, err)
	} else {
		log.Debugf("successfully saved definition %s", name)
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
