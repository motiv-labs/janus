package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var configFile string

func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "janus",
		Short: "Janus is an API Gateway",
		Long: `
This is a lightweight API Gateway and Management Platform that enables you
to control who accesses your API, when they access it and how they access it.
API Gateway will also record detailed analytics on how your users are interacting
with your API and when things go wrong.
Complete documentation is available at https://hellofresh.gitbooks.io/janus`,
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is $PWD/janus.toml)")

	cmd.AddCommand(NewCheckCmd(ctx))
	cmd.AddCommand(NewVersionCmd(ctx))
	cmd.AddCommand(NewServerStartCmd(ctx))

	return cmd
}
