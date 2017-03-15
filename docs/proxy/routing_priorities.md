### Routing priorities

An API may define matching rules based on its `hosts`, `listen_path`, and `methods`
fields. For Janus to match an incoming request to an API, all existing fields
must be satisfied. However, Janus allows for quite some flexibility by allowing
two or more APIs to be configured with fields containing the same values - when
this occurs, Janus applies a priority rule.

The rule is that : **when evaluating a request, Janus will first try
to match the APIs with the most rules**.

For example, two APIs are configured like this:

```json
{
    "name": "API 1",
    "proxy": {
        "listen_path": "/",
        "upstream_url": "http://my-api.com",
        "hosts": ["example.com"]
    }
},
{
    "name": "API 2",
    "proxy": {
        "listen_path": "/",
        "upstream_url": "http://my-api.com",
        "hosts": ["example.com"],
        "methods": ["POST"]
    }
}
```

api-2 has a `hosts` field **and** a `methods` field, so it will be
evaluated first by Janus. By doing so, we avoid api-1 "shadowing" calls
intended for api-2.

Thus, this request will match api-1:

```http
GET / HTTP/1.1
Host: example.com
```

And this request will match api-2:

```http
POST / HTTP/1.1
Host: example.com
```

Following this logic, if a third API was to be configured with a `hosts` field,
a `methods` field, and a `listen_path` field, it would be evaluated first by Janus.
