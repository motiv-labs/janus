package main

import (
	"os"

	"github.com/hellofresh/janus/cmd"
	log "github.com/sirupsen/logrus"
)

var version = "0.0.0-dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error(err.Error())
		os.Exit(1)
	}
}
