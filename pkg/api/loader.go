package api

import (
	"github.com/NYTimes/gziphandler"
	log "github.com/Sirupsen/logrus"
	"github.com/hellofresh/janus/pkg/cors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/oauth"
	"github.com/hellofresh/janus/pkg/proxy"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/ulule/limiter"
)

// Loader is responsible for loading all apis form a datastore and configure them in a register
type Loader struct {
	register *proxy.Register
	storage  store.Store
	authRepo oauth.Repository
}

// NewLoader creates a new instance of the api manager
func NewLoader(register *proxy.Register, storage store.Store, authRepo oauth.Repository) *Loader {
	return &Loader{register, storage, authRepo}
}

// LoadDefinitions will connect and download ApiDefintions from a Mongo DB instance.
func (m *Loader) LoadDefinitions(repo APISpecRepository) {
	specs := m.getAPISpecs(repo)
	m.RegisterApis(specs)
}

// RegisterApis load application middleware
func (m *Loader) RegisterApis(apiSpecs []*Spec) {
	log.Debug("Loading API configurations")

	for _, referenceSpec := range apiSpecs {
		m.RegisterApi(referenceSpec)
	}
}

func (m *Loader) RegisterApi(referenceSpec *Spec) {
	//Validates the proxy
	active := proxy.Validate(referenceSpec.Proxy)
	if false == referenceSpec.Active {
		log.WithField("listen_path", referenceSpec.Proxy.ListenPath).Debug("API is not active, skiping...")
		active = false
	}

	if active {
		var handlers []router.Constructor
		if referenceSpec.RateLimit.Enabled {
			rate, err := limiter.NewRateFromFormatted(referenceSpec.RateLimit.Limit)
			if err != nil {
				panic(err)
			}

			limiterStore, err := m.storage.ToLimiterStore(referenceSpec.Slug)
			if err != nil {
				panic(err)
			}

			handlers = append(handlers, limiter.NewHTTPMiddleware(limiter.NewLimiter(limiterStore, rate)).Handler)
			handlers = append(handlers, middleware.NewRateLimitLogger().Handler)
		} else {
			log.Debug("Rate limit is not enabled")
		}

		if referenceSpec.CorsMeta.Enabled {
			handlers = append(handlers, cors.NewMiddleware(referenceSpec.CorsMeta, false).Handler)
		} else {
			log.Debug("CORS is not enabled")
		}

		if referenceSpec.UseOauth2 {
			handlers = append(handlers, NewKeyExistsMiddleware(referenceSpec).Handler)
		} else {
			log.Debug("OAuth2 is not enabled")
		}

		if referenceSpec.UseCompression {
			handlers = append(handlers, gziphandler.GzipHandler)
		} else {
			log.Debug("Compression is not enabled")
		}

		m.register.Add(proxy.NewRoute(referenceSpec.Proxy, handlers...))
		log.WithField("listen_path", referenceSpec.Proxy.ListenPath).Debug("API registered")
	} else {
		log.WithField("listen_path", referenceSpec.Proxy.ListenPath).Error("Listen path is invalid or not active, skipping...")
	}
}

//getAPISpecs Load application specs from datasource
func (m *Loader) getAPISpecs(repo APISpecRepository) []*Spec {
	definitions, err := repo.FindAll()
	if err != nil {
		log.Panic(err)
	}

	var specs []*Spec
	for _, definition := range definitions {
		spec, err := m.makeSpec(definition)
		if nil != err {
			continue
		}

		specs = append(specs, spec)
	}

	return specs
}

func (m *Loader) makeSpec(definition *Definition) (*Spec, error) {
	spec := new(Spec)
	spec.Definition = definition
	if definition.UseOauth2 {
		manager, err := m.getManager(definition.OAuthServerSlug)
		if nil != err {
			log.WithError(err).Error("OAuth Configuration for this API is incorrect, skipping...")
			return nil, err
		}
		spec.Manager = manager
	}

	return spec, nil
}

func (m *Loader) getManager(oAuthServerSlug string) (oauth.Manager, error) {
	oauthServer, err := m.authRepo.FindBySlug(oAuthServerSlug)
	if nil != err {
		return nil, err
	}

	managerType, err := oauth.ParseType(oauthServer.TokenStrategy.Name)
	if nil != err {
		return nil, err
	}

	return oauth.NewManagerFactory(m.storage, oauthServer.TokenStrategy.Settings).Build(managerType)
}
