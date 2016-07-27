package main

import (
	log "github.com/Sirupsen/logrus"
)

type Logger struct {
	*log.Entry
}

func createContextualLogger(spec *APISpec) *Logger {
	fields := log.Fields{
		"id":   spec.ID,
		"name": spec.Name,
	}

	return &Logger{log.WithFields(fields)}
}
