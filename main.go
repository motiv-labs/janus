package main

import (
	"os"

	"github.com/hellofresh/janus/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd := cmd.NewRootCmd()

	if err := cmd.Execute(); err != nil {
		log.WithError(err).Error(err.Error())
		os.Exit(1)
	}
}
