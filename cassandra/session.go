package cassandra

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
)
const (
	// Cassandra cluster host
	ClusterHostName = "db"
	// System keyspace
	SystemKeyspace = "system"
	// Github taxi dispatcher keyspace
	AppKeyspace = "janus"
	// default timeout
	Timeout = 300
)

// SessionHolder holds our connection to Cassandra
var sessionHolder wrapper.Holder

func GetSession() wrapper.SessionInterface {
	return sessionHolder.GetSession()
}

func SetSessionHolder(holder wrapper.Holder) {
	sessionHolder = holder
}
