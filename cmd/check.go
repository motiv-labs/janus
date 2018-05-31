package cmd

import (
	"context"

	"github.com/hellofresh/janus/pkg/config"
	"github.com/spf13/cobra"
)

// NewCheckCmd creates a new check command
func NewCheckCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "check [config-file]",
		Short: "Check the validity of a given Janus configuration file. (default /etc/janus/janus.toml)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  RunCheck,
	}
}

// RunCheck is the run command to check Janus configurations
func RunCheck(cmd *cobra.Command, args []string) error {
	_, err := config.Load(args[0])
	if err != nil {
		return err
	}

	cmd.Printf("The configuration file is valid\n")
	return nil
}
