package main

import (
	"github.com/spf13/cobra"
)

// RunVersion is the run command to check Janus version
func RunVersion(cmd *cobra.Command, args []string) {
	cmd.Printf("Janus v%s", version)
}
