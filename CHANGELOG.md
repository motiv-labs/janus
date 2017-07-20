# Unrealeased

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

- Added Open Tracing support. Available tracers are Google Cloud Platform, Zipkin and appdash.

# 2.0.0

## Changed

- Split the application in two different ports, an administrative port (defaults to `8081`) and proxies port (defaults to `8080`). This way we avoid route collision with the admin routes and also we don't need to load tons of middlewares for the admin routes that are not necessary.
- Now the docker image is super tiny, less then 14mb when decompressed.
- API Definition and OAuth Server Definition don't depend on an ID anymore, now the name becomes the unique identifier. This works both in MongoDB and file based configurations.
- Handled 404 in a more elegant way

## Added

- Added possibility to create configurations using YAML, JSON, TOML or environemnt variables.
- Added a host matcher middleware.
