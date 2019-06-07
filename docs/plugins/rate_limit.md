# Rate Limiting

Rate limit how many HTTP requests a developer can make in a given period of seconds, minutes, hours, days, months or years.

## Configuration

The plain rate limit config:

```json
"rate_limit": {
    "enabled": true,
    "config": {
        "limit": "10-S",
        "policy": "local",
        "trust_forward_headers": false
    }
}
```

| Configuration        | Description                                                                                                                                                                                                                                                 |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| limit                 | Defines the limit rule for the proxy. i.e. 5 reqs/second: `5-S`, 10 reqs/minute: `10-M`, 1000 reqs/hour: `1000-H`                                                                                                                                           |
| policy                | The rate-limiting policies to use for retrieving and incrementing the limits. Available values are `local` (counters will be stored locally in-memory on the node) and `redis` (counters are stored on a Redis server and will be shared across the nodes). |
| redis.dsn             | The DSN for the redis instance/cluster to be used                                                                                                                                                                                                           |
| redis.prefix          | A prefix to be used on redis keys. It defaults to `limiter`                                                                                                                                                                                                 |                                                        |
| trust_forward_headers | If set to True, `X-Forwarded-For` and `X-Real-IP` headers will be used instead of the source ip. Defaults to False.                                                                                                                                         |

## Headers sent to the client

When this plugin is enabled, Janus will send some additional headers back to the client telling how many requests are available and what are the limits allowed, for example:

```
X-Ratelimit-Limit: 10
X-Ratelimit-Remaining: 9
X-Ratelimit-Reset: 1491383478
```

If any of the limits configured is being reached, the plugin will return a HTTP/1.1 `429` status code to the client with the following plain text body:

```
Limit exceeded
```

# Implementation considerations

The plugin supports 3 policies, which each have their specific pros and cons.

| Policy | Pros                                                      | Cons                                                                                                                                |
|--------|-----------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| redis  | accurate, lesser performance impact than a cluster policy | extra redis installation required, bigger performance impact than a local policy                                                    |
| local  | minimal performance impact                                | less accurate, and unless a consistent-hashing load balancer is used in front of Janus, it diverges when scaling the number of nodes |

There are 2 use cases that are most common:

### 1. every transaction counts. 
These are for example transactions with financial consequences. Here the highest level of accuracy is required.

### 2. backend protection. 
This is where accuracy is not as relevant, but it is merely used to protect backend services from overload. Either by specific users, or to protect against an attack in general.

> NOTE: the redis policy does not support the Sentinel protocol for high available master-slave architectures. When using rate-limiting for general protection the chances of both redis being down and the system being under attack are rather small. Check with your own use case wether you can handle this (small) risk.
