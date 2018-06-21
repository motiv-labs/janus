# Circuit Breaker

Janus has a circuit breaker plugin that can be configured for each endpoint. You can check 
our [example](https://github.com/hellofresh/janus/tree/master/examples/plugin-cb) on how
to use the plugin.

## Configuration

The plain cb config:

```json
{
    "name" : "cb",
    "enabled" : true,
    "config" : {
        "name": "my-circuit-breaker",
        "timeout" : 1000,
        "max_concurrent_requests": 100,
        "error_percent_threshold": 50,
        "request_volume_threshold": 20,
        "sleep_window": 5000,
        "predicate": "statusCode == 0 || statusCode >= 500"
    }
}
```

Configuration | Description
:---|:---|
| name                        | Circuit Breaker name to group stats |
| timeout                     | Timeout that the CB will wait till the request responds |
| max_concurrent_requests     | How many commands of the same type can run at the same time |
| error_percent_threshold     | Causes circuits to open once the rolling measure of errors exceeds this percent of requests |
| request_volume_threshold    | Is the minimum number of requests needed before a circuit can be tripped due to health |
| sleep_window                | Is how long, in milliseconds, to wait after a circuit opens before testing for recovery |
| predicate                   | The rule that we will check to define if the request was successful or not. You have access to `statusCode` and all the `request` object. Defaults to `statusCode == 0 || statusCode >= 500` |
