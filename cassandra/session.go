package cassandra

import (
	"github.com/motiv-labs/cassandra"
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
	//span := opentracing.StartSpan("GetSession", opentracing.ChildOf(parentSpan.Context()))
	//defer span.Finish()
	//span.SetTag("Package", "cassandra")
	return sessionHolder.GetSession(nil)
}

func SetSessionHolder(holder cassandra.Holder) {
	//span := opentracing.StartSpan("SetSessionHolder", opentracing.ChildOf(parentSpan.Context()))
	//defer span.Finish()
	//span.SetTag("Package", "cassandra")
	sessionHolder = holder
}
