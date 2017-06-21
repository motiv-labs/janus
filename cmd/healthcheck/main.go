package main

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "8080", "port on localhost to check")
	flag.Parse()

	resp, err := http.Get("http://127.0.0.1:" + *port + "/status") // Note pointer dereference using *

	// If there is an error or non-200 status, exit with 1 signaling unsuccessful check.
	if err != nil || resp.StatusCode != 200 {
		os.Exit(1)
	}
	os.Exit(0)
}
