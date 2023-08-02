package wrapper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

// Schema file to create keyspace if required
const (
	schemaDefaultPath                  = "/usr/local/bin"
	schemaDefaultFileName              = "schema.sql"
	defaultCassandraClusterConsistency = gocql.Quorum
	defaultUsername                    = ""
	defaultPassword                    = ""
	defaultSSLCert                     = ""

	envCassandraClusterConsistency = "CLUSTER_CONSISTENCY"
	envCassandraSchemaPath         = "CASSANDRA_SCHEMA_PATH"
	envCassandraSchemaFileName     = "CASSANDRA_SCHEMA_FILE_NAME"
	envUsername                    = "CASSANDRA_USERNAME"
	envPassword                    = "CASSANDRA_PASSWORD"
	envSSLCert                     = "CASSANDRA_SSL_CERT"
)

var schemaPath = "/usr/local/bin"
var schemaFileName = "schema.sql"
var clusterConsistency = gocql.Quorum
var username = ""
var password = ""
var sslCert = ""

// Package level initialization.
//
// init functions are automatically executed when the programs starts
func init() {

	// reading and setting up environment variables
	schemaPath = getenv(envCassandraSchemaPath, schemaDefaultPath)
	schemaFileName = getenv(envCassandraSchemaFileName, schemaDefaultFileName)
	clusterConsistency = checkConsistency(getenv(envCassandraClusterConsistency, defaultCassandraClusterConsistency.String()))
	username = getenvnolog(envUsername, defaultUsername)
	password = getenvnolog(envPassword, defaultPassword)
	sslCert = getenvnolog(envSSLCert, defaultSSLCert)

	log.Debugf("Got schema path: %s", schemaPath)
	log.Debugf("Got schema file name: %s", schemaFileName)
	log.Debugf("Got cluster consistency: %s", clusterConsistency)
	log.Debugf("Got username: %s", username)
}

// sessionInitializer is an initializer for a cassandra session
type sessionInitializer struct {
	clusterHostName     string
	clusterHostUsername string
	clusterHostPassword string
	clusterHostSSLCert  string
	keyspace            string
	consistency         gocql.Consistency
}

// sessionHolder stores a cassandra session
type sessionHolder struct {
	session SessionInterface
}

// New return a cassandra session Initializer
func New(clusterHostName, keyspace string) Initializer {
	log.Debugf("in new")

	return sessionInitializer{
		clusterHostName:     clusterHostName,
		clusterHostUsername: username,
		clusterHostPassword: password,
		clusterHostSSLCert:  sslCert,
		keyspace:            keyspace,
		consistency:         clusterConsistency,
	}
}

// Initialize waits for a Cassandra session, initializes Cassandra keyspace and creates tables if required.
// NOTE: Needs to be called only once on the app startup, won't fail if it is called multiple times but is not necessary.
//
// Params:
//
//	clusterHostName: Cassandra cluster host
//	systemKeyspace: System keyspace
//	appKeyspace: Application keyspace
//	connectionTimeout: timeout to get the connection
func Initialize(clusterHostName, systemKeyspace, appKeyspace string, connectionTimeout time.Duration) {
	log.Debug("Setting up cassandra db")
	connectionHolder, err := loop(connectionTimeout, New(clusterHostName, systemKeyspace), "cassandra-db")
	if err != nil {
		log.Fatalf("error connecting to Cassandra db: %v", err)
		panic(err)
	}
	defer connectionHolder.CloseSession()

	log.Debug("Setting up cassandra keyspace")
	err = createAppKeyspaceIfRequired(clusterHostName, systemKeyspace, appKeyspace)
	if err != nil {
		log.Fatalf("error creating keyspace for Cassandra db: %v", err)
		panic(err)
	}

	log.Info("Cassandra keyspace has been set up")
}

// NewSession starts a new cassandra session for the given keyspace
// NOTE: It is responsibility of the caller to close this new session.
//
// Returns a session Holder for the session, or an error if can't start the session
func (i sessionInitializer) NewSession() (Holder, error) {
	session, err := newKeyspaceSession(i.clusterHostName, i.keyspace,
		i.clusterHostUsername, i.clusterHostPassword, i.clusterHostSSLCert, i.consistency,
		600*time.Millisecond)
	if err != nil {
		log.Errorf("error starting Cassandra session for the cluster hostname: %s and keyspace: %s - %v",
			i.clusterHostName, i.keyspace, err)
		return nil, err
	}
	sessionRetry := sessionRetry{session}
	connectionHolder := sessionHolder{sessionRetry}
	return connectionHolder, nil
}

