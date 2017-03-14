# Proxy Reference

Janus listens for traffic on four ports, which by default are:

`:8080` on which Janus listens for incoming HTTP traffic from your clients, and forwards it to your upstream services.
`:8443` on which Janus listens for incoming HTTPS traffic. This port has a similar behavior as the `:8080` port, except that it expects HTTPS traffic only. This port can be disabled via the configuration file.
`:8081` on which the [Admin API](admin_api.md) used to configure Janus listens.
`:8444` on which the [Admin API](admin_api.md) listens for HTTPS traffic.

### Table of Contents

- [Terminology][proxy-terminology]
- [Overview][proxy-overview]
- [Reminder: How to add an API to Janus][proxy-reminder]
- [Routing capabilities][proxy-routing-capabilities]
    - [Request Host header][proxy-request-host-header]
        - [Using wildcard hostnames][proxy-using-wildcard-hostnames]
        - [The `preserve_host_header` property][proxy-preserve-host-property]
    - [Request URI][proxy-request-uri]
        - [The `strip_listen_path` property][proxy-strip-uri-property]
    - [Request HTTP method][proxy-request-http-method]
- [Routing priorities][proxy-routing-priorities]
- [Proxying behavior][proxy-proxying-behavior]
    - [1. Load balancing][proxy-load-balancing]
    - [2. Plugins execution][proxy-plugins-execution]
    - [3. Proxying & upstream timeouts][proxy-proxying-upstream-timeouts]
    - [4. Response][proxy-response]
- [Configuring a fallback API][proxy-configuring-a-fallback-api]
- [Configuring SSL for an API][proxy-configuring-ssl-for-an-api]
    - [The `https_only` property][proxy-the-https-only-property]
    - [The `http_if_terminated` property][proxy-the-http-if-terminated-property]
- [Proxy WebSocket traffic][proxy-websocket]
- [Conclusion][proxy-conclusion]

[proxy-terminology]: #terminology
[proxy-overview]: #overview
[proxy-reminder]: #reminder-how-to-add-an-api-to-Janus
[proxy-routing-capabilities]: #routing-capabilities
[proxy-request-host-header]: #request-host-header
[proxy-using-wildcard-hostnames]: #using-wildcard-hostnames
[proxy-preserve-host-property]: #the-preserve_host_header-property
[proxy-request-uri]: #request-uri
[proxy-strip-uri-property]: #the-strip_listen_path-property
[proxy-request-http-method]: #request-http-method
[proxy-routing-priorities]: #routing-priorities
[proxy-proxying-behavior]: #proxying-behavior
[proxy-load-balancing]: #1-load-balancing
[proxy-plugins-execution]: #2-plugins-execution
[proxy-proxying-upstream-timeouts]: #3-proxying-amp-upstream-timeouts
[proxy-response]: #proxy-response
[proxy-configuring-a-fallback-api]: #configuring-a-fallback-api
[proxy-configuring-ssl-for-an-api]: #configuring-ssl-for-an-api
[proxy-the-https-only-property]: #the-https_only-property
[proxy-the-http-if-terminated-property]: #the-http_if_terminated-property
[proxy-websocket]: #proxy-websocket-traffic
[proxy-conclusion]: #conclusion

### Terminology

`API`: This term refers to the API entity of Janus. You configure your APIs, that point to your own upstream services, through the Admin API.
`Middleware`: This refers to Janus "middleware", which are pieces of business logic that run in the proxying lifecycle. Middleware can be configured through the Admin API - either globally (all incoming traffic) or on a per-API basis.
`Client`: Refers to the downstream client making requests to Janus's proxy port.
`Upstream service`: Refers to your own API/service sitting behind Janus, to which client requests are forwarded.

### Overview

From a high level perspective, Janus will listen for HTTP traffic on its configured proxy port (`8080` by default), recognize which upstream service is being requested, run the configured middlewares for that API, and forward the HTTP request upstream to your own API or service.

When a client makes a request to the proxy port, Janus will decide to which upstream service or API to route (or forward) the incoming request, depending on the API configuration in Janus, which is managed via the Admin API. You can configure APIs with various properties, but the three relevant ones for routing incoming traffic are hosts, uris, and methods.

If Janus cannot determine to which upstream API a given request should be routed, Janus will respond with:

```http
HTTP/1.1 404 Not Found
Content-Type: application/json
Server: Janus/<x.x.x>

{
    "error": "no API found with those values"
}
```

### Routing capabilities

Let's now discuss how Janus matches a request to the configured hosts, uris and methods properties (or fields) of your API. Note that all three of these fields are optional, but at least one of them must be specified. For a client request to match an API:

The request must include all of the configured fields
The values of the fields in the request must match at least one of the configured values (While the field configurations accepts one or more values, a request needs only one of the values to be considered a match)
Let's go through a few examples. Consider an API configured like this:

