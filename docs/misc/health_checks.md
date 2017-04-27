### Health Checks

Health checks can be added to each API definition by simply setting these properties:

*URL*: The url to be checked

You will be able to check all your health checks on the admin REST endpoint `/status`. The response
will be something like this:

```
[
  {
    "name": "configurations",
    "check": {
      "status": "healthy",
      "subject": "fully operating"
    }
  },
  {
    "name": "example",
    "check": {
      "status": "partially_healthy",
      "subject": "rabbitMQ connection is not respoding"
    }
  },
  {
    "name": "foo",
    "check": {
      "status": "unhealthy",
      "subject": "postgres database down"
    }
  }
]
```

Each one of the services must provide an endpoint that Janus can use to check how is the service performing.
The response code that the endpoint returns will define if the service is *healthy*, *partially healthy* or *unhealthy*

| Code           | Description               |
|----------------|---------------------------|
| 200 - 399      | Service fully working     |
| 400 - 499      | Service partially working |
| 500 >          | Service not working       |