// GetSession returns the stored cassandra session
func (holder sessionHolder) GetSession() SessionInterface {
	return holder.session
}

// CloseSession closes the cassandra session
func (holder sessionHolder) CloseSession() {
	holder.session.Close()
}

// newKeyspaceSession returns a new session for the given keyspace
func newKeyspaceSession(clusterHostName, keyspace, username, password, sslCert string, clusterConsistency gocql.Consistency, clusterTimeout time.Duration) (*gocql.Session, error) {
	log.Infof("Creating new cassandra session for cluster hostname: %s and keyspace: %s", clusterHostName, keyspace)
	cluster := gocql.NewCluster(clusterHostName)
	cluster.Keyspace = keyspace
	cluster.Timeout = clusterTimeout
	if username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}
	if sslCert != "" {
		cluster.SslOpts = &gocql.SslOptions{
			CaPath: sslCert,
		}
	}
	cluster.Consistency = clusterConsistency
	return cluster.CreateSession()
}

// createAppKeyspaceIfRequired creates the keyspace for the app if it doesn't exist
func createAppKeyspaceIfRequired(clusterHostName, systemKeyspace, appKeyspace string) error {
	// Getting the schema file if exist
	stmtList, err := getStmtsFromFile(path.Join(schemaPath, schemaFileName))
	if err != nil {
		return err
	}
	if stmtList == nil { // Didn't fail but returned nil, probably the file does not exist
		return nil
	}

	log.Info("about to create a session with a 5 minute timeout to allow for all schema creation")
	session, err := newKeyspaceSession(clusterHostName, systemKeyspace, username, password, sslCert, clusterConsistency, 5*time.Minute)
	if err != nil {
		return err
	}
	currentKeyspace := systemKeyspace

	var sessionList []*gocql.Session
	defer func() {
		for _, s := range sessionList {
			if s != nil && !s.Closed() {
				s.Close()
			}
		}
	}()

	log.Debugf("Creating new keyspace if required: %s", appKeyspace)

	for _, stmt := range stmtList {
		log.Debugf("Executing statement: %s", stmt)
		// New session for use statement
		newKeyspace, isCaseSensitive := getKeyspaceNameFromUseStmt(stmt)
		if newKeyspace != "" {
			if (isCaseSensitive && newKeyspace != currentKeyspace) || (!isCaseSensitive &&
				strings.ToLower(newKeyspace) != strings.ToLower(currentKeyspace)) {
				log.Infof("about to create a session with a 5 minute timeout to set keyspace: %s", newKeyspace)
				session, err = newKeyspaceSession(clusterHostName, newKeyspace, username, password, sslCert, clusterConsistency, 5*time.Minute) //5 minutes
				if err != nil {
					return err
				}
				currentKeyspace = newKeyspace
				sessionList = append(sessionList, session)
				log.Debugf("Changed to new keyspace: %s", newKeyspace)
			}
			continue
		}

		// execute statement
		err = session.Query(stmt).Exec()
		if err != nil {
			log.Errorf("statement error: %v", err)
			return err
		}
		log.Debug("Statement executed")
	}

	log.Debugf("app keyspace set to: %s", appKeyspace)
	return nil
}

// getStmtsFromFile extracts CQL statements from the file
func getStmtsFromFile(fileName string) ([]string, error) {
	// Verify first if the file exist
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) { // Does not exist
			log.Warnf("no schema file [%s] found initializing Cassandra.", fileName)
			return nil, nil
		}
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	pattern := regexp.MustCompile(`(?ms)([^"';]*?)("(?:[^"]|"")*"|'(?:[^']|'')*'|\$\$.*?\$\$|(/\*.*?\*/)|((?:--|//).*?$)|;\n?)`)

	var stmtList []string

	i := 0
	contentLength := len(content)
	var stmt bytes.Buffer
	for i < contentLength {
		subIndexes := pattern.FindSubmatchIndex(content[i:])
		if len(subIndexes) > 0 {
			end := subIndexes[1]
			stmt.Write(getMatch(content, i, subIndexes, 2, 3))
			stmtTail := getMatch(content, i, subIndexes, 4, 5)
			comment := getMatch(content, i, subIndexes, 6, 7)
			lineComment := getMatch(content, i, subIndexes, 8, 9)
			if comment == nil && lineComment == nil {
				if stmtTail != nil && string(bytes.TrimSpace(stmtTail)) == ";" {
					stmtList = append(stmtList, stmt.String())
					stmt.Reset()
				} else {
					stmt.Write(stmtTail)
				}
			}
			i = i + end
		} else {
			break
		}
	}

	return stmtList, nil

}