```json
{
    "name": "My API",
    "slug": "my-api",
    "hosts": ["example.com", "service.com"],
    "proxy": {
        "listen_path": "/foo/*",
        "target_url": "http://my-api.com",
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


### Request Host header

Routing a request based on its Host header is the most straightforward way to proxy traffic through Janus, as this is the intended usage of the HTTP Host header. Janus makes it easy to do so via the hosts field of the API entity.

`hosts` accepts multiple values, which must be in an array format when specifying them via the Admin API:

```json
{
    "hosts": ["my-api.com", "example.com", "service.com"]
}
```

To satisfy the hosts condition of this API, any incoming request from a client must now have its Host header set to one of:

```http
Host: my-api.com
```

or:

```http
Host: example.com
```

or:

```http
Host: service.com
```

#### Using wildcard hostnames

To provide flexibility, Janus allows you to specify hostnames with wildcards in the hosts field. Wildcard hostnames allow any matching Host header to satisfy the condition, and thus match a given API.

Wildcard hostnames must contain only one asterisk at the leftmost or rightmost label of the domain. Examples:

`*.example.org` would allow Host values such as `a.example.com` and `x.y.example.com` to match.
`example.*` would allow Host values such as `example.com` and `example.org` to match.
A complete example would look like this:

```json
{
    "slug": "my-api",
    "hosts": ["*.example.com", "service.com"]
}
```

Which would allow the following requests to match this API:

```http
GET / HTTP/1.1
Host: an.example.com
```

```http
GET / HTTP/1.1
Host: service.com
```
[Back to TOC](#table-of-contents)

#### The `preserve_host_header` property

When proxying, Janus's default behavior is to set the upstream request's Host header to the hostname of the API's `target_url` property. The `preserve_host_header` field accepts a boolean flag instructing Janus not to do so.

For example, when the `preserve_host_header` property is not changed and an API is configured like this:

```json
{
    "name": "My API",
    "slug": "my-api",
    "hosts": ["service.com"],
    "proxy": {
        "listen_path": "/foo/*",
        "target_url": "http://my-api.com",
        "methods": ["GET"]
    }
}
```

A possible request from a client to Janus could be:

```http
GET / HTTP/1.1
Host: service.com
```

Janus would extract the Host header value from the the hostname of the API's `target_url` field, and would send the following request to your upstream service:

```http
GET / HTTP/1.1
Host: my-api.com
```

However, by explicitly configuring your API with `preserve_host_header=true`:

```json
{
    "name": "My API",
    "slug": "my-api",
    "hosts": ["example.com", "service.com"],
    "proxy": {
        "listen_path": "/foo/*",
        "target_url": "http://my-api.com",
        "methods": ["GET"],
        "preserve_host_header": true
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

[Back to TOC](#table-of-contents)

#### Request URI

Another way for Janus to route a request to a given upstream service is to
specify a request URI via the `proxy.listen_path` property. To satisfy this field's
condition, a client request's URI **must** be prefixed with one of the values
of the `proxy.listen_path` field.

For example, in an API configured like this:

```json
{
    "name": "My API",
    "slug": "my-api",
    "proxy": {
        "listen_path": "/hello/*",
        "target_url": "http://my-api.com",
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

[Back to TOC](#table-of-contents)

##### The `strip_listen_path` property

It may be desirable to specify a URI prefix to match an API, but not
include it in the upstream request. To do so, use the `strip_listen_path` boolean
property by configuring an API like this:

```json
{
    "name": "My API",
    "slug": "my-api",
    "proxy": {
        "strip_listen_path" : true,
        "listen_path": "/service/*",
        "target_url": "http://my-api.com",
        "methods": ["GET"]
    }
}
```

Enabling this flag instructs Janus that when proxying this API, it should **not**
include the matching URI prefix in the upstream request's URI. For example, the
following client's request to the API configured as above:

```http
GET /service/path/to/resource HTTP/1.1
Host: my-api.com
```

Will cause Janus to send the following request to your upstream service:

```http
GET /path/to/resource HTTP/1.1
Host: my-api.com
```

[Back to TOC](#table-of-contents)

##### The `append_listen_path` property

You might also want to always append the `listen_path` to the upstream `target_url`. 
To do so, use the `append_listen_path` boolean property by configuring an API like this:

```json
{
    "name": "My API",
    "slug": "my-api",
    "proxy": {
        "append_listen_path" : true,
        "listen_path": "/service/*",
        "target_url": "http://my-api.com/example",
    }
}
```

Enabling this flag instructs Janus that when proxying this API, it should **always**
include the matching URI prefix in the upstream request's URI. For example, the
following client's request to the API configured as above:

```http
GET /service/path/to/resource HTTP/1.1
Host: my-api.com
```

Will cause Janus to send the following request to your upstream service:

```http
GET /example/service/path/to/resource HTTP/1.1
Host: my-api.com
```

[Back to TOC](#table-of-contents)

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
    "slug": "my-api",
    "proxy": {
        "strip_listen_path" : true,
        "listen_path": "/hello/*",
        "target_url": "http://my-api.com",
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

[Back to TOC](#table-of-contents)

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
    "slug": "api-1",
    "proxy": {
        "listen_path": "/",
        "target_url": "http://my-api.com",
        "hosts": ["example.com"]
    }
},
{
    "name": "API 2",
    "slug": "api-2",
    "proxy": {
        "listen_path": "/",
        "target_url": "http://my-api.com",
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

[Back to TOC](#table-of-contents)

### Conclusion

Through this guide, we hope you gained knowledge of the underlying proxying
mechanism of Janus, from how is a request matched to an API, to how to allow for
using the WebSocket protocol or setup SSL for an API.

[Back to TOC](#table-of-contents)
