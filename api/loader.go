package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/etcinit/speedbump"
	"github.com/hellofresh/janus/cors"
	"github.com/hellofresh/janus/limitter"
	"github.com/hellofresh/janus/middleware"
	"github.com/hellofresh/janus/oauth"
	"github.com/hellofresh/janus/proxy"
	"github.com/hellofresh/janus/router"
	"gopkg.in/redis.v3"
)

type Loader struct {
	registerChan *proxy.RegisterChan
	redisClient  *redis.Client
	accessor     *middleware.DatabaseAccessor
	manager      *oauth.Manager
	debug        bool
}

// NewLoader creates a new instance of the api manager
func NewLoader(registerChan *proxy.RegisterChan, redisClient *redis.Client, accessor *middleware.DatabaseAccessor, manager *oauth.Manager, debug bool) *Loader {
	return &Loader{registerChan, redisClient, accessor, manager, debug}
}

// Load loads all api specs from a datasource
func (m *Loader) Load() {
	specs := m.getAPISpecs()
	m.LoadApps(specs)
}

// LoadApps load application middleware
func (m *Loader) LoadApps(apiSpecs []*Spec) {
	log.Debug("Loading API configurations")

	for _, referenceSpec := range apiSpecs {
		var skip bool

		//Validates the proxy
		skip = proxy.Validate(referenceSpec.Proxy)
		if false == referenceSpec.Active {
			log.Debug("API is not active, skiping...")
			skip = false
		}

		if skip {
			hasher := speedbump.PerSecondHasher{}
			limit := referenceSpec.RateLimit.Limit
			limiter := speedbump.NewLimiter(m.redisClient, hasher, limit)

			var handlers []router.Constructor
			if referenceSpec.RateLimit.Enabled {
				handlers = append(handlers, limitter.NewRateLimitMiddleware(limiter, hasher, limit).Handler)
			} else {
				log.Debug("Rate limit is not enabled")
			}

			if referenceSpec.CorsMeta.Enabled {
				handlers = append(handlers, cors.NewMiddleware(referenceSpec.CorsMeta, m.debug).Handler)
			} else {
				log.Debug("CORS is not enabled")
			}

			if referenceSpec.UseOauth2 {
				handlers = append(handlers, oauth.NewKeyExistsMiddleware(m.manager).Handler)
			} else {
				log.Debug("OAuth2 is not enabled")
			}

			m.registerChan.One <- proxy.NewRoute(referenceSpec.Proxy, handlers...)
			log.Debug("Proxy registered")
		} else {
			log.Error("Listen path is empty, skipping...")
		}
	}
}

//getAPISpecs Load application specs from datasource
func (m *Loader) getAPISpecs() []*Spec {
	log.Debug("Using App Configuration from Mongo DB")
	repo, err := NewMongoAppRepository(m.accessor.Session.DB(""))
	if err != nil {
		log.Panic(err)
	}

	definitions, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var APISpecs = []*Spec{}
	for _, definition := range definitions {
		newAppSpec := Spec{}
		newAppSpec.Definition = definition
		APISpecs = append(APISpecs, &newAppSpec)
	}

	return APISpecs
}
