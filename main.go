package main

import (
	"github.com/hellofresh/janus/cmd"
	log "github.com/sirupsen/logrus"
)

var version = "0.0.0-dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Could not run command")
	}
}
