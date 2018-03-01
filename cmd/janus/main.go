package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile  string
	versionFlag bool

	// Root command
	rootCmd = &cobra.Command{
		Use:   "janus",
		Short: "Janus is an API Gateway",
		Long: `
This is a lightweight API Gateway and Management Platform that enables you
to control who accesses your API, when they access it and how they access it.
API Gateway will also record detailed analytics on how your users are interacting
with your API and when things go wrong.
Complete documentation is available at https://hellofresh.gitbooks.io/janus`,
		Run: RunServer,
	}
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file (default is $PWD/janus.toml)")

	rootCmd.AddCommand(NewCheckCmd())
	rootCmd.AddCommand(NewVersionCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("Something went wrong")
	}
}
