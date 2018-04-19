# Retry

The retry plugin allows you to configure retry rules for your proxy. This enables you to be more resilient for any network or any other kind of failure.

## Configuration

The plain retry config:

```json
{
    "name" : "retry",
    "enabled" : false,
    "config" : {
        "attempts" : 3,
        "backoff": "1s"
    }
}
```

| Configuration | Description        |
| attempts      | Number of attempts |
| backoff       | Time that we should wait to retry. This must be given in the [ParseDuration](https://golang.org/pkg/time/#ParseDuration) format. Defaults to `1s` |
| predicate     | The rule that we will check to define if the request was successful or not. You have access to `statusCode` and all the `request` object. Defaults to `statusCode == 0 || statusCode >= 500` |
