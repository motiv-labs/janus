package wrapper

import (
	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	defaultCassandraRetryAttempts           = "3"
	defaultCassandraSecondsToSleepIncrement = "1"

	envCassandraAttempts                = "CASSANDRA_RETRY_ATTEMPTS"
	envCassandraSecondsToSleepIncrement = "CASSANDRA_SECONDS_SLEEP_INCREMENT"
)

var cassandraRetryAttempts = 3
var cassandraSecondsToSleepIncrement = 1

// Package level initialization.
//
// init functions are automatically executed when the programs starts
func init() {
	cassandraRetryAttempts, err := strconv.Atoi(getenv(envCassandraAttempts, defaultCassandraRetryAttempts))
	if err != nil {
		log.Errorf("error trying to get CASSANDRA_RETRY_ATTEMPTS value: %s",
			getenv(envCassandraAttempts, defaultCassandraRetryAttempts))
		cassandraRetryAttempts = 3
	}

	cassandraSecondsToSleep, err := strconv.Atoi(getenv(envCassandraSecondsToSleepIncrement, defaultCassandraSecondsToSleepIncrement))
	if err != nil {
		log.Errorf("error trying to get CASSANDRA_SECONDS_SLEEP value: %s",
			getenv(envCassandraSecondsToSleepIncrement, defaultCassandraSecondsToSleepIncrement))
		cassandraSecondsToSleepIncrement = 1
	}

	log.Debugf("got cassandraRetryAttempts: %d", cassandraRetryAttempts)
	log.Debugf("got cassandraSecondsToSleepIncrement: %d", cassandraSecondsToSleep)
}

// queryRetry is an implementation of QueryInterface
type queryRetry struct {
	goCqlQuery *gocql.Query
}

// iterRetry is an implementation of IterInterface
type iterRetry struct {
	goCqlIter *gocql.Iter
}

// Exec wrapper to retry around gocql Exec(). We have a retry approach in place + incremental approach used. For example:
// First time it will wait 1 second, second time 2 seconds, ... It will depend on the values for retries and seconds to wait.
func (q queryRetry) Exec() error {
	log.Debug("running queryRetry Exec() method")

	retryAttempts := cassandraRetryAttempts
	secondsToSleep := 0

	var err error

	attempts := 1
	for attempts <= retryAttempts {
		//we will try to run the method several times until attempts is met
		err = q.goCqlQuery.Exec()
		if err != nil {
			log.Warnf("error when running Exec(): %v, attempt: %d / %d", err, attempts, retryAttempts)

			// incremental sleep
			secondsToSleep = secondsToSleep + cassandraSecondsToSleepIncrement

			log.Warnf("sleeping for %d second", secondsToSleep)

			time.Sleep(time.Duration(secondsToSleep) * time.Second)
		} else {
			// in case the error is nil, we stop and return
			return err
		}

		attempts = attempts + 1
	}

	return err
}

// Scan wrapper to retry around gocql Scan(). We have a retry approach in place + incremental approach used. For example:
// First time it will wait 1 second, second time 2 seconds, ... It will depend on the values for retries and seconds to wait.
func (q queryRetry) Scan(dest ...interface{}) error {
	log.Debug("running queryRetry Scan() method")

	retries := cassandraRetryAttempts
	secondsToSleep := 0

	var err error

	attempts := 1
	for attempts <= retries {
		//we will try to run the method several times until attempts is met
		err = q.goCqlQuery.Scan(dest...)
		if err != nil {
			log.Warnf("error when running Scan(): %v, attempt: %d / %d", err, attempts, retries)

			// incremental sleep
			secondsToSleep = secondsToSleep + cassandraSecondsToSleepIncrement

			log.Warnf("sleeping for %d second", secondsToSleep)

			log.Warnf("sleeping for %d second", secondsToSleep)
			time.Sleep(time.Duration(secondsToSleep) * time.Second)
		} else {
			// in case the error is nil, we stop and return
			return err
		}

		attempts = attempts + 1
	}

	return err
}

// Iter just a wrapper to be able to call this method
func (q queryRetry) Iter() IterInterface {
	log.Debug("running queryRetry Iter() method")

	return iterRetry{goCqlIter: q.goCqlQuery.Iter()}
}

// PageState just a wrapper to be able to call this method
func (q queryRetry) PageState(state []byte) QueryInterface {
	log.Debug("running queryRetry PageState() method")

	return queryRetry{goCqlQuery: q.goCqlQuery.PageState(state)}
}

// PageSize just a wrapper to be able to call this method
func (q queryRetry) PageSize(n int) QueryInterface {
	log.Debug("running queryRetry PageSize() method")

	return queryRetry{goCqlQuery: q.goCqlQuery.PageSize(n)}
}

//
func (i iterRetry) Scan(dest ...interface{}) bool {
	log.Debug("running iterRetry Scan() method")

	return i.goCqlIter.Scan(dest...)
}

// WillSwitchPage is just a wrapper to be able to call this method
func (i iterRetry) WillSwitchPage() bool {
	log.Debug("running iterRetry Close() method")

	return i.goCqlIter.WillSwitchPage()
}

// PageState is just a wrapper to be able to call this method
func (i iterRetry) PageState() []byte {
	log.Debug("running iterRetry PageState() method")

	return i.goCqlIter.PageState()
}

// Close is just a wrapper to be able to call this method
func (i iterRetry) Close() error {
	log.Debug("running iterRetry Close() method")

	return i.goCqlIter.Close()
}

// ScanAndClose is a wrapper to retry around the gocql Scan() and Close().
// We have a retry approach in place + incremental approach used. For example:
// First time it will wait 1 second, second time 2 seconds, ... It will depend on the values for retries
// and seconds to wait.
func (i iterRetry) ScanAndClose(handle func() bool, dest ...interface{}) error {

	retries := cassandraRetryAttempts
	secondsToSleep := 0

	var err error

	attempts := 1
	for attempts <= retries {

		// Scan consumes the next row of the iterator and copies the columns of the
		// current row into the values pointed at by dest.
		for i.goCqlIter.Scan(dest...) {
			if !handle() {
				break
			}
		}

		// we will try to run the method several times until attempts is met
		if err = i.goCqlIter.Close(); err != nil {

			log.Warnf("error when running Close(): %v, attempt: %d / %d", err, attempts, retries)

			// incremental sleep
			secondsToSleep += cassandraSecondsToSleepIncrement

			log.Warnf("sleeping for %d second", secondsToSleep)

			time.Sleep(time.Duration(secondsToSleep) * time.Second)
		} else {
			// in case the error is nil, we stop and return
			return err
		}

		attempts++
	}

	return err
}
