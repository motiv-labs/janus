package cassandra

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
)

const (
	// ClusterHostName representing Cassandra cluster host
	ClusterHostName = "db"
	// SystemKeyspace is system keyspace
	SystemKeyspace = "system"
	// AppKeyspace Github taxi dispatcher keyspace
	AppKeyspace = "janus"
	// Timeout represents default timeout
	Timeout = 300
)

// SessionHolder holds our connection to Cassandra
var sessionHolder wrapper.Holder

// GetSession returns session
func GetSession() wrapper.SessionInterface {
	return sessionHolder.GetSession()
}

// SetSessionHolder setter
func SetSessionHolder(holder wrapper.Holder) {
	sessionHolder = holder
}
