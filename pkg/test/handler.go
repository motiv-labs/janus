package test

import (
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
)

// Ping is a test handler
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("OK\n"))
}

// FailWith is a test handler that fails
func FailWith(statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})
}

// RecoveryHandler represents the recovery handler
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	errors.Handler(w, r, err)
}
