### Circuit breaker

Circuit breaker are enabled by default for each endpoint and the settings can be adjusted globally, but also per API definition.

#### Configuration
Set the following environment variables to overwrite the settings globally:

*CB_TIMEOUT* CB_TIMEOUT is how long to wait for command to complete, in milliseconds
*CB_MAX_CONCURRENT* CB_MAX_CONCURRENT is how many commands of the same type can run at the same time
*CB_VOLUME_THRESHOLD* CB_VOLUME_THRESHOLD is the minimum number of requests needed before a circuit can be tripped due to health
*CB_SLEEP_WINDOW* CB_SLEEP_WINDOW is how long, in milliseconds, to wait after a circuit opens before testing for recovery
*CB_ERROR_PRECENT_THRESHOLD* CB_ERROR_PRECENT_THRESHOLD DefaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
*CB_DASHBOARD_ENABLED* CB_DASHBOARD_ENABLED enables a streaming endpoint which can be consumed by the hystrix-dashboard
*CB_DASHBOARD_PORT* CB_DASHBOARD_PORT port for the streaming endpoint

To overwrite the settings per endpoint add the `circuit_breaker` config to your endpoint definition:
```json
{
  "name":"example",
  "active":true,
  "proxy":{
    ...
  },
  "circuit_breaker":{
    "timeout": 1000,
    "max_concurrent_requests": 10,
    "request_volume_threshold": 20,
    "sleep_window": 5000,
    "error_percent_threshold": 50
  }
}
```

#### Metrics
Janus uses the hystrix-go implementation of the circuit breaker pattern, which automatically reports failures and success to statsd.

#### References
[circuit-breaker](https://martinfowler.com/bliki/CircuitBreaker.html)
[hystrix](https://github.com/Netflix/Hystrix/wiki#what)
[hystrix-go](https://github.com/afex/hystrix-go)
[hystrix-dashboard (deprecated)](https://github.com/Netflix-Skunkworks/hystrix-dashboard)
