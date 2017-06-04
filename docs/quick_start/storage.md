# Storage

Some Janus functionality requires key/value storage. These functionalities are:

* [Rate Limit](../plugins/rate_limit.md) Plugin
* [OAuth 2.0](../auth/oauth.md) Token validation strategy (when set to `storage`)
* Proxy changes propagation across all running Janus instances (when storage type is `redis`, see bellow)

## Configuring storage

Storage can be configured with either `STORAGE_DSN` environment variable or config file value, depending on
config file format. E.g. for `toml` config file format storage configuration looks like:

```toml
[storage]
  dsn = "<storage dsn>"
```

## Storage types

Janus supports the following storage types:

* `in memory` - DSN format is `memory://localhost`, scheme is the only part that matters here.
* `redis` - DSN format is `redis://<host>:<port>[/?prefix=<prefix>]`, prefix is optional
  and is set to `janus` by default, unless you set other value explicitly.
