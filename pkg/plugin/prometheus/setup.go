package prometheus

import (
	"fmt"
	"net/http"

	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	adminRouter router.Router
)

func init() {
	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)

	plugin.RegisterPlugin("prometheus", plugin.Plugin{
		Action: setupPrometheus,
	})
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return fmt.Errorf("Could not convert event to startup type")
	}
	if adminRouter == nil {
		return errors.New(http.StatusNotFound, "invalid admin router given")
	}
	if e.Config.Prometheus.Enabled {
		adminRouter.GET(e.Config.Prometheus.MetricPath, prometheus.Handler().ServeHTTP)
	}
	return nil
}

func onAdminAPIStartup(event interface{}) error {
	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return fmt.Errorf("Could not convert event to admin startup type")
	}
	adminRouter = e.Router
	return nil
}

func setupPrometheus(route *proxy.Route, rawConfig plugin.Config) error {
	route.AddInbound(func(handle http.Handler) http.Handler {
		return prometheus.InstrumentHandlerWithOpts(prometheus.SummaryOpts{
			Subsystem: "janus",
			ConstLabels: prometheus.Labels{
				"listen_path":  route.Proxy.ListenPath,
				"upstream_url": route.Proxy.UpstreamURL,
			},
		}, handle)
	})
	return nil
}
