package prometheus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hellofresh/janus/pkg/plugin"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	onStartupEventChan      chan *plugin.OnStartup
	onAdminStartupEventChan chan *plugin.OnAdminAPIStartup
)

func init() {
	onStartupEventChan = make(chan *plugin.OnStartup)
	onAdminStartupEventChan = make(chan *plugin.OnAdminAPIStartup)

	go onReady()

	plugin.RegisterEventHook(plugin.StartupEvent, onStartup)
	plugin.RegisterEventHook(plugin.AdminAPIStartupEvent, onAdminAPIStartup)

	plugin.RegisterPlugin("prometheus", plugin.Plugin{
		Action: setupPrometheus,
	})
}

func onReady() {
	var (
		startupEvent      *plugin.OnStartup
		startupAdminEvent *plugin.OnAdminAPIStartup
	)
	for i := 0; i < 2; i++ {
		select {
		case startupEvent = <-onStartupEventChan:
		case startupAdminEvent = <-onAdminStartupEventChan:
		case <-time.NewTimer(time.Second * 5).C:
			log.Errorf("prometheus plugin reach timeout of 5s waiting for startup & admin event")
			return
		}
	}

	if cfg := startupEvent.Config.Prometheus; cfg.Enabled {
		startupAdminEvent.Router.GET(cfg.MetricPath, prometheus.Handler().ServeHTTP)
	}
}

func onStartup(event interface{}) error {
	e, ok := event.(plugin.OnStartup)
	if !ok {
		return fmt.Errorf("Could not convert event to startup type")
	}
	onStartupEventChan <- &e
	return nil
}

func onAdminAPIStartup(event interface{}) error {
	e, ok := event.(plugin.OnAdminAPIStartup)
	if !ok {
		return fmt.Errorf("Could not convert event to admin startup type")
	}
	onAdminStartupEventChan <- &e
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
