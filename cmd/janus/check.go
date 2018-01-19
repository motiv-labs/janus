package main

import (
	"github.com/hellofresh/janus/pkg/config"
	"github.com/spf13/cobra"
)

// NewCheckCmd creates a new check command
func NewCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check [config-file]",
		Short: "Check the validity of a given Janus configuration file. (default /etc/janus/janus.toml)",
		Args:  cobra.MinimumNArgs(1),
		Run:   RunCheck,
	}
}

// RunCheck is the run command to check Janus configurations
func RunCheck(cmd *cobra.Command, args []string) {
	_, err := config.Load(args[0])
	if nil != err {
		cmd.Printf("An error occurred: %s", err.Error())
		return
	}

	cmd.Printf("The configuration file is valid")
}
