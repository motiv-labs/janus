package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "0.0.0-dev"

// NewVersionCmd creates a new version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the version information",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			RunVersion()
		},
	}
}

// RunVersion runs the command to print the current version
func RunVersion() {
	log.Infof("janus %s", version)
}
