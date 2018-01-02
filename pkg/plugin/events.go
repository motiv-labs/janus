package plugin

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/notifier"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	stats "github.com/hellofresh/stats-go"
	mgo "gopkg.in/mgo.v2"
)

// Define the event names for the startup and shutdown events
const (
	StartupEvent         string = "startup"
	AdminAPIStartupEvent string = "admin_startup"

	ReloadEvent   string = "reload"
	ShutdownEvent string = "shutdown"
)

// OnStartup represents a event that happens when Janus starts up on the main process
type OnStartup struct {
	Notifier     notifier.Notifier
	StatsClient  stats.Client
	Repository   api.Repository
	MongoSession *mgo.Session
	Register     *proxy.Register
	Config       *config.Specification
}

// OnReload represents a event that happens when Janus hot reloads it's configurations
type OnReload struct {
	Register   *proxy.Register
	Repository api.Repository
}

// OnAdminAPIStartup represents a event that happens when Janus starts up the admin API
type OnAdminAPIStartup struct {
	Router router.Router
}
