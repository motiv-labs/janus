#### The `preserve_host` property

When proxying, Janus's default behavior is to set the upstream request's Host header to the hostname of the API's `upstream_url` property. The `preserve_host` field accepts a boolean flag instructing Janus not to do so.

For example, when the `preserve_host` property is not changed and an API is configured like this:

```json
{
    "name": "My API",
    "hosts": ["service.com"],
    "proxy": {
        "listen_path": "/foo/*",
        "upstream_url": "http://my-api.com",
        "methods": ["GET"]
    }
}
```

A possible request from a client to Janus could be:

```http
GET / HTTP/1.1
Host: service.com
```

Janus would extract the Host header value from the the hostname of the API's `upstream_url` field, and would send the following request to your upstream service:

```http
GET / HTTP/1.1
Host: my-api.com
```

However, by explicitly configuring your API with `preserve_host=true`:

```json
{
    "name": "My API",
    "hosts": ["example.com", "service.com"],
    "proxy": {
        "listen_path": "/foo/*",
        "upstream_url": "http://my-api.com",
        "methods": ["GET"],
        "preserve_host": true
    }
}
```

And assuming the same request from the client:

```http
GET / HTTP/1.1
Host: service.com
```

Janus would preserve the Host on the client request and would send the following request to your upstream service:

```http
GET / HTTP/1.1
Host: service.com
```
