### Health Checks

Health checks can be added to each API definition by simply setting these properties:

*url*: The url to be checked
*timeout*: A timeout in seconds for the health check. If the timeout is reached an error is returned.

You will be able to see all your health checks on the admin REST endpoint `/status`. 
If everything is ok you will see something like this:

```json
{
    "status": "OK",
    "timestamp": "2017-06-21T13:06:50.546685883+02:00"
}
```

If you have any problems you'll see something like this:

```json
{
    "status": "Partially Available",
    "timestamp": "2017-06-21T14:44:38.782346389+02:00",
    "failures": {
        "example": "example is not available at the moment"
    }
}
```

Each one of the services must provide an endpoint that Janus can use to check how is the service performing.
The response code that the endpoint returns will define if the service is *healthy*, *partially healthy* or *unhealthy*

| Code           | Description               |
|----------------|---------------------------|
| 200 - 399      | Service fully working     |
| 400 - 499      | Service partially working |
| 500 >          | Service not working       |
