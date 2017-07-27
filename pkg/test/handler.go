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

// RecoveryHandler represents the recovery handler
func RecoveryHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	errors.Handler(w, err)
}
