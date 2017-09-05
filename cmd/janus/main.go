package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version     string
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

	checkCmd = &cobra.Command{
		Use:   "check [config-file]",
		Short: "Check the validity of a given Janus configuration file. (default /etc/janus/janus.toml)",
		Run:   RunCheck,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print Janus's version",
		Run:   RunVersion,
	}
)

func init() {
	cobra.OnInitialize(func() {
		if versionFlag {
			fmt.Println("Janus v" + version)
			os.Exit(0)
		}
	})

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Source of a configuration file")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print application version")
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
