package appdash

import (
	"net"
	"net/http"
	"net/url"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"sourcegraph.com/sourcegraph/appdash"
	appdashtracer "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

// Server represents a new appdash server (local or remote)
type Server struct {
	collectorDSN    string
	appdashHTTPAddr string
}

// NewServer creates a new instance of Server
func NewServer(collectorDSN string, appdashHTTPAddr string) *Server {
	return &Server{collectorDSN: collectorDSN, appdashHTTPAddr: appdashHTTPAddr}
}

// Listen starts the appdash server collector and UI
func (s *Server) Listen() error {
	memStore := appdash.NewMemoryStore()
	store := &appdash.RecentStore{
		MinEvictAge: 20 * time.Second,
		DeleteStore: memStore,
	}

	s.listenCollector(s.collectorDSN, store)
	s.listenWebUI(s.appdashHTTPAddr, store, memStore)

	return nil
}

// GetTracer returns an open tracing compatible tracer
func (s *Server) GetTracer() opentracing.Tracer {
	collector := appdash.NewRemoteCollector(s.collectorDSN)

	// Here we use the local collector to create a new opentracing.Tracer
	return appdashtracer.NewTracer(collector)
}

func (s *Server) listenCollector(collectorDSN string, store appdash.Store) error {
	l, err := net.Listen("tcp", collectorDSN)
	if err != nil {
		return err
	}

	collectorPort := l.Addr().(*net.TCPAddr).Port
	log.Info("Appdash collector listening on tcp:%d", collectorPort)
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	return nil
}

func (s *Server) listenWebUI(appdashHTTPAddr string, store appdash.Store, queryer appdash.Queryer) error {
	appdashURLStr := "http://localhost" + appdashHTTPAddr
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		return err
	}

	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		return err
	}

	tapp.Store = store
	tapp.Queryer = queryer

	log.Infof("Appdash web UI running at %s", appdashHTTPAddr)

	go func() {
		log.Fatal(http.ListenAndServe(appdashHTTPAddr, tapp))
	}()

	return nil
}
