# Unreleased
- None


# 3.8.9
- Added `CONN_PURGE_INTERVAL` environment variable as a way to prevent stale http keep-alive connections

# 3.8.8

## Added
- Rate limiter configuration to respect `X-Forwarded-For` and `X-Real-IP` headers

## Changed
- Rate limiter plugin now ignores `X-Forwarded-For` and `X-Real-IP` headers by default

# 3.8.7

## Added

- Url parameters can be used in the target definition. Thanks to @Serjick
- Request-ID to error handler logs
- Additional attributes to tracing spans

## Changed

- Log writer initialized earlier
- Use in-memory repository for basic auth plugin when Mongo is not available
- Use `gofrs/uuid` instead of `satori/go.uuid`
- Rate limiter respects `X-Forwarded-For` and `X-Real-IP` HTTP headers

## Fixed

- Circuit breaker plugin statsd collector prefix
 
# 3.8.6

## Updated

- `http_server_request_latency` to include HTTP method key

# 3.8.5

## Fixed
- Fixed plugin configuration not being validated

## Updated
- Added stats and tracing support with opencensus

## Removed
- Tracing support via opentracing.io

# 3.8.4

## Fixed
- Fixed configuration listener that made API stuck

# 3.8.3

## Added
- Support b3 http propagation format for jaeger

## Fixed
- Race condition on application start. Reported on #348

## Updated
- Added more debug information to recovery handler to track application errors
- New `options_passthrough` parameter for CORS plugin. Thanks to @locker1776

# 3.8.1

## Fixed
- Open tracing error and http status code tags were not being set during tracing

# 3.8.0

## Added
- New Retry plugin: you can now configure your endpoints to have a retry in case of a failed request
- New `read`, `write` and `idle` timeouts for Janus server global configurations
- New `dial` and `response_header` timeouts that can be set per endpoint
- New `/debug/pprof` endpoint (handlers from `net/http/pprof`) on API port for debugging and profiling (can be enabled with `start` command flags)
- Alias `rr` for roundrobin balancer
- Add request id as a tag into tracing for seamlessly correlation in tracing UI

## Fixed
- Fixed bug when using the configuration file in a linux/64 system

## Updated
- Added `name` parameter for `cb` (Circuit Breaker) plugin to set group explicitly

## Removed
- Redis is not necessary anymore for the cluster to work
- Removed proxy definition property `enable_load_balancing` as it was not being used

# 3.7.0

## Added

- Leeway support for JWT date fields validation
- Support for zero weight when using the weight algorithm for balancing
- New header `X-Request-Id` that makes sure it create a new id for each request. It also ties it up with open tracing

## Fixed

- Fixed oauth rate limit reported on #276

## Removed

- `Upstream_URL` support is removed, see the [Upgrade Notes](docs/upgrade/3.7.x.md)

# 3.6.0

## Added

- Extra JWT metrics for token validation success and error

## Fixed

- Fixed a bug for the `oauth servers` when rows were empty it was returning `null` on the json response

## Updated

