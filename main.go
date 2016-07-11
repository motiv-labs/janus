package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/kataras/iris"
	"github.com/hellofresh/api-gateway/storage"
	"strings"
)

var APILoader = APIDefinitionLoader{}

// Specification for basic configurations
type Specification struct {
	storage.Database
	Port  int                 `envconfig:"PORT"`
	Debug bool                `envconfig:"DEBUG"`
}

func loadConfigEnv() Specification {
	var s Specification
	err := envconfig.Process("", &s)

	if err != nil {
		log.Fatal(err.Error())
	}

	return s
}

// initializeDatabase Initialize DB connection
func initializeDatabase(dbConfig storage.Database) *storage.DatabaseAccessor {
	accessor, err := storage.NewServer(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	//Use the middleware
	iris.Use(NewDatabase(*accessor))

	return accessor
}

func loadAPIEndpoints(proxyRegister *ProxyRegister) {
	group := iris.Party("/apis")
	{
		appHandler := AppsAPI{proxyRegister: proxyRegister}
		group.API("/", appHandler)
	}
}

func getAPISpecs(accessor *storage.DatabaseAccessor, dbConfig storage.Database) []*APIDefinition {
	var APISpecs []*APIDefinition

	log.Debug("Using App Configuration from Mongo DB")
	APISpecs = APILoader.LoadDefinitionsFromDatastore(accessor.Session, dbConfig)

	return APISpecs;
}

func loadApps(apiSpecs []*APIDefinition) {
	log.Info("Loading API configurations.")

	proxyRegister := NewProxyRegister()

	for _, referenceSpec := range apiSpecs {
		var skip bool

		if referenceSpec.Proxy.ListenPath == "" {
			log.Error("Listen path is empty, skipping API ID: ", referenceSpec.ID)
			skip = true
		}

		if strings.Contains(referenceSpec.Proxy.ListenPath, " ") {
			log.Error("Listen path contains spaces, is invalid, skipping API ID: ", referenceSpec.ID)
			skip = true
		}

		if !skip {
			proxyRegister.Register(referenceSpec.Proxy)
			cb := NewCircuitBreaker(referenceSpec)

			log.Debug("Proxy registered")
		}
	}
}

func main() {
	log.SetLevel(log.DebugLevel)

	s := loadConfigEnv()
	accessor := initializeDatabase(s.Database)
	defer accessor.Close()

	specs := getAPISpecs(accessor, s.Database)
	loadApps(specs)

	proxyRegister := NewProxyRegister()
	loadAPIEndpoints(proxyRegister)

	iris.Listen(fmt.Sprintf(":%v", s.Port))
}
