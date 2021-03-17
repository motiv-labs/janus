package cassandra

import (
	"github.com/motiv-labs/cassandra"
	"github.com/opentracing/opentracing-go"
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
var sessionHolder cassandra.Holder

func GetSession() cassandra.SessionInterface {
	span := opentracing.StartSpan("GetSession")
	defer span.Finish()
	span.SetTag("Package", "cassandra")
	return sessionHolder.GetSession(span)
}

func SetSessionHolder(holder cassandra.Holder) {
	sessionHolder = holder
}
