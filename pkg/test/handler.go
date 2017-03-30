package test

import "net/http"

// Ping is a test handler
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte("OK\n"))
}
