package cmd

import (
	"context"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
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
)

// ServerStartOptions are the command flags
type ServerStartOptions struct {
	profilingEnabled bool
	profilingPublic  bool
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context) *cobra.Command {
	opts := &ServerStartOptions{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a Janus web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts)
		},
	}

	cmd.PersistentFlags().BoolVarP(&opts.profilingEnabled, "profiling-enabled", "", false, "Enable profiler, will be available on API port at /debug/pprof path")
	cmd.PersistentFlags().BoolVarP(&opts.profilingPublic, "profiling-public", "", false, "Allow accessing profiler endpoint w/out authentication")

	return cmd
}

// RunServerStart is the run command to start Janus
func RunServerStart(ctx context.Context, opts *ServerStartOptions) error {
	log.WithField("version", version).Info("Janus starting...")

	initConfig()
	initLog()
	initStatsClient()
	initStatsExporter()
	initTracingExporter()

	tracingFactory := opentracing.New(globalConfig.Tracing)
	tracingFactory.Setup()

	defer tracingFactory.Close()
	defer statsClient.Close()
	defer globalConfig.Log.Flush()

	repo, err := api.BuildRepository(globalConfig.Database.DSN, globalConfig.Cluster.UpdateFrequency)
	if err != nil {
		return errors.Wrap(err, "could not build a repository for the database")
	}
	defer repo.Close()

	svr := server.New(
		server.WithGlobalConfig(globalConfig),
		server.WithMetricsClient(statsClient),
		server.WithProvider(repo),
		server.WithProfiler(opts.profilingEnabled, opts.profilingPublic),
	)

	ctx = ContextWithSignal(ctx)
	svr.StartWithContext(ctx)
	defer svr.Close()

	svr.Wait()
	log.Info("Shutting down")

	return nil
}
