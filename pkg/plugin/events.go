package plugin

import (
	"github.com/hellofresh/janus/cassandra/wrapper"
	"github.com/hellofresh/stats-go/client"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
)

// Define the event names for the startup and shutdown events
const (
	StartupEvent         = "startup"
	AdminAPIStartupEvent = "admin_startup"

	ReloadEvent   = "reload"
	ShutdownEvent = "shutdown"
	SetupEvent    = "setup"
)

// OnStartup represents a event that happens when Janus starts up on the main process
type OnStartup struct {
	StatsClient   client.Client
	MongoDB       *mongo.Database
	Cassandra 	  wrapper.Holder
	Register      *proxy.Register
	Config        *config.Specification
	Configuration []*api.Definition
}

// OnReload represents a event that happens when Janus hot reloads it's configurations
type OnReload struct {
	Configurations []*api.Definition
}

// OnAdminAPIStartup represents a event that happens when Janus starts up the admin API
type OnAdminAPIStartup struct {
	Router router.Router
}
