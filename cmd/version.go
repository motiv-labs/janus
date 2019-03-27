package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var version = "0.0.0-dev"

// NewVersionCmd creates a new version command
func NewVersionCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "Print the version information",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("janus %s\n", version)
		},
	}
}
