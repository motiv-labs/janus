
### Routing capabilities

Let's now discuss how Janus matches a request to the configured hosts, uris and methods properties (or fields) of your API. Note that all three of these fields are optional, but at least one of them must be specified. For a client request to match an API:

The request must include all of the configured fields
The values of the fields in the request must match at least one of the configured values (While the field configurations accepts one or more values, a request needs only one of the values to be considered a match)
Let's go through a few examples. Consider an API configured like this:

```json
{
    "name": "My API",
    "hosts": ["example.com", "service.com"],
    "proxy": {
        "listen_path": "/foo/*",
        "upstream_url": "http://my-api.com",
        "methods": ["GET"]
    }
}
```

Some of the possible requests matching this API could be:

```http
GET /foo HTTP/1.1
Host: example.com
```

```http
GET /foo HTTP/1.1
Host: service.com
```

```http
GET /foo/hello/world HTTP/1.1
Host: example.com
```

All three of these requests satisfy all the conditions set in the API definition.

However, the following requests would not match the configured conditions:

```http
GET / HTTP/1.1
Host: example.com
```

```http
POST /bar HTTP/1.1
Host: example.com
```

```http
GET /foo HTTP/1.1
Host: foo.com
```

All three of these requests satisfy only two of configured conditions. The first request's URI is not a match for any of the configured uris, same for the second request's HTTP method, and the third request's Host header.

Now that we understand how the hosts, uris, and methods properties work together, let's explore each property individually.
