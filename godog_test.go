package janus

import (
	"flag"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/features/bootstrap"
	"github.com/hellofresh/janus/pkg/config"
)

var runGoDogTests bool

func init() {
	flag.BoolVar(&runGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	flag.Bool("random", false, "Randomize features/scenarios execution. Flag is passed to godog")
	flag.Bool("stop-on-failure", false, "Stop processing on first failed scenario.. Flag is passed to godog")
	flag.Parse()
}

func FeatureContext(s *godog.Suite) {
	c, err := config.Load("")
	if nil != err {
		log.WithError(err).Panic("Error initializing statsd client")
	}

	bootstrap.RegisterRequestContext(s, c.Port, c.Web.Port)
}

func TestMain(m *testing.M) {
	if !runGoDogTests {
		os.Exit(0)
	}

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: "progress",
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
