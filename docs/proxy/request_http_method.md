#### Request HTTP method

Client requests can also be routed depending on their
HTTP method by specifying the `methods` field. By default, Janus will route a
request to an API regardless of its HTTP method. But when this field is set,
only requests with the specified HTTP methods will be matched.

This field also accepts multiple values. Here is an example of an API allowing
routing via `GET` and `HEAD` HTTP methods:

```json
{
    "name": "My API",
    "proxy": {
        "strip_path" : true,
        "listen_path": "/hello/*",
        "upstream_url": "http://my-api.com",
        "methods": ["GET", "HEAD"]
    }
}
```

Such an API would be matched with the following requests:

```http
GET / HTTP/1.1
Host:
```

```http
HEAD /resource HTTP/1.1
Host:
```

But would not match a `POST` or `DELETE` request. This allows for much more
granularity when configuring APIs and Middlewares. For example, one could imagine
two APIs pointing to the same upstream service: one API allowing unlimited
unauthenticated `GET` requests, and a second API allowing only authenticated
and rate-limited `POST` requests (by applying the authentication and rate
limiting plugins to such requests).
