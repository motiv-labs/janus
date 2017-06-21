package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	port := flag.String("port", "8081", "port on localhost to check")
	flag.Parse()

	url := "http://127.0.0.1:" + *port + "/status"
	log.WithField("url", url).Info("Checking...")
	resp, err := http.Get(url)

	// If there is an error or non-200 status, exit with 1 signaling unsuccessful check.
	if err != nil || resp.StatusCode != 200 {
		log.WithError(err).Error("An error happened when checking Janus status")
		os.Exit(1)
	}

	log.Info("All good!")
	os.Exit(0)
}