- Bumped [stats-go](https://github.com/hellofresh/stats-go) to current latest stable version (0.6.3) - this changes stats DSN config value format, see [`stats-go`](https://github.com/hellofresh/stats-go#instance-creation) docs for details

# 3.5.0

## Added

- Check GitHub permissions. Sets `is_admin` into the jwt token when the chosen provider is Github
- Jaeger support as distributed tracing backend
- Added Proxy Listen Path validation to prevent `chi` from panicking in case of invalid listen path
- Added load balancing for upstream targets. Now you can add multiple upstream targets and Janus will balance the requests.
- Added support for url parameters both in listen path and upstreams.

## Fixed

- Monitor health check endpoints only of active proxies. Reported on #203
- Fix hot reload was not working when using in memory storage implementation
- Fix oauth servers post endpoint incorrect behaviour. Reported on #234
- Add constant time compare to basic auth password. Reported on #194

## Removed

- Appdash support

## Updated

- THe docker image does not depend on a github release anymore

## Deprecated

- `upstream_url` is now deprecated in favor of using the `upstreams` object. This will allow Janus to balance requests if you have more than one upstream target.

# 3.3.0

## Added

- Added response transformer plugin
- Added basic auth plugin
- Added github login for the Admin API

## Updated

- Changed our dependency management tool from glide to Dep

## Fixed

- Fixed problems when using -c flag to specify a configuration file
- Fixed oAuth2 introspection token strategy when configuring an oauth server

# 3.2.1

## Added

- Added request body limit plugin
- Track application start/restart with stats metrics `<prefix>.app.init.<host>.<app-file>`

## Fixed

- Concurrent map writes in [stats-go](https://github.com/hellofresh/stats-go/pull/15)
- Non sampled spans recording in [gcloud-opentracing](https://github.com/hellofresh/gcloud-opentracing/pull/1)

# 3.2.0

## Added
- Added support for JWT signature validation chain for `jwt` token strategy
- Added support for OAuth2 `introspection` token strategy
- Added rate limit configurations for all endpoints of an OAuth2 server

## Removed
- Dropped support for `storage` token strategy

# 3.1.0

## Changed

- Moved Concourse CI scripts to another repo
- Changed health check JSON output to be in alignment with [health-go](https://github.com/hellofresh/health-go)
- Logging configuring is now handled by [logging-go](https://github.com/hellofresh/logging-go), so more logging options now
- Bumped Chi router to 3.0, see [changelog](https://github.com/go-chi/chi/blob/master/CHANGELOG.md) if you're using parametrised urls

## Added

- Added [plugin to transform](./docs/plugins/request_transformer.md) a request to an upstream. You can now modify headers and query string before the request is sent
- Added godog for behaviour tests
- Allow insecure upstream SSL certificate
- Added health-check statement on the Dockerfile. This will allow you to deploy the container to swarm/kubernetes/ecs and have it checked the `/status` endpoint.

# 3.0.0

## Changed

- Using viper to load the API definitions when using file based configurations. This allows you to configure your API definitions in YAML, JSON and TOML.
- The underling router was changed from [httptreemux](https://github.com/dimfeld/httptreemux) to [Chi](https://github.com/pressly/chi).
- Proper Mux reload when an API or OAuth server is changed

## Added

- Adds the ability to hot reload proxy definitions. To enable this feature you MUST use Redis as your datastore. If you use `in memory` storage this feature will not be enabled.
- Added the ability to enable or disable plugins per API definitions. This will bring us a lot of flexibility in developing new plugins and hooking them up. This feature is a BC and we should upgrade the major version because of that.
- Added health checks to any API definition

## Fixed

- Rate limit bug that was around for quite a while.
- Problems when creating a new API definition

# 2.2.0

## Changed

- Now the docker image is super tiny, less then 14mb when decompressed.
- Using commands to start Janus. This way we can improve the organization on how we want the binary to work. Also, this will allow us to probably move towards an ideal solution for hot reload of configs.

## Added

- Added coveralls as our coverage tool.
- Added plugins specifically for the round tripper. This allows us to decouple the token logic from the tripper.

# 2.1.0

## Changed

- The CI pipeline now bumps the patch version automatically.
- Updated docker compose to use the TOML config file
- Replaced the statsd implementation for our stats-go package

## Added

- Added Open Tracing support. Available tracers are Google Cloud Platform and Jaeger.

# 2.0.0

## Changed

- Split the application in two different ports, an administrative port (defaults to `8081`) and proxies port (defaults to `8080`). This way we avoid route collision with the admin routes and also we don't need to load tons of middlewares for the admin routes that are not necessary.
- Now the docker image is super tiny, less then 14mb when decompressed.
- API Definition and OAuth Server Definition don't depend on an ID anymore, now the name becomes the unique identifier. This works both in MongoDB and file based configurations.
- Handled 404 in a more elegant way

## Added

- Added possibility to create configurations using YAML, JSON, TOML or environemnt variables.
- Added a host matcher middleware.
