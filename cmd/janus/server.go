package main

import (
	"context"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/opentracing"
	"github.com/hellofresh/janus/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	// this is needed to call the init function on each plugin
	_ "github.com/hellofresh/janus/pkg/plugin/basic"
	_ "github.com/hellofresh/janus/pkg/plugin/bodylmt"
	_ "github.com/hellofresh/janus/pkg/plugin/cb"
	_ "github.com/hellofresh/janus/pkg/plugin/compression"
	_ "github.com/hellofresh/janus/pkg/plugin/cors"
	_ "github.com/hellofresh/janus/pkg/plugin/oauth2"
	_ "github.com/hellofresh/janus/pkg/plugin/rate"
	_ "github.com/hellofresh/janus/pkg/plugin/requesttransformer"
	_ "github.com/hellofresh/janus/pkg/plugin/responsetransformer"
	_ "github.com/hellofresh/janus/pkg/plugin/retry"

	// dynamically registered auth providers
	_ "github.com/hellofresh/janus/pkg/jwt/basic"
	_ "github.com/hellofresh/janus/pkg/jwt/github"

	// internal plugins
	_ "github.com/hellofresh/janus/pkg/web"
)

// RunServer is the run command to start Janus
func RunServer(cmd *cobra.Command, args []string) {
	log.WithField("version", version).Info("Janus starting...")

	initConfig()
	initLog()
	initStatsd()

	tracingFactory := opentracing.New(globalConfig.Tracing)
	tracingFactory.Setup()

	defer tracingFactory.Close()
	defer statsClient.Close()
	defer globalConfig.Log.Flush()

	repo, err := api.BuildRepository(globalConfig.Database.DSN, globalConfig.Cluster.UpdateFrequency)
	if err != nil {
		log.WithError(err).Fatal("Could not build a repository for the database")
	}
	defer repo.Close()

	svr := server.New(
		server.WithGlobalConfig(globalConfig),
		server.WithMetricsClient(statsClient),
		server.WithProvider(repo),
	)

	ctx := ContextWithSignal(context.Background())
	svr.StartWithContext(ctx)
	defer svr.Close()

	svr.Wait()
	log.Info("Shutting down")
}
