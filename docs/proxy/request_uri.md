#### Request URI

Another way for Janus to route a request to a given upstream service is to
specify a request URI via the `proxy.listen_path` property. To satisfy this field's
condition, a client request's URI **must** be prefixed with one of the values
of the `proxy.listen_path` field.

For example, in an API configured like this:

```json
{
    "name": "My API",
    "proxy": {
        "listen_path": "/hello/*",
        "upstream_url": "http://my-api.com",
        "methods": ["GET"],
    }
}
```

The following requests would match the configured API:

```http
GET /hello HTTP/1.1
Host: my-api.com
```

```http
GET /hello/resource?param=value HTTP/1.1
Host: my-api.com
```

```http
GET /hello/world/resource HTTP/1.1
Host: anything.com
```

For each of these requests, Janus detects that their URI is prefixed with one of
the API's `proxy.listen_path` values. By default, Janus would then forward the request
upstream with the untouched, **same URI**.

When proxying with URIs prefixes, **the longest URIs get evaluated first**.
This allow you to define two APIs with two URIs: `/service` and
`/service/resource`, and ensure that the former does not "shadow" the latter.