// getMatch returns the matched substring if there's a match, nil otherwise
func getMatch(src []byte, base int, match []int, start int, end int) []byte {
	if match[start] >= 0 {
		return src[base+match[start] : base+match[end]]
	}

	return nil
}

// getKeyspaceNameFromUseStmt return keyspace name for use statement
func getKeyspaceNameFromUseStmt(stmt string) (string, bool) {
	pattern := regexp.MustCompile(`(?ms)[Uu][Ss][Ee]\s+("(?:[^"]|"")+"|\w+)`)
	if pattern.MatchString(stmt) {
		match := pattern.FindStringSubmatch(stmt)
		if len(match) > 1 {
			keyspace := match[1]
			caseSensitive := false
			if strings.HasPrefix(keyspace, "\"") && strings.HasSuffix(keyspace, "\"") {
				keyspace = strings.Trim(keyspace, "\"")
				caseSensitive = true
			}
			return keyspace, caseSensitive
		}
	}
	return "", false
}

// Loop is a loop that tries to get a connection until a timeout is reached
//
// Params:
//
//	timeout: timeout to get the connection
//	initializer : initializer to start the session
//	connectionHost : name of host for the connection
//
// Returns a session Holder to store the session, or an error if the timeout was reached
func loop(timeout time.Duration, initializer Initializer, connectionHost string) (Holder, error) {
	log.Debugf("Connection loop to connect to %s, timeout to use: %s", connectionHost, timeout)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("connection to %s failed after %s timeout", connectionHost, timeout)

		case <-ticker.C:
			log.Infof("Trying to connect to: %s", connectionHost)
			connectionHolder, err := initializer.NewSession()
			if err == nil {
				log.Infof("Successful connection to: %s", connectionHost)
				return connectionHolder, nil
			}
			log.Infof("Trying to connect to %s, failed attempt: %v", connectionHost, err)
		}
	}

}

// getenv get a string value from an environment variable
// or return the given default value if the environment variable is not set
//
// Params:
//
//	envVariable : environment variable
//	defaultValue : value to return if environment variable is not set
//
// Returns the string value for the specified variable
func getenv(envVariable string, defaultValue string) string {

	log.Debugf("Setting value for: %s", envVariable)
	returnValue := defaultValue
	log.Debugf("Default value for %s : %s", envVariable, defaultValue)
	envStr := os.Getenv(envVariable)
	if envStr != "" {
		returnValue = envStr
		log.Debugf("Init value for %s set to: %s", envVariable, envStr)
	}

	return returnValue
}

// getenvnolog get a string value from an environment variable
// or return the given default value if the environment variable is not set
//
// Params:
//
//	envVariable : environment variable
//	defaultValue : value to return if environment variable is not set
//
// Returns the string value for the specified variable
func getenvnolog(envVariable string, defaultValue string) string {

	log.Debugf("Setting value for: %s", envVariable)
	returnValue := defaultValue
	log.Debugf("Default value for %s : %s", envVariable, defaultValue)
	envStr := os.Getenv(envVariable)
	if envStr != "" {
		returnValue = envStr
	}

	return returnValue
}

func checkConsistency(envVar string) gocql.Consistency {
	switch strings.ToLower(envVar) {
	case "any":
		log.Debugf("consistency set to any")
		return gocql.Any
	case "one":
		log.Debugf("consistency set to one")
		return gocql.One
	case "two":
		log.Debugf("consistency set to two")
		return gocql.Two
	case "three":
		log.Debugf("consistency set to three")
		return gocql.Three
	case "quorum":
		log.Debugf("consistency set to quorum")
		return gocql.Quorum
	case "all":
		log.Debugf("consistency set to all")
		return gocql.All
	case "localquorum":
		log.Debugf("consistency set to local quorum")
		return gocql.LocalQuorum
	case "eachquorum":
		log.Debugf("consistency set to each quorum")
		return gocql.EachQuorum
	case "localone":
		log.Debugf("consistency set to local one")
		return gocql.LocalOne
	default:
		log.Debugf("consistency set to %s", clusterConsistency.String())
		return clusterConsistency
	}
}
