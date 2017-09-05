package main

import (
	"github.com/hellofresh/janus/pkg/config"
	"github.com/spf13/cobra"
)

// RunCheck is the run command to check Janus configurations
func RunCheck(cmd *cobra.Command, args []string) {
	_, err := config.Load(args[0])
	if nil != err {
		cmd.Printf("An error occurred: %s", err.Error())
		return
	}

	cmd.Printf("The configuration file is valid")
}
