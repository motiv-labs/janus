package opentracing

import "github.com/sirupsen/logrus"

type jaegerLoggerAdapter struct {
	log *logrus.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.log.Error(msg)
}

// Infof adapts infof messages to logrus
func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.log.Infof(msg, args...)
}
