package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
)

func main() {
	var RootCmd = &cobra.Command{
		Use:   "janus",
		Short: "Janus is an API Gateway",
		Long: `This is a lightweight API Gateway and Management Platform that enables you 
				to control who accesses your API, when they access it and how they access it. 
				API Gateway will also record detailed analytics on how your users are interacting 
				with your API and when things go wrong.
                Complete documentation is available at https://hellofresh.gitbooks.io/janus`,
		Run: RunServer,
	}
	RootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Source of a configuration file")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
